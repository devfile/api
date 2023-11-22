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

package overrides

import (
	"bytes"
	"fmt"
	"go/ast"

	"go/printer"
	"go/token"
	"io"
	"regexp"
	"strings"

	"golang.org/x/tools/go/ast/astutil"

	"github.com/devfile/api/generator/genutils"
	"github.com/go-toolsmith/astcopy"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"

	"github.com/elliotchance/orderedmap"
)

//go:generate go run -mod=mod sigs.k8s.io/controller-tools/cmd/helpgen@v0.6.2 generate:headerFile=../header.go.txt,year=2020 paths=.

// +controllertools:marker:generateHelp:category=Overrides

// FieldOverridesInclude drives whether a field should be overriden in devfile parent or plugins
type FieldOverridesInclude struct {
	// Omit indicates that this field cannot be overridden at all.
	Omit bool `marker:",optional"`
	// OmmitInPlugin indicates that this field cannot be overridden in a devfile plugin.
	OmitInPlugin bool `marker:",optional"`
	// Description indicates the description that should be added as Go documentation on the generated structs.
	Description string `marker:",optional"`
}

var (
	overridesFieldMarker = markers.Must(markers.MakeDefinition("devfile:overrides:include", markers.DescribesField, FieldOverridesInclude{}))
	overridesTypeMarker  = markers.Must(markers.MakeDefinition("devfile:overrides:generate", markers.DescribesType, struct{}{}))
)

// +controllertools:marker:generateHelp

// Generator generates additional GO code for the overriding of elements in devfile parent or plugins.
type Generator struct {

	// IsForPluginOverrides indicates that the generated code should be done for plugin overrides.
	// When false, the parent overrides are generated
	IsForPluginOverrides bool `marker:"isForPluginOverrides,optional"`

	suffix            string
	rootTypeToProcess typeToProcess
}

// RegisterMarkers registers the markers of the Generator
func (Generator) RegisterMarkers(into *markers.Registry) error {
	if err := markers.RegisterAll(into, overridesFieldMarker, overridesTypeMarker); err != nil {
		return err
	}
	into.AddHelp(overridesFieldMarker, FieldOverridesInclude{}.Help())
	into.AddHelp(overridesTypeMarker, markers.SimpleHelp("Overrides", "indicates that a type should be selected to create Overrides for it"))
	return genutils.RegisterUnionMarkers(into)
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
	for _, root := range ctx.Roots {
		ctx.Checker.Check(root)

		root.NeedTypesInfo()

		var rootStructToOverride *markers.TypeInfo
		packageTypes := map[string]*markers.TypeInfo{}
		if err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
			if info.Markers.Get(overridesTypeMarker.Name) != nil {
				if rootStructToOverride == nil {
					rootStructToOverride = info
				} else {
					root.AddError(fmt.Errorf("Marker %v should be added to only one Struct type, but was added on %v and %v",
						overridesTypeMarker.Name,
						rootStructToOverride.Name,
						info.Name,
					))
				}
			}
			packageTypes[info.RawSpec.Name.Name] = info
		}); err != nil {
			root.AddError(err)
			return nil
		}

		if rootStructToOverride == nil {
			root.AddError(fmt.Errorf("Marker %v should be added to at least one Struct type",
				overridesTypeMarker.Name,
			))
			return nil
		}

		config := printer.Config{
			Tabwidth: 2,
			Mode:     printer.UseSpaces,
		}

		g.suffix = "ParentOverride"
		if g.IsForPluginOverrides {
			g.suffix = "PluginOverride"
		}
		g.rootTypeToProcess = typeToProcess{
			OverrideTypeName: g.suffix + "s",
			TypeInfo:         rootStructToOverride,
			MandatoryKey:     "",
		}

		overrides := g.process(root, packageTypes)

		fileNamePart := "parent_overrides"
		if g.IsForPluginOverrides {
			fileNamePart = "plugin_overrides"
		}

		genutils.WriteFormattedSourceFile(fileNamePart, ctx, root, func(buf *bytes.Buffer) {
			buf.WriteString(`
import (
	attributes "github.com/devfile/api/v2/pkg/attributes"
)

`)
			config.Fprint(buf, root.Fset, overrides)
			buf.WriteString(`
func (overrides ` + g.rootTypeToProcess.OverrideTypeName + `) isOverride() {}
`)
		})
	}

	return nil
}

