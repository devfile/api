package v1alpha2

// DockerfileLocationType describes the type of
// the location for the Dockerfile outerloop build.
// Only one of the following location type may be specified.
// +kubebuilder:validation:Enum=Uri;Registry;Git
type DockerfileLocationType string

const (
	UriLikeDockerfileLocationType      DockerfileLocationType = "Uri"
	RegistryLikeDockerfileLocationType DockerfileLocationType = "Registry"
	GitLikeDockerfileLocationType      DockerfileLocationType = "Git"
)

// Dockerfile Image type to specify the outerloop build using a Dockerfile
type DockerfileImage struct {
	BaseImage          `json:",inline"`
	DockerfileLocation `json:",inline"`
	Dockerfile         `json:",inline"`
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

	// Dockerfile's Devfile Registry source
	// +optional
	Registry *DockerfileDevfileRegistrySource `json:"registry,omitempty"`

	// Dockerfile's Git source
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
	RootRequired *bool `json:"rootRequired,omitempty"`
}

type DockerfileDevfileRegistrySource struct {
	// Id in a devfile registry that contains a Dockerfile
	Id string `json:"id"`

	// Devfile Registry URL to pull the Dockerfile from when using the Devfile Registry as Dockerfile src.
	// To ensure the Dockerfile gets resolved consistently in different environments,
	// it is recommended to always specify the `devfileRegistryUrl` when `Id` is used.
	// +optional
	DevfileRegistryUrl string `json:"devfileRegistryUrl,omitempty"`
}

type DockerfileGitProjectSource struct {
	GitProjectSource `json:",inline"`

	// Location of the Dockerfile in the Git repository when using git as Dockerfile src.
	// +optional
	GitLocation string `json:"gitLocation,omitempty"`
}
