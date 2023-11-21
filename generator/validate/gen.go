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

//go:generate go run -mod=mod sigs.k8s.io/controller-tools/cmd/helpgen@v0.6.2 generate:headerFile=../header.go.txt,year=2020 paths=.

// +controllertools:marker:generateHelp

// Generator validates the consistency of the API GO code.
//
// Validity checks are related to unions, patchStrategy, and optional fields.
type Generator struct{}

// RegisterMarkers registers the markers of the Generator
func (Generator) RegisterMarkers(into *markers.Registry) error {
	err := genutils.RegisterUnionMarkers(into)
	if err != nil {
		return err
	}
	return crdmarkers.Register(into)
}

func (Generator) CheckFilter() loader.NodeFilter {
	return func(node ast.Node) bool {
		// ignore interfaces
		_, isIface := node.(*ast.InterfaceType)
		return !isIface
	}
}

// Generate validates the source code
func (g Generator) Generate(ctx *genall.GenerationContext) error {
	for _, root := range ctx.Roots {
		ctx.Checker.Check(root)

		root.NeedTypesInfo()

		packageTypes := map[string]*markers.TypeInfo{}
		if err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
			packageTypes[info.Name] = info
		}); err != nil {
			root.AddError(err)
			return nil
		}

		for _, typeToCheck := range packageTypes {
			checkUnion(typeToCheck, root, packageTypes)
		}

	}

	return nil
}

func checkUnion(typeToCheck *markers.TypeInfo, root *loader.Package, packageTypes map[string]*markers.TypeInfo) {
	if typeToCheck.Markers.Get(genutils.UnionMarker.Name) == nil {
		return
	}

	unionMembers := []string{}
	var unionDiscriminatorField *markers.FieldInfo
	for index, unionField := range typeToCheck.Fields {
		if unionField.Markers.Get(genutils.UnionDiscriminatorMarker.Name) != nil {
			if unionDiscriminatorField != nil {
				root.AddError(loader.ErrFromNode(fmt.Errorf(
					"union `%v` should have only 1 union discriminator, but has 2: `%v` and `%v`",
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
			"union `%v` should have 1 union discriminator. See here for details: https://github.com/kubernetes/enhancements/blob/master/keps/sig-api-machinery/20190325-unions.md#proposal",
			typeToCheck.Name), typeToCheck.RawSpec))
		return
	}

	if unionDiscriminatorField.Markers.Get("optional") == nil {
		root.AddError(loader.ErrFromNode(fmt.Errorf(
			"in union `%v` the union discriminator `%v` should have the `+optional` comment marker",
			typeToCheck.Name,
			unionDiscriminatorField.Name), unionDiscriminatorField.RawField))
	}
	if !strings.Contains(unionDiscriminatorField.Tag.Get("json"), ",omitempty") {
		root.AddError(loader.ErrFromNode(fmt.Errorf(
			"in union `%v` the union discriminator `%v` should contain the `omitempty` option in its `json` tag, since it is expected to be an optional field",
			typeToCheck.Name,
			unionDiscriminatorField.Name), unionDiscriminatorField.RawField))
	}

	wrongTypeError := loader.ErrFromNode(fmt.Errorf(
		"in union `%v` the union discriminator `%v` should have a `string` type, or a type (defined in the same package) whose underlying type is a string",
		typeToCheck.Name,
		unionDiscriminatorField.Name), unionDiscriminatorField.RawField)

	discriminatorTypeRef, typeFound := root.TypesInfo.Types[unionDiscriminatorField.RawField.Type]
	if !typeFound {
		root.AddError(wrongTypeError)
		return
	}

	underlyingType, underlyingTypeIsBasic := discriminatorTypeRef.Type.Underlying().(*types.Basic)
	if !underlyingTypeIsBasic || underlyingType.Kind() != types.String {
		root.AddError(wrongTypeError)
		return
	}

	getEnumMarker := func(markers markers.MarkerValues) *crdmarkers.Enum {
		enumMarkerIf := markers.Get("kubebuilder:validation:Enum")
		if enumMarkerIf == nil {
			return nil
		}
		enumMarker := enumMarkerIf.(crdmarkers.Enum)
		return &enumMarker
	}

	var fieldEnumMarkerPtr *crdmarkers.Enum = getEnumMarker(unionDiscriminatorField.Markers)
	var typeEnumMarkerPtr *crdmarkers.Enum = nil
	var enumTypeInfo *markers.TypeInfo = nil

	if enumTypeNamed, isNamed := discriminatorTypeRef.Type.(*types.Named); isNamed {
		enumTypeInfo = packageTypes[enumTypeNamed.Obj().Name()]
		if enumTypeInfo != nil {
			typeEnumMarkerPtr = getEnumMarker(enumTypeInfo.Markers)
		} else {
			root.AddError(loader.ErrFromNode(fmt.Errorf(
				"in union `%v` the union discriminator `%v` should have a string-based type defined in the same package",
				typeToCheck.Name,
				unionDiscriminatorField.Name), unionDiscriminatorField.RawField))
		}
	}

	if fieldEnumMarkerPtr != nil && typeEnumMarkerPtr != nil {
		root.AddError(loader.ErrFromNode(fmt.Errorf(
			"type `%v` should not define a `kubebuilder:validation:Enum` annotation, since it is already defined in the discriminator `%v` of union `%v`",
			enumTypeInfo.Name,
			unionDiscriminatorField.Name,
			typeToCheck.Name), enumTypeInfo.RawSpec))
	}

	if fieldEnumMarkerPtr == nil && typeEnumMarkerPtr == nil {
		root.AddError(loader.ErrFromNode(fmt.Errorf(
			"in union `%v` the union discriminator `%v` should specify the allowed union values through a comment marker (kubebuilder:validation:Enum=...)",
			typeToCheck.Name,
			unionDiscriminatorField.Name), unionDiscriminatorField.RawField))
		return
	}

	checkEnumValues := func(enumMarkerPtr *crdmarkers.Enum, node ast.Node) {
		for index, enumValue := range *enumMarkerPtr {
			stringValue, isString := enumValue.(string)
			if !isString || stringValue != unionMembers[index] {
				root.AddError(loader.ErrFromNode(fmt.Errorf(
					"in union `%v` the union discriminator `%v` should allow the following values through the `kubebuilder:validation:Enum` comment marker: "+strings.Join(unionMembers, ";"),
					typeToCheck.Name,
					unionDiscriminatorField.Name), node))
				break
			}
		}
	}
	if fieldEnumMarkerPtr != nil {
		checkEnumValues(fieldEnumMarkerPtr, unionDiscriminatorField.RawField)
	} else {
		checkEnumValues(typeEnumMarkerPtr, enumTypeInfo.RawSpec)
	}
}
