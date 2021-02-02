package validation

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
)

var buildGroup = v1alpha2.BuildCommandGroupKind
var runGroup = v1alpha2.RunCommandGroupKind

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
func generateDummyApplyCommand(name, component string, group *v1alpha2.CommandGroup) v1alpha2.Command {
	return v1alpha2.Command{
		Id: name,
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

// generateDummyCompositeCommand returns a dummy VscodeLaunch command for testing
func generateDummyVscodeLaunchCommand(name string, commandLocation v1alpha2.VscodeConfigurationCommandLocation, group *v1alpha2.CommandGroup) v1alpha2.Command {
	return v1alpha2.Command{
		Id: name,
		CommandUnion: v1alpha2.CommandUnion{
			VscodeLaunch: &v1alpha2.VscodeConfigurationCommand{
				BaseCommand: v1alpha2.BaseCommand{
					Group: group,
				},
				VscodeConfigurationCommandLocation: commandLocation,
			},
		},
	}
}

// generateDummyCompositeCommand returns a dummy VscodeTask command for testing
func generateDummyVscodeTaskCommand(name string, commandLocation v1alpha2.VscodeConfigurationCommandLocation, group *v1alpha2.CommandGroup) v1alpha2.Command {
	return v1alpha2.Command{
		Id: name,
		CommandUnion: v1alpha2.CommandUnion{
			VscodeTask: &v1alpha2.VscodeConfigurationCommand{
				BaseCommand: v1alpha2.BaseCommand{
					Group: group,
				},
				VscodeConfigurationCommandLocation: commandLocation,
			},
		},
	}
}

func TestValidateCommands(t *testing.T) {

	component := "alias1"

	components := []v1alpha2.Component{
		generateDummyContainerComponent(component, nil, nil, nil),
	}

	uriCommandLocation := v1alpha2.VscodeConfigurationCommandLocation{
		Uri: "/some/path",
	}
	inValidUriCommandLocation := v1alpha2.VscodeConfigurationCommandLocation{
		Uri: "http//wronguri",
	}
	inlinedCommandLocation := v1alpha2.VscodeConfigurationCommandLocation{
		Inlined: "inlined code",
	}

	duplicateKeyErr := "duplicate key: somecommand1"
	noDefaultCmdErr := ".*there should be exactly one default command, currently there is no default command"
	multipleDefaultCmdErr := ".*there should be exactly one default command, currently there is more than one default command"
	invalidURIErr := "invalid URI for request"
	invalidCmdErr := ".*command does not map to a container component"

	tests := []struct {
		name     string
		commands []v1alpha2.Command
		wantErr  *string
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
				generateDummyCompositeCommand("composite1", []string{"somecommand1", "somecommand2"}, &v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: true}),
			},
		},
		{
			name: "Duplicate commands",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("somecommand1", component, nil),
				generateDummyExecCommand("somecommand1", component, nil),
			},
			wantErr: &duplicateKeyErr,
		},
		{
			name: "Duplicate commands, different types",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("somecommand1", component, nil),
				generateDummyCompositeCommand("somecommand1", []string{"fakecommand"}, &v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: true}),
			},
			wantErr: &duplicateKeyErr,
		},
		{
			name: "Different command types belonging to the same group but no default",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("somecommand1", component, &v1alpha2.CommandGroup{Kind: buildGroup}),
				generateDummyCompositeCommand("somecommand2", []string{"somecommand1"}, &v1alpha2.CommandGroup{Kind: buildGroup}),
			},
			wantErr: &noDefaultCmdErr,
		},
		{
			name: "Different command types belonging to the same group with more than one default",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("somecommand1", component, &v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: true}),
				generateDummyExecCommand("somecommand3", component, &v1alpha2.CommandGroup{Kind: buildGroup}),
				generateDummyCompositeCommand("somecommand2", []string{"somecommand1"}, &v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: true}),
			},
			wantErr: &multipleDefaultCmdErr,
		},
		{
			name: "Valid VscodeTask command with URI",
			commands: []v1alpha2.Command{
				generateDummyVscodeTaskCommand("somevscodetask", uriCommandLocation, nil),
			},
		},
		{
			name: "Valid VscodeLaunch command with Inlined",
			commands: []v1alpha2.Command{
				generateDummyVscodeLaunchCommand("somevscodelaunch", inlinedCommandLocation, nil),
			},
		},
		{
			name: "Invalid VscodeLaunch command with wrong URI",
			commands: []v1alpha2.Command{
				generateDummyVscodeLaunchCommand("somevscodelaunch", inValidUriCommandLocation, nil),
			},
			wantErr: &invalidURIErr,
		},
		{
			name: "Invalid Apply command with wrong component",
			commands: []v1alpha2.Command{
				generateDummyApplyCommand("command", "invalidComponent", nil),
			},
			wantErr: &invalidCmdErr,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCommands(tt.commands, components)
			if tt.wantErr != nil && assert.Error(t, err) {
				assert.Regexp(t, *tt.wantErr, err.Error(), "Error message should match")
			} else {
				assert.NoError(t, err, "Expected error to be nil")
			}
		})
	}
}

