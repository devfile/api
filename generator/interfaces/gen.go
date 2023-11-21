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

package interfaces

import (
	"bytes"
	"go/ast"
	"go/printer"
	"strings"

	"sigs.k8s.io/controller-tools/pkg/loader"

	"github.com/devfile/api/generator/genutils"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/markers"

	"github.com/elliotchance/orderedmap"
)

//go:generate go run -mod=mod sigs.k8s.io/controller-tools/cmd/helpgen@v0.6.2 generate:headerFile=../header.go.txt,year=2020 paths=.

var (
	toplevelListMarker = markers.Must(markers.MakeDefinition("devfile:toplevellist", markers.DescribesField, struct{}{}))
)

// +controllertools:marker:generateHelp

// Generator generates GO source code required for the API
//
// Generated source code mainly consists in interface implementations to manage unions and keyed top-level lists
type Generator struct{}

// RegisterMarkers registers the markers of the Generator
func (Generator) RegisterMarkers(into *markers.Registry) error {
	if err := markers.RegisterAll(into, toplevelListMarker); err != nil {
		return err
	}
	into.AddHelp(toplevelListMarker,
		markers.SimpleHelp("Devfile", "indicates that a given field of the Devfile body structure is a top-level list that should be managed through strategic merge patch during parent of plugin overriding."))
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

		unions := orderedmap.NewOrderedMap()
		toplevelListContainers := orderedmap.NewOrderedMap()
		keyed := orderedmap.NewOrderedMap()
		if err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
			if info.Markers.Get(genutils.UnionMarker.Name) != nil {
				unions.Set(info.Name, info)
				return
			}
			for i, field := range info.Fields {
				if field.Markers.Get(toplevelListMarker.Name) != nil {
					toplevelListContainers.Set(info.Name, info)
					if arrayType, isArrayType := field.RawField.Type.(*ast.ArrayType); isArrayType {
						if ident, isIdent := arrayType.Elt.(*ast.Ident); isIdent {
							keyed.Set(ident.Name, &info.Fields[i])
						}
					}
				}
			}
		}); err != nil {
			root.AddError(err)
			return nil
		}

		genutils.WriteFormattedSourceFile("keyed_definitions", ctx, root, func(buf *bytes.Buffer) {
			for elt := keyed.Front(); elt != nil; elt = elt.Next() {
				typeName := elt.Key.(string)
				field := elt.Value.(*markers.FieldInfo)
				mergeKey := strings.Title(genutils.GetPatchMergeKey(field))
				buf.WriteString(`
func (keyed ` + typeName + `) Key() string {
	return keyed.` + mergeKey + `
}
`)
			}
		})

		genutils.WriteFormattedSourceFile("toplevellistcontainer_definitions", ctx, root, func(buf *bytes.Buffer) {
			for elt := toplevelListContainers.Front(); elt != nil; elt = elt.Next() {
				typeName := elt.Key.(string)
				theType := elt.Value.(*markers.TypeInfo)
				buf.WriteString(`
func (container ` + typeName + `) GetToplevelLists() TopLevelLists {
	return TopLevelLists{`)
				for _, field := range theType.Fields {
					if field.Markers.Get(toplevelListMarker.Name) != nil {
						buf.WriteString(`
		"` + field.Name + `": extractKeys(container.` + field.Name + `),`)
					}
				}
				buf.WriteString(`
	}
}
`)
			}
		})

		genutils.WriteFormattedSourceFile("union_definitions", ctx, root, func(buf *bytes.Buffer) {
			buf.WriteString(`
import (
	"reflect"
)
`)
			for elt := unions.Front(); elt != nil; elt = elt.Next() {
				typeName := elt.Key.(string)
				theType := elt.Value.(*markers.TypeInfo)
				visitorName := typeName + "Visitor"
				visitorType := strings.ToLower(string(typeName[0])) + typeName[1:]
				discriminatorName := ""
				fieldMap := orderedmap.NewOrderedMap()
				for _, field := range theType.Fields {
					if field.Markers.Get(genutils.UnionDiscriminatorMarker.Name) != nil {
						discriminatorName = field.Name
					} else {
						buf := new(bytes.Buffer)
						printer.Fprint(buf, root.Fset, field.RawField.Type)
						fieldMap.Set(field.Name, buf.String())
					}
				}

				buf.WriteString(`
var ` + visitorType + ` reflect.Type = reflect.TypeOf(` + visitorName + `{})

func (union ` + typeName + `) Visit(visitor ` + visitorName + `) error {
	return visitUnion(union, visitor)
}
func (union *` + typeName + `) discriminator() *string {
	return (*string)(&union.` + discriminatorName + `)
}
func (union *` + typeName + `) Normalize() error {
	return normalizeUnion(union, ` + visitorType + `)
}
func (union *` + typeName + `) Simplify() {
	simplifyUnion(union, ` + visitorType + `)
}

// +k8s:deepcopy-gen=false
type ` + visitorName + ` struct {`)

				for elt := fieldMap.Front(); elt != nil; elt = elt.Next() {
					fieldName := elt.Key.(string)
					fieldType := elt.Value.(string)
					buf.WriteString(`
	` + fieldName + ` func(` + fieldType + `) error`)
				}
				buf.WriteString(`
}
`)
			}
		})
	}

	return nil
}
