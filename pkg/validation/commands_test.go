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

func TestValidateCommands(t *testing.T) {

	component := "alias1"

	components := []v1alpha2.Component{
		generateDummyContainerComponent(component, nil, nil, nil),
	}

	tests := []struct {
		name      string
		exec      []v1alpha2.Command
		composite []v1alpha2.Command
		wantErr   bool
	}{
		{
			name: "Case 1: Valid Exec Command",
			exec: []v1alpha2.Command{
				generateDummyExecCommand("command", component, v1alpha2.CommandGroup{Kind: runGroup}),
			},
			wantErr: false,
		},
		{
			name: "Case 2: Valid Composite Command",
			exec: []v1alpha2.Command{
				generateDummyExecCommand("somecommand1", component, v1alpha2.CommandGroup{}),
				generateDummyExecCommand("somecommand2", component, v1alpha2.CommandGroup{}),
			},
			composite: []v1alpha2.Command{
				generateDummyCompositeCommand("composite1", []string{"somecommand1", "somecommand2"}, v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: true}),
			},
			wantErr: false,
		},
		{
			name: "Case 3: Duplicate commands",
			exec: []v1alpha2.Command{
				generateDummyExecCommand("somecommand1", component, v1alpha2.CommandGroup{}),
				generateDummyExecCommand("somecommand1", component, v1alpha2.CommandGroup{}),
			},
			wantErr: true,
		},
		{
			name: "Case 4: Duplicate commands, different types",
			exec: []v1alpha2.Command{
				generateDummyExecCommand("somecommand1", component, v1alpha2.CommandGroup{}),
			},
			composite: []v1alpha2.Command{
				generateDummyCompositeCommand("somecommand1", []string{"fakecommand"}, v1alpha2.CommandGroup{Kind: buildGroup, IsDefault: true}),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCommands(append(tt.composite, tt.exec...), components)
			if !tt.wantErr == (err != nil) {
				t.Errorf("TestValidateAction unexpected error: %v", err)
				return
			}
		})
	}

}

func TestValidateExecCommand(t *testing.T) {

	component := "alias1"
	invalidComponent := "garbagealias"

	components := []v1alpha2.Component{
		generateDummyContainerComponent(component, nil, nil, nil),
	}

	tests := []struct {
		name    string
		exec    v1alpha2.Command
		wantErr bool
	}{
		{
			name:    "Case 1: Valid Exec Command",
			exec:    generateDummyExecCommand("command", component, v1alpha2.CommandGroup{Kind: runGroup}),
			wantErr: false,
		},
		{
			name:    "Case 2: Invalid Exec Command with missing component",
			exec:    generateDummyExecCommand("command", "", v1alpha2.CommandGroup{Kind: runGroup}),
			wantErr: true,
		},
		{
			name:    "Case 3: Valid Exec Command with invalid component",
			exec:    generateDummyExecCommand("command", invalidComponent, v1alpha2.CommandGroup{Kind: runGroup}),
			wantErr: true,
		},
		{
			name:    "Case 4: valid Exec Command with Group nil",
			exec:    generateDummyExecCommand("command", component, v1alpha2.CommandGroup{}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateExecCommand(tt.exec, components)
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
