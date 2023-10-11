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

// parentCommandAdded adds a new command to the test schema data
func (testDevFile *TestDevfile) parentCommandAdded(command schema.CommandParentOverride) {
	LogInfoMessage(fmt.Sprintf("parent command added Id: %s", command.Id))
	testDevFile.SchemaDevFile.Parent.Commands = append(testDevFile.SchemaDevFile.Parent.Commands, command)
}

// addParentEnv creates and returns a specifed number of env attributes in a schema structure
func addParentEnv(numEnv int) []schema.EnvVarParentOverride {
	commandEnvs := make([]schema.EnvVarParentOverride, numEnv)
	for i := 0; i < numEnv; i++ {
		commandEnvs[i].Name = "Name_" + GetRandomString(5, false)
		commandEnvs[i].Value = "Value_" + GetRandomString(5, false)
		LogInfoMessage(fmt.Sprintf("Add Parent Env: %s", commandEnvs[i]))
	}
	return commandEnvs
}

// addParentGroup creates and returns a group in a schema structure
func (testDevFile *TestDevfile) addParentGroup() *schema.CommandGroupParentOverride {

	commandGroup := schema.CommandGroupParentOverride{}
	commandGroup.Kind = schema.CommandGroupKindParentOverride(GetRandomValue(GroupKinds).String())
	kind := schema.CommandGroupKind(commandGroup.Kind)
	LogInfoMessage(fmt.Sprintf("parent group Kind: %s, default already set %t", commandGroup.Kind, testDevFile.GroupDefaults[kind]))
	// Ensure only one and at least one of each type are labelled as default
	if !testDevFile.GroupDefaults[schema.CommandGroupKind(kind)] {
		testDevFile.GroupDefaults[schema.CommandGroupKind(kind)] = true
		commandGroup.IsDefault = &isTrue
	} else {
		commandGroup.IsDefault = &isFalse
	}
	LogInfoMessage(fmt.Sprintf("parent group isDefault: %t", *commandGroup.IsDefault))
	return &commandGroup
}

// AddParentCommand creates a command of a specified type in a schema structure and populates it with random values
func (testDevFile *TestDevfile) AddParentCommand(commandType schema.CommandType) schema.CommandParentOverride {

	var command *schema.CommandParentOverride
	switch commandType {
	case schema.ExecCommandType:
		command = testDevFile.createParentExecCommand()
		testDevFile.SetParentExecCommandValues(command)
	case schema.CompositeCommandType:
		command = testDevFile.createParentCompositeCommand()
		testDevFile.SetParentCompositeCommandValues(command)
	case schema.ApplyCommandType:
		command = testDevFile.createParentApplyCommand()
		testDevFile.SetParentApplyCommandValues(command)
	}
	return *command
}

// createParentExecCommand creates and returns an empty exec command in a schema structure
func (testDevFile *TestDevfile) createParentExecCommand() *schema.CommandParentOverride {

	LogInfoMessage("Create a parent exec command :")
	command := schema.CommandParentOverride{}
	command.Id = GetRandomUniqueString(8, true)
	LogInfoMessage(fmt.Sprintf("command Id: %s", command.Id))
	command.Exec = &schema.ExecCommandParentOverride{}
	testDevFile.parentCommandAdded(command)
	return &command

}

// SetParentExecCommandValues randomly sets/updates exec command attributes to random values
func (testDevFile *TestDevfile) SetParentExecCommandValues(command *schema.CommandParentOverride) {

	execCommand := command.Exec

	// exec command must be mentioned by a container component
	execCommand.Component = testDevFile.GetParentContainerName()

	execCommand.CommandLine = GetRandomString(4, false) + " " + GetRandomString(4, false)
	LogInfoMessage(fmt.Sprintf("....... commandLine: %s", execCommand.CommandLine))

	// If group already leave it to make sure defaults are not deleted or added
	if execCommand.Group == nil {
		if GetRandomDecision(2, 1) {
			execCommand.Group = testDevFile.addParentGroup()
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
		execCommand.Env = addParentEnv(GetRandomNumber(1, 4))
	} else {
		execCommand.Env = nil
	}
	LogInfoMessage(fmt.Sprintf("parent command updated Id: %s", command.Id))

}

// createParentCompositeCommand creates an empty composite command in a schema structure
func (testDevFile *TestDevfile) createParentCompositeCommand() *schema.CommandParentOverride {

	LogInfoMessage("Create a parent composite command :")
	command := schema.CommandParentOverride{}
	command.Id = GetRandomUniqueString(8, true)
	LogInfoMessage(fmt.Sprintf("command Id: %s", command.Id))
	command.Composite = &schema.CompositeCommandParentOverride{}
	testDevFile.parentCommandAdded(command)

	return &command
}

// SetParentCompositeCommandValues randomly sets/updates composite command attributes to random values
func (testDevFile *TestDevfile) SetParentCompositeCommandValues(command *schema.CommandParentOverride) {

	compositeCommand := command.Composite
	numCommands := GetRandomNumber(1, 3)

	for i := 0; i < numCommands; i++ {
		execCommand := testDevFile.AddParentCommand(schema.ExecCommandType)
		compositeCommand.Commands = append(compositeCommand.Commands, execCommand.Id)
		LogInfoMessage(fmt.Sprintf("....... command %d of %d : %s", i, numCommands, execCommand.Id))
	}

	// If group already exists - leave it to make sure defaults are not deleted or added
	if compositeCommand.Group == nil {
		if GetRandomDecision(2, 1) {
			compositeCommand.Group = testDevFile.addParentGroup()
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

	LogInfoMessage(fmt.Sprintf("parent command updated Id: %s", command.Id))
}

// createParentApplyCommand creates an apply command in a schema structure
func (testDevFile *TestDevfile) createParentApplyCommand() *schema.CommandParentOverride {

	LogInfoMessage("Create a parent apply command :")
	command := schema.CommandParentOverride{}
	command.Id = GetRandomUniqueString(8, true)
	LogInfoMessage(fmt.Sprintf("command Id: %s", command.Id))
	command.Apply = &schema.ApplyCommandParentOverride{}
	testDevFile.parentCommandAdded(command)
	return &command
}

// SetApplyCommandValues randomly sets/updates apply command attributes to random values
func (testDevFile *TestDevfile) SetParentApplyCommandValues(command *schema.CommandParentOverride) {
	applyCommand := command.Apply

	applyCommand.Component = testDevFile.GetParentContainerName()

	if GetRandomDecision(2, 1) {
		applyCommand.Group = testDevFile.addParentGroup()
	}

	if GetBinaryDecision() {
		applyCommand.Label = GetRandomString(63, false)
		LogInfoMessage(fmt.Sprintf("....... label: %s", applyCommand.Label))
	}

	LogInfoMessage(fmt.Sprintf("parent command updated Id: %s", command.Id))
}
