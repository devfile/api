//
//
// Copyright Red Hat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

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

// +k8s:openapi-gen=true
// +union
type Component struct {
	// Type of component
	//
	// +unionDiscriminator
	// +optional
	ComponentType ComponentType `json:"componentType,omitempty"`

	// Allows adding and configuring workspace-related containers
	// +optional
	Container *ContainerComponent `json:"container,omitempty"`

	// Allows specifying the definition of a volume
	// shared by several other components
	// +optional
	Volume *VolumeComponent `json:"volume,omitempty"`

	// Allows importing a plugin.
	//
	// Plugins are mainly imported devfiles that contribute components, commands
	// and events as a consistent single unit. They are defined in either YAML files
	// following the devfile syntax,
	// or as `DevWorkspaceTemplate` Kubernetes Custom Resources
	// +optional
	Plugin *PluginComponent `json:"plugin,omitempty"`

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

	// Custom component whose logic is implementation-dependant
	// and should be provided by the user
	// possibly through some dedicated controller
	// +optional
	Custom *CustomComponent `json:"custom,omitempty"`
}

type CustomComponent struct {
	// Mandatory name that allows referencing the component
	// in commands, or inside a parent
	Name string `json:"name"`

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

// +k8s:openapi-gen=true
// +union
type PluginComponentsOverride struct {
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
