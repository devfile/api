package v1alpha1

//Dockerfile component of devfile
type Dockerfile struct {
	// Mandatory name that allows referencing the Volume component in Container volume mounts or inside a parent
	Name string `json:"name"`

	// Mandatory path to source code
	Source *Source `json:"source"`

	// Mandatory path to dockerfile
	DockerfileLocation string `json:"dockerfileLocation"`

	// Mandatory destination to registry to push built image
	Destination string `json:"destination,omitempty"`
}


//Source within dockerfile component
type Source struct {
	// Mandatory path to local source directory folder
	SourceDir string `json:"sourceDir"`

	// Mandatory path to source repository hosted locally or on cloud
	Location string `json:"location"`
}
