package validation

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

var buildGroup = v1alpha2.BuildCommandGroupKind
var runGroup = v1alpha2.RunCommandGroupKind

// generateDummyExecCommand returns a dummy exec command for testing
func generateDummyExecCommand(name, component string, group v1alpha2.CommandGroup) v1alpha2.Command {
	return v1alpha2.Command{
		Id: name,
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				LabeledCommand: v1alpha2.LabeledCommand{
					BaseCommand: v1alpha2.BaseCommand{
						Group: &group,
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
func generateDummyApplyCommand(name, component string, group v1alpha2.CommandGroup) v1alpha2.Command {
	return v1alpha2.Command{
		Id: name,
		CommandUnion: v1alpha2.CommandUnion{
			Apply: &v1alpha2.ApplyCommand{
				LabeledCommand: v1alpha2.LabeledCommand{
					BaseCommand: v1alpha2.BaseCommand{
						Group: &group,
					},
				},
				Component: component,
			},
		},
	}
}

// generateDummyCompositeCommand returns a dummy composite command for testing
func generateDummyCompositeCommand(name string, commands []string, group v1alpha2.CommandGroup) v1alpha2.Command {
	return v1alpha2.Command{
		Id: name,
		CommandUnion: v1alpha2.CommandUnion{
			Composite: &v1alpha2.CompositeCommand{
				LabeledCommand: v1alpha2.LabeledCommand{
					BaseCommand: v1alpha2.BaseCommand{
						Group: &group,
					},
				},
				Commands: commands,
			},
		},
	}
}

// generateDummyCompositeCommand returns a dummy VscodeLaunch command for testing
func generateDummyVscodeLaunchCommand(name string, commandLocation v1alpha2.VscodeConfigurationCommandLocation, group v1alpha2.CommandGroup) v1alpha2.Command {
	return v1alpha2.Command{
		Id: name,
		CommandUnion: v1alpha2.CommandUnion{
			VscodeLaunch: &v1alpha2.VscodeConfigurationCommand{
				BaseCommand: v1alpha2.BaseCommand{
					Group: &group,
				},
				VscodeConfigurationCommandLocation: commandLocation,
			},
		},
	}
}

// generateDummyCompositeCommand returns a dummy VscodeTask command for testing
func generateDummyVscodeTaskCommand(name string, commandLocation v1alpha2.VscodeConfigurationCommandLocation, group v1alpha2.CommandGroup) v1alpha2.Command {
	return v1alpha2.Command{
		Id: name,
		CommandUnion: v1alpha2.CommandUnion{
			VscodeTask: &v1alpha2.VscodeConfigurationCommand{
				BaseCommand: v1alpha2.BaseCommand{
					Group: &group,
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

	tests := []struct {
		name     string
		commands []v1alpha2.Command
		wantErr  bool
	}{
		{
			name: "Case 1: Valid Exec Command",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("command", component, v1alpha2.CommandGroup{Kind: runGroup}),
			},
			wantErr: false,
		},
		{
			name: "Case 2: Valid Composite Command",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("somecommand1", component, v1alpha2.CommandGroup{}),
				generateDummyExecCommand("somecommand2", component, v1alpha2.CommandGroup{}),
				generateDummyCompositeCommand("composite1", []string{"somecommand1", "somecommand2"}, v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: true}),
			},
			wantErr: false,
		},
		{
			name: "Case 3: Duplicate commands",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("somecommand1", component, v1alpha2.CommandGroup{}),
				generateDummyExecCommand("somecommand1", component, v1alpha2.CommandGroup{}),
			},
			wantErr: true,
		},
		{
			name: "Case 4: Duplicate commands, different types",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("somecommand1", component, v1alpha2.CommandGroup{}),
				generateDummyCompositeCommand("somecommand1", []string{"fakecommand"}, v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: true}),
			},
			wantErr: true,
		},
		{
			name: "Case 5: Multiple same group kind but no default",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("somecommand1", component, v1alpha2.CommandGroup{Kind: buildGroup}),
				generateDummyCompositeCommand("somecommand2", []string{"somecommand1"}, v1alpha2.CommandGroup{Kind: buildGroup}),
			},
			wantErr: true,
		},
		{
			name: "Case 6: Multiple same group kind with more than 1 default",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("somecommand1", component, v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: true}),
				generateDummyExecCommand("somecommand3", component, v1alpha2.CommandGroup{Kind: buildGroup}),
				generateDummyCompositeCommand("somecommand2", []string{"somecommand1"}, v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: true}),
			},
			wantErr: true,
		},
		{
			name: "Case 7: Valid VscodeTask command with Uri",
			commands: []v1alpha2.Command{
				generateDummyVscodeTaskCommand("somevscodetask", uriCommandLocation, v1alpha2.CommandGroup{}),
			},
			wantErr: false,
		},
		{
			name: "Case 8: Valid VscodeLaunch command with Inlined",
			commands: []v1alpha2.Command{
				generateDummyVscodeLaunchCommand("somevscodelaunch", inlinedCommandLocation, v1alpha2.CommandGroup{}),
			},
			wantErr: false,
		},
		{
			name: "Case 9: Invalid VscodeLaunch command with wrong Uri",
			commands: []v1alpha2.Command{
				generateDummyVscodeLaunchCommand("somevscodelaunch", inValidUriCommandLocation, v1alpha2.CommandGroup{}),
			},
			wantErr: true,
		},
		{
			name: "Case 10: Invalid apply command with wrong component",
			commands: []v1alpha2.Command{
				generateDummyApplyCommand("command", "invalidComponent", v1alpha2.CommandGroup{}),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCommands(tt.commands, components)
			if !tt.wantErr == (err != nil) {
				t.Errorf("TestValidateAction unexpected error: %v", err)
				return
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

	tests := []struct {
		name    string
		command v1alpha2.Command
		wantErr bool
	}{
		{
			name:    "Case 1: Valid Exec Command",
			command: generateDummyExecCommand("command", component, v1alpha2.CommandGroup{Kind: runGroup}),
			wantErr: false,
		},
		{
			name:    "Case 2: Invalid Exec Command with missing component",
			command: generateDummyExecCommand("command", "", v1alpha2.CommandGroup{Kind: runGroup}),
			wantErr: true,
		},
		{
			name:    "Case 3: Valid Exec Command with invalid component",
			command: generateDummyExecCommand("command", invalidComponent, v1alpha2.CommandGroup{Kind: runGroup}),
			wantErr: true,
		},
		{
			name:    "Case 4: Valid Exec Command with Group nil",
			command: generateDummyExecCommand("command", component, v1alpha2.CommandGroup{}),
			wantErr: false,
		},
		{
			name:    "Case 5: Valid Apply Command",
			command: generateDummyApplyCommand("command", component, v1alpha2.CommandGroup{}),
			wantErr: false,
		},
		{
			name:    "Case 6: Invalid Apply Command with wrong component",
			command: generateDummyApplyCommand("command", invalidComponent, v1alpha2.CommandGroup{}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCommandComponent(tt.command, components)
			if !tt.wantErr == (err != nil) {
				t.Errorf("TestValidateAction unexpected error: %v", err)
				return
			}
		})
	}
}

func TestValidateCompositeCommand(t *testing.T) {

	component := "alias1"
	id := []string{"command1", "command2", "command3", "command4", "command5"}
	validExecCommands := []v1alpha2.Command{
		generateDummyExecCommand(id[0], component, v1alpha2.CommandGroup{Kind: runGroup}),
		generateDummyExecCommand(id[1], component, v1alpha2.CommandGroup{Kind: buildGroup}),
		generateDummyExecCommand(id[2], component, v1alpha2.CommandGroup{Kind: runGroup}),
	}
	components := []v1alpha2.Component{
		generateDummyContainerComponent(component, nil, nil, nil),
	}

	tests := []struct {
		name              string
		compositeCommands []v1alpha2.Command
		execCommands      []v1alpha2.Command
		wantErr           bool
	}{
		{
			name: "Case 1: Valid Composite Command",
			compositeCommands: []v1alpha2.Command{
				generateDummyCompositeCommand(id[3], []string{id[0], id[1], id[2]}, v1alpha2.CommandGroup{Kind: buildGroup}),
			},
			execCommands: validExecCommands,
			wantErr:      false,
		},
		{
			name: "Case 2: Invalid composite command, references non-existent command",
			compositeCommands: []v1alpha2.Command{
				generateDummyCompositeCommand(id[3], []string{id[0], "fakecommand", id[2]}, v1alpha2.CommandGroup{Kind: buildGroup}),
			},
			execCommands: validExecCommands,
			wantErr:      true,
		},
		{
			name: "Case 3: Invalid composite command, references itself",
			compositeCommands: []v1alpha2.Command{
				generateDummyCompositeCommand(id[3], []string{id[0], id[3], id[2]}, v1alpha2.CommandGroup{Kind: buildGroup}),
			},
			execCommands: validExecCommands,
			wantErr:      true,
		},
		{
			name: "Case 4: Invalid composite command, indirectly references itself",
			compositeCommands: []v1alpha2.Command{
				generateDummyCompositeCommand(id[3], []string{id[4], id[3], id[2]}, v1alpha2.CommandGroup{Kind: buildGroup}),
				generateDummyCompositeCommand(id[4], []string{id[0], id[3], id[2]}, v1alpha2.CommandGroup{Kind: buildGroup}),
			},
			execCommands: validExecCommands,
			wantErr:      true,
		},
		{
			name: "Case 5: Invalid composite command, points to invalid exec command",
			compositeCommands: []v1alpha2.Command{
				generateDummyCompositeCommand(id[3], []string{id[0], id[1]}, v1alpha2.CommandGroup{Kind: buildGroup}),
			},
			execCommands: []v1alpha2.Command{
				generateDummyExecCommand(id[0], component, v1alpha2.CommandGroup{Kind: runGroup}),
				generateDummyExecCommand(id[1], "some-fake-component", v1alpha2.CommandGroup{Kind: buildGroup}),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		commandMap := getCommandsMap(append(tt.execCommands, tt.compositeCommands...))

		t.Run(tt.name, func(t *testing.T) {
			cmd := tt.compositeCommands[0]
			parentCommands := make(map[string]string)

			err := validateCompositeCommand(&cmd, parentCommands, commandMap, components)
			if !tt.wantErr == (err != nil) {
				t.Errorf("TestValidateAction unexpected error: %v", err)
				return
			}
		})
	}
}

func TestValidateGroup(t *testing.T) {

	component := "alias1"

	tests := []struct {
		name     string
		commands []v1alpha2.Command
		wantErr  bool
	}{
		{
			name: "Case 1: Two default run commands",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("run command", component, v1alpha2.CommandGroup{Kind: runGroup, IsDefault: true}),
				generateDummyExecCommand("customcommand", component, v1alpha2.CommandGroup{Kind: runGroup, IsDefault: true}),
			},
			wantErr: true,
		},
		{
			name: "Case 2: No default for more than one build commands",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("build command", component, v1alpha2.CommandGroup{Kind: buildGroup}),
				generateDummyExecCommand("build command 2", component, v1alpha2.CommandGroup{Kind: buildGroup}),
			},
			wantErr: true,
		},
		{
			name: "Case 3: One command does not need default",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("test command", component, v1alpha2.CommandGroup{Kind: buildGroup}),
			},
			wantErr: false,
		},
		{
			name: "Case 4: One command can have default",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("test command", component, v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: true}),
			},
			wantErr: false,
		},
		{
			name: "Case 5: Composite commands in group",
			commands: []v1alpha2.Command{
				generateDummyExecCommand("build command", component, v1alpha2.CommandGroup{Kind: buildGroup}),
				generateDummyExecCommand("build command 2", component, v1alpha2.CommandGroup{Kind: buildGroup}),
				generateDummyCompositeCommand("composite1", []string{"build command", "build command 2"}, v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: true}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGroup(tt.commands)
			if !tt.wantErr && err != nil {
				t.Errorf("TestValidateGroup unexpected error: %v", err)
			}
		})
	}
}