// typeToProcess contains all required information about the how to process a given type.
// A list of `typeToProcess` instances can be returned the `createOverride` function
// since it is when processing a type that we possibly encounter new types to process.
type typeToProcess struct {
	OverrideTypeName   string
	TypeInfo           *markers.TypeInfo
	MandatoryKey       string
	DropEnumAnnotation bool
}

func (g Generator) process(root *loader.Package, packageTypes map[string]*markers.TypeInfo) []ast.Decl {
	toProcess := []typeToProcess{g.rootTypeToProcess}
	processed := orderedmap.NewOrderedMap()
	for len(toProcess) > 0 {
		nextOne := toProcess[0]
		toProcess = toProcess[1:]
		if _, isAlreadyProcessed := processed.Get(nextOne.TypeInfo.Name); isAlreadyProcessed &&
			nextOne.MandatoryKey == "" {
			continue
		}

		newOverride, newTypesToOverride, errors := g.createOverride(nextOne, packageTypes)
		processed.Set(nextOne.TypeInfo.Name, newOverride)
		for _, err := range errors {
			root.AddError(loader.ErrFromNode(err, nextOne.TypeInfo.RawSpec))
		}
		toProcess = append(toProcess, newTypesToOverride...)
	}

	overrides := []ast.Decl{}
	for elt := processed.Front(); elt != nil; elt = elt.Next() {
		overrides = append(overrides, elt.Value.(ast.Decl))
	}
	return overrides
}

// fieldChange provides the required information about how overrides generation should handle a given field
type fieldChange struct {
	fieldInfo      markers.FieldInfo
	overrideMarker FieldOverridesInclude
}

