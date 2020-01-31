package v1alpha1

import runtime "k8s.io/apimachinery/pkg/runtime"

type Events struct {
	// +optional
	Workspace *WorkspaceEvents `json:"workspace,omitempty"`
	Editor    *EditorEvents `json:"editor,omitempty"`
	Custom    []CustomEvent `json:"custom,omitempty"`
}

type WorkspaceEvents struct {
	// +optional
	PostCreate *WorkspaceEvent `json:"postCreate,omitempty"` 

	// +optional
	PreStart   *WorkspaceEvent `json:"preStart,omitempty"`

	// +optional
	PostStart  *WorkspaceEvent `json:"preStart,omitempty"`

	// +optional
	PreStop    *WorkspaceEvent `json:"preStop,omitempty"`

	// +optional
	PostStop   *WorkspaceEvent `json:"postStop,omitempty"`

	// +optional
	PostDelete *WorkspaceEvent `json:"postDelete,omitempty`
}

type BaseEvent struct {
	// The alias of the command to trigger
	Command string `json:"command"`
}

type WorkspaceEvent struct {
	BaseEvent `json:",inline"`
}

type EditorEvents struct {
	// +optional
	PostFirstOpen  *EditorEvent `json:"postFirstOpen,omitempty"`

	// +optional
	PostOpen       *EditorEvent `json:"postOpen,omitempty"`

	// +optional
	PostClone      *EditorEvent    `json:"postClone,omitempty"` 
}

type EditorEvent struct {
	BaseEvent `json:",inline"`
}

type CustomEvent struct {
	BaseEvent  `json:",inline"`
	EventClass string `json:"commandClass"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:EmbeddedResource
	EmbeddedResource runtime.RawExtension `json:"embeddedResource"`
}
