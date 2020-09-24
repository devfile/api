package v1alpha2

// Devfile describes the structure of a cloud-native workspace and development environment.
// +devfile:jsonschema:generate:omitCustomUnionMembers=true
type Devfile struct {
	// Devfile schema version
	// +kubebuilder:validation:Pattern=^([2-9][0-9]*)\.([0-9]+)\.([0-9]+)(\-[0-9a-z-]+(\.[0-9a-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$
	SchemaVersion string `json:"schemaVersion"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	// Optional metadata
	Metadata DevfileMetadata `json:"metadata,omitempty"`

	DevWorkspaceTemplateSpec `json:",inline"`
}

type DevfileMetadata struct {
	// Optional devfile name
	Name string `json:"name,omitempty"`

	// Optional semver-compatible version
	// +kubebuilder:validation:Pattern=^([0-9]+)\.([0-9]+)\.([0-9]+)(\-[0-9a-z-]+(\.[0-9a-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$
	Version string `json:"version,omitempty"`
}
