package v1alpha1

type RegistryLocation struct {
	Id string  `json:"id"`

	// +optional
	RegistryUrl string `json:"registryUrl,omitempty"`
}

// PluginLocationType describes the type of location where the plugin definition can be fetched from.
// Only one of the following values may be specified.
// +kubebuilder:validation:Enum= Registry;Url
type PluginLocationType string

const (
	RegistryPluginLocationType    PluginLocationType = "Registry"
	UrlPluginLocationType    PluginLocationType = "Url"
)

// +k8s:openapi-gen=true
// +union
type ChePluginLocation struct {
	// Type of plugin location
	// +
	// +unionDiscriminator
	// +optional
	LocationType PluginLocationType `json:"locationType"`

	// Location in a plugin registry
	// +optional
	Registry *RegistryLocation `json:"registry,omitempty"`

	// Location defined as an URL
	// +optional
	Url string `json:"url,omitempty"`
}

type PluginLikeComponent struct {
	BaseComponent `json:",inline"`
	MemoryLimit  string `json:"memoryLimit,omitempty"`
	Location ChePluginLocation `json:",inline"`
}

type ChePluginComponent struct {
	PluginLikeComponent `json:",inline"`
}

type CheEditorComponent struct {
	PluginLikeComponent `json:",inline"`
}
