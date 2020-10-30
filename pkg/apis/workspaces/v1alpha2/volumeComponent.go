package v1alpha2

// Component that allows the developer to declare and configure a volume into his workspace
type VolumeComponent struct {
	BaseComponent `json:",inline" yaml:",inline"`
	Volume        `json:",inline" yaml:",inline"`
}

// Volume that should be mounted to a component container
type Volume struct {
	// +optional
	// Size of the volume
	Size string `json:"size,omitempty" yaml:"size,omitempty"`
}
