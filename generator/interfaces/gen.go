package interfaces

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/printer"
	"io"
	"strings"

	"github.com/devfile/api/generator/genutils"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"

	"github.com/elliotchance/orderedmap"
)

//go:generate go run sigs.k8s.io/controller-tools/cmd/helpgen generate:headerFile=../header.go.txt,year=2020 paths=.

var (
	toplevelListMarker = markers.Must(markers.MakeDefinition("devfile:toplevellist", markers.DescribesField, struct{}{}))
)

type copyStructs struct {
	StructName string
	Fields     []string
}

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

// Generate generates the artifacts
func (g Generator) Generate(ctx *genall.GenerationContext) error {
	for _, root := range ctx.Roots {
		ctx.Checker.Check(root, func(node ast.Node) bool {
			// ignore interfaces
			_, isIface := node.(*ast.InterfaceType)
			return !isIface
		})

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

		g.writeFoFile("keyed_definitions", ctx, root, func(buf *bytes.Buffer) {
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

		g.writeFoFile("toplevellistcontainer_definitions", ctx, root, func(buf *bytes.Buffer) {
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

		g.writeFoFile("union_definitions", ctx, root, func(buf *bytes.Buffer) {
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

func (g Generator) writeFoFile(filename string, ctx *genall.GenerationContext, root *loader.Package, writeContents func(*bytes.Buffer)) error {
	buf := new(bytes.Buffer)
	buf.WriteString(`
package ` + root.Name + `
`)

	writeContents(buf)

	outContents, err := format.Source(buf.Bytes())
	if err != nil {
		root.AddError(err)
		return err
	}

	fullname := "zz_generated." + filename + ".go"
	outputFile, err := ctx.Open(root, fullname)
	if err != nil {
		root.AddError(err)
		return err
	}
	defer outputFile.Close()
	n, err := outputFile.Write(outContents)
	if err != nil {
		root.AddError(err)
		return err
	}
	if n < len(outContents) {
		root.AddError(io.ErrShortWrite)
		return err
	}
	return nil
}
