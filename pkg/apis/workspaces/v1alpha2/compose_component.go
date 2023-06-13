package v1alpha2

type ComposeFileComponentLocationType string

const (
	UriComposeComponentLocationType ComposeFileComponentLocationType = "Uri"
)

// +union
type ComposeFileComponentLocation struct {
	// Type of Compose file Component Location
	// +
	// +unionDiscriminator
	// +optional
	LocationType ComposeFileComponentLocationType `json:"locationType,omitempty"`

	// Location uri of the docker-compose file

	Uri string `json:"uri,omitempty"`
}

type ComposeLikeComponent struct {
	BaseComponent                `json:",inline"`
	ComposeFileComponentLocation `json:",inline"`
}

// Component allows the developer to reuse an existing Compose file
type ComposeComponent struct {
	ComposeLikeComponent `json:",inline"`
}
