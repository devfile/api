package variables

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
)

func TestValidateAndReplaceEndpoint(t *testing.T) {

	tests := []struct {
		name         string
		testFile     string
		outputFile   string
		variableFile string
		wantErr      bool
	}{
		{
			name:         "Good Substitution",
			testFile:     "test-fixtures/components/endpoint.yaml",
			outputFile:   "test-fixtures/components/endpoint-output.yaml",
			variableFile: "test-fixtures/variables/variables-referenced.yaml",
			wantErr:      false,
		},
		{
			name:         "Invalid Reference",
			testFile:     "test-fixtures/components/endpoint.yaml",
			outputFile:   "test-fixtures/components/endpoint.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testEndpoint := v1alpha2.Endpoint{}
			readFileToStruct(t, tt.testFile, &testEndpoint)
			testEndpointArr := []v1alpha2.Endpoint{testEndpoint}

			testVariable := make(map[string]string)
			readFileToStruct(t, tt.variableFile, &testVariable)

			err := validateAndReplaceForEndpoint(testVariable, testEndpointArr)
			_, ok := err.(*InvalidKeysError)
			if tt.wantErr && !ok {
				t.Errorf("Expected InvalidKeysError error from test but got %+v", err)
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else {
				expectedEndpoint := v1alpha2.Endpoint{}
				readFileToStruct(t, tt.outputFile, &expectedEndpoint)
				expectedEndpointArr := []v1alpha2.Endpoint{expectedEndpoint}
				assert.Equal(t, expectedEndpointArr, testEndpointArr, "The two values should be the same.")
			}
		})
	}
}
