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

package validation

import (
	"github.com/devfile/api/v2/pkg/attributes"
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

var buildGroup = v1alpha2.BuildCommandGroupKind
var runGroup = v1alpha2.RunCommandGroupKind
var isTrue bool = true

// generateDummyExecCommand returns a dummy exec command for testing
func generateDummyExecCommand(name, component string, group *v1alpha2.CommandGroup) v1alpha2.Command {
	return v1alpha2.Command{
		Id: name,
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				LabeledCommand: v1alpha2.LabeledCommand{
					BaseCommand: v1alpha2.BaseCommand{
						Group: group,
					},
				},
				CommandLine: "command",
				Component:   component,
				WorkingDir:  "workDir",
			},
		},
	}
}

// generateDummyExecCommand returns a dummy apply command for testing
func generateDummyApplyCommand(name, component string, group *v1alpha2.CommandGroup, cmdAttributes attributes.Attributes) v1alpha2.Command {
	return v1alpha2.Command{
		Attributes: cmdAttributes,
		Id:         name,
		CommandUnion: v1alpha2.CommandUnion{
			Apply: &v1alpha2.ApplyCommand{
				LabeledCommand: v1alpha2.LabeledCommand{
					BaseCommand: v1alpha2.BaseCommand{
						Group: group,
					},
				},
				Component: component,
			},
		},
	}
}

// generateDummyCompositeCommand returns a dummy composite command for testing
func generateDummyCompositeCommand(name string, commands []string, group *v1alpha2.CommandGroup) v1alpha2.Command {
	return v1alpha2.Command{
		Id: name,
		CommandUnion: v1alpha2.CommandUnion{
			Composite: &v1alpha2.CompositeCommand{
				LabeledCommand: v1alpha2.LabeledCommand{
					BaseCommand: v1alpha2.BaseCommand{
						Group: group,
					},
				},
				Commands: commands,
			},
		},
	}
}

func TestValidateCommands(t *testing.T) {

	component := "alias1"

	components := []v1alpha2.Component{
		generateDummyContainerComponent(component, nil, nil, nil, v1alpha2.Annotation{}, false),
	}

	duplicateKeyErr := "duplicate key: somecommand1"
	noDefaultCmdErr := ".*there should be exactly one default command, currently there is no default command"
	multipleDefaultCmdErr := ".*there should be exactly one default command, currently there are multiple default commands"
	invalidCmdErr := ".*command does not map to a valid component"
	nonExistCmdInComposite := "the command .* mentioned in the composite command does not exist in the devfile"

	parentOverridesFromMainDevfile := attributes.Attributes{}.PutString(ImportSourceAttribute,
		"uri: http://127.0.0.1:8080").PutString(ParentOverrideAttribute, "main devfile")
	invalidCmdErrWithImportAttributes := ".*command does not map to a valid component, imported from uri: http://127.0.0.1:8080, in parent overrides from main devfile"

	tests := []struct {
		name     string
		commands []v1alpha2.Command
		wantErr  []string
	}{
		{
			name: "Valid Exec Command",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("command", component, &v1alpha2.CommandGroup{Kind: runGroup}),
			},
		},
		{
			name: "Valid Composite Command",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("somecommand1", component, nil),
				generateDummyExecCommand("somecommand2", component, nil),
				generateDummyCompositeCommand("composite1", []string{"somecommand1", "somecommand2"}, &v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: &isTrue}),
			},
		},
		{
			name: "Duplicate commands",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("somecommand1", component, nil),
				generateDummyExecCommand("somecommand1", component, nil),
			},
			wantErr: []string{duplicateKeyErr},
		},
		{
			name: "Multiple errors: Duplicate commands, non-exist command in composite command",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("somecommand1", component, nil),
				generateDummyCompositeCommand("somecommand1", []string{"fakecommand"}, &v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: &isTrue}),
			},
			wantErr: []string{duplicateKeyErr, nonExistCmdInComposite},
		},
		{
			name: "Different command types belonging to the same group but no default",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("somecommand1", component, &v1alpha2.CommandGroup{Kind: buildGroup}),
				generateDummyCompositeCommand("somecommand2", []string{"somecommand1"}, &v1alpha2.CommandGroup{Kind: buildGroup}),
			},
			wantErr: []string{noDefaultCmdErr},
		},
		{
			name: "Different command types belonging to the same group with more than one default",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("somecommand1", component, &v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: &isTrue}),
				generateDummyExecCommand("somecommand3", component, &v1alpha2.CommandGroup{Kind: buildGroup}),
				generateDummyCompositeCommand("somecommand2", []string{"somecommand1"}, &v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: &isTrue}),
			},
			wantErr: []string{multipleDefaultCmdErr},
		},
		{
			name: "Invalid Apply command with wrong component",
			commands: []v1alpha2.Command{
				generateDummyApplyCommand("command", "invalidComponent", nil, attributes.Attributes{}),
			},
			wantErr: []string{invalidCmdErr},
		},
		{
			name: "Valid Composite command with one subcommand referencing the other",
			commands: []v1alpha2.Command{
				generateDummyCompositeCommand("composite-1", []string{"composite-a", "composite-b"}, nil),
				generateDummyExecCommand("basic-exec", component, nil),
				generateDummyCompositeCommand("composite-b", []string{"composite-a"}, nil),
				generateDummyCompositeCommand("composite-a", []string{"basic-exec"}, nil),
			},
		},
		{
			name: "Invalid command with import source attribute",
			commands: []v1alpha2.Command{
				generateDummyApplyCommand("command", "invalidComponent", nil, parentOverridesFromMainDevfile),
			},
			wantErr: []string{invalidCmdErrWithImportAttributes},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCommands(tt.commands, components)

			if merr, ok := err.(*multierror.Error); ok && tt.wantErr != nil {
				assert.Equal(t, len(tt.wantErr), len(merr.Errors), "Error list length should match")
				for i := 0; i < len(merr.Errors); i++ {
					assert.Regexp(t, tt.wantErr[i], merr.Errors[i].Error(), "Error message should match")
				}
			} else {
				assert.Equal(t, nil, err, "Error should be nil")
			}
		})
	}
}

