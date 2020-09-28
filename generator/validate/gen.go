package interfaces

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/devfile/api/generator/genutils"
	"github.com/iancoleman/strcase"
	crdmarkers "sigs.k8s.io/controller-tools/pkg/crd/markers"
	"sigs.k8s.io/controller-tools/pkg/genall"
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
				// Check Union
				unionMembers := []string{}
				var unionDiscriminatorField *markers.FieldInfo
				for index, unionField := range typeToCheck.Fields {
					if unionField.Markers.Get(genutils.UnionDiscriminatorMarker.Name) != nil {
						if unionDiscriminatorField != nil {
							root.AddError(fmt.Errorf(
								"Union %v should have only 1 union discriminator, but has 2: %v and %v",
								typeToCheck.Name,
								unionDiscriminatorField.Name,
								unionField.Name))
						}
						unionDiscriminatorField = &(typeToCheck.Fields[index])
					} else {
						unionMembers = append(unionMembers, unionField.Name)
					}
				}
				if unionDiscriminatorField == nil {
					root.AddError(fmt.Errorf(
						"Union %v should have 1 union discriminator. See here for details: https://github.com/kubernetes/enhancements/blob/master/keps/sig-api-machinery/20190325-unions.md#proposal",
						typeToCheck.Name))
				}

				if unionDiscriminatorField.Markers.Get("optional") == nil {
					root.AddError(fmt.Errorf(
						"In union %v the union discriminator %v should have the `+optional` comment marker",
						typeToCheck.Name,
						unionDiscriminatorField.Name))
				}
				if unionDiscriminatorField.Tag.Get("json") == "" {
					root.AddError(fmt.Errorf(
						"In union %v the union discriminator %v should have the `+optional` comment marker",
						typeToCheck.Name,
						unionDiscriminatorField.Name))
				}
				// Simplify unions writing by adding the kubebuiler enum annotation ourselves in the CRD + Json schema generation.
				// => Nothing more to check here, well, yes, still check thatno enum is added manually at least.
			}
		}

	}

	return nil
}
