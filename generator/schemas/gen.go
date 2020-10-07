package schemas

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"path/filepath"

	"go/ast"

	"encoding/json"

	"github.com/devfile/api/generator/genutils"
	"github.com/iancoleman/strcase"
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/controller-tools/pkg/crd"
	"sigs.k8s.io/controller-tools/pkg/loader"

	crdmarkers "sigs.k8s.io/controller-tools/pkg/crd/markers"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/markers"

	"github.com/coreos/go-semver/semver"
	"gomodules.xyz/orderedmap"
)

//go:generate go run sigs.k8s.io/controller-tools/cmd/helpgen generate:headerFile=../header.go.txt,year=2020 paths=.

// +controllertools:marker:generateHelp:category=Devfile

// GenerateJSONSchema drives whether a Json schema should be generated from this GO Struct type
type GenerateJSONSchema struct {

	// OmitCustomUnionMembers indicates that the Json schema gnerated from this type should omit Custom union members.
	OmitCustomUnionMembers bool `marker:",optional"`

	// Title indicates the content ot the Json Schema `title` attribute
	Title string `marker:",optional"`
}

var (
	jsonschemaVersionMarker  = markers.Must(markers.MakeDefinition("devfile:jsonschema:version", markers.DescribesPackage, ""))
	jsonschemaGenerateMarker = markers.Must(markers.MakeDefinition("devfile:jsonschema:generate", markers.DescribesType, GenerateJSONSchema{}))
)

// +controllertools:marker:generateHelp

// Generator generates JSON schemas from the GO source code of the Kubernetes API
//
// A JSON Schema is generated for each GO structure that had the `devfile:jsonschema:generate` annotation.
// The semver-compatible version of JSON Schemas is defined by the `devfile:jsonschema:version` annotation on the package. Typically in the `doc.go` file.
type Generator struct{}

// RegisterMarkers registers the markers of the Generator
func (Generator) RegisterMarkers(into *markers.Registry) error {
	if err := markers.RegisterAll(into, jsonschemaVersionMarker, jsonschemaGenerateMarker); err != nil {
		return err
	}
	if err := crdmarkers.Register(into); err != nil {
		return err
	}
	into.AddHelp(jsonschemaGenerateMarker, GenerateJSONSchema{}.Help())
	into.AddHelp(jsonschemaVersionMarker,
		markers.SimpleHelp("Devfile", "defines the semver-compatible version of the Json schemas that will be generated from the K8S API"))
	return genutils.RegisterUnionMarkers(into)
}

