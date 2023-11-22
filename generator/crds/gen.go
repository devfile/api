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

package crds

import (
	"fmt"
	"go/ast"

	"github.com/devfile/api/generator/genutils"

	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-tools/pkg/crd"
	crdmarkers "sigs.k8s.io/controller-tools/pkg/crd/markers"
	"sigs.k8s.io/controller-tools/pkg/genall"
	"sigs.k8s.io/controller-tools/pkg/loader"
	"sigs.k8s.io/controller-tools/pkg/markers"
)

//go:generate go run -mod=mod sigs.k8s.io/controller-tools/cmd/helpgen@v0.6.2 generate:headerFile=../header.go.txt,year=2020 paths=.

// +controllertools:marker:generateHelp

// Generator generates CustomResourceDefinition YAML manifests for each root Kubernetes resource.
//
// Currently this generates v1 and v1beta1 CRDs for the `DevWorkspace` and `DevWorkspaceTemplate` resources.
type Generator struct{}

func (Generator) CheckFilter() loader.NodeFilter {
	return func(node ast.Node) bool {
		return true
	}
}

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
			root.AddError(fmt.Errorf("the package should have a valid 'groupName' annotation"))
			return nil
		}

		unionDiscriminators := unionDiscriminatorsByGV[groupVersion]

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
		apiVersions := []string{}
		for _, apiVersion := range crdRaw.Spec.Versions {
			apiVersions = append(apiVersions, apiVersion.Name)
			unionDiscriminators := unionDiscriminatorsByGV[groupKind.WithVersion(apiVersion.Name).GroupVersion()]
			genutils.AddUnionOneOfConstraints(apiVersion.Schema.OpenAPIV3Schema, unionDiscriminators, false)
		}

		latestAPIVersion := genutils.LatestKubeLikeVersion(apiVersions)

		for pkg, gv := range parser.GroupVersions {
			if gv.Group != groupKind.Group || gv.Version != latestAPIVersion {
				continue
			}

			typeIdent := crd.TypeIdent{Package: pkg, Name: groupKind.Kind}
			typeInfo := parser.Types[typeIdent]
			if typeInfo == nil {
				continue
			}

			for _, markerVals := range typeInfo.Markers {
				for _, val := range markerVals {
					crdMarker, isCrdResourceMarker := val.(crdmarkers.Resource)
					if !isCrdResourceMarker {
						continue
					}
					if err := crdMarker.ApplyToCRD(&crdRaw.Spec, latestAPIVersion); err != nil {
						pkg.AddError(loader.ErrFromNode(err /* an okay guess */, typeInfo.RawSpec))
					}
				}
			}
		}

		for i, ver := range crdVersions {
			copiedCrd := crdRaw.DeepCopy()

			// drop defaults in v1beta1 since they are not supported there
			if crdVersions[i] == "v1beta1" {
				for _, apiVersion := range copiedCrd.Spec.Versions {
					genutils.EditJSONSchema(
						apiVersion.Schema.OpenAPIV3Schema,
						func(schema *apiext.JSONSchemaProps) (newVisitor genutils.Visitor, stop bool) {
							if schema != nil {
								schema.Default = nil
							}
							return
						})
				}
			}
			extCrd, err := crd.AsVersion(*copiedCrd, schema.GroupVersion{Group: apiext.SchemeGroupVersion.Group, Version: ver})
			if err != nil {
				return err
			}

			var fileName string
			if i == 0 {
				fileName = fmt.Sprintf("%s_%s.yaml", crdRaw.Spec.Group, crdRaw.Spec.Names.Plural)
			} else {
				fileName = fmt.Sprintf("%s_%s.%s.yaml", crdRaw.Spec.Group, crdRaw.Spec.Names.Plural, crdVersions[i])
			}
			if err := ctx.WriteYAML(fileName, extCrd); err != nil {
				return err
			}
		}
	}

	return nil
}
