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
	ImportReferenceType ImportReferenceType `json:"importReferenceType,omitempty"`

	// URI Reference of a Devfile yaml file, can be a full URL
	// or a relative URI with the current devfile as the base URI
	// +optional
	Uri string `json:"uri,omitempty"`

	// Id in a registry that contains a Devfile yaml file
	// +optional
	Id string `json:"id,omitempty"`

	// Reference to a Kubernetes CRD of type DevWorkspaceTemplate
	// +optional
	Kubernetes *KubernetesCustomResourceImportReference `json:"kubernetes,omitempty"`
}

type KubernetesCustomResourceImportReference struct {
	Name string `json:"name"`

	// +optional
	Namespace string `json:"namespace,omitempty"`
}

type ImportReference struct {
	ImportReferenceUnion `json:",inline"`
	// Registry URL to pull devfile from with the specified Id in the import reference.
	// It is recommended to always have the regsitryURL specified if Id is provided to reference a parent devfile
	// It to provides a well-defined parent and ensures it gets resolved consistently in different envrionments.
	// +optional
	RegistryUrl string `json:"registryUrl,omitempty"`
}
