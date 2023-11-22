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

package getters

import (
	"bytes"
	"fmt"
	"go/ast"
	"strconv"

	"github.com/devfile/api/generator/genutils"
	"github.com/elliotchance/orderedmap"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

//go:generate go run -mod=mod sigs.k8s.io/controller-tools/cmd/helpgen@v0.6.2 generate:headerFile=../header.go.txt,year=2021 paths=.

var (
	// GetterTypeMarker is associated with a type that's used as the pointer receiver of the getter method
	GetterTypeMarker = markers.Must(markers.MakeDefinition("devfile:getter:generate", markers.DescribesType, struct{}{}))
	// DefaultFieldMarker is associated with a boolean pointer field to indicate the default boolean value
	DefaultFieldMarker = markers.Must(markers.MakeDefinition("devfile:default:value", markers.DescribesField, ""))
)

// +controllertools:marker:generateHelp

// Generator generates getter methods that are used to return values for the boolean pointer fields.
//
// The pointer receiver is determined from the `devfile:getter:generate` annotated type.  The method will return the value of the
// field if it's been set, otherwise it will return the default value specified by the devfile:default:value annotation.
type Generator struct{}

// RegisterMarkers registers the markers of the Generator
func (Generator) RegisterMarkers(into *markers.Registry) error {
	if err := markers.RegisterAll(into, GetterTypeMarker, DefaultFieldMarker); err != nil {
		return err
	}
	into.AddHelp(GetterTypeMarker,
		markers.SimpleHelp("Devfile", "indicates the type that's used as the pointer receiver of the getter method"))
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

// getterInfo stores the info to generate the getter method
type getterInfo struct {
	funcName   string
	defaultVal string
}

// Generate generates the artifacts
func (g Generator) Generate(ctx *genall.GenerationContext) error {
	for _, root := range ctx.Roots {
		ctx.Checker.Check(root)
		root.NeedTypesInfo()

		typesToProcess := orderedmap.NewOrderedMap()
		if err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
			if info.Markers.Get(GetterTypeMarker.Name) != nil {
				var getters []getterInfo
				for _, field := range info.Fields {
					defaultVal := field.Markers.Get(DefaultFieldMarker.Name)
					if defaultVal != nil {
						if _, err := strconv.ParseBool(defaultVal.(string)); err != nil {
							root.AddError(fmt.Errorf("devfile:default:value marker specified on %s/%s does not have a true or false value.  Value is %s", info.Name, field.Name, defaultVal.(string)))
						}

						//look for boolean pointers
						if ptr, isPtr := field.RawField.Type.(*ast.StarExpr); isPtr {
							if ident, ok := ptr.X.(*ast.Ident); ok {
								if ident.Name == "bool" {
									getters = append(getters, getterInfo{
										field.Name,
										defaultVal.(string),
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
				if len(getters) > 0 {
					typesToProcess.Set(info, getters)
				} else {
					root.AddError(fmt.Errorf("type %s does not have the field marker, devfile:default:value specified on a boolean pointer field", info.Name))
				}
				return
			}

		}); err != nil {
			root.AddError(err)
			return nil
		}

		genutils.WriteFormattedSourceFile("getters", ctx, root, func(buf *bytes.Buffer) {
			for elt := typesToProcess.Front(); elt != nil; elt = elt.Next() {
				cmd := elt.Key.(*markers.TypeInfo)
				fields := elt.Value.([]getterInfo)
				for _, getter := range fields {
					fName := getter.funcName
					defaultVal := getter.defaultVal
					getterMethod := fmt.Sprintf(`
// Get%[1]s returns the value of the boolean property.  If unset, it's the default value specified in the devfile:default:value marker
func (in *%[2]s) Get%[1]s() bool {
return getBoolOrDefault(in.%[1]s, %[3]s)}`, fName, cmd.Name, defaultVal)
					buf.WriteString(getterMethod)
				}
			}

			internalHelper := `

func getBoolOrDefault(input *bool, defaultVal bool) bool {
	if input != nil {
		return *input 
	} 
	return defaultVal }`
			buf.WriteString(internalHelper)
		})
	}

	return nil
}
