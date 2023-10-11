//
//
// Copyright Red Hat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"fmt"

	schema "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

var isTrue bool = true
var isFalse bool = false

// commandAdded adds a new command to the test schema data and notifies the follower
func (testDevFile *TestDevfile) commandAdded(command schema.Command) {
	LogInfoMessage(fmt.Sprintf("command added Id: %s", command.Id))
	testDevFile.SchemaDevFile.Commands = append(testDevFile.SchemaDevFile.Commands, command)
	if testDevFile.Follower != nil {
		testDevFile.Follower.AddCommand(command)
	}
}

// commandUpdated notifies the follower of the command which has been updated
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
	commandGroup.Kind = schema.CommandGroupKind(GetRandomValue(GroupKinds).String())
	LogInfoMessage(fmt.Sprintf("group Kind: %s, default already set %t", commandGroup.Kind, testDevFile.GroupDefaults[commandGroup.Kind]))
	// Ensure only one and at least one of each type are labelled as default
	if !testDevFile.GroupDefaults[commandGroup.Kind] {
		testDevFile.GroupDefaults[commandGroup.Kind] = true
		commandGroup.IsDefault = &isTrue
	} else {
		commandGroup.IsDefault = &isFalse
	}
	LogInfoMessage(fmt.Sprintf("group isDefault: %t", *commandGroup.IsDefault))
	return &commandGroup
}

// AddCommand creates a command of a specified type in a schema structure and pupulates it with random attributes
func (testDevFile *TestDevfile) AddCommand(commandType schema.CommandType) schema.Command {

	var command *schema.Command
	switch commandType {
	case schema.ExecCommandType:
		command = testDevFile.createExecCommand()
		testDevFile.SetExecCommandValues(command)
	case schema.CompositeCommandType:
		command = testDevFile.createCompositeCommand()
		testDevFile.SetCompositeCommandValues(command)
	case schema.ApplyCommandType:
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

// SetExecCommandValues randomly sets/updates exec command attributes to random values
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

	value := GetBinaryDecision()
	execCommand.HotReloadCapable = &value
	LogInfoMessage(fmt.Sprintf("....... HotReloadCapable: %t", *execCommand.HotReloadCapable))

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

// SetCompositeCommandValues randomly sets/updates composite command attributes to random values
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
		compositeCommand.Parallel = &isTrue
		LogInfoMessage(fmt.Sprintf("....... Parallel: %t", *compositeCommand.Parallel))
	}

	testDevFile.commandUpdated(*command)
}

// SetCompositeCommandCommands set the commands in a composite command to a specific type
func (testDevFile *TestDevfile) SetCompositeCommandCommands(command *schema.Command, commandType schema.CommandType) {
	compositeCommand := command.Composite
	compositeCommand.Commands = nil
	numCommands := GetRandomNumber(1, 3)
	for i := 0; i < numCommands; i++ {
		command := testDevFile.AddCommand(commandType)
		compositeCommand.Commands = append(compositeCommand.Commands, command.Id)
		LogInfoMessage(fmt.Sprintf("....... command %d of %d : %s", i, numCommands, command.Id))
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

// SetApplyCommandValues randomly sets/updates apply command attributes to random values
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