func (g Generator) createOverride(newTypeToProcess typeToProcess, packageTypes map[string]*markers.TypeInfo) (ast.Decl, []typeToProcess, []error) {
	errors := []error{}
	var alreadyFoundType *ast.TypeSpec = nil
	fieldChanges := map[token.Pos]fieldChange{}

	typeToOverride := newTypeToProcess.TypeInfo
	if typeToOverride.Fields != nil {
		for _, field := range typeToOverride.Fields {
			fieldPos := field.RawField.Pos()
			if !fieldPos.IsValid() {
				errors = append(errors,
					fmt.Errorf("Field %v in type %v doesn't have a valid position in the source file",
						field.Name,
						typeToOverride.Name,
					))
				continue
			}
			overridesMarker := FieldOverridesInclude{}
			if markerEntry := field.Markers.Get(overridesFieldMarker.Name); markerEntry != nil {
				overridesMarker = markerEntry.(FieldOverridesInclude)
			}
			fieldChanges[fieldPos] = fieldChange{
				fieldInfo:      field,
				overrideMarker: overridesMarker,
			}
		}
	}

	overrideGenDecl := astcopy.GenDecl(typeToOverride.RawDecl)
	if typeToOverride.Markers.Get(overridesTypeMarker.Name) != nil {
		overrideGenDecl.Doc = updateComments(overrideGenDecl, overrideGenDecl.Doc, `.*`, ` *\+`+overridesTypeMarker.Name+` *`)
	}
	if newTypeToProcess.DropEnumAnnotation {
		overrideGenDecl.Doc = updateComments(
			overrideGenDecl, overrideGenDecl.Doc,
			`.*`,
			` *`+regexp.QuoteMeta("+kubebuilder:validation:Enum=")+`.*`,
		)
	}

	overrideGenDecl.Doc = updateComments(
		overrideGenDecl, overrideGenDecl.Doc,
		`.*`,
		` *`+regexp.QuoteMeta("+devfile:getter:generate")+`.*`,
	)

	overrideGenDecl.Doc = updateComments(
		overrideGenDecl, overrideGenDecl.Doc,
		`.*`,
		` *`+regexp.QuoteMeta("+devfile:jsonschema:generate")+` *`,
	)

	if newTypeToProcess == g.rootTypeToProcess {
		overrideGenDecl.Doc = updateComments(
			overrideGenDecl, overrideGenDecl.Doc,
			`.*`,
			``,
			"+devfile:jsonschema:generate",
		)
	}

	moreTypesToAdd := []typeToProcess{}
	overrideGenDecl = astutil.Apply(overrideGenDecl,
		func(cursor *astutil.Cursor) bool {
			if typeSpec, isTypeSpec := cursor.Node().(*ast.TypeSpec); isTypeSpec {
				if alreadyFoundType != nil {
					errors = append(errors,
						fmt.Errorf("types %v and %v are defined in the same type definition - please avoid defining several types in the same type definition",
							alreadyFoundType.Name,
							typeSpec.Name,
						))
					return false
				}
				typeSpec.Name.Name = newTypeToProcess.OverrideTypeName
			}
			if astField, isField := cursor.Node().(*ast.Field); isField {
				if newTypeToProcess == g.rootTypeToProcess &&
					cursor.Index() == 0 {
					cursor.InsertBefore(&ast.Field{
						Type: &ast.Ident{Name: "OverridesBase"},
						Tag:  &ast.BasicLit{Kind: token.STRING, Value: "`json:\",inline\"`"},
					})
				}

				fieldChange := fieldChanges[cursor.Node().Pos()]
				field := fieldChange.fieldInfo

				overridesMarker := fieldChange.overrideMarker

				shouldSkip := func(overridesMarker FieldOverridesInclude) bool {
					if overridesMarker.Omit ||
						(overridesMarker.OmitInPlugin && g.IsForPluginOverrides) {
						return true
					}
					return false
				}

				if shouldSkip(overridesMarker) {
					cursor.Delete()
					return true
				}

				if overridesMarker.Description != "" {
					astField.Doc = updateComments(
						astField, astField.Doc,
						` *\+[^ ]+.*`,
						` *\+`+overridesFieldMarker.Name+`.*`,
						overridesMarker.Description,
						"Overriding is done according to K8S strategic merge patch standard rules.",
					)
				}

				if field.Name != newTypeToProcess.MandatoryKey {
					// Make the field optional by default, unless typeToProcess contains a MandatoryKey nonempty field
					jsonTag := field.Tag.Get("json")
					if jsonTag != "" &&
						jsonTag != "-" &&
						!strings.Contains(jsonTag, ",inline") &&
						!strings.Contains(jsonTag, ",omitempty") {
						newJSONTag := jsonTag + ",omitempty"
						astField.Tag.Value = strings.Replace(astField.Tag.Value, `json:"`+jsonTag+`"`, `json:"`+newJSONTag+`"`, 1)
						astField.Doc = updateComments(
							astField, astField.Doc,
							`.*`,
							` *`+regexp.QuoteMeta("+optional")+`.*`,
							` +optional`,
						)
					}
				}

				// Remove the `default` directives for overrides, since it doesn't make sense.
				astField.Doc = updateComments(
					astField, astField.Doc,
					`.*`,
					` *`+regexp.QuoteMeta("+kubebuilder:default")+` *=.*`,
				)

				//remove the +devfile:default:values for overrides
				astField.Doc = updateComments(
					astField, astField.Doc,
					`.*`,
					` *`+regexp.QuoteMeta("+devfile:default:value")+` *=.*`,
				)

				processFieldType := func(ident *ast.Ident) *typeToProcess {
					typeToOverride, existsInPackage := packageTypes[ident.Name]
					if !existsInPackage {
						return nil
					}
					ident.Name = ident.Name + g.suffix
					return &typeToProcess{
						OverrideTypeName: ident.Name,
						TypeInfo:         typeToOverride,
					}
				}

				var fieldTypeToProcess *typeToProcess

				switch fieldType := astField.Type.(type) {
				case *ast.ArrayType:
					switch elementType := fieldType.Elt.(type) {
					case *ast.Ident:
						fieldTypeToProcess = processFieldType(elementType)
						if fieldTypeToProcess != nil {
							fieldTypeToProcess.MandatoryKey = strings.Title(genutils.GetPatchMergeKey(&field))
						}
					}
				case *ast.Ident:
					fieldTypeToProcess = processFieldType(fieldType)
					if field.Markers.Get(genutils.UnionDiscriminatorMarker.Name) != nil {
						enumValues := []string{}
						for _, f := range typeToOverride.Fields {
							pos := f.RawField.Pos()
							fieldChange := fieldChanges[f.RawField.Pos()]
							if pos != cursor.Node().Pos() &&
								!shouldSkip(fieldChange.overrideMarker) {
								enumValues = append(enumValues, fieldChange.fieldInfo.Name)
							}
						}
						kubebuilderAnnotation := "+kubebuilder:validation:Enum=" + strings.Join(enumValues, ";")
						astField.Doc = updateComments(
							astField, astField.Doc,
							`.*`,
							` *`+regexp.QuoteMeta("+kubebuilder:validation:Enum=")+`.*`,
							kubebuilderAnnotation)
						if fieldTypeToProcess != nil {
							fieldTypeToProcess.DropEnumAnnotation = true
						}
					}
				case *ast.StarExpr:
					switch elementType := fieldType.X.(type) {
					case *ast.Ident:
						fieldTypeToProcess = processFieldType(elementType)
					}
				case *ast.MapType:
					switch elementType := fieldType.Key.(type) {
					case *ast.Ident:
						fieldTypeToProcess = processFieldType(elementType)
					}
					switch elementType := fieldType.Value.(type) {
					case *ast.Ident:
						fieldTypeToProcess = processFieldType(elementType)
					}
				default:
				}

				if fieldTypeToProcess != nil {
					moreTypesToAdd = append(moreTypesToAdd, *fieldTypeToProcess)
				}
			}
			return true
		},
		func(*astutil.Cursor) bool { return true },
	).(*ast.GenDecl)

	return overrideGenDecl, moreTypesToAdd, errors
}

