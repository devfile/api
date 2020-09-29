//
// Copyright (c) 2019-2020 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//

package v1alpha2

// K8sLikeComponentLocationType describes the type of
// the location the configuration is fetched from.
// Only one of the following component type may be specified.
// +kubebuilder:validation:Enum=Uri;Inlined
type K8sLikeComponentLocationType string

const (
	UriK8sLikeComponentLocationType     K8sLikeComponentLocationType = "Uri"
	InlinedK8sLikeComponentLocationType K8sLikeComponentLocationType = "Inlined"
)

// +union
type K8sLikeComponentLocation struct {
	// Type of Kubernetes-like location
	// +
	// +unionDiscriminator
	// +optional
	LocationType K8sLikeComponentLocationType `json:"locationType,omitempty"`

	// Location in a file fetched from a uri.
	// +optional
	Uri string `json:"uri,omitempty"`

	// Inlined manifest
	// +optional
	Inlined string `json:"inlined,omitempty"`
}

type K8sLikeComponent struct {
	BaseComponent            `json:",inline"`
	K8sLikeComponentLocation `json:",inline"`
	Endpoints                []Endpoint `json:"endpoints,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
}

// Component that allows partly importing Kubernetes resources into the workspace POD
type KubernetesComponent struct {
	K8sLikeComponent `json:",inline"`
}

// Component that allows partly importing Openshift resources into the workspace POD
type OpenshiftComponent struct {
	K8sLikeComponent `json:",inline"`
}
