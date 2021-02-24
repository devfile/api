package attributes

import (
	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
)

// ValidateAndReplaceForParent validates the parent data for global attribute references(except parent overrides) and replaces them with the attribute value
func ValidateAndReplaceForParent(attributes apiAttributes.Attributes, parent *v1alpha2.Parent) error {
	var err error

	if parent != nil {
		switch {
		case parent.Id != "":
			// Validate parent id
			if parent.Id, err = validateAndReplaceDataWithAttribute(parent.Id, attributes); err != nil {
				return err
			}
		case parent.Uri != "":
			// Validate parent uri
			if parent.Uri, err = validateAndReplaceDataWithAttribute(parent.Uri, attributes); err != nil {
				return err
			}
		case parent.Kubernetes != nil:
			// Validate parent kubernetes name
			if parent.Kubernetes.Name, err = validateAndReplaceDataWithAttribute(parent.Kubernetes.Name, attributes); err != nil {
				return err
			}

			// Validate parent kubernetes namespace
			if parent.Kubernetes.Namespace, err = validateAndReplaceDataWithAttribute(parent.Kubernetes.Namespace, attributes); err != nil {
				return err
			}
		}

		// Validate parent registry url
		if parent.RegistryUrl, err = validateAndReplaceDataWithAttribute(parent.RegistryUrl, attributes); err != nil {
			return err
		}

		// Note: No need to substitute parent overrides at this point. Call global attribute validation/substitution
		// after merging the flattened parent devfile to the main devfile. Parent's global attribute key can
		// be overridden in parent overrides or the alternative is to mention the attribute as a main devfile
		// global attribute if parent devfile does not have a global attribute
	}

	return nil
}
