package v1alpha1

// Component that allows the developer to declare and configure a volume into his workspace
type VolumeComponent struct {
	BaseComponent `json:",inline"`
	Volume        `json:",inline"`
}

// Volume that should be mounted to a component container
type Volume struct {
	// Mandatory name that allows referencing the Volume component
	// in Container volume mounts or inside a parent
	Name string `json:"name"`

	// +optional
	// Size of the volume
	Size string `json:"size,omitempty"`
}
