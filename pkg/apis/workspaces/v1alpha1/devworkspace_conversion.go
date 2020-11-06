package v1alpha1

import (
	"github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// Spokes for conversion have to satisfy the Convertible interface.
var _ conversion.Convertible = (*DevWorkspace)(nil)

func (src *DevWorkspace) ConvertTo(destRaw conversion.Hub) error {
	dest := destRaw.(*v1alpha2.DevWorkspace)
	return convertDevWorkspaceTo_v1alpha2(src, dest)
}

func (dest *DevWorkspace) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.DevWorkspace)
	return convertDevWorkspaceFrom_v1alpha2(src, dest)
}
