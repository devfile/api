package v1alpha1

type PluginComponent struct {
	BaseComponent   `json:",inline"`
	ImportReference `json:",inline"`
	PluginOverrides `json:",inline"`

	// +optional
	// Optional name that allows referencing the component
	// in commands, or inside a parent
	// If omitted it will be infered from the location (uri or registryEntry)
	Name string `json:"name,omitempty"`
}
