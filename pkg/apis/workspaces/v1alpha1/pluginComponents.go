package v1alpha1

type PluginOverrides struct {
	// Overrides of components encapsulated in a plugin.
	// Overriding is done using a strategic merge
	// +optional
	Components []ComponentOverride `json:"components,omitempty"`

	// Overrides of commands encapsulated in a plugin.
	// Overriding is done using a strategic merge
	// +optional
	Commands []Command `json:"commands,omitempty"`
}

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
