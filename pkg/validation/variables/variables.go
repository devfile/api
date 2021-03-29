package variables

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

// ValidateAndReplaceGlobalVariable validates the workspace template spec data for global variable references and replaces them with the variable value
func ValidateAndReplaceGlobalVariable(workspaceTemplateSpec *v1alpha2.DevWorkspaceTemplateSpec) error {

	var err error

	if workspaceTemplateSpec != nil {
		// Validate the components and replace for global variable
		if err = ValidateAndReplaceForComponents(workspaceTemplateSpec.Variables, workspaceTemplateSpec.Components); err != nil {
			return err
		}

		// Validate the commands and replace for global variable
		if err = ValidateAndReplaceForCommands(workspaceTemplateSpec.Variables, workspaceTemplateSpec.Commands); err != nil {
			return err
		}

		// Validate the projects and replace for global variable
		if err = ValidateAndReplaceForProjects(workspaceTemplateSpec.Variables, workspaceTemplateSpec.Projects); err != nil {
			return err
		}

		// Validate the starter projects and replace for global variable
		if err = ValidateAndReplaceForStarterProjects(workspaceTemplateSpec.Variables, workspaceTemplateSpec.StarterProjects); err != nil {
			return err
		}
	}

	return nil
}

var globalVariableRegex = regexp.MustCompile(`\{\{(.*?)\}\}`)

// validateAndReplaceDataWithVariable validates the string for a global variable and replaces it. An error
// is returned if the string references an invalid global variable key
func validateAndReplaceDataWithVariable(val string, variables map[string]string) (string, error) {
	matches := globalVariableRegex.FindAllStringSubmatch(val, -1)
	for _, match := range matches {
		varValue, ok := variables[match[1]]
		if !ok {
			return "", fmt.Errorf("Variable with key %q does not exist", match[1])
		}
		val = strings.Replace(val, match[0], varValue, -1)
	}

	return val, nil
}
