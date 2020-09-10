package v1alpha2

// +k8s:openapi-gen=true
// +union
type BuildGuidances struct {

	// +optional
	// Allows specifying a dockerfile to initiate build
	Dockerfile *Dockerfile `json:"dockerfile,omitempty"`

	// +optional
	// Allows specifying a builder image to initiate s2i (SourceToImage) build
	SourceToImage *SourceToImage `json:"s2i,omitempty"`
}

type Dockerfile struct {

	// Name that allows referencing a build guidance
	Name string `json:"name"`

	// Dockerfile location which can be an URL or a path relative to buildContext
	DockerfileLocation string `json:"dockerfileLocation"`

	// +optional
	// Path of source directory to establish build context.  Default to the top level directory.
	BuildContext string `json:"buildContext"`

	// +optional
	// Optional flag that specifies whether unprivileged builder pod is required.  Default is false.
	Rootless bool `json:"rootless,omitempty"`
}

type SourceToImage struct {

	// Name that allows referencing a build guidance
	Name string `json:"name"`

	// Namespace where builder image is present
	BuilderImageNamespace string `json:"builderImageNamespace"`

	// Builder image name with tag
	BuilderImageStreamTag string `json:"builderImageStreamTag"`

	// +optional
	// Script URL to override default scripts provided by builder image
	ScriptLocation string `json:"scriptLocation,omitempty"`

	// +optional
	// Flag that indicates whether to perform incremental builds or no
	IncrementalBuild bool `json:"incrementalBuild,omitempty"`
}
