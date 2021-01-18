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
	// Persistent defines whether to use persistent storage for this volume. Defaults
	// to true. When set to false, storage is tied to a pod lifecycle and is erased on
	// delete.
	// +optional
	// +kubebuilder:default=true
	Persistent bool `json:"persistent,omitempty"`
	// ReadOnly specifies whether the volume should be mounted without write capabilities.
	// Defaults to false.
	// +optional
	ReadOnly bool `json:"readonly,omitempty"`
	// External defines information about volumes that exist outside of the current workspace.
	// They are not created or deleted while processing a devfile but are still mounted into
	// component containers. When left empty, it is assumed that a new volume is to be created.
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
	// kubebuilder:validation:Enum="storage,configmap,secret"
	Type ExistingVolumeType `json:"type"`
}

// ExistingVolumeType defines the type of an external Volume
type ExistingVolumeType string

const (
	// StorageVolumeType specifies persistent storage, e.g. a PersistentVolumeClaim
	StorageVolumeType ExistingVolumeType = "storage"
	// ConfigmapVolumeType specifies a configmap
	ConfigmapVolumeType ExistingVolumeType = "configmap"
	// SecretVolumeType specifies a secret
	SecretVolumeType ExistingVolumeType = "secret"
)
