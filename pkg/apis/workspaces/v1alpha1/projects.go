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

type Project struct {
	// Project name
	Name string `json:"name"`

	// Path relative to the root of the projects to which this project should be cloned into. This is a unix-style relative path (i.e. uses forward slashes). The path is invalid if it is absolute or tries to escape the project root through the usage of '..'. If not specified, defaults to the project name.
	// +optional
	ClonePath string `json:"clonePath,omitempty"`

	ProjectSource `json:",inline"`
}
type StarterProject struct {
	Project `json:",inline"`

	// Description of a starter project
	// +optional
	Description string `json:"description,omitempty"`
}

// ProjectSourceType describes the type of Project sources.
// Only one of the following project sources may be specified.
// If none of the following policies is specified, the default one
// is AllowConcurrent.
// +kubebuilder:validation:Enum=Git;Github;Zip;Custom
type ProjectSourceType string

const (
	GitProjectSourceType    ProjectSourceType = "Git"
	GitHubProjectSourceType ProjectSourceType = "Github"
	ZipProjectSourceType    ProjectSourceType = "Zip"
	CustomProjectSourceType ProjectSourceType = "Custom"
)

// +k8s:openapi-gen=true
// +union
type ProjectSource struct {
	// Type of project source
	// +
	// +unionDiscriminator
	// +optional
	SourceType ProjectSourceType `json:"sourceType,omitempty"`

	// Project's Git source
	// +optional
	Git *GitProjectSource `json:"git,omitempty"`

	// Project's GitHub source
	// +optional
	Github *GithubProjectSource `json:"github,omitempty"`

	// Project's Zip source
	// +optional
	Zip *ZipProjectSource `json:"zip,omitempty"`

	// Project's Custom source
	// +optional
	Custom *CustomProjectSource `json:"custom,omitempty"`
}

type CommonProjectSource struct {
	// Part of project to populate in the working directory.
	// +optional
	SparseCheckoutDir string `json:"sparseCheckoutDir,omitempty"`
}

type CustomProjectSource struct {
	ProjectSourceClass string `json:"projectSourceClass"`
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:EmbeddedResource
	EmbeddedResource runtime.RawExtension `json:"embeddedResource"`
}

type ZipProjectSource struct {
	CommonProjectSource `json:",inline"`

	// Zip project's source location address. Should be file path of the archive, e.g. file://$FILE_PATH
	// +required
	Location string `json:"location,omitempty"`
}

type GitLikeProjectSource struct {
	CommonProjectSource `json:",inline"`

	// Defines from what the project should be checked out. Required if there are more than one remote configured
	// +optional
	CheckoutFrom *CheckoutFrom `json:"checkoutFrom,omitempty"`

	// The remotes map which should be initialized in the git project. Must have at least one remote configured
	// +optional
	Remotes map[string]string `json:"remotes,omitempty"`
}

type CheckoutFrom struct {
	// The revision to checkout from. Should be branch name, tag or commit id.
	// Default branch is used if missing or specified revision is not found.
	// +optional
	Revision string `json:"revision,omitempty"`
	// The remote name should be used as init. Required if there are more than one remote configured
	// +optional
	Remote string `json:"remote,omitempty"`
}

type GitProjectSource struct {
	GitLikeProjectSource `json:",inline"`
}

type GithubProjectSource struct {
	GitLikeProjectSource `json:",inline"`
}
