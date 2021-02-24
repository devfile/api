package attributes

import (
	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
)

// validateAndReplaceForEndpoint validates the endpoint data for global attribute references and replaces them with the attribute value
func validateAndReplaceForEndpoint(attributes apiAttributes.Attributes, endpoints []v1alpha2.Endpoint) error {

	for i := range endpoints {
		var err error

		// Validate endpoint path
		if endpoints[i].Path, err = validateAndReplaceDataWithAttribute(endpoints[i].Path, attributes); err != nil {
			return err
		}
	}

	return nil
}
