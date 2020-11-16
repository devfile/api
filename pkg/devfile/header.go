package devfile

// DevfileHeader describes the structure of the devfile-specific top-level fields
// that are not part of the K8S API structures
type DevfileHeader struct {
	// Devfile schema version
	// +kubebuilder:validation:Pattern=^([2-9])\.([0-9]+)\.([0-9]+)(\-[0-9a-z-]+(\.[0-9a-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$
	SchemaVersion string `json:"schemaVersion"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	// Optional metadata
	Metadata DevfileMetadata `json:"metadata,omitempty"`
}

type DevfileMetadata struct {
	// Optional devfile name
	// +optional
	Name string `json:"name,omitempty"`

	// Optional semver-compatible version
	// +optional
	// +kubebuilder:validation:Pattern=^([0-9]+)\.([0-9]+)\.([0-9]+)(\-[0-9a-z-]+(\.[0-9a-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$
	Version string `json:"version,omitempty"`

	// Optional devfile metadata artifact
	// +optional
	Artifact []MetadataArtifact `json:"artifact,omitempty"`
}

// MetadataArtifact contains additional information that is not defined in devfile schema
type MetadataArtifact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
