package attributes

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
	"github.com/stretchr/testify/assert"
)

func TestValidateAndReplaceEndpoint(t *testing.T) {

	tests := []struct {
		name          string
		testFile      string
		outputFile    string
		attributeFile string
		wantErr       bool
	}{
		{
			name:          "Good Substitution",
			testFile:      "test-fixtures/components/endpoint.yaml",
			outputFile:    "test-fixtures/components/endpoint-output.yaml",
			attributeFile: "test-fixtures/attributes/attributes-referenced.yaml",
			wantErr:       false,
		},
		{
			name:          "Invalid Reference",
			testFile:      "test-fixtures/components/endpoint.yaml",
			attributeFile: "test-fixtures/attributes/attributes-notreferenced.yaml",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testEndpoint := v1alpha2.Endpoint{}
			readFileToStruct(t, tt.testFile, &testEndpoint)
			testEndpointArr := []v1alpha2.Endpoint{testEndpoint}

			testAttribute := apiAttributes.Attributes{}
			readFileToStruct(t, tt.attributeFile, &testAttribute)

			err := validateAndReplaceForEndpoint(testAttribute, testEndpointArr)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error from test but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else if err == nil {
				expectedEndpoint := v1alpha2.Endpoint{}
				readFileToStruct(t, tt.outputFile, &expectedEndpoint)
				expectedEndpointArr := []v1alpha2.Endpoint{expectedEndpoint}
				assert.Equal(t, expectedEndpointArr, testEndpointArr, "The two values should be the same.")
			}
		})
	}
}
