package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DevWorkspaceSpec defines the desired state of DevWorkspace
// +k8s:openapi-gen=true
type DevWorkspaceSpec struct {
	Started bool          `json:"started"`
	RoutingClass string   `json:"routingClass,omitempty"`
	Template DevWorkspaceTemplateSpec `json:"template,omitempty"`
}

// DevWorkspaceStatus defines the observed state of DevWorkspace
// +k8s:openapi-gen=true
type DevWorkspaceStatus struct {
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
