package attributes

import (
	"regexp"
	"strings"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
)

// ValidateAndReplaceGlobalAttribute validates the workspace template spec data for global attribute references and replaces them with the attribute value
func ValidateAndReplaceGlobalAttribute(workspaceTemplateSpec *v1alpha2.DevWorkspaceTemplateSpec) error {

	var err error

	if workspaceTemplateSpec != nil {
		// Validate the components and replace for global attribute
		if err = ValidateAndReplaceForComponents(workspaceTemplateSpec.Attributes, workspaceTemplateSpec.Components); err != nil {
			return err
		}

		// Validate the commands and replace for global attribute
		if err = ValidateAndReplaceForCommands(workspaceTemplateSpec.Attributes, workspaceTemplateSpec.Commands); err != nil {
			return err
		}

		// Validate the projects and replace for global attribute
		if err = ValidateAndReplaceForProjects(workspaceTemplateSpec.Attributes, workspaceTemplateSpec.Projects); err != nil {
			return err
		}

		// Validate the starter projects and replace for global attribute
		if err = ValidateAndReplaceForStarterProjects(workspaceTemplateSpec.Attributes, workspaceTemplateSpec.StarterProjects); err != nil {
			return err
		}
	}

	return nil
}

var globalAttributeRegex = regexp.MustCompile(`\{\{(.*?)\}\}`)

// validateAndReplaceDataWithAttribute validates the string for a global attribute and replaces it. An error
// is returned if the string references an invalid global attribute key
func validateAndReplaceDataWithAttribute(val string, attributes apiAttributes.Attributes) (string, error) {
	matches := globalAttributeRegex.FindAllStringSubmatch(val, -1)
	for _, match := range matches {
		var err error
		attrValue := attributes.GetString(match[1], &err)
		if err != nil {
			return "", err
		}
		val = strings.Replace(val, match[0], attrValue, -1)
	}

	return val, nil
}
