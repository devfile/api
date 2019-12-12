package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DevWorkspaceSpec defines the desired state of DevWorkspace
// +k8s:openapi-gen=true
type DevWorkspaceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Started bool          `json:"started"`
	EndpointsClass string  `json:"endpointsClass,omitempty"`
	Template DevWorkspaceTemplateSpec `json:"template,omitempty"`
}

// Structure of the workspace. This is also the specification of a workspace template.
type DevWorkspaceTemplateSpec struct {
	// Predefined, ready-to-use, workspace-related commands
	Commands          []Command      `json:"commands,omitempty"`
	
	// Projects worked on in the workspace, containing names and sources locations
	Projects          []Project      `json:"projects,omitempty"`
	
	// List of the workspace components, such as editor and plugins,
	// user-provided containers, or other types of components
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge,retainKeys
	Components        []Component `json:"components,omitempty"`
}

// DevWorkspaceStatus defines the observed state of DevWorkspace
// +k8s:openapi-gen=true
type DevWorkspaceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

	// Id of the workspace
	WorkspaceId string `json:"workspaceId"`
	// URL at which the Worksace Editor can be joined
	MainIdeUrl string `json:"mainIdeUrl,omitempty"`
	// AdditionalInfo
	AdditionalInfo map[string]string `json:"additionalInfo,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevWorkspace is the Schema for the devworkspaces API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=devworkspaces,scope=Namespaced
type DevWorkspace struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DevWorkspaceSpec   `json:"spec,omitempty"`
	Status DevWorkspaceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevWorkspaceList contains a list of DevWorkspace
type DevWorkspaceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DevWorkspace `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DevWorkspace{}, &DevWorkspaceList{})
}