func TestValidateCommandComponent(t *testing.T) {

	containerComponent := "alias1"
	kubeComponent := "alias2"
	openshiftComponent := "alias3"
	imageComponent := "alias4"
	volumeComponent := "alias5"
	nonexistComponent := "garbagealias"

	components := []v1alpha2.Component{
		generateDummyContainerComponent(containerComponent, nil, nil, nil, v1alpha2.Annotation{}, false),
		generateDummyKubernetesComponent(kubeComponent, nil, ""),
		generateDummyOpenshiftComponent(openshiftComponent, nil, ""),
		generateDummyImageComponent(imageComponent, v1alpha2.DockerfileSrc{}),
		generateDummyVolumeComponent(volumeComponent, ""),
	}

	invalidCmdErr := ".*command does not map to a valid component"

	tests := []struct {
		name    string
		command v1alpha2.Command
		wantErr *string
	}{
		{
			name:    "Valid Exec Command",
			command: generateDummyExecCommand("command", containerComponent, &v1alpha2.CommandGroup{Kind: runGroup}),
		},
		{
			name:    "Invalid Exec Command with missing component",
			command: generateDummyExecCommand("command", "", &v1alpha2.CommandGroup{Kind: runGroup}),
			wantErr: &invalidCmdErr,
		},
		{
			name:    "Exec Command with non-exist component",
			command: generateDummyExecCommand("command", nonexistComponent, &v1alpha2.CommandGroup{Kind: runGroup}),
			wantErr: &invalidCmdErr,
		},
		{
			name:    "Exec Command with image component",
			command: generateDummyExecCommand("command", imageComponent, &v1alpha2.CommandGroup{Kind: runGroup}),
			wantErr: &invalidCmdErr,
		},
		{
			name:    "Exec Command with kubernetes component",
			command: generateDummyExecCommand("command", kubeComponent, &v1alpha2.CommandGroup{Kind: runGroup}),
			wantErr: &invalidCmdErr,
		},
		{
			name:    "Exec Command with openshift component",
			command: generateDummyExecCommand("command", openshiftComponent, &v1alpha2.CommandGroup{Kind: runGroup}),
			wantErr: &invalidCmdErr,
		},
		{
			name:    "Exec Command with volume component",
			command: generateDummyExecCommand("command", volumeComponent, &v1alpha2.CommandGroup{Kind: runGroup}),
			wantErr: &invalidCmdErr,
		},
		{
			name:    "Valid Exec Command with Group nil",
			command: generateDummyExecCommand("command", containerComponent, nil),
		},
		{
			name:    "Valid Apply Command with container component",
			command: generateDummyApplyCommand("command", containerComponent, nil, attributes.Attributes{}),
		},
		{
			name:    "Valid Apply Command with image component",
			command: generateDummyApplyCommand("command", imageComponent, nil, attributes.Attributes{}),
		},
		{
			name:    "Valid Apply Command with kubernetes component",
			command: generateDummyApplyCommand("command", kubeComponent, nil, attributes.Attributes{}),
		},
		{
			name:    "Valid Apply Command with openshift component",
			command: generateDummyApplyCommand("command", openshiftComponent, nil, attributes.Attributes{}),
		},
		{
			name:    "Apply Command with non-exist component",
			command: generateDummyApplyCommand("command", nonexistComponent, &v1alpha2.CommandGroup{Kind: runGroup}, attributes.Attributes{}),
			wantErr: &invalidCmdErr,
		},
		{
			name:    "Apply Command with volume component",
			command: generateDummyApplyCommand("command", volumeComponent, &v1alpha2.CommandGroup{Kind: runGroup}, attributes.Attributes{}),
			wantErr: &invalidCmdErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCommandComponent(tt.command, components)
			if tt.wantErr != nil && assert.Error(t, err) {
				assert.Regexp(t, *tt.wantErr, err.Error(), "Error message should match")
			} else {
				assert.NoError(t, err, "Expected error to be nil")
			}
		})
	}
}

