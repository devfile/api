package v1alpha1

// Structure of the workspace. This is also the specification of a workspace template.
// +k8s:openapi-gen=true
type DevWorkspaceTemplateSpec struct {
	// Predefined, ready-to-use, workspace-related commands
	Commands          []Command      `json:"commands,omitempty"`

	// Projects worked on in the workspace, containing names and sources locations
	Projects          []Project      `json:"projects,omitempty"`
	
	// List of the workspace components, such as editor and plugins,
	// user-provided containers, or other types of components
	// +optional
	Components        []Component `json:"components,omitempty"`
}
