package v1alpha2

// ImportReferenceType describes the type of location
// from where the referenced template structure should be retrieved.
// Only one of the following parent locations may be specified.
// +kubebuilder:validation:Enum=Uri;Id;Kubernetes
type ImportReferenceType string

const (
	UriImportReferenceType        ImportReferenceType = "Uri"
	IdImportReferenceType         ImportReferenceType = "Id"
	KubernetesImportReferenceType ImportReferenceType = "Kubernetes"
)

// Location from where the an import reference is retrieved
// +union
type ImportReferenceUnion struct {
	// type of location from where the referenced template structure should be retrieved
	// +
	// +unionDiscriminator
	// +optional
	ImportReferenceType ImportReferenceType `json:"importReferenceType,omitempty" yaml:"importReferenceType,omitempty"`

	// Uri of a Devfile yaml file
	// +optional
	Uri string `json:"uri,omitempty" yaml:"uri,omitempty"`

	// Id in a registry that contains a Devfile yaml file
	// +optional
	Id string `json:"id,omitempty" yaml:"id,omitempty"`

	// Reference to a Kubernetes CRD of type DevWorkspaceTemplate
	// +optional
	Kubernetes *KubernetesCustomResourceImportReference `json:"kubernetes,omitempty" yaml:"kubernetes,omitempty"`
}

type KubernetesCustomResourceImportReference struct {
	Name string `json:"name" yaml:"name"`

	// +optional
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

type ImportReference struct {
	ImportReferenceUnion `json:",inline" yaml:",inline"`
	// +optional
	RegistryUrl string `json:"registryUrl,omitempty" yaml:"registryUrl,omitempty"`
}
