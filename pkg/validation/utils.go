package validation

import (
	"net/url"
	"strings"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

// attribute keys for imported and overridden elements
// the value of those keys is the resource information
const (
	// attribute key of the imported element resource information
	ImportSourceAttribute = "library.devfile.io/imported-from"
	// attribute key of the parent overridden element resource information
	ParentOverrideAttribute = "library.devfile.io/parent-override-from"
	// attribute key of the plugin overridden element resource information
	PluginOverrideAttribute = "library.devfile.io/plugin-override-from"
)

// getCommandsMap iterates through the commands and returns a map of command
func getCommandsMap(commands []v1alpha2.Command) map[string]v1alpha2.Command {
	commandMap := make(map[string]v1alpha2.Command, len(commands))

	for _, command := range commands {
		command.Id = strings.ToLower(command.Id)
		commandMap[command.Id] = command
	}

	return commandMap
}

// ValidateURI checks if the string is with valid uri format, return error if not valid
func ValidateURI(uri string) error {
	if strings.HasPrefix(uri, "http") {
		if _, err := url.ParseRequestURI(uri); err != nil {
			return err
		}
	} else if _, err := url.Parse(uri); err != nil {
		return err
	}

	return nil
}
