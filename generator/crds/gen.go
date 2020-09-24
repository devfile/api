package crds

import (
	"bytes"
	"fmt"
	"go/format"
	"io"

	"github.com/devfile/api/generator/genutils"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-tools/pkg/crd"

	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	crdmarkers "sigs.k8s.io/controller-tools/pkg/crd/markers"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

//go:generate go run sigs.k8s.io/controller-tools/cmd/helpgen generate:headerFile=../header.go.txt,year=2020 paths=.

// +controllertools:marker:generateHelp

// Generator generates CustomResourceDefinition YAML manifests for each root Kubernetes resource.
//
// Currently this generates v1 and v1beta1 CRDs for the `DevWorkspace` and `DevWorkspaceTemplate` resources.
type Generator struct{}

// RegisterMarkers registers the markers of the Generator
func (Generator) RegisterMarkers(into *markers.Registry) error {
	if err := genutils.RegisterUnionMarkers(into); err != nil {
		return err
	}
	return crdmarkers.Register(into)
}

// Generate generates the artifacts
func (g Generator) Generate(ctx *genall.GenerationContext) error {
	parser := &crd.Parser{
		Collector:           ctx.Collector,
		Checker:             ctx.Checker,
		AllowDangerousTypes: false,
	}

	crd.AddKnownTypes(parser)
	for _, root := range ctx.Roots {
		parser.NeedPackage(root)
	}

	metav1Pkg := crd.FindMetav1(ctx.Roots)
	if metav1Pkg == nil {
		// no objects in the roots, since nothing imported metav1
		return nil
	}

	kubeKinds := crd.FindKubeKinds(parser, metav1Pkg)
	if len(kubeKinds) == 0 {
		// no objects in the roots
		return nil
	}

	unionDiscriminatorsByGV := map[schema.GroupVersion][]markers.FieldInfo{}

	for _, root := range ctx.Roots {
		packageMarkers, err := markers.PackageMarkers(ctx.Collector, root)
		if err != nil {
			root.AddError(err)
			return nil
		}

		var groupVersion schema.GroupVersion
		switch groupName := packageMarkers.Get("groupName").(type) {
		case string:
			groupVersion = schema.GroupVersion{
				Group:   groupName,
				Version: root.Name,
			}
		default:
			root.AddError(fmt.Errorf("The package should have a valid 'groupName' annotation"))
			return nil
		}

		unionDiscriminators, found := unionDiscriminatorsByGV[groupVersion]
		if !found {
			unionDiscriminators = []markers.FieldInfo{}
		}

		if err := markers.EachType(ctx.Collector, root, func(info *markers.TypeInfo) {
			if info.Markers.Get(genutils.UnionMarker.Name) != nil {
				for _, field := range info.Fields {
					if field.Markers.Get(genutils.UnionDiscriminatorMarker.Name) != nil {
						unionDiscriminators = append(unionDiscriminators, field)
					}
				}
				return
			}
		}); err != nil {
			root.AddError(err)
			return nil
		}
		unionDiscriminatorsByGV[groupVersion] = unionDiscriminators
	}

	crdVersions := []string{"v1", "v1beta1"}

	for groupKind := range kubeKinds {
		parser.NeedCRDFor(groupKind, nil)
		crdRaw := parser.CustomResourceDefinitions[groupKind]
		for _, apiVersion := range crdRaw.Spec.Versions {
			unionDiscriminators := unionDiscriminatorsByGV[groupKind.WithVersion(apiVersion.Name).GroupVersion()]
			genutils.AddUnionOneOfConstraints(apiVersion.Schema.OpenAPIV3Schema, unionDiscriminators, false)
		}

		versionedCRDs := make([]interface{}, len(crdVersions))
		for i, ver := range crdVersions {
			conv, err := crd.AsVersion(crdRaw, schema.GroupVersion{Group: apiext.SchemeGroupVersion.Group, Version: ver})
			if err != nil {
				return err
			}
			versionedCRDs[i] = conv
		}

		for i, crd := range versionedCRDs {
			var fileName string
			if i == 0 {
				fileName = fmt.Sprintf("%s_%s.yaml", crdRaw.Spec.Group, crdRaw.Spec.Names.Plural)
			} else {
				fileName = fmt.Sprintf("%s_%s.%s.yaml", crdRaw.Spec.Group, crdRaw.Spec.Names.Plural, crdVersions[i])
			}
			if err := ctx.WriteYAML(fileName, crd); err != nil {
				return err
			}
		}
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