// Generate generates the artifacts
func (g Generator) Generate(ctx *genall.GenerationContext) error {
	parser := &crd.Parser{
		Collector:           ctx.Collector,
		Checker:             ctx.Checker,
		AllowDangerousTypes: false,
	}

	crd.AddKnownTypes(parser)
	metav1PackageOverride := parser.PackageOverrides["k8s.io/apimachinery/pkg/apis/meta/v1"]
	parser.PackageOverrides["k8s.io/apimachinery/pkg/apis/meta/v1"] = func(p *crd.Parser, pkg *loader.Package) {
		metav1PackageOverride(p, pkg)
		delete(p.Schemata, crd.TypeIdent{Name: "ObjectMeta", Package: pkg})
	}

	for _, root := range ctx.Roots {
		ctx.Checker.Check(root, func(node ast.Node) bool {
			// ignore interfaces
			_, isIface := node.(*ast.InterfaceType)
			return !isIface
		})

		root.NeedTypesInfo()

		parser.NeedPackage(root)

		unionDiscriminators := []markers.FieldInfo{}
		jsonschemaRequested := []*markers.TypeInfo{}

		if err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
			if info.Markers.Get(genutils.UnionMarker.Name) != nil {
				for _, field := range info.Fields {
					if field.Markers.Get(genutils.UnionDiscriminatorMarker.Name) != nil {
						unionDiscriminators = append(unionDiscriminators, field)
					}
				}
				return
			}
			if info.Markers.Get(jsonschemaGenerateMarker.Name) != nil {
				jsonschemaRequested = append(jsonschemaRequested, info)
				return
			}
		}); err != nil {
			root.AddError(err)
			return nil
		}

		if len(jsonschemaRequested) == 0 {
			return nil
		}

		var devfileSchemaVersion *semver.Version
		packageMarkers, err := markers.PackageMarkers(ctx.Collector, root)
		if err != nil {
			return err
		}
		switch markerValue := packageMarkers.Get(jsonschemaVersionMarker.Name).(type) {
		case string:
			devfileSchemaVersion, err = semver.NewVersion(markerValue)
			if err != nil {
				root.AddError(fmt.Errorf("In order to generate Json schemas from the K8S API, you should provide a valid semver-compatible devfile version in the +devfile:jsonschema:version comment marker of the K8S API package (in the doc.go file)"))
				return nil
			}
		default:
			root.AddError(fmt.Errorf("In order to generate Json schemas from the K8S API, you should annotate the K8S API package (in the doc.go file) with the +devfile:jsonschema:version comment marker"))
			return nil
		}

		for _, typeToProcess := range jsonschemaRequested {
			typeIdent := crd.TypeIdent{
				Package: root,
				Name:    typeToProcess.Name,
			}
			schemaGenerateMarker := typeToProcess.Markers.Get(jsonschemaGenerateMarker.Name).(GenerateJSONSchema)
			parser.NeedFlattenedSchemaFor(typeIdent)
			currentJSONSchema, found := parser.FlattenedSchemata[typeIdent]
			if !found {
				root.AddError(fmt.Errorf("Json schema for type " + typeIdent.Package.Name + "/" + typeIdent.Name + " could not be generated"))
				continue
			}

			fieldsToSkip := []string{}
			if schemaGenerateMarker.OmitCustomUnionMembers {
				fieldsToSkip = append(fieldsToSkip, "Custom")
			}

			genutils.AddUnionOneOfConstraints(&currentJSONSchema, unionDiscriminators, true, fieldsToSkip...)

			// Fix descriptions to have them Markdown compatible
			genutils.EditJSONSchema(&currentJSONSchema, func(schema *apiext.JSONSchemaProps) (newVisitor genutils.Visitor, stop bool) {
				if schema == nil || schema.Description == "" {
					return
				}
				regex, err := regexp.Compile(" \\n ([^-])")
				if err != nil {
					root.AddError(err)
					return
				}
				schema.Description = strings.ReplaceAll(schema.Description, " \t", "\n")
				schema.Description = strings.ReplaceAll(schema.Description, " \n - ", "\n- ")
				schema.Description = regex.ReplaceAllString(schema.Description, "\n\n$1")
				return
			})

			// Add the additionalProperies required to reflect the expected behavior from the K8S API,
			// (preserve-unknown-fields false by default unless on the Devfile metadata)
			genutils.EditJSONSchema(&currentJSONSchema, func(schema *apiext.JSONSchemaProps) (newVisitor genutils.Visitor, stop bool) {
				if schema == nil ||
					schema.Type != "object" ||
					schema.AdditionalProperties != nil {
					return
				}
				schema.AdditionalProperties = &apiext.JSONSchemaPropsOrBool{
					Allows: schema.XPreserveUnknownFields != nil && *schema.XPreserveUnknownFields,
				}
				return
			})

			// Remove Kubernetes extensions from the generated Json Schema
			genutils.EditJSONSchema(&currentJSONSchema, func(schema *apiext.JSONSchemaProps) (newVisitor genutils.Visitor, stop bool) {
				if schema == nil {
					return
				}
				schema.XEmbeddedResource = false
				schema.XPreserveUnknownFields = nil
				return
			})

			if schemaGenerateMarker.Title == "" {
				schemaGenerateMarker.Title = typeToProcess.Name + " schema - Version " + devfileSchemaVersion.String()
			}

			(&currentJSONSchema).Title = schemaGenerateMarker.Title

			jsonSchema, err := json.MarshalIndent(&currentJSONSchema, "", "  ")
			if err != nil {
				return err
			}

			genutils.EditJSONSchema(
				&currentJSONSchema,
				func(schema *apiext.JSONSchemaProps) (newVisitor genutils.Visitor, stop bool) {
					if schema != nil {
						schema.Default = nil
					}
					return
				})

			ideTargetedSchemasExplanation :=
				`IDE-targeted variants of the schemas provide the following difference compared to the main schemas:
- They contain additional non-standard ` + "`markdownDescription`" + ` attributes that are used by IDEs such a VSCode
to provide markdown-rendered documentation hovers. 
- They don't contain ` + "`default`" + ` attributes, since this triggers unwanted addition of defaulted fields during completion in IDEs.`

			(&currentJSONSchema).Title = (&currentJSONSchema).Title + " - IDE-targeted variant"
			(&currentJSONSchema).Description = (&currentJSONSchema).Description + "\n\n" + ideTargetedSchemasExplanation

			ideTargetedJsonSchema, err := json.MarshalIndent(&currentJSONSchema, "", "  ")
			if err != nil {
				return err
			}
			ideTargetedJsonSchemaMap := orderedmap.New()
			json.Unmarshal(ideTargetedJsonSchema, ideTargetedJsonSchemaMap)
			addMarkdownDescription(ideTargetedJsonSchemaMap)
			ideTargetedJsonSchema, err = json.MarshalIndent(ideTargetedJsonSchemaMap, "", "  ")

			schemaBaseName := strcase.ToKebab(typeToProcess.Name)
			schemaFolder := "latest"
			folderForIdeTargetedSchemas := filepath.Join(schemaFolder, "ide-targeted")
			schemaFileName := schemaBaseName + ".json"
			err = writeFile(ctx, schemaFolder, schemaFileName, jsonSchema)
			if err != nil {
				root.AddError(err)
				return nil
			}
			err = writeFile(ctx, folderForIdeTargetedSchemas, "Readme.md", []byte(ideTargetedSchemasExplanation))
			if err != nil {
				root.AddError(err)
				return nil
			}
			err = writeFile(ctx, folderForIdeTargetedSchemas, schemaFileName, ideTargetedJsonSchema)
			if err != nil {
				root.AddError(err)
				return nil
			}
			err = writeFile(ctx, schemaFolder, "jsonSchemaVersion.txt", []byte(devfileSchemaVersion.String()))
			if err != nil {
				root.AddError(err)
				return nil
			}
			err = writeFile(ctx, schemaFolder, "k8sApiVersion.txt", []byte(root.Name))
			if err != nil {
				root.AddError(err)
				return nil
			}
		}
	}

	return nil
}

