package validation

import (
	"strings"
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

func TestIsEventValid(t *testing.T) {

	containers := []string{"container1", "container2"}
	execCommands := []v1alpha2.Command{
		generateDummyExecCommand("command1", containers[0], v1alpha2.CommandGroup{}),
		generateDummyExecCommand("command2", containers[1], v1alpha2.CommandGroup{}),
	}
	compCommands := []v1alpha2.Command{
		generateDummyCompositeCommand("composite1", []string{"command1", "command2"}, v1alpha2.CommandGroup{}),
	}
	commandMap := getCommandsMap(append(execCommands, compCommands...))

	tests := []struct {
		name       string
		eventType  string
		eventNames []string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:      "Case 1: Valid events",
			eventType: "preStart",
			eventNames: []string{
				"command1",
				"composite1",
			},
			wantErr: false,
		},
		{
			name:      "Case 2: Invalid events with wrong mapping to devfile command",
			eventType: "preStart",
			eventNames: []string{
				"command12345iswrong",
				"composite1",
			},
			wantErr:    true,
			wantErrMsg: "does not map to a valid devfile command",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := isEventValid(tt.eventNames, tt.eventType, commandMap)
			if err != nil && !tt.wantErr {
				t.Errorf("TestIsEventValid error: %v", err)
			} else if err != nil && tt.wantErr {
				if !strings.Contains(err.Error(), tt.wantErrMsg) {
					t.Errorf("TestIsEventValid error mismatch - %s; does not contain: %s", err.Error(), tt.wantErrMsg)
				}
			}
		})
	}

}