// writeFormatted outputs the given code, after gofmt-ing it.  If we couldn't gofmt,
// we write the unformatted code for debugging purposes.
func (g Generator) writeOut(ctx *genall.GenerationContext, root *loader.Package, outBytes []byte) {
	fileToWrite := "zz_generated.parent_overrides.go"
	if g.IsForPluginOverrides {
		fileToWrite = "zz_generated.plugin_overrides.go"
	}
	outputFile, err := ctx.Open(root, fileToWrite)
	if err != nil {
		root.AddError(err)
		return
	}
	defer outputFile.Close()
	n, err := outputFile.Write(outBytes)
	if err != nil {
		root.AddError(err)
		return
	}
	if n < len(outBytes) {
		root.AddError(io.ErrShortWrite)
	}
}

// updateComments defines, through regexps, which comment lines should be kept and which should be dropped,
// It also provides additional comment lines that will be *prepended* to the existing comment lines.
// In both regexps and additional lines, the comment prefix `//` should be omitted.
func updateComments(commentedNode ast.Node, commentGroup *ast.CommentGroup, keepRegexp string, dropRegexp string, additionalLines ...string) *ast.CommentGroup {
	if commentGroup == nil {
		commentGroup = &ast.CommentGroup{}
	}
	commentsToKeep := []*ast.Comment{}
	for _, comment := range commentGroup.List {
		if keep, _ := regexp.MatchString(`^ *//`+keepRegexp+`$`, comment.Text); keep {
			if drop, _ := regexp.MatchString(`^ *//`+dropRegexp+`$`, comment.Text); !drop {
				comment.Slash = token.NoPos
				commentsToKeep = append(commentsToKeep, comment)
			}
		}
	}
	commentGroup.List = []*ast.Comment{}
	for _, line := range additionalLines {
		commentGroup.List = append(commentGroup.List, &ast.Comment{Text: "// " + line})
	}
	commentGroup.List = append(commentGroup.List, commentsToKeep...)
	if len(commentGroup.List) > 0 {
		commentGroup.List[0].Slash = commentedNode.Pos() - 1
	}
	return commentGroup
}