func writeFile(ctx *genall.GenerationContext, schemaFolder, schemaFileName string, jsonSchema []byte) error {
	err := doWriteFile(ctx, schemaFolder, schemaFileName, jsonSchema)
	if pathError, isPathError := err.(*os.PathError); isPathError &&
		errors.Is(err, os.ErrNotExist) {
		os.MkdirAll(filepath.Dir(pathError.Path), os.ModePerm)
		err = writeFile(ctx, schemaFolder, schemaFileName, jsonSchema)
	}
	return err
}

func doWriteFile(ctx *genall.GenerationContext, schemaFolder, schemaFileName string, jsonSchema []byte) error {
	writer, err := ctx.Open(nil, filepath.Join(schemaFolder, schemaFileName))
	if err != nil {
		return err
	}
	defer writer.Close()
	_, err = writer.Write(jsonSchema)
	return err
}

func addMarkdownDescription(orderedMap *orderedmap.OrderedMap) {
	if orderedMap == nil {
		return
	}
	description, descriptionExists := orderedMap.Get("description")
	if descriptionExists {
		orderedMap.Set("markdownDescription", description)
	}

	theType, typeExists := orderedMap.Get("type")
	if typeExists {
		switch theType {
		case "object":
			walkProperties(getChildMap(orderedMap, "properties"))
		case "array":
			addMarkdownDescription(getChildMap(orderedMap, "items"))
		}
	}

}

func walkProperties(properties *orderedmap.OrderedMap) {
	if properties == nil {
		return
	}
	for _, propertyName := range properties.Keys() {
		propIf, exists := properties.Get(propertyName)
		if exists {
			if property, isOrderedMap := propIf.(*orderedmap.OrderedMap); isOrderedMap {
				addMarkdownDescription(property)
			}
		}
	}
}

func getChildMap(orderedMap *orderedmap.OrderedMap, childName string) *orderedmap.OrderedMap {
	childIf, childExists := orderedMap.Get(childName)
	if childExists {
		if child, isOrderedMap := childIf.(*orderedmap.OrderedMap); isOrderedMap {
			return child
		}
	}
	return nil
}