func TestValidateCompositeCommand(t *testing.T) {

	component := "alias1"
	validExecCommands := []v1alpha2.Command{
		generateDummyExecCommand("command1", component, &v1alpha2.CommandGroup{Kind: runGroup}),
		generateDummyExecCommand("command2", component, &v1alpha2.CommandGroup{Kind: buildGroup}),
		generateDummyExecCommand("command3", component, &v1alpha2.CommandGroup{Kind: runGroup}),
	}
	components := []v1alpha2.Component{
		generateDummyContainerComponent(component, nil, nil, nil, v1alpha2.Annotation{}, false),
	}

	invalidCmdErr := ".*command does not map to a valid component"
	missingCmdErr := ".*the command .* mentioned in the composite command does not exist in the devfile"
	selfRefCmdErr := ".*composite command cannot reference itself"
	indirectRefCmdErr := "composite command cannot indirectly reference itself"

	tests := []struct {
		name                 string
		commands             []v1alpha2.Command
		testCompositeCommand string
		wantErr              *string
	}{
		{
			name: "Valid Composite Command",
			commands: append(validExecCommands,
				generateDummyCompositeCommand("command4", []string{"command1", "command2", "command3"}, &v1alpha2.CommandGroup{Kind: buildGroup})),
			testCompositeCommand: "command4",
		},
		{
			name: "Invalid composite command, references non-existent command",
			commands: append(validExecCommands,
				generateDummyCompositeCommand("command4", []string{"command1", "fakecommand", "command3"}, &v1alpha2.CommandGroup{Kind: buildGroup})),
			testCompositeCommand: "command4",
			wantErr:              &missingCmdErr,
		},
		{
			name: "Invalid composite command, references itself",
			commands: append(validExecCommands,
				generateDummyCompositeCommand("command4", []string{"command1", "command4", "command3"}, &v1alpha2.CommandGroup{Kind: buildGroup})),
			testCompositeCommand: "command4",
			wantErr:              &selfRefCmdErr,
		},
		{
			name: "Invalid composite command, indirectly references itself",
			commands: append(validExecCommands,
				generateDummyCompositeCommand("command4", []string{"command5", "command3"}, &v1alpha2.CommandGroup{Kind: buildGroup}),
				generateDummyCompositeCommand("command5", []string{"command1", "command4", "command3"}, &v1alpha2.CommandGroup{Kind: buildGroup})),
			testCompositeCommand: "command4",
			wantErr:              &indirectRefCmdErr,
		},
		{
			name: "Invalid composite command, points to invalid exec command",
			commands: []v1alpha2.Command{
				generateDummyCompositeCommand("command4", []string{"command1", "command2"}, &v1alpha2.CommandGroup{Kind: buildGroup}),
				generateDummyExecCommand("command1", component, &v1alpha2.CommandGroup{Kind: runGroup}),
				generateDummyExecCommand("command2", "some-fake-component", &v1alpha2.CommandGroup{Kind: buildGroup}),
			},
			testCompositeCommand: "command4",
			wantErr:              &invalidCmdErr,
		},
		{
			name: "Valid Composite command with one subcommand referencing the other",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("basic-exec", component, nil),
				generateDummyCompositeCommand("composite-a", []string{"basic-exec"}, nil),
				generateDummyCompositeCommand("composite-b", []string{"composite-a"}, nil),
				generateDummyCompositeCommand("composite-1", []string{"composite-a", "composite-b"}, nil),
			},
			testCompositeCommand: "composite-1",
		},
	}
	for _, tt := range tests {
		commandMap := getCommandsMap(tt.commands)

		t.Run(tt.name, func(t *testing.T) {
			cmd := commandMap[tt.testCompositeCommand]
			parentCommands := make(map[string]string)

			err := validateCompositeCommand(&cmd, parentCommands, commandMap, components)
			if tt.wantErr != nil && assert.Error(t, err) {
				assert.Regexp(t, *tt.wantErr, err.Error(), "Error message should match")
			} else {
				assert.NoError(t, err, "Expected error to be nil")
			}
		})
	}
}

