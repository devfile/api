package validation

import (
	"strings"
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

func TestValidateEvents(t *testing.T) {

	containers := []string{"container1", "container2"}

	commands := []v1alpha2.Command{
		generateDummyApplyCommand("apply1", containers[0], v1alpha2.CommandGroup{}),
		generateDummyApplyCommand("apply2", containers[0], v1alpha2.CommandGroup{}),
		generateDummyExecCommand("exec1", containers[1], v1alpha2.CommandGroup{}),
		generateDummyExecCommand("exec2", containers[1], v1alpha2.CommandGroup{}),
		generateDummyCompositeCommand("compositeOnlyApply", []string{"apply1", "apply2"}, v1alpha2.CommandGroup{}),
		generateDummyCompositeCommand("compositeOnlyExec", []string{"exec1", "exec2"}, v1alpha2.CommandGroup{}),
		generateDummyCompositeCommand("compositeExecApply", []string{"exec1", "apply1"}, v1alpha2.CommandGroup{}),
	}

	tests := []struct {
		name       string
		events     v1alpha2.Events
		wantErr    bool
		wantErrMsg string
	}{
		{
			name: "Case 1: Valid preStart events - Apply and Composite Apply Commands",
			events: v1alpha2.Events{
				WorkspaceEvents: v1alpha2.WorkspaceEvents{
					PreStart: []string{
						"apply1",
						"apply2",
						"compositeOnlyApply",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Case 2: Invalid postStart and postStop events",
			events: v1alpha2.Events{
				WorkspaceEvents: v1alpha2.WorkspaceEvents{
					PostStart: []string{
						"apply1",
						"exec2",
						"compositeExecApply",
					},
					PostStop: []string{
						"apply12",
						"exec2",
						"compositeExecApply",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEvents(tt.events, commands)
			if err != nil && !tt.wantErr {
				t.Errorf("TestValidateEvents error: %v", err)
			}
		})
	}
}

func TestIsEventValid(t *testing.T) {

	containers := []string{"container1", "container2"}

	commands := []v1alpha2.Command{
		generateDummyApplyCommand("apply1", containers[0], v1alpha2.CommandGroup{}),
		generateDummyApplyCommand("apply2", containers[0], v1alpha2.CommandGroup{}),
		generateDummyExecCommand("exec1", containers[1], v1alpha2.CommandGroup{}),
		generateDummyExecCommand("exec2", containers[1], v1alpha2.CommandGroup{}),
		generateDummyCompositeCommand("compositeOnlyApply", []string{"apply1", "apply2"}, v1alpha2.CommandGroup{}),
		generateDummyCompositeCommand("compositeOnlyExec", []string{"exec1", "exec2"}, v1alpha2.CommandGroup{}),
		generateDummyCompositeCommand("compositeExecApply", []string{"exec1", "apply1"}, v1alpha2.CommandGroup{}),
	}

	tests := []struct {
		name       string
		eventType  string
		eventNames []string
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:      "Case 1: Valid preStart events - Apply and Composite Apply Commands",
			eventType: preStart,
			eventNames: []string{
				"apply1",
				"apply2",
				"compositeOnlyApply",
			},
			wantErr: false,
		},
		{
			name:      "Case 2: Invalid postStop events - Non-Apply and Composite Apply Commands",
			eventType: postStop,
			eventNames: []string{
				"exec2",
				"apply2",
				"compositeOnlyApply",
			},
			wantErr: true,
		},
		{
			name:      "Case 3: Invalid postStop events - Apply and Composite Mixed Commands",
			eventType: postStop,
			eventNames: []string{
				"apply1",
				"apply2",
				"compositeExecApply",
			},
			wantErr: true,
		},
		{
			name:      "Case 4: Valid postStart events - Exec and Composite Exec Commands",
			eventType: postStart,
			eventNames: []string{
				"exec1",
				"exec2",
				"compositeOnlyExec",
			},
			wantErr: false,
		},
		{
			name:      "Case 5: Invalid postStart events - Exec and Composite Mixed Commands",
			eventType: postStart,
			eventNames: []string{
				"exec1",
				"exec2",
				"compositeExecApply",
			},
			wantErr: true,
		},
		{
			name:      "Case 6: Invalid preStop events - Non-Exec and Composite Exec Commands",
			eventType: preStop,
			eventNames: []string{
				"exec1",
				"apply2",
				"compositeOnlyExec",
			},
			wantErr: true,
		},
		{
			name:      "Case 7: Invalid events - Missing event",
			eventType: preStop,
			eventNames: []string{
				"exec1",
				"apply2isInvalid",
				"compositeOnlyExec",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commandMap := getCommandsMap(commands)
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
