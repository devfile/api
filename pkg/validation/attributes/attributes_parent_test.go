package attributes

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
	"github.com/stretchr/testify/assert"
)

func TestValidateAndReplaceParent(t *testing.T) {

	tests := []struct {
		name          string
		testFile      string
		outputFile    string
		attributeFile string
		wantErr       bool
	}{
		{
			name:          "Good Uri Substitution",
			testFile:      "test-fixtures/parent/parent-uri.yaml",
			outputFile:    "test-fixtures/parent/parent-uri-output.yaml",
			attributeFile: "test-fixtures/attributes/attributes-referenced.yaml",
			wantErr:       false,
		},
		{
			name:          "Good Id Substitution",
			testFile:      "test-fixtures/parent/parent-id.yaml",
			outputFile:    "test-fixtures/parent/parent-id-output.yaml",
			attributeFile: "test-fixtures/attributes/attributes-referenced.yaml",
			wantErr:       false,
		},
		{
			name:          "Good Kube Substitution",
			testFile:      "test-fixtures/parent/parent-kubernetes.yaml",
			outputFile:    "test-fixtures/parent/parent-kubernetes-output.yaml",
			attributeFile: "test-fixtures/attributes/attributes-referenced.yaml",
			wantErr:       false,
		},
		{
			name:          "Invalid Reference",
			testFile:      "test-fixtures/parent/parent-id.yaml",
			attributeFile: "test-fixtures/attributes/attributes-notreferenced.yaml",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParent := v1alpha2.Parent{}
			readFileToStruct(t, tt.testFile, &testParent)

			testAttribute := apiAttributes.Attributes{}
			readFileToStruct(t, tt.attributeFile, &testAttribute)

			err := ValidateAndReplaceForParent(testAttribute, &testParent)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error from test but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else if err == nil {
				expectedParent := v1alpha2.Parent{}
				readFileToStruct(t, tt.outputFile, &expectedParent)
				assert.Equal(t, expectedParent, testParent, "The two values should be the same.")
			}
		})
	}
}