func TestValidateGroup(t *testing.T) {

	component := "alias1"

	noDefaultCmdErr := ".*there should be exactly one default command, currently there is no default command"
	multipleDefaultError := ".*there should be exactly one default command, currently there are multiple default commands"
	multipleDefaultCmdErr := multipleDefaultError + "; command: run command; command: customcommand"

	parentOverridesFromMainDevfile := attributes.Attributes{}.PutString(ImportSourceAttribute,
		"uri: http://127.0.0.1:8080").PutString(ParentOverrideAttribute, "main devfile")
	multipleDefaultCmdErrWithImportAttributes := multipleDefaultError +
		"; command: run command; command: customcommand, imported from uri: http://127.0.0.1:8080, in parent overrides from main devfile"

	tests := []struct {
		name     string
		commands []v1alpha2.Command
		group    v1alpha2.CommandGroupKind
		wantErr  *string
	}{
		{
			name: "Two default run commands",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("run command", component, &v1alpha2.CommandGroup{Kind: runGroup, IsDefault: &isTrue}),
				generateDummyExecCommand("customcommand", component, &v1alpha2.CommandGroup{Kind: runGroup, IsDefault: &isTrue}),
			},
			group:   runGroup,
			wantErr: &multipleDefaultCmdErr,
		},
		{
			name: "Two default run commands with import source attribute",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("run command", component, &v1alpha2.CommandGroup{Kind: runGroup, IsDefault: &isTrue}),
				generateDummyApplyCommand("customcommand", component, &v1alpha2.CommandGroup{Kind: runGroup, IsDefault: &isTrue}, parentOverridesFromMainDevfile),
			},
			group:   runGroup,
			wantErr: &multipleDefaultCmdErrWithImportAttributes,
		},
		{
			name: "No default for more than one build commands",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("build command", component, &v1alpha2.CommandGroup{Kind: buildGroup}),
				generateDummyExecCommand("build command 2", component, &v1alpha2.CommandGroup{Kind: buildGroup}),
			},
			group:   buildGroup,
			wantErr: &noDefaultCmdErr,
		},
		{
			name: "One command does not need default",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("test command", component, &v1alpha2.CommandGroup{Kind: buildGroup}),
			},
			group: buildGroup,
		},
		{
			name: "One command can have default",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("test command", component, &v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: &isTrue}),
			},
			group: buildGroup,
		},
		{
			name: "Composite commands in group",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("build command", component, &v1alpha2.CommandGroup{Kind: buildGroup}),
				generateDummyExecCommand("build command 2", component, &v1alpha2.CommandGroup{Kind: buildGroup}),
				generateDummyCompositeCommand("composite1", []string{"build command", "build command 2"}, &v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: &isTrue}),
			},
			group: buildGroup,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGroup(tt.commands, tt.group)
			if tt.wantErr != nil && assert.Error(t, err) {
				assert.Regexp(t, *tt.wantErr, err.Error(), "Error message should match")
			} else {
				assert.NoError(t, err, "Expected error to be nil")
			}
		})
	}
}
