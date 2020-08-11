package v1alpha1

//Dockerfile component of devfile
type Dockerfile struct {
	// Mandatory name that allows referencing the Volume component in Container volume mounts or inside a parent
	Name string `json:"name"`

	// path to source code, if empty - source is assumed as the directory having devfile
	Source *Source `json:"source,omitempty"`

	// Mandatory path to dockerfile
	DockerfileLocation string `json:"dockerfileLocation"`

	// destination to registry to push built image
	Destination string `json:"destination,omitempty"`

	// field indicating whether rootless/unprivileged builder container is required
	Rootless bool `json:"rootless,omitempty"`
}

//Source within dockerfile component
type Source struct {
	// path to local source directory folder
	SourceDir string `json:"sourceDir,omitempty"`

	// path to source repository hosted locally or on cloud
	Location string `json:"location,omitempty"`
}
