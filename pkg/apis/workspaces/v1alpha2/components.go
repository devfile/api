package v1alpha2

import runtime "k8s.io/apimachinery/pkg/runtime"

// ComponentType describes the type of component.
// Only one of the following component type may be specified.
// +kubebuilder:validation:Enum=Container;Kubernetes;Openshift;Volume;Plugin;Custom
type ComponentType string

const (
	ContainerComponentType  ComponentType = "Container"
	KubernetesComponentType ComponentType = "Kubernetes"
	OpenshiftComponentType  ComponentType = "Openshift"
	PluginComponentType     ComponentType = "Plugin"
	VolumeComponentType     ComponentType = "Volume"
	CustomComponentType     ComponentType = "Custom"
)

// Workspace component: Anything that will bring additional features / tooling / behaviour / context
// to the workspace, in order to make working in it easier.
type BaseComponent struct {
}

//+k8s:openapi-gen=true
type Component struct {
	// Mandatory name that allows referencing the component
	// from other elements (such as commands) or from an external
	// devfile that may reference this component through a parent or a plugin.
	Name           string `json:"name"`
	ComponentUnion `json:",inline"`
}

// +union
type ComponentUnion struct {
	// Type of component
	//
	// +unionDiscriminator
	// +optional
	ComponentType ComponentType `json:"componentType,omitempty"`

	// Allows adding and configuring workspace-related containers
	// +optional
	Container *ContainerComponent `json:"container,omitempty"`

	// Allows importing into the workspace the Kubernetes resources
	// defined in a given manifest. For example this allows reusing the Kubernetes
	// definitions used to deploy some runtime components in production.
	//
	// +optional
	Kubernetes *KubernetesComponent `json:"kubernetes,omitempty"`

	// Allows importing into the workspace the OpenShift resources
	// defined in a given manifest. For example this allows reusing the OpenShift
	// definitions used to deploy some runtime components in production.
	//
	// +optional
	Openshift *OpenshiftComponent `json:"openshift,omitempty"`

	// Allows importing a plugin.
	//
	// Plugins are mainly imported devfiles that contribute components, commands
	// and events as a consistent single unit. They are defined in either YAML files
	// following the devfile syntax,
	// or as `DevWorkspaceTemplate` Kubernetes Custom Resources
	// +optional
	// +devfile:overrides:include:omitInPlugin=true
	Plugin *PluginComponent `json:"plugin,omitempty"`

	// Allows specifying the definition of a volume
	// shared by several other components
	// +optional
	Volume *VolumeComponent `json:"volume,omitempty"`

	// Custom component whose logic is implementation-dependant
	// and should be provided by the user
	// possibly through some dedicated controller
	// +optional
	// +devfile:overrides:include:omit=true
	Custom *CustomComponent `json:"custom,omitempty"`
}

type CustomComponent struct {
	// Class of component that the associated implementation controller
	// should use to process this command with the appropriate logic
	ComponentClass string `json:"componentClass"`

	// Additional free-form configuration for this custom component
	// that the implementation controller will know how to use
	//
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:EmbeddedResource
	EmbeddedResource runtime.RawExtension `json:"embeddedResource"`
}

// PluginComponentsOverrideType describes the type of components
// that can be overriden for a plugin.
// Only one of the following component type may be specified.
// +kubebuilder:validation:Enum=Container;Kubernetes;Openshift;Volume
type PluginComponentsOverrideType string

const (
	ContainerPluginComponentsOverrideType  PluginComponentsOverrideType = "Container"
	KubernetesPluginComponentsOverrideType PluginComponentsOverrideType = "Kubernetes"
	OpenshiftPluginComponentsOverrideType  PluginComponentsOverrideType = "Openshift"
	VolumePluginComponentsOverrideType     PluginComponentsOverrideType = "Volume"
)

//+k8s:openapi-gen=true
type PluginComponentsOverride struct {
	// Mandatory name that allows referencing the Volume component
	// in Container volume mounts or inside a parent
	Name                          string `json:"name"`
	PluginComponentsOverrideUnion `json:",inline"`
}

// +union
type PluginComponentsOverrideUnion struct {
	// Type of component override for a plugin
	//
	// +unionDiscriminator
	// +optional
	ComponentType PluginComponentsOverrideType `json:"componentType,omitempty"`

	// Configuration overriding for a Container component in a plugin
	// +optional
	Container *ContainerComponent `json:"container,omitempty"`

	// Configuration overriding for a Volume component in a plugin
	// +optional
	Volume *VolumeComponent `json:"volume,omitempty"`

	// Configuration overriding for a Kubernetes component in a plugin
	// +optional
	Kubernetes *KubernetesComponent `json:"kubernetes,omitempty"`

	// Configuration overriding for an OpenShift component in a plugin
	// +optional
	Openshift *OpenshiftComponent `json:"openshift,omitempty"`
}
