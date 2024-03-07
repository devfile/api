package v1alpha1

// ComposeFileComponentLocationType describes the type of
// the location where the docker-compose file is fetched from.
// Only one of the following types can be specified.
// +kubebuilder:validation:Enum=Uri;Inlined
type ComposeFileComponentLocationType string

const ( 
	UriComposeComponentLocationType ComposeFileComponentLocationType = "Uri"
	InlinedComposeComponentLocationType ComposeFileComponentLocationType = "Inlined"
)
// +union
type ComposeFileComponentLocation struct{
	// Type of Compose File Component Location
	// + 
	// +unionDiscriminator
	// +optional
	LocationType ComposeFileComponentLocationType `json:"locationType,omitempty"`

	// Location uri of the docker-compose file
	// +optional
	Uri string `json:"uri,omitempty"`

	// Inlined docker-compose file
	// +optional
	Inlined string `json:"inlined,omitempty"`
}

type ComposeLikeComponent struct {
	BaseComponent `json:",inline"`
	ComposeFileComponentLocation `json:",inline"`

	// Mandatory name that allows referencing the component
	// in commands, or inside a parent
	Name string `json:"name"`
}

// Component allows the developer to reuse an existing docker-compose file
type ComposeComponent struct{
	ComposeLikeComponent `json:",inline"`
}