package v1alpha1

import runtime "k8s.io/apimachinery/pkg/runtime"

// ComponentType describes the type of component.
// Only one of the following component type may be specified.
// +kubebuilder:validation:Enum= Container;Kubernetes;Openshift;CheEditor;Volume;ChePlugin;Custom
type ComponentType string

const (
	ContainerComponentType    ComponentType = "Container"
	KubernetesComponentType ComponentType = "Kubernetes"
	OpenshiftComponentType    ComponentType = "Openshift"
	CheEditorComponentType ComponentType = "CheEditor"
	ChePluginComponentType ComponentType = "ChePlugin"
	VolumeComponentType    ComponentType = "Volume"
	CustomComponentType ComponentType = "Custom"
)



// Workspace component: Anything that will bring additional features / tooling / behaviour / context
// to the workspace, in order to make working in it easier.
type BaseComponent struct {
}

type Component struct {
	PolymorphicComponent `json:",inline"`
}

// +k8s:openapi-gen=true
// +union
type PolymorphicComponent struct {
	// Type of project source
	// +
	// +unionDiscriminator
	// +optional
	Type ComponentType `json:"type"`

	// Container component
	// +optional
	Container *ContainerComponent `json:"container,omitempty"`

	// Volume component
	// +optional
	Volume *VolumeComponent `json:"volume,omitempty"`

	// CheEditor component
	// +optional
	CheEditor *CheEditorComponent `json:"cheEditor,omitempty"`

	// ChePlugin component
	// +optional
	ChePlugin *ChePluginComponent `json:"chePlugin,omitempty"`

	// Kubernetes component
	// +optional
	Kubernetes *KubernetesComponent `json:"kubernetes,omitempty"`

	// Openshift component
	// +optional
	Openshift *OpenshiftComponent `json:"openshift,omitempty"`

	// Custom component
	// +optional
	Custom *CustomComponent `json:"custom,omitempty"`
}

type CustomComponent struct {
	Name string `json:"name"`
	ComponentClass string `json:"componentClass"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:EmbeddedResource
	EmbeddedResource runtime.RawExtension `json:"embeddedResource"`
}