func TestValidateCommandComponent(t *testing.T) {

	component := "alias1"
	invalidComponent := "garbagealias"

	components := []v1alpha2.Component{
		generateDummyContainerComponent(component, nil, nil, nil),
	}

	invalidCmdErr := ".*command does not map to a container component"

	tests := []struct {
		name    string
		command v1alpha2.Command
		wantErr *string
	}{
		{
			name:    "Valid Exec Command",
			command: generateDummyExecCommand("command", component, &v1alpha2.CommandGroup{Kind: runGroup}),
		},
		{
			name:    "Invalid Exec Command with missing component",
			command: generateDummyExecCommand("command", "", &v1alpha2.CommandGroup{Kind: runGroup}),
			wantErr: &invalidCmdErr,
		},
		{
			name:    "Valid Exec Command with invalid component",
			command: generateDummyExecCommand("command", invalidComponent, &v1alpha2.CommandGroup{Kind: runGroup}),
			wantErr: &invalidCmdErr,
		},
		{
			name:    "Valid Exec Command with Group nil",
			command: generateDummyExecCommand("command", component, nil),
		},
		{
			name:    "Valid Apply Command",
			command: generateDummyApplyCommand("command", component, nil),
		},
		{
			name:    "Invalid Apply Command with wrong component",
			command: generateDummyApplyCommand("command", invalidComponent, &v1alpha2.CommandGroup{Kind: runGroup}),
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
		generateDummyContainerComponent(component, nil, nil, nil),
	}

	invalidCmdErr := ".*command does not map to a container component"
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
	multipleDefaultCmdErr := ".*there should be exactly one default command, currently there is more than one default command"

	tests := []struct {
		name     string
		commands []v1alpha2.Command
		wantErr  *string
	}{
		{
			name: "Two default run commands",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("run command", component, &v1alpha2.CommandGroup{Kind: runGroup, IsDefault: true}),
				generateDummyExecCommand("customcommand", component, &v1alpha2.CommandGroup{Kind: runGroup, IsDefault: true}),
			},
			wantErr: &multipleDefaultCmdErr,
		},
		{
			name: "No default for more than one build commands",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("build command", component, &v1alpha2.CommandGroup{Kind: buildGroup}),
				generateDummyExecCommand("build command 2", component, &v1alpha2.CommandGroup{Kind: buildGroup}),
			},
			wantErr: &noDefaultCmdErr,
		},
		{
			name: "One command does not need default",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("test command", component, &v1alpha2.CommandGroup{Kind: buildGroup}),
			},
		},
		{
			name: "One command can have default",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("test command", component, &v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: true}),
			},
		},
		{
			name: "Composite commands in group",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("build command", component, &v1alpha2.CommandGroup{Kind: buildGroup}),
				generateDummyExecCommand("build command 2", component, &v1alpha2.CommandGroup{Kind: buildGroup}),
				generateDummyCompositeCommand("composite1", []string{"build command", "build command 2"}, &v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: true}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGroup(tt.commands)
			if tt.wantErr != nil && assert.Error(t, err) {
				assert.Regexp(t, *tt.wantErr, err.Error(), "Error message should match")
			} else {
				assert.NoError(t, err, "Expected error to be nil")
			}
		})
	}
}
