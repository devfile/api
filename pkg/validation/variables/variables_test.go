package variables

import (
	"io/ioutil"
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestValidateGlobalVariableBasic(t *testing.T) {

	tests := []struct {
		name       string
		testFile   string
		outputFile string
		wantErr    bool
	}{
		{
			name:       "Successful global variable substitution",
			testFile:   "test-fixtures/all/devfile-good.yaml",
			outputFile: "test-fixtures/all/devfile-good-output.yaml",
			wantErr:    false,
		},
		{
			name:     "Invalid Reference",
			testFile: "test-fixtures/all/devfile-bad.yaml",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDWT := v1alpha2.DevWorkspaceTemplateSpec{}
			readFileToStruct(t, tt.testFile, &testDWT)

			err := ValidateAndReplaceGlobalVariable(&testDWT)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error from test but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else if err == nil {
				expectedDWT := v1alpha2.DevWorkspaceTemplateSpec{}
				readFileToStruct(t, tt.outputFile, &expectedDWT)
				assert.Equal(t, expectedDWT, testDWT, "The two values should be the same.")
			}
		})
	}
}

func TestValidateAndReplaceDataWithVariable(t *testing.T) {

	invalidVariableErr := ".*Variable with key .* does not exist.*"

	tests := []struct {
		name       string
		testString string
		variables  map[string]string
		wantValue  string
		wantErr    *string
	}{
		{
			name:       "Valid variable reference",
			testString: "image-{{version}}:{{tag}}-14",
			variables: map[string]string{
				"version": "1.x.x",
				"tag":     "dev",
			},
			wantValue: "image-1.x.x:dev-14",
			wantErr:   nil,
		},
		{
			name:       "Invalid variable reference",
			testString: "image-{{version}}:{{invalid}}-14",
			variables: map[string]string{
				"version": "1.x.x",
				"tag":     "dev",
			},
			wantErr: &invalidVariableErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, err := validateAndReplaceDataWithVariable(tt.testString, tt.variables)
			if tt.wantErr != nil && assert.Error(t, err) {
				assert.Regexp(t, *tt.wantErr, err.Error(), "Error message should match")
			} else {
				assert.NoError(t, err, "Expected error to be nil")
				if gotValue != tt.wantValue {
					assert.Equal(t, tt.wantValue, gotValue, "Return value should match")
				}
			}
		})
	}
}

func readFileToStruct(t *testing.T, path string, into interface{}) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test file from %s: %s", path, err.Error())
	}
	err = yaml.Unmarshal(bytes, into)
	if err != nil {
		t.Fatalf("Failed to unmarshal file into struct: %s", err.Error())
	}
}
