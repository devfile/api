package v1alpha1

// +k8s:openapi-gen=true
// +union
type K8sLikeComponentLocation struct {
	// Type of Kubernetes-like location
	// +
	// +unionDiscriminator
	// +optional
	LocationType string `json:"locationType"`

	// Location in a plugin registry
	// +optional
	Url string `json:"url,omitempty"`

	// Reference to the plugin definition
	// +optional
	Inlined string `json:"inlined,omitempty"`
}

type K8sLikeComponent struct {
	BaseComponent                          `json:",inline"`
	K8sLikeComponentLocation               `json:",inline"`
	// Mandatory name that allows referencing the component
	// in commands, or inside a parent
	Name string `json:"name"`
}

// Component that allows partly importing Kubernetes resources into the workspace POD
type KubernetesComponent struct {
	K8sLikeComponent `json:",inline"`
}

// Component that allows partly importing Openshift resources into the workspace POD
type OpenshiftComponent struct {
	K8sLikeComponent `json:",inline"`
}
