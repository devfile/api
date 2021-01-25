package v1alpha2

// Component that allows the developer to declare and configure a volume into his workspace
type VolumeComponent struct {
	BaseComponent `json:",inline"`
	Volume        `json:",inline"`
}

// Volume that should be mounted to a component container
type Volume struct {
	// +optional
	// Size of the volume
	Size string `json:"size,omitempty"`

	// External defines information about volumes that exist outside of the current workspace.
	// They are not created or deleted while processing a devfile but are still mounted into
	// component containers. When left empty, it is assumed that a new volume is to be created.
	//
	// Note: External volumes should be used with care, as they make devfiles less portable. It
	// is assumed that external volumes exist already.
	// +optional
	External ExistingVolumeRef `json:"external,omitempty"`
}

// ExistingVolumeRef is a refernce to a volume that exists outside the lifecycle of components
type ExistingVolumeRef struct {
	// Name defines the name of the resource
	Name string `json:"name"`
	// Type defines the type of the resource:
	//
	// - `storage` specifies that this volume refers to a PersistentVolumeClaim
	// - `configmap` specifies that this volume refers to a ConfigMap
	// - `secret` specifies that this volume refers to a Secret
	// kubebuilder:validation:Enum="persistent,configmap,secret"
	Type ExistingVolumeType `json:"type"`
}

// ExistingVolumeType defines the type of an external Volume
type ExistingVolumeType string

const (
	// PersistentVolumeType specifies persistent storage, e.g. a PersistentVolumeClaim
	PersistentVolumeType ExistingVolumeType = "persistent"
	// ConfigmapVolumeType specifies a configmap
	ConfigmapVolumeType ExistingVolumeType = "configmap"
	// SecretVolumeType specifies a secret
	SecretVolumeType ExistingVolumeType = "secret"
)
