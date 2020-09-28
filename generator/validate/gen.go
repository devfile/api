package validate

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/devfile/api/generator/genutils"
	crdmarkers "sigs.k8s.io/controller-tools/pkg/crd/markers"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

//go:generate go run sigs.k8s.io/controller-tools/cmd/helpgen generate:headerFile=../header.go.txt,year=2020 paths=.

// +controllertools:marker:generateHelp

// Generator validates the consistency of the API GO code.
//
// Validaty check are related to unions, patchStrategy, and optional fields.
type Generator struct{}

// RegisterMarkers registers the markers of the Generator
func (Generator) RegisterMarkers(into *markers.Registry) error {
	err := genutils.RegisterUnionMarkers(into)
	if err != nil {
		return err
	}
	return crdmarkers.Register(into)
}

// Generate generates the artifacts
func (g Generator) Generate(ctx *genall.GenerationContext) error {
	for _, root := range ctx.Roots {
		ctx.Checker.Check(root, func(node ast.Node) bool {
			// ignore interfaces
			_, isIface := node.(*ast.InterfaceType)
			return !isIface
		})

		root.NeedTypesInfo()

		packageTypes := map[string]*markers.TypeInfo{}
		if err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
			packageTypes[info.Name] = info
		}); err != nil {
			root.AddError(err)
			return nil
		}

		for _, typeToCheck := range packageTypes {
			if typeToCheck.Markers.Get(genutils.UnionMarker.Name) != nil {
				// Check Unions
				unionMembers := []string{}
				var unionDiscriminatorField *markers.FieldInfo
				for index, unionField := range typeToCheck.Fields {
					if unionField.Markers.Get(genutils.UnionDiscriminatorMarker.Name) != nil {
						if unionDiscriminatorField != nil {
							root.AddError(loader.ErrFromNode(fmt.Errorf(
								"Union `%v` should have only 1 union discriminator, but has 2: `%v` and `%v`",
								typeToCheck.Name,
								unionDiscriminatorField.Name,
								unionField.Name), unionField.RawField))
						}
						unionDiscriminatorField = &(typeToCheck.Fields[index])
					} else {
						unionMembers = append(unionMembers, unionField.Name)
					}
				}
				if unionDiscriminatorField == nil {
					root.AddError(loader.ErrFromNode(fmt.Errorf(
						"Union `%v` should have 1 union discriminator. See here for details: https://github.com/kubernetes/enhancements/blob/master/keps/sig-api-machinery/20190325-unions.md#proposal",
						typeToCheck.Name), typeToCheck.RawSpec))
				} else {
					if unionDiscriminatorField.Markers.Get("optional") == nil {
						root.AddError(loader.ErrFromNode(fmt.Errorf(
							"In union `%v` the union discriminator `%v` should have the `+optional` comment marker",
							typeToCheck.Name,
							unionDiscriminatorField.Name), unionDiscriminatorField.RawField))
					}
					if !strings.Contains(unionDiscriminatorField.Tag.Get("json"), ",omitempty") {
						root.AddError(loader.ErrFromNode(fmt.Errorf(
							"In union `%v` the union discriminator `%v` should contain the `omitempty` option in its `json` tag, since it is expected to be an optional field",
							typeToCheck.Name,
							unionDiscriminatorField.Name), unionDiscriminatorField.RawField))
					}

					wrongTypeError := loader.ErrFromNode(fmt.Errorf(
						"In union `%v` the union discriminator `%v` should have a `string` type, or a type (defined in the same package) whose underlying type is a string",
						typeToCheck.Name,
						unionDiscriminatorField.Name), unionDiscriminatorField.RawField)

					discriminatorTypeRef, typeFound := root.TypesInfo.Types[unionDiscriminatorField.RawField.Type]
					if !typeFound {
						root.AddError(wrongTypeError)
					}

					underlyingType, underlyingTypeIsBasic := discriminatorTypeRef.Type.Underlying().(*types.Basic)
					if !underlyingTypeIsBasic || underlyingType.Kind() != types.String {
						root.AddError(wrongTypeError)
					}

					getEnumValidation := func(markers markers.MarkerValues, node ast.Node) *crdmarkers.Enum {
						enumMarkerIf := markers.Get("kubebuilder:validation:Enum")
						if enumMarkerIf != nil {
							if enumMarker, isEnum := enumMarkerIf.(crdmarkers.Enum); isEnum {
								return &enumMarker
							}
							root.AddError(loader.ErrFromNode(fmt.Errorf(
								"In union `%v` the union discriminator `%v` should only have the kubebuilder marker that defines the allowed union values as an Enum (kubebuilder:validation:Enum=....)",
								typeToCheck.Name,
								unionDiscriminatorField.Name), node))
						}
						return nil
					}

					enumMarkerPtr := getEnumValidation(unionDiscriminatorField.Markers, unionDiscriminatorField.RawField)
					var enumAstNode ast.Node = unionDiscriminatorField.RawField

					if _, isBasic := discriminatorTypeRef.Type.(*types.Basic); !isBasic {
						if typeNamed, isNamed := discriminatorTypeRef.Type.(*types.Named); isNamed {
							enumTypeInfo, isInSamePackage := packageTypes[typeNamed.Obj().Name()]
							if !isInSamePackage {
								root.AddError(loader.ErrFromNode(fmt.Errorf(
									"In union `%v` the union discriminator `%v` should have a string-based type defined in the same package",
									typeToCheck.Name,
									unionDiscriminatorField.Name), unionDiscriminatorField.RawField))
							} else {
								typeEnumMarkerPtr := getEnumValidation(enumTypeInfo.Markers, enumTypeInfo.RawSpec)
								if typeEnumMarkerPtr != nil {
									if enumMarkerPtr != nil {
										root.AddError(loader.ErrFromNode(fmt.Errorf(
											"Type `%v` should not define a `kubebuilder:validation:Enum` annotation, since it is already defined in the discriminator `%v` of union `%v`",
											enumTypeInfo.Name,
											unionDiscriminatorField.Name,
											typeToCheck.Name), enumTypeInfo.RawSpec))
									} else {
										enumMarkerPtr = typeEnumMarkerPtr
										enumAstNode = enumTypeInfo.RawSpec
									}
								}
							}
						}
					}

					if enumMarkerPtr == nil {
						root.AddError(loader.ErrFromNode(fmt.Errorf(
							"In union `%v` the union discriminator `%v` should specify the allowed union values through an annotation (kubebuilder:validation:Enum=....)",
							typeToCheck.Name,
							unionDiscriminatorField.Name), unionDiscriminatorField.RawField))
					} else {
						for index, enumValue := range *enumMarkerPtr {
							stringValue, isString := enumValue.(string)
							if !isString || stringValue != unionMembers[index] {
								root.AddError(loader.ErrFromNode(fmt.Errorf(
									"In union `%v` the union discriminator `%v` should allow the following values through the `kubebuilder:validation:Enum` anntation: "+strings.Join(unionMembers, ";"),
									typeToCheck.Name,
									unionDiscriminatorField.Name), enumAstNode))
								break
							}
						}
					}
				}
			}
		}

	}

	return nil
}
