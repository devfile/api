package v1alpha1

import runtime "k8s.io/apimachinery/pkg/runtime"

type Project struct {
	// Project name
	Name          string `json:"name"`

	// Path relative to the root of the projects to which this project should be cloned into. This is a unix-style relative path (i.e. uses forward slashes). The path is invalid if it is absolute or tries to escape the project root through the usage of '..'. If not specified, defaults to the project name.
	// +optional
	ClonePath     string `json:"clonePath,omitempty"`

	ProjectSource `json:",inline"`
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
	SourceType ProjectSourceType `json:"sourceType"`

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
	// Project's source location address. Should be URL for git and github located projects, or; file:// for zip
	Location string `json:"location"`

	// Part of project to populate in the working directory.
	// +optional
	SparseCheckoutDir string `json:"sparseCheckoutDir,omitmepty"`
}

type CustomProjectSource struct {
	ProjectSourceClass string `json:"projectSourceClass"`
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:EmbeddedResource
	EmbeddedResource runtime.RawExtension `json:"embeddedResource"`
}

type ZipProjectSource struct {
	CommonProjectSource `json:",inline"`
}

type GitLikeProjectSource struct {
	CommonProjectSource `json:",inline"`

	// The tag or commit id to reset the checked out branch to
	// +optional
	StartPoint string `json:"startPoint,omitempty"`

	// The branch to check
	// +optional
	Branch string `json:"branch,omitempty"`
}

type GitProjectSource struct {
	GitLikeProjectSource `json:",inline"`
}

type GithubProjectSource struct {
	GitLikeProjectSource `json:",inline"`
}
