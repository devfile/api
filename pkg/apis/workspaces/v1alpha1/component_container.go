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

// Component that allows the developer to add a configured container into his workspace
type ContainerComponent struct {
	BaseComponent `json:",inline"`
	Container     `json:",inline"`
	Endpoints     []Endpoint `json:"endpoints,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
}

type Container struct {
	Name string `json:"name"`

	Image string `json:"image,omitempty"`

	// +optional
	// Environment variables used in this container
	Env []EnvVar `json:"env,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// +optional
	// List of volumes mounts that should be mounted is this container.
	VolumeMounts []VolumeMount `json:"volumeMounts,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// +optional
	MemoryLimit string `json:"memoryLimit,omitempty"`

	// The command to run in the dockerimage component instead of the default one provided in the image.
	//
	// Defaults to an empty array, meaning use whatever is defined in the image.
	// +optional
	Command []string `json:"command,omitempty" patchStrategy:"replace"`

	// The arguments to supply to the command running the dockerimage component. The arguments are supplied either to the default command provided in the image or to the overridden command.
	//
	// Defaults to an empty array, meaning use whatever is defined in the image.
	// +optional
	Args []string `json:"args,omitempty" patchStrategy:"replace"`

	// +optional
	MountSources bool `json:"mountSources,omitempty"`

	// Optional specification of the path in the container where
	// project sources should be transferred/mounted when `mountSources` is `true`.
	// When omitted, the value of the `PROJECTS_ROOT` environment variable is used.
	// +optional
	SourceMapping string `json:"sourceMapping,omitempty"`

	// Specify if a container should run in its own separated pod,
	// instead of running as part of the main development environment pod.
	//
	// Default value is `false`
	// +optional
	DedicatedPod bool `json:"dedicatedPod,omitempty"`
}

type EnvVar struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

// Volume that should be mounted to a component container
type VolumeMount struct {
	// The volume mount name is the name of an existing `Volume` component.
	// If several containers mount the same volume name
	// then they will reuse the same volume and will be able to access to the same files.
	Name string `json:"name"`

	// The path in the component container where the volume should be mounted.
	// If not path is mentioned, default path is the is `/<name>`.
	// +optional
	Path string `json:"path,omitempty"`
}
