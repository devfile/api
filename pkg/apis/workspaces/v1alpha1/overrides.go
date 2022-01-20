package v1alpha1

type OverridesBase struct {
	// Overrides of commands encapsulated in a parent devfile or a plugin.
	// Overriding is done using a strategic merge patch
	// +optional
	Commands []Command `json:"commands,omitempty" patchStrategy:"merge" patchMergeKey:"id"`

	// Not implemented for now
	// additional directives to drive the strategic merge patch
	// OverrideDirectives []OverrideDirective `json:"overrideDirectives,omitempty"`
}

type Overrides struct {
	OverridesBase `json:",inline"`

	// Overrides of projects encapsulated in a parent devfile.
	// Overriding is done using a strategic merge patch.
	// +optional
	Projects []Project `json:"projects,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// Overrides of starter projects encapsulated in a parent devfile
	// Overriding is done using a strategic merge patch.
	// +optional
	StarterProjects []StarterProject `json:"starterProjects,omitempty"`

	// Overrides of components encapsulated in a parent devfile.
	// Overriding is done using a strategic merge patch
	// +optional
	Components []Component `json:"components,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
}

type PluginOverrides struct {
	OverridesBase `json:",inline"`

	// Overrides of components encapsulated in a plugin.
	// Overriding is done using a strategic merge patch.
	// A plugin cannot override embedded plugin components.
	// +optional
	Components []PluginComponentsOverride `json:"components,omitempty"`
}
