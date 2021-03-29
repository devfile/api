package variables

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
)

func TestValidateAndReplaceExecCommand(t *testing.T) {

	tests := []struct {
		name         string
		testFile     string
		outputFile   string
		variableFile string
		wantErr      bool
	}{
		{
			name:         "Good Substitution",
			testFile:     "test-fixtures/commands/exec.yaml",
			outputFile:   "test-fixtures/commands/exec-output.yaml",
			variableFile: "test-fixtures/variables/variables-referenced.yaml",
			wantErr:      false,
		},
		{
			name:         "Invalid Reference",
			testFile:     "test-fixtures/commands/exec.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testExecCommand := v1alpha2.ExecCommand{}
			readFileToStruct(t, tt.testFile, &testExecCommand)

			testVariable := make(map[string]string)
			readFileToStruct(t, tt.variableFile, &testVariable)

			err := validateAndReplaceForExecCommand(testVariable, &testExecCommand)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error from test but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else if err == nil {
				expectedExecCommand := v1alpha2.ExecCommand{}
				readFileToStruct(t, tt.outputFile, &expectedExecCommand)
				assert.Equal(t, expectedExecCommand, testExecCommand, "The two values should be the same.")
			}
		})
	}
}

func TestValidateAndReplaceCompositeCommand(t *testing.T) {

	tests := []struct {
		name         string
		testFile     string
		outputFile   string
		variableFile string
		wantErr      bool
	}{
		{
			name:         "Good Substitution",
			testFile:     "test-fixtures/commands/composite.yaml",
			outputFile:   "test-fixtures/commands/composite-output.yaml",
			variableFile: "test-fixtures/variables/variables-referenced.yaml",
			wantErr:      false,
		},
		{
			name:         "Invalid Reference",
			testFile:     "test-fixtures/commands/composite.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCompositeCommand := v1alpha2.CompositeCommand{}
			readFileToStruct(t, tt.testFile, &testCompositeCommand)

			testVariable := make(map[string]string)
			readFileToStruct(t, tt.variableFile, &testVariable)

			err := validateAndReplaceForCompositeCommand(testVariable, &testCompositeCommand)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error from test but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else if err == nil {
				expectedCompositeCommand := v1alpha2.CompositeCommand{}
				readFileToStruct(t, tt.outputFile, &expectedCompositeCommand)
				assert.Equal(t, expectedCompositeCommand, testCompositeCommand, "The two values should be the same.")
			}
		})
	}
}

func TestValidateAndReplaceApplyCommand(t *testing.T) {

	tests := []struct {
		name         string
		testFile     string
		outputFile   string
		variableFile string
		wantErr      bool
	}{
		{
			name:         "Good Substitution",
			testFile:     "test-fixtures/commands/apply.yaml",
			outputFile:   "test-fixtures/commands/apply-output.yaml",
			variableFile: "test-fixtures/variables/variables-referenced.yaml",
			wantErr:      false,
		},
		{
			name:         "Invalid Reference",
			testFile:     "test-fixtures/commands/apply.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testApplyCommand := v1alpha2.ApplyCommand{}
			readFileToStruct(t, tt.testFile, &testApplyCommand)

			testVariable := make(map[string]string)
			readFileToStruct(t, tt.variableFile, &testVariable)

			err := validateAndReplaceForApplyCommand(testVariable, &testApplyCommand)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error from test but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else if err == nil {
				expectedApplyCommand := v1alpha2.ApplyCommand{}
				readFileToStruct(t, tt.outputFile, &expectedApplyCommand)
				assert.Equal(t, expectedApplyCommand, testApplyCommand, "The two values should be the same.")
			}
		})
	}
}
