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

// Structure of the workspace. This is also the specification of a workspace template.
// +k8s:openapi-gen=true
type DevWorkspaceTemplateSpec struct {
	// Parent workspace template
	// +optional
	Parent *Parent `json:"parent,omitempty"`

	DevWorkspaceTemplateSpecContent `json:",inline"`
}

type DevWorkspaceTemplateSpecContent struct {
	// Predefined, ready-to-use, workspace-related commands
	// +optional
	//
	Commands []Command `json:"commands,omitempty" patchStrategy:"merge" patchMergeKey:"id"`

	// Bindings of commands to events.
	// Each command is referred-to by its name.
	// +optional
	Events *Events `json:"events,omitempty"`

	// Projects worked on in the workspace, containing names and sources locations
	// +optional
	Projects []Project `json:"projects,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// StarterProjects is a project that can be used as a starting point when bootstrapping new projects
	// +optional
	StarterProjects []StarterProject `json:"starterProjects,omitempty"`

	// List of the workspace components, such as editor and plugins,
	// user-provided containers, or other types of components
	// +optional
	Components []Component `json:"components,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
}
