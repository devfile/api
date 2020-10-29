package v1alpha1

import (
	"errors"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// Spokes for conversion have to satisfy the Convertible interface.
var _ conversion.Convertible = (*DevWorkspaceTemplate)(nil)

func (src *DevWorkspaceTemplate) ConvertTo(dstRaw conversion.Hub) error {
	return errors.New("Unimplemented")
}

func (dst *DevWorkspaceTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	return errors.New("Unimplemented")
}
