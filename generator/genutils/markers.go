package genutils

import (
	"sigs.k8s.io/controller-tools/pkg/markers"
)

var (
	// UnionMarker is the definition of the union marker, as defined in https://github.com/kubernetes/enhancements/blob/master/keps/sig-api-machinery/20190325-unions.md#proposal
	UnionMarker = markers.Must(markers.MakeDefinition("union", markers.DescribesType, struct{}{}))
	// UnionDiscriminatorMarker is the definition of the union discriminator marker, as defined in https://github.com/kubernetes/enhancements/blob/master/keps/sig-api-machinery/20190325-unions.md#proposal
	UnionDiscriminatorMarker = markers.Must(markers.MakeDefinition("unionDiscriminator", markers.DescribesField, struct{}{}))
)

// RegisterUnionMarkers registers the `union` and `unionDiscriminator` markers
func RegisterUnionMarkers(into *markers.Registry) error {
	if err := markers.RegisterAll(into, UnionMarker, UnionDiscriminatorMarker); err != nil {
		return err
	}
	into.AddHelp(UnionMarker,
		markers.SimpleHelp("Devfile", "indicates that a given Struct type is a K8S union, and its fields (apart from the discriminator) are mutually exclusive. K8S unions are described here: https://github.com/kubernetes/enhancements/blob/master/keps/sig-api-machinery/20190325-unions.md#proposal"))
	into.AddHelp(UnionDiscriminatorMarker,
		markers.SimpleHelp("Devfile", "indicates that a given field of an union Struct type is the union discriminator. K8S unions are described here: https://github.com/kubernetes/enhancements/blob/master/keps/sig-api-machinery/20190325-unions.md#proposal"))
	return nil
}
