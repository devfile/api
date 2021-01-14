package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

// ValidateCommands validates the devfile commands:
// 1. if there are commands with duplicate IDs, an error is returned
// 2. checks if its either a valid command
// 3. checks if commands belonging to a specific group obeys the rule of 1 default command
func ValidateCommands(commands []v1alpha2.Command, components []v1alpha2.Component) (err error) {
	processedCommands := make(map[string]string, len(commands))
	groupCommandMap := make(map[v1alpha2.CommandGroup][]v1alpha2.Command)
	commandMap := getCommandsMap(commands)

	for _, command := range commands {
		// Check if the command is in the list of already processed commands
		// If there's a hit, it means more than one command share the same ID and we should error out
		if isInt(command.Id) {
			return &InvalidNameOrIdError{id: command.Id, resourceType: "command"}
		}
		if _, exists := processedCommands[command.Id]; exists {
			return &InvalidCommandError{commandId: command.Id, reason: "duplicate commands present with the same id"}
		}
		processedCommands[command.Id] = command.Id

		parentCommands := make(map[string]string)
		err = validateCommand(command, parentCommands, commandMap, components)
		if err != nil {
			return err
		}

		commandGroup := *getGroup(command)
		if !reflect.DeepEqual(commandGroup, v1alpha2.CommandGroup{}) {
			groupCommandMap[commandGroup] = append(groupCommandMap[commandGroup], command)
		}
	}

	groupErrors := ""
	for group, commands := range groupCommandMap {
		if err = validateGroup(commands); err != nil {
			groupErrors += fmt.Sprintf("\ncommand group %s error - %s", group.Kind, err.Error())
		}
	}

	if len(groupErrors) > 0 {
		err = fmt.Errorf("%s", groupErrors)
	}

	return err
}

// validateCommand validates a given devfile command
func validateCommand(command v1alpha2.Command, parentCommands map[string]string, devfileCommands map[string]v1alpha2.Command, components []v1alpha2.Component) (err error) {

	switch {
	case command.Composite != nil:
		return validateCompositeCommand(&command, parentCommands, devfileCommands, components)
	case command.Exec != nil || command.Apply != nil:
		return validateCommandComponent(command, components)
	case command.VscodeLaunch != nil:
		if command.VscodeLaunch.Uri != "" {
			return ValidateURI(command.VscodeLaunch.Uri)
		}
	case command.VscodeTask != nil:
		if command.VscodeTask.Uri != "" {
			return ValidateURI(command.VscodeTask.Uri)
		}
	default:
		err = fmt.Errorf("command %s type is invalid", command.Id)
	}

	return err
}

// validateGroup validates commands belonging to a specific group kind. If there are multiple commands belonging to the same group:
// 1. without any default, err out
// 2. with more than one default, err out
func validateGroup(commands []v1alpha2.Command) (err error) {
	defaultCommandCount := 0

	if len(commands) > 1 {
		for _, command := range commands {
			if getGroup(command).IsDefault {
				defaultCommandCount++
			}
		}
	} else {
		// if there is only one command, it is the default command for the group
		defaultCommandCount = 1
	}

	if defaultCommandCount == 0 {
		return fmt.Errorf("there should be exactly one default command, currently there is no default command")
	} else if defaultCommandCount > 1 {
		return fmt.Errorf("there should be exactly one default command, currently there is more than one default command")
	}

	return nil
}

// getGroup returns the group the command belongs to
func getGroup(command v1alpha2.Command) *v1alpha2.CommandGroup {
	switch {
	case command.Composite != nil:
		return command.Composite.Group
	case command.Exec != nil:
		return command.Exec.Group
	case command.Apply != nil:
		return command.Apply.Group
	case command.VscodeLaunch != nil:
		return command.VscodeLaunch.Group
	case command.VscodeTask != nil:
		return command.VscodeTask.Group
	case command.Custom != nil:
		return command.Custom.Group

	default:
		return nil
	}
}

// validateCommandComponent validates the given exec or apply command, the command should map to a valid container component
func validateCommandComponent(command v1alpha2.Command, components []v1alpha2.Component) (err error) {

	if command.Exec == nil && command.Apply == nil {
		return &InvalidCommandError{commandId: command.Id, reason: "should be of type exec or apply"}
	}

	var commandComponent string
	if command.Exec != nil {
		commandComponent = command.Exec.Component
	} else if command.Apply != nil {
		commandComponent = command.Apply.Component
	}

	// must map to a container component
	isComponentValid := false
	for _, component := range components {
		if component.Container != nil && commandComponent == component.Name {
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

		err := validateCommand(subCommand, parentCommands, devfileCommands, components)
		if err != nil {
			return err
		}

		// if subCommand.Composite != nil {
		// 	// Recursively validate the composite subcommand
		// 	err := validateCompositeCommand(&subCommand, parentCommands, devfileCommands, components)
		// 	if err != nil {
		// 		// Don't wrap the error message here to make the error message more readable to the user
		// 		return err
		// 	}
		// } else {
		// 	err := validateCommandComponent(subCommand, components)
		// 	if err != nil {
		// 		return err
		// 	}
		// }
	}
	return nil
}
