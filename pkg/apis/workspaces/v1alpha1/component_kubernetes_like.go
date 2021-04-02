package v1alpha1

// K8sLikeComponentLocationType describes the type of
// the location the configuration is fetched from.
// Only one of the following component type may be specified.
// +kubebuilder:validation:Enum=Uri;Inlined
type K8sLikeComponentLocationType string

const (
	UriK8sLikeComponentLocationType     K8sLikeComponentLocationType = "Uri"
	InlinedK8sLikeComponentLocationType K8sLikeComponentLocationType = "Inlined"
)

// +k8s:openapi-gen=true
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

	// Mandatory name that allows referencing the component
	// in commands, or inside a parent
	Name string `json:"name"`

	Endpoints []Endpoint `json:"endpoints,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
}

// Component that allows partly importing Kubernetes resources into the workspace POD
type KubernetesComponent struct {
	K8sLikeComponent `json:",inline"`
}

// Component that allows partly importing Openshift resources into the workspace POD
type OpenshiftComponent struct {
	K8sLikeComponent `json:",inline"`
}
