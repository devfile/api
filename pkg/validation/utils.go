package validation

import (
	"net/url"
	"strings"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

const (
	ImportSourceAttribute   = "library.devfile.io/imported-from"
	ParentOverrideAttribute = "library.devfile.io/parent-override-from"
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
