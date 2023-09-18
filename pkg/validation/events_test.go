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
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/api/v2/pkg/attributes"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

func TestValidateEvents(t *testing.T) {

	containers := []string{"container1", "container2"}

	commands := []v1alpha2.Command{
		generateDummyApplyCommand("apply1", containers[0], nil, attributes.Attributes{}),
		generateDummyApplyCommand("apply2", containers[0], nil, attributes.Attributes{}),
		generateDummyExecCommand("exec1", containers[1], nil),
		generateDummyExecCommand("exec2", containers[1], nil),
		generateDummyCompositeCommand("compositeOnlyApply", []string{"apply1", "apply2"}, nil),
		generateDummyCompositeCommand("compositeOnlyExec", []string{"exec1", "exec2"}, nil),
		generateDummyCompositeCommand("compositeExecApply", []string{"exec1", "apply1"}, nil),
	}

	preStartPostStopErr := ".*does not map to a valid devfile command.*\n.*should either map to an apply command or a composite command with apply commands.*"
	postStartPreStopErr := ".*does not map to a valid devfile command.*\n.*should either map to an exec command or a composite command with exec commands.*"

	tests := []struct {
		name    string
		events  v1alpha2.Events
		wantErr []string
	}{
		{
			name: "Valid preStart events - Apply and Composite Apply Commands",
			events: v1alpha2.Events{
				DevWorkspaceEvents: v1alpha2.DevWorkspaceEvents{
					PreStart: []string{
						"apply1",
						"apply2",
						"compositeOnlyApply",
					},
				},
			},
		},
		{
			name: "Invalid exec commands in preStart",
			events: v1alpha2.Events{
				DevWorkspaceEvents: v1alpha2.DevWorkspaceEvents{
					PreStart: []string{
						"apply12",
						"exec2",
						"compositeExecApply",
					},
				},
			},
			wantErr: []string{preStartPostStopErr},
		},
		{
			name: "Invalid exec commands in postStop",
			events: v1alpha2.Events{
				DevWorkspaceEvents: v1alpha2.DevWorkspaceEvents{
					PostStop: []string{
						"apply12",
						"exec2",
						"compositeExecApply",
					},
				},
			},
			wantErr: []string{preStartPostStopErr},
		},
		{
			name: "Invalid apply commands in postStart",
			events: v1alpha2.Events{
				DevWorkspaceEvents: v1alpha2.DevWorkspaceEvents{
					PostStart: []string{
						"apply12",
						"exec2",
						"compositeExecApply",
					},
				},
			},
			wantErr: []string{postStartPreStopErr},
		},
		{
			name: "Invalid apply commands in preStop",
			events: v1alpha2.Events{
				DevWorkspaceEvents: v1alpha2.DevWorkspaceEvents{
					PreStop: []string{
						"apply12",
						"exec2",
						"compositeExecApply",
					},
				},
			},
			wantErr: []string{postStartPreStopErr},
		},
		{
			name: "Multiple errors: Invalid exec commands in postStop, Invalid apply commands in preStop",
			events: v1alpha2.Events{
				DevWorkspaceEvents: v1alpha2.DevWorkspaceEvents{
					PreStop: []string{
						"apply12",
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
			wantErr: []string{postStartPreStopErr, preStartPostStopErr},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEvents(tt.events, commands)

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

func TestIsEventValid(t *testing.T) {

	containers := []string{"container1", "container2"}

	commands := []v1alpha2.Command{
		generateDummyApplyCommand("apply1", containers[0], nil, attributes.Attributes{}),
		generateDummyApplyCommand("apply2", containers[0], nil, attributes.Attributes{}),
		generateDummyExecCommand("exec1", containers[1], nil),
		generateDummyExecCommand("exec2", containers[1], nil),
		generateDummyCompositeCommand("compositeOnlyApply", []string{"apply1", "apply2"}, nil),
		generateDummyCompositeCommand("compositeOnlyExec", []string{"exec1", "exec2"}, nil),
		generateDummyCompositeCommand("compositeExecApply", []string{"exec1", "apply1"}, nil),
	}

	missingCmdErr := "does not map to a valid devfile command"
	applyCmdErr := "should either map to an apply command or a composite command with apply commands"
	execCmdErr := "should either map to an exec command or a composite command with exec commands"

	tests := []struct {
		name       string
		eventType  string
		eventNames []string
		wantErr    *string
	}{
		{
			name:      "Valid preStart events - Apply and Composite Apply Commands",
			eventType: preStart,
			eventNames: []string{
				"apply1",
				"apply2",
				"compositeOnlyApply",
			},
		},
		{
			name:      "Invalid postStop events - Invalid Exec Command",
			eventType: postStop,
			eventNames: []string{
				"exec2",
				"apply2",
				"compositeOnlyApply",
			},
			wantErr: &applyCmdErr,
		},
		{
			name:      "Invalid postStop events - Invalid Composite Command with Exec & Apply Subcommands",
			eventType: postStop,
			eventNames: []string{
				"apply1",
				"apply2",
				"compositeExecApply",
			},
			wantErr: &applyCmdErr,
		},
		{
			name:      "Valid postStart events - Exec and Composite Exec Commands",
			eventType: postStart,
			eventNames: []string{
				"exec1",
				"exec2",
				"compositeOnlyExec",
			},
		},
		{
			name:      "Invalid postStart events - Invalid Composite Commands with Exec & Apply Subcommands",
			eventType: postStart,
			eventNames: []string{
				"exec1",
				"exec2",
				"compositeExecApply",
			},
			wantErr: &execCmdErr,
		},
		{
			name:      "Invalid preStop events - Invalid Apply Command",
			eventType: preStop,
			eventNames: []string{
				"exec1",
				"apply2",
				"compositeOnlyExec",
			},
			wantErr: &execCmdErr,
		},
		{
			name:      "Invalid events - Missing event",
			eventType: preStop,
			eventNames: []string{
				"exec1",
				"apply2isInvalid",
				"compositeOnlyExec",
			},
			wantErr: &missingCmdErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commandMap := getCommandsMap(commands)
			err := isEventValid(tt.eventNames, tt.eventType, commandMap)

			if tt.wantErr != nil && assert.Error(t, err) {
				assert.Regexp(t, *tt.wantErr, err.Error(), "Error message should match")
			} else {
				assert.NoError(t, err, "Expected error to be nil")
			}
		})
	}
}
