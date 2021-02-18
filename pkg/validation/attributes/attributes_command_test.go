package attributes

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
	"github.com/stretchr/testify/assert"
)

func TestValidateExecCommand(t *testing.T) {

	tests := []struct {
		name       string
		testFile   string
		expected   v1alpha2.ExecCommand
		attributes apiAttributes.Attributes
		wantErr    bool
	}{
		{
			name:     "Good Substitution",
			testFile: "test-fixtures/commands/exec.yaml",
			expected: v1alpha2.ExecCommand{
				CommandLine: "tail -f /dev/null",
				WorkingDir:  "FOO",
				Component:   "BAR",
				LabeledCommand: v1alpha2.LabeledCommand{
					Label: "1",
				},
				Env: []v1alpha2.EnvVar{
					{
						Name:  "FOO",
						Value: "BAR",
					},
				},
			},
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"version": "1",
				"devnull": "/dev/null",
				"bar":     "BAR",
				"foo":     "FOO",
			}, nil),
			wantErr: false,
		},
		{
			name:     "Invalid Reference",
			testFile: "test-fixtures/commands/exec.yaml",
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"foo": "FOO",
			}, nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testExecCommand := v1alpha2.ExecCommand{}

			readFileToStruct(t, tt.testFile, &testExecCommand)

			err := validateExecCommand(tt.attributes, &testExecCommand)
			if tt.wantErr == (err == nil) {
				t.Errorf("error: %v", err)
				return
			} else if err == nil {
				assert.Equal(t, tt.expected, testExecCommand, "The two values should be the same.")
			}
		})
	}
}

func TestValidateCompositeCommand(t *testing.T) {

	tests := []struct {
		name       string
		testFile   string
		expected   v1alpha2.CompositeCommand
		attributes apiAttributes.Attributes
		wantErr    bool
	}{
		{
			name:     "Good Substitution",
			testFile: "test-fixtures/commands/composite.yaml",
			expected: v1alpha2.CompositeCommand{
				LabeledCommand: v1alpha2.LabeledCommand{
					Label: "1",
				},
				Commands: []string{
					"FOO",
					"BAR",
				},
			},
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"version": "1",
				"foo":     "FOO",
			}, nil),
			wantErr: false,
		},
		{
			name:     "Invalid Reference",
			testFile: "test-fixtures/commands/composite.yaml",
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"foo": "FOO",
			}, nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCompositeCommand := v1alpha2.CompositeCommand{}

			readFileToStruct(t, tt.testFile, &testCompositeCommand)

			err := validateCompositeCommand(tt.attributes, &testCompositeCommand)
			if tt.wantErr == (err == nil) {
				t.Errorf("error: %v", err)
				return
			} else if err == nil {
				assert.Equal(t, tt.expected, testCompositeCommand, "The two values should be the same.")
			}
		})
	}
}

func TestValidateApplyCommand(t *testing.T) {

	tests := []struct {
		name       string
		testFile   string
		expected   v1alpha2.ApplyCommand
		attributes apiAttributes.Attributes
		wantErr    bool
	}{
		{
			name:     "Good Substitution",
			testFile: "test-fixtures/commands/apply.yaml",
			expected: v1alpha2.ApplyCommand{
				LabeledCommand: v1alpha2.LabeledCommand{
					Label: "1",
				},
				Component: "FOO",
			},
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"version": "1",
				"foo":     "FOO",
			}, nil),
			wantErr: false,
		},
		{
			name:     "Invalid Reference",
			testFile: "test-fixtures/commands/apply.yaml",
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"foo": "FOO",
			}, nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testApplyCommand := v1alpha2.ApplyCommand{}

			readFileToStruct(t, tt.testFile, &testApplyCommand)

			err := validateApplyCommand(tt.attributes, &testApplyCommand)
			if tt.wantErr == (err == nil) {
				t.Errorf("error: %v", err)
				return
			} else if err == nil {
				assert.Equal(t, tt.expected, testApplyCommand, "The two values should be the same.")
			}
		})
	}
}
