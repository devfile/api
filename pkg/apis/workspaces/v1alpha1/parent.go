package v1alpha1

type Parent struct {
	ParentLocation `json:",inline"`
}

// ParentLocationType describes the type of location
// from where the parent workspace structure should be retrieved.
// Only one of the following parent locations may be specified.
// +kubebuilder:validation:Enum=Uri;RegistryEntry;Kubernetes
type ParentLocationType string

const (
	UriParentLocationType    ParentLocationType = "Uri"
	RegistryEntryParentLocationType ParentLocationType = "RegistryEntry"
	KubernetesParentLocationType    ParentLocationType = "Kubernetes"
)

// Location from where the parent workspace structure is retrieved
// +k8s:openapi-gen=true
// +union
type ParentLocation struct {
	// Type of parent location
	// +
	// +unionDiscriminator
	// +optional
	LocationType ParentLocationType `json:"locationType"`

	// Uri of a Devfile yaml file
	// +optional
	Uri string `json:"uri,omitempty"`

	// Entry in a registry (base URL + ID) that contains a Devfile yaml file  
	// +optional
	RegistryEntry *RegistryEntryParentLocation `json:"registryEntry,omitempty"`

	// Reference to a Kubernetes CRD of type DevWorkspaceTemplate
	// +optional
	Kubernetes *KubernetesCustomResourceParentLocation `json:"kubernetes,omitempty"`
}

type RegistryEntryParentLocation struct {
	Id string  `json:"id"`

	// +optional
	baseUrl string `json:"baseUrl,omitempty"`
}

type KubernetesCustomResourceParentLocation struct {
	Name string  `json:"name"`

	// +optional
	Namespace string `json:"namespace,omitempty"`
}
