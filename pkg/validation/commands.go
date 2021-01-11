package validation

import (
	"fmt"
	"strings"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

// ValidateCommands validates the devfile commands:
// 1. if there are commands with duplicate IDs, an error is returned
// 2. checks if its either a valid exec or composite command
func ValidateCommands(commands []v1alpha2.Command, components []v1alpha2.Component) (err error) {
	processedCommands := make(map[string]string, len(commands))
	commandMap := getCommandsMap(commands)

	for _, command := range commands {
		// Check if the command is in the list of already processed commands
		// If there's a hit, it means more than one command share the same ID and we should error out
		commandID := strings.ToLower(command.Id)
		if _, exists := processedCommands[commandID]; exists {
			return &InvalidCommandError{commandId: command.Id, reason: "duplicate commands present with the same id"}
		}
		processedCommands[commandID] = commandID

		err = validateCommand(command, commandMap, components)
		if err != nil {
			return err
		}
	}

	return
}

// validateCommand validates a given devfile command
func validateCommand(command v1alpha2.Command, devfileCommands map[string]v1alpha2.Command, components []v1alpha2.Component) (err error) {

	// If the command is a composite command, need to validate that it is valid
	if command.Composite != nil {
		parentCommands := make(map[string]string)
		return validateCompositeCommand(&command, parentCommands, devfileCommands, components)
	}

	return validateExecCommand(command, components)
}

// validateExecCommand validates the given exec command, the command should:
// 1. have a component
// 2. have a command line
// 3. map to a valid container component
func validateExecCommand(command v1alpha2.Command, components []v1alpha2.Component) (err error) {

	if command.Exec == nil {
		return &InvalidCommandError{commandId: command.Id, reason: "should be of type exec"}
	}

	// must map to a container component
	isComponentValid := false
	for _, component := range components {
		if component.Container != nil && command.Exec.Component == component.Name {
			isComponentValid = true
		}
	}
	if !isComponentValid {
		return &InvalidCommandError{commandId: command.Id, reason: "command does not map to a container component"}
	}

	return
}

// validateCompositeCommand checks that the specified composite command is valid. The command:
// 1. should not reference itself via s subcommand
// 2. should not indirectly reference itself via a subcommand which is a composite command
// 3. should reference a valid devfile command
// 4. should have a valid exec sub command
func validateCompositeCommand(command *v1alpha2.Command, parentCommands map[string]string, devfileCommands map[string]v1alpha2.Command, components []v1alpha2.Component) error {

	// Store the command ID in a map of parent commands
	parentCommands[command.Id] = command.Id

	if command.Composite == nil {
		return &InvalidCommandError{commandId: command.Id, reason: "should be of type composite"}
	}

	// Loop over the commands and validate that each command points to a command that's in the devfile
	for _, cmd := range command.Composite.Commands {
		if strings.ToLower(cmd) == command.Id {
			return &InvalidCommandError{commandId: command.Id, reason: "composite command cannot reference itself"}
		}

		// Don't allow commands to indirectly reference themselves, so check if the command equals any of the parent commands in the command tree
		_, ok := parentCommands[strings.ToLower(cmd)]
		if ok {
			return &InvalidCommandError{commandId: command.Id, reason: "composite command cannot indirectly reference itself"}
		}

		subCommand, ok := devfileCommands[strings.ToLower(cmd)]
		if !ok {
			return &InvalidCommandError{commandId: command.Id, reason: fmt.Sprintf("the command %q mentioned in the composite command does not exist in the devfile", cmd)}
		}

		if subCommand.Composite != nil {
			// Recursively validate the composite subcommand
			err := validateCompositeCommand(&subCommand, parentCommands, devfileCommands, components)
			if err != nil {
				// Don't wrap the error message here to make the error message more readable to the user
				return err
			}
		} else {
			err := validateExecCommand(subCommand, components)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
