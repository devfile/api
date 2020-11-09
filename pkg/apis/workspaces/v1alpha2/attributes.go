package v1alpha2

import (
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// Attributes provides a way to add a map of arbitrary YAML/JSON
// objects.
type Attributes map[string]apiext.JSON
