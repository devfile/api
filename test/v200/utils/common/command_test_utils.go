package common

import (
	"fmt"

	schema "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

// commandAdded adds a new command to the test schema data and notifies the follower
func (testDevFile *TestDevfile) commandAdded(command schema.Command) {
	LogInfoMessage(fmt.Sprintf("command added Id: %s", command.Id))
	testDevFile.SchemaDevFile.Commands = append(testDevFile.SchemaDevFile.Commands, command)
	if testDevFile.Follower != nil {
		testDevFile.Follower.AddCommand(command)
	}
}

// commandUpdated and notifies the follower of the command which has been updated
func (testDevFile *TestDevfile) commandUpdated(command schema.Command) {
	LogInfoMessage(fmt.Sprintf("command updated Id: %s", command.Id))
	if testDevFile.Follower != nil {
		testDevFile.Follower.UpdateCommand(command)
	}
}

// addEnv creates and returns a specifed number of env attributes in a schema structure
func addEnv(numEnv int) []schema.EnvVar {
	commandEnvs := make([]schema.EnvVar, numEnv)
	for i := 0; i < numEnv; i++ {
		commandEnvs[i].Name = "Name_" + GetRandomString(5, false)
		commandEnvs[i].Value = "Value_" + GetRandomString(5, false)
		LogInfoMessage(fmt.Sprintf("Add Env: %s", commandEnvs[i]))
	}
	return commandEnvs
}

// addAttributes creates returns a specifed number of attributes in a schema structure
func addAttributes(numAtrributes int) map[string]string {
	attributes := make(map[string]string)
	for i := 0; i < numAtrributes; i++ {
		AttributeName := "Name_" + GetRandomString(6, false)
		attributes[AttributeName] = "Value_" + GetRandomString(6, false)
		LogInfoMessage(fmt.Sprintf("Add attribute : %s = %s", AttributeName, attributes[AttributeName]))
	}
	return attributes
}

// addGroup creates and returns a group in a schema structure
func (testDevFile *TestDevfile) addGroup() *schema.CommandGroup {

	commandGroup := schema.CommandGroup{}
	commandGroup.Kind = GetRandomGroupKind()
	LogInfoMessage(fmt.Sprintf("group Kind: %s, default already set %t", commandGroup.Kind, testDevFile.GroupDefaults[commandGroup.Kind]))
	// Ensure only one and at least one of each type are labelled as default
	if !testDevFile.GroupDefaults[commandGroup.Kind] {
		testDevFile.GroupDefaults[commandGroup.Kind] = true
		commandGroup.IsDefault = true
	} else {
		commandGroup.IsDefault = false
	}
	LogInfoMessage(fmt.Sprintf("group isDefault: %t", commandGroup.IsDefault))
	return &commandGroup
}

// AddCommand creates a command of a specified type in a schema structure and pupulates it with random attributes
func (testDevFile *TestDevfile) AddCommand(commandType schema.CommandType) schema.Command {

	var command *schema.Command
	if commandType == schema.ExecCommandType {
		command = testDevFile.createExecCommand()
		testDevFile.SetExecCommandValues(command)
	} else if commandType == schema.CompositeCommandType {
		command = testDevFile.createCompositeCommand()
		testDevFile.SetCompositeCommandValues(command)
	} else if commandType == schema.ApplyCommandType {
		command = testDevFile.createApplyCommand()
		testDevFile.SetApplyCommandValues(command)
	}
	return *command
}

// createExecCommand creates and returns an empty exec command in a schema structure
func (testDevFile *TestDevfile) createExecCommand() *schema.Command {

	LogInfoMessage("Create an exec command :")
	command := schema.Command{}
	command.Id = GetRandomUniqueString(8, true)
	LogInfoMessage(fmt.Sprintf("command Id: %s", command.Id))
	command.Exec = &schema.ExecCommand{}
	testDevFile.commandAdded(command)
	return &command

}

// SetExecCommandValues randomly sets exec command attribute to random values
func (testDevFile *TestDevfile) SetExecCommandValues(command *schema.Command) {

	execCommand := command.Exec

	// exec command must be mentioned by a container component
	execCommand.Component = testDevFile.GetContainerName()

	execCommand.CommandLine = GetRandomString(4, false) + " " + GetRandomString(4, false)
	LogInfoMessage(fmt.Sprintf("....... commandLine: %s", execCommand.CommandLine))

	// If group already leave it to make sure defaults are not deleted or added
	if execCommand.Group == nil {
		if GetRandomDecision(2, 1) {
			execCommand.Group = testDevFile.addGroup()
		}
	}

	if GetBinaryDecision() {
		execCommand.Label = GetRandomString(12, false)
		LogInfoMessage(fmt.Sprintf("....... label: %s", execCommand.Label))
	} else {
		execCommand.Label = ""
	}

	if GetBinaryDecision() {
		execCommand.WorkingDir = "./tmp"
		LogInfoMessage(fmt.Sprintf("....... WorkingDir: %s", execCommand.WorkingDir))
	} else {
		execCommand.WorkingDir = ""
	}

	execCommand.HotReloadCapable = GetBinaryDecision()
	LogInfoMessage(fmt.Sprintf("....... HotReloadCapable: %t", execCommand.HotReloadCapable))

	if GetBinaryDecision() {
		execCommand.Env = addEnv(GetRandomNumber(1, 4))
	} else {
		execCommand.Env = nil
	}
	testDevFile.commandUpdated(*command)

}

// createCompositeCommand creates an empty composite command in a schema structure
func (testDevFile *TestDevfile) createCompositeCommand() *schema.Command {

	LogInfoMessage("Create a composite command :")
	command := schema.Command{}
	command.Id = GetRandomUniqueString(8, true)
	LogInfoMessage(fmt.Sprintf("command Id: %s", command.Id))
	command.Composite = &schema.CompositeCommand{}
	testDevFile.commandAdded(command)

	return &command
}

// SetCompositeCommandValues randomly sets composite command attribute to random values
func (testDevFile *TestDevfile) SetCompositeCommandValues(command *schema.Command) {

	compositeCommand := command.Composite
	numCommands := GetRandomNumber(1, 3)

	for i := 0; i < numCommands; i++ {
		execCommand := testDevFile.AddCommand(schema.ExecCommandType)
		compositeCommand.Commands = append(compositeCommand.Commands, execCommand.Id)
		LogInfoMessage(fmt.Sprintf("....... command %d of %d : %s", i, numCommands, execCommand.Id))
	}

	// If group already exists - leave it to make sure defaults are not deleted or added
	if compositeCommand.Group == nil {
		if GetRandomDecision(2, 1) {
			compositeCommand.Group = testDevFile.addGroup()
		}
	}

	if GetBinaryDecision() {
		compositeCommand.Label = GetRandomString(12, false)
		LogInfoMessage(fmt.Sprintf("....... label: %s", compositeCommand.Label))
	}

	if GetBinaryDecision() {
		compositeCommand.Parallel = true
		LogInfoMessage(fmt.Sprintf("....... Parallel: %t", compositeCommand.Parallel))
	}

	testDevFile.commandUpdated(*command)
}

// createApplyCommand creates an apply command in a schema structure
func (testDevFile *TestDevfile) createApplyCommand() *schema.Command {

	LogInfoMessage("Create a apply command :")
	command := schema.Command{}
	command.Id = GetRandomUniqueString(8, true)
	LogInfoMessage(fmt.Sprintf("command Id: %s", command.Id))
	command.Apply = &schema.ApplyCommand{}
	testDevFile.commandAdded(command)
	return &command
}

// SetApplyCommandValues randomly sets apply command attributes to random values
func (testDevFile *TestDevfile) SetApplyCommandValues(command *schema.Command) {
	applyCommand := command.Apply

	applyCommand.Component = testDevFile.GetContainerName()

	if GetRandomDecision(2, 1) {
		applyCommand.Group = testDevFile.addGroup()
	}

	if GetBinaryDecision() {
		applyCommand.Label = GetRandomString(63, false)
		LogInfoMessage(fmt.Sprintf("....... label: %s", applyCommand.Label))
	}

	testDevFile.commandUpdated(*command)
}
