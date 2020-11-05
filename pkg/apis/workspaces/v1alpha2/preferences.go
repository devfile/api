package v1alpha2

import (
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type Preference struct {
	// Mandatory key that uniquely references the preference
	// especially from an external defvile that may override this preference
	// through a parent or a plugin.
	Name                     string `json:"name"`
	// Map of implementation-dependant string-based free-form attributes.
	// +optional
	// +devfile:overrides:include:omit=true
	Attributes map[string]string `json:"attributes,omitempty"`
	PreferenceLocationUnion `json:",inline"`
}

type YamlPreference map[string]apiext.JSON

// +union
type PreferenceLocationUnion struct {
	// Type of preference
	//
	// +unionDiscriminator
	// +optional
	// +kubebuilder:validation:Enum=Yaml;Inline;Uri
	PreferenceType string `json:"preferenceType,omitempty"`

	// Free-form Yaml preference
	// +optional
	Yaml YamlPreference `json:"yaml,omitempty"`

	// Opaque raw string preference
	// +optional
	Inline string `json:"inline,omitempty"`

	// uri where the preferences string should be loaded from
	// +optional
	Uri string `json:"uri,omitempty"`
}
