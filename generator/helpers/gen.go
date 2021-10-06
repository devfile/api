package helpers

import (
	"bytes"
	"fmt"
	"github.com/devfile/api/generator/genutils"
	"github.com/elliotchance/orderedmap"
	"go/ast"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

//go:generate go run sigs.k8s.io/controller-tools/cmd/helpgen generate:headerFile=../header.go.txt,year=2020 paths=.

var (
	// HelperTypeMarker is associated with a type that's used as a parameter to the helper function
	HelperTypeMarker = markers.Must(markers.MakeDefinition("devfile:helper:generate", markers.DescribesType, struct{}{}))
	// DefaultFieldMarker is associated with a boolean pointer field to indicate what the default value should be
	DefaultFieldMarker = markers.Must(markers.MakeDefinition("devfile:default:value", markers.DescribesField, ""))
)

// +controllertools:marker:generateHelp

// Generator generates helper functions that are used to return values for boolean pointer fields.
//
// The helper function takes as a parameter, the `devfile:default:generate` annotated type and returns the value of the
// field if it's been set, otherwise it will return the default value specified by the devfile:default:value annotation.
type Generator struct{}

// RegisterMarkers registers the markers of the Generator
func (Generator) RegisterMarkers(into *markers.Registry) error {
	if err := markers.RegisterAll(into, HelperTypeMarker, DefaultFieldMarker); err != nil {
		return err
	}
	into.AddHelp(HelperTypeMarker,
		markers.SimpleHelp("Devfile", "indicates the type that's passed in as a parameter to the generated helper functions"))
	into.AddHelp(DefaultFieldMarker,
		markers.SimpleHelp("Devfile", "indicates the default value of a boolean pointer field"))
	return genutils.RegisterUnionMarkers(into)

}

func (Generator) CheckFilter() loader.NodeFilter {
	return func(node ast.Node) bool {
		// ignore interfaces
		_, isIface := node.(*ast.InterfaceType)
		return !isIface
	}
}

// helperInfo stores the info to generate the helper function
type helperInfo struct {
	funcName   string
	defaultVal string
	returnType string
}

// Generate generates the artifacts
func (g Generator) Generate(ctx *genall.GenerationContext) error {
	for _, root := range ctx.Roots {
		ctx.Checker.Check(root)
		root.NeedTypesInfo()

		typesToProcess := orderedmap.NewOrderedMap()
		if err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
			if info.Markers.Get(HelperTypeMarker.Name) != nil {
				var helpers []helperInfo
				for _, field := range info.Fields {
					defaultVal := field.Markers.Get(DefaultFieldMarker.Name)
					if defaultVal != nil {
						//look for boolean pointers
						if ptr, isPtr := field.RawField.Type.(*ast.StarExpr); isPtr {
							if ident, ok := ptr.X.(*ast.Ident); ok {
								if ident.Name == "bool" {
									helpers = append(helpers, helperInfo{
										field.Name,
										defaultVal.(string),
										ident.Name,
									})
								} else {
									root.AddError(fmt.Errorf("devfile:default:value marker is specified on %s/%s which is not a boolean pointer", info.Name, field.Name))
								}
							}
						} else {
							root.AddError(fmt.Errorf("devfile:default:value marker is specified on %s/%s which is not a boolean pointer", info.Name, field.Name))
						}

					}
				}
				if len(helpers) > 0 {
					typesToProcess.Set(info, helpers)
				} else {
					root.AddError(fmt.Errorf("type %s does not have the field marker, devfile:default:value specified on a boolean pointer field", info.Name))
				}
				return
			}

		}); err != nil {
			root.AddError(err)
			return nil
		}

		genutils.WriteFormattedSourceFile("helpers", ctx, root, func(buf *bytes.Buffer) {
			for elt := typesToProcess.Front(); elt != nil; elt = elt.Next() {
				param := elt.Key.(*markers.TypeInfo)
				fields := elt.Value.([]helperInfo)
				for _, helper := range fields {
					fName := helper.funcName
					returnType := helper.returnType

					helperFunction := `
//` + fName + ` returns the value of the boolean property.  If unset, it's the default value specified in the devfile:default:value marker` + `
func ` + fName + `(in *` + param.Name + `) ` + returnType + ` { 
	if in.` + fName + ` !=nil{ 
		return *in.` + fName + `}`

					defaultVal := helper.defaultVal
					if fName != "MountSources" {
						buf.WriteString(helperFunction + `
							return ` + defaultVal + ` }`)
					} else {
						buf.WriteString(helperFunction + `else {
							if DedicatedPod(in) { return false } 
								return ` + defaultVal + ` }}`)
					}
				}
			}
		})

	}

	return nil
}
