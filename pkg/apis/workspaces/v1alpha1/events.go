package v1alpha1

type Events struct {
	WorkspaceEvents `json:",inline"`
}

type WorkspaceEvents struct {
	// +optional
	PreStart []string `json:"preStart,omitempty"`

	// +optional
	PostStart []string `json:"postStart,omitempty"`

	// +optional
	PreStop []string `json:"preStop,omitempty"`

	// +optional
	PostStop []string `json:"postStop,omitempty"`
}
