//
//
// Copyright Red Hat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package schemas

import (
	"errors"
	"fmt"
	"go/ast"
	"os"
	"regexp"
	"strings"

	"path/filepath"

	"encoding/json"

	"github.com/devfile/api/generator/genutils"
	"github.com/iancoleman/strcase"
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-tools/pkg/crd"
	"sigs.k8s.io/controller-tools/pkg/loader"

	crdmarkers "sigs.k8s.io/controller-tools/pkg/crd/markers"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/markers"

	"github.com/coreos/go-semver/semver"
	"gomodules.xyz/orderedmap"
)

//go:generate go run -mod=mod sigs.k8s.io/controller-tools/cmd/helpgen@v0.6.2 generate:headerFile=../header.go.txt,year=2020 paths=.

// +controllertools:marker:generateHelp:category=Devfile

// GenerateJSONSchema drives whether a Json schema should be generated from this GO Struct type
type GenerateJSONSchema struct {

	// OmitCustomUnionMembers indicates that the Json schema generated from this type should omit Custom union members.
	OmitCustomUnionMembers bool `marker:",optional"`

	// OmitPluginUnionMembers indicates that the Json schema generated from this type should omit Plugin component union members.
	OmitPluginUnionMembers bool `marker:",optional"`

	ShortenEndpointNameLength bool `marker:",optional"`

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

type toGenerate struct {
	groupName            string
	version              string
	devfileSchemaVersion *semver.Version
	unionDiscriminators  []markers.FieldInfo
	jsonschemaRequested  []*markers.TypeInfo
}

func (Generator) CheckFilter() loader.NodeFilter {
	return func(node ast.Node) bool {
		// ignore interfaces
		_, isIface := node.(*ast.InterfaceType)
		return !isIface
	}
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

	toGenerateByPackage := map[*loader.Package]toGenerate{}

	apiVersionsByAPIGroup := map[string][]string{}
	schemaVersionsByGV := map[schema.GroupVersion]*semver.Version{}
	packageByGV := map[schema.GroupVersion]*loader.Package{}

	for _, root := range ctx.Roots {
		forRoot := toGenerate{
			version: root.Name,
		}

		ctx.Checker.Check(root)

		root.NeedTypesInfo()

		parser.NeedPackage(root)

		if err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
			if info.Markers.Get(genutils.UnionMarker.Name) != nil {
				for _, field := range info.Fields {
					if field.Markers.Get(genutils.UnionDiscriminatorMarker.Name) != nil {
						forRoot.unionDiscriminators = append(forRoot.unionDiscriminators, field)
					}
				}
				return
			}
			if info.Markers.Get(jsonschemaGenerateMarker.Name) != nil {
				forRoot.jsonschemaRequested = append(forRoot.jsonschemaRequested, info)
				return
			}
		}); err != nil {
			root.AddError(err)
			return nil
		}

		if len(forRoot.jsonschemaRequested) == 0 {
			continue
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
			forRoot.devfileSchemaVersion = devfileSchemaVersion
		default:
			root.AddError(fmt.Errorf("In order to generate Json schemas from the K8S API, you should annotate the K8S API package (in the doc.go file) with the +devfile:jsonschema:version comment marker"))
			return nil
		}

		switch groupName := packageMarkers.Get("groupName").(type) {
		case string:
			forRoot.groupName = groupName
		default:
			root.AddError(fmt.Errorf("The package should have a valid 'groupName' annotation"))
			return nil
		}
		groupVersion := schema.GroupVersion{
			Group:   forRoot.groupName,
			Version: forRoot.version,
		}

		apiVersionsByAPIGroup[groupVersion.Group] = append(apiVersionsByAPIGroup[groupVersion.Group], groupVersion.Version)
		schemaVersionsByGV[groupVersion] = forRoot.devfileSchemaVersion
		packageByGV[groupVersion] = root
		toGenerateByPackage[root] = forRoot
	}

	for groupName, apiVersions := range apiVersionsByAPIGroup {
		genutils.SortKubeLikeVersion(apiVersions)

		var currentSchemaVersion *semver.Version
		var currentAPIVersion string
		for _, apiVersion := range apiVersions {
			groupVersion := schema.GroupVersion{Group: groupName, Version: apiVersion}
			schemaVersion := schemaVersionsByGV[groupVersion]
			if currentSchemaVersion != nil && schemaVersion != nil {
				if schemaVersion.Compare(*currentSchemaVersion) <= 0 {
					packageByGV[groupVersion].AddError(
						fmt.Errorf(`The schema versions should be incremented on each increment of the corresponding K8S apiVersion.
This is not the case in the "%s' API group:
  '%s' K8S apiVersion => '%s' Json schema version
  '%s' K8S apiVersion => '%s' Json schema version`,
							groupName, currentAPIVersion, currentSchemaVersion.String(), apiVersion, schemaVersion.String()))
					return nil
				}
			}
			currentAPIVersion = apiVersion
			currentSchemaVersion = schemaVersion
		}
	}

	for root, toDo := range toGenerateByPackage {
		for _, typeToProcess := range toDo.jsonschemaRequested {
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
			if schemaGenerateMarker.OmitPluginUnionMembers {
				fieldsToSkip = append(fieldsToSkip, "Plugin")
			}

			genutils.AddUnionOneOfConstraints(&currentJSONSchema, toDo.unionDiscriminators, true, fieldsToSkip...)

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

			// Add the additionalProperties required to reflect the expected behavior from the K8S API,
			// (preserve-unknown-fields false by default)
			genutils.EditJSONSchema(&currentJSONSchema, func(schema *apiext.JSONSchemaProps) (newVisitor genutils.Visitor, stop bool) {
				if schema == nil ||
					schema.Type != "object" ||
					schema.AdditionalProperties != nil {
					return
				}
				schema.AdditionalProperties = &apiext.JSONSchemaPropsOrBool{
					// Allows when schema does not describe any property or has preserveUnknownFields
					Allows: len(schema.Properties) == 0 ||
						schema.XPreserveUnknownFields != nil && *schema.XPreserveUnknownFields,
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
				schemaGenerateMarker.Title = typeToProcess.Name + " schema - Version " + toDo.devfileSchemaVersion.String()
			}

			(&currentJSONSchema).Title = schemaGenerateMarker.Title

			// Update endpoint name length limit to 15 chars in devfile spec, if ShortenEndpointNameLength is specified
			// To fix issue: https://github.com/devfile/api/issues/700, but also to hold backward compatibility for devworkspace
			genutils.EditJSONSchema(&currentJSONSchema, func(schema *apiext.JSONSchemaProps) (newVisitor genutils.Visitor, stop bool) {
				if schema == nil {
					return
				}
				if schema.Type != "object" {
					return
				}
				if len(schema.Properties) == 0 {
					return
				}
				for propName, prop := range schema.Properties {
					if propName == "endpoints" && prop.Items != nil && prop.Items.Schema != nil {
						for endpointPropName, endpointProp := range prop.Items.Schema.Properties {
							if endpointPropName == "name" {
								if schemaGenerateMarker.ShortenEndpointNameLength {
									*endpointProp.MaxLength = int64(15)
									break
								}
							}
						}
					}
				}
				return
			})

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
			err = json.Unmarshal(ideTargetedJsonSchema, ideTargetedJsonSchemaMap)
			if err != nil {
				return err
			}
			addMarkdownDescription(ideTargetedJsonSchemaMap)
			ideTargetedJsonSchema, err = json.MarshalIndent(ideTargetedJsonSchemaMap, "", "  ")

			schemaBaseName := strcase.ToKebab(typeToProcess.Name)
			schemaFolder := "latest"
			if toDo.version != genutils.LatestKubeLikeVersion(apiVersionsByAPIGroup[toDo.groupName]) {
				schemaFolder = toDo.version
			}
			folderForIdeTargetedSchemas := filepath.Join(schemaFolder, "ide-targeted")
			schemaFileName := schemaBaseName + ".json"
			err = writeFile(ctx, schemaFolder, schemaFileName, jsonSchema)
			if err != nil {
				root.AddError(err)
				return nil
			}
			err = writeFile(ctx, folderForIdeTargetedSchemas, "README.md", []byte(ideTargetedSchemasExplanation))
			if err != nil {
				root.AddError(err)
				return nil
			}
			err = writeFile(ctx, folderForIdeTargetedSchemas, schemaFileName, ideTargetedJsonSchema)
			if err != nil {
				root.AddError(err)
				return nil
			}
			err = writeFile(ctx, schemaFolder, "jsonSchemaVersion.txt", []byte(toDo.devfileSchemaVersion.String()))
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
