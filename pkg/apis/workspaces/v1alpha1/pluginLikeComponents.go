package v1alpha1

type RegistryEntryPluginLocation struct {
	Id string  `json:"id"`

	// +optional
	BaseUrl string `json:"baseUrl,omitempty"`
}

// PluginLocationType describes the type of location where the plugin definition can be fetched from.
// Only one of the following values may be specified.
// +kubebuilder:validation:Enum= RegistryEntry;Uri
type PluginLocationType string

const (
	RegistryEntryPluginLocationType    PluginLocationType = "RegistryEntry"
	UriPluginLocationType    PluginLocationType = "Uri"
)

// +k8s:openapi-gen=true
// +union
type ChePluginLocation struct {
	// Type of plugin location
	// +
	// +unionDiscriminator
	// +optional
	LocationType PluginLocationType `json:"locationType"`

	// Location of an entry inside a plugin registry
	// +optional
	RegistryEntry *RegistryEntryPluginLocation `json:"registryEntry,omitempty"`

	// Location defined as an URI
	// +optional
	Uri string `json:"uri,omitempty"`
}

type PluginLikeComponent struct {
	BaseComponent `json:",inline"`
	MemoryLimit  string `json:"memoryLimit,omitempty"`
	ChePluginLocation `json:",inline"`
}

type ChePluginComponent struct {
	PluginLikeComponent `json:",inline"`
}

type CheEditorComponent struct {
	PluginLikeComponent `json:",inline"`
}
