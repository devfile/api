package v1alpha2

// BuildGuidanceType describes the type of build guidance.
// Only one of the following build guidance type may be specified.
// +kubebuilder:validation:Enum=Dockerfile;SourceToImage
type BuildGuidanceType string

const (
	DockerfileBuildGuidanceType    BuildGuidanceType = "Dockerfile"
	SourceToImageBuildGuidanceType BuildGuidanceType = "SourceToImage"
)

//+k8s:openapi-gen=true
type BuildGuidance struct {
	// Mandatory name that allows referencing the buid guidance from other elements or from
	// an external devfile that may reference this build guidance through a parent or a plugin.
	Name               string `json:"name"`
	BuildGuidanceUnion `json:",inline"`
}

// +union
type BuildGuidanceUnion struct {
	// Type of build guidance
	//
	// +unionDiscriminator
	// +optional
	BuildGuidanceType `json:"buildGuidanceType,omitempty"`

	// +optional
	// Allows specifying a dockerfile to initiate build
	Dockerfile *Dockerfile `json:"dockerfile,omitempty"`

	// +optional
	// Allows specifying a builder image to initiate s2i (SourceToImage) build
	SourceToImage *SourceToImage `json:"s2i,omitempty"`
}

type Dockerfile struct {
	// Dockerfile location which can be an URL or a path relative to buildContext
	DockerfileLocation string `json:"dockerfileLocation"`

	// +optional
	// Path of source directory to establish build context.  Default to the top level directory.
	BuildContext string `json:"buildContext,omitempty"`

	// +optional
	// Optional flag that specifies whether unprivileged builder pod is required.  Default is false.
	Rootless bool `json:"rootless,omitempty"`
}

type SourceToImage struct {
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
