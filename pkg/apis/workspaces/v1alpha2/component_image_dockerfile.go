package v1alpha2

// DockerfileLocationType describes the type of
// the location for the Dockerfile outerloop build.
// Only one of the following location type may be specified.
// +kubebuilder:validation:Enum=Uri;Id;Git
type DockerfileLocationType string

const (
	UriLikeDockerfileLocationType DockerfileLocationType = "Uri"
	IdLikeDockerfileLocationType  DockerfileLocationType = "Id"
	GitLikeDockerfileLocationType DockerfileLocationType = "Git"
)

// Dockerfile Image type to specify the outerloop build using a Dockerfile
type DockerfileImage struct {
	BaseImage          `json:",inline"`
	DockerfileLocation `json:",inline"`
	Dockerfile         `json:",inline"`

	// Registry URL to pull the Dockerfile from when using id as Dockerfile src.
	// To ensure the dockerfile gets resolved consistently in different environments,
	// it is recommended to always specify the `regsitryURL` when `Id` is used.
	// +optional
	RegistryUrl string `json:"registryUrl,omitempty"`
}

// +union
type DockerfileLocation struct {
	// Type of Dockerfile location
	// +
	// +unionDiscriminator
	// +optional
	LocationType DockerfileLocationType `json:"locationType,omitempty"`

	// URI Reference of a Dockerfile.
	// It can be a full URL or a relative URI from the current devfile as the base URI.
	// +optional
	Uri string `json:"uri,omitempty"`

	// Id in a registry that contains a Dockerfile
	// +optional
	Id string `json:"id,omitempty"`

	// Project's Git source
	// +optional
	Git *DockerfileGitProjectSource `json:"git,omitempty"`
}

type Dockerfile struct {
	// Path of source directory to establish build context. Defaults to ${PROJECT_ROOT} in the container
	// +optional
	BuildContext string `json:"buildContext,omitempty"`

	// The arguments to supply to the dockerfile build.
	// +optional
	Args []string `json:"args,omitempty" patchStrategy:"replace"`

	// Specify if a privileged builder pod is required.
	//
	// Default value is `false`
	// +optional
	RootRequired bool `json:"rootRequired,omitempty"`
}

type DockerfileGitProjectSource struct {
	GitProjectSource `json:",inline"`

	// Location of the Dockerfile in the Git repository when using git as Dockerfile src.
	// +optional
	GitLocation string `json:"gitLocation,omitempty"`
}
