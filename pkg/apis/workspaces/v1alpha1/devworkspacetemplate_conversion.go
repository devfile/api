package v1alpha1

import (
	"github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// Spokes for conversion have to satisfy the Convertible interface.
var _ conversion.Convertible = (*DevWorkspaceTemplate)(nil)

func (src *DevWorkspaceTemplate) ConvertTo(destRaw conversion.Hub) error {
	dest := destRaw.(*v1alpha2.DevWorkspaceTemplate)
	return convertDevWorkspaceTemplateTo_v1alpha2(src, dest)

}

func (dest *DevWorkspaceTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.DevWorkspaceTemplate)
	return convertDevWorkspaceTemplateFrom_v1alpha2(src, dest)
}
