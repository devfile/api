package v1alpha2

type ComposeFileComponentLocationType string

const (
	UriComposeComponentLocationType     ComposeFileComponentLocationType = "Uri"
	InlinedComposeComponentLocationType ComposeFileComponentLocationType = "Inlined"
)

type ComposeFileComponentLocation struct {

	// Type of Compose file Component Location

	LocationType ComposeFileComponentLocationType `json:"locationType,omitempty"`

	// Location uri of the docker-compose file

	Uri string `json:"uri,omitempty"`

	// Inline docker-file

	Inlined string `json:"inlined,omitempty"`
}

type ComposeLikeComponent struct {
	BaseComponent                `json:",inline"`
	ComposeFileComponentLocation `json:",inline"`
}

// Component allows the developer to reuse an existing Compose file
type ComposeComponent struct {
	ComposeLikeComponent `json:",inline"`
}
