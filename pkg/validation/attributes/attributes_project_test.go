package attributes

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
	"github.com/stretchr/testify/assert"
)

func TestValidateAndReplaceProjects(t *testing.T) {

	tests := []struct {
		name          string
		testFile      string
		outputFile    string
		attributeFile string
		wantErr       bool
	}{
		{
			name:          "Good Substitution",
			testFile:      "test-fixtures/projects/project.yaml",
			outputFile:    "test-fixtures/projects/project-output.yaml",
			attributeFile: "test-fixtures/attributes/attributes-referenced.yaml",
			wantErr:       false,
		},
		{
			name:          "Invalid Reference",
			testFile:      "test-fixtures/projects/project.yaml",
			attributeFile: "test-fixtures/attributes/attributes-notreferenced.yaml",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testProject := v1alpha2.Project{}
			readFileToStruct(t, tt.testFile, &testProject)
			testProjectArr := []v1alpha2.Project{testProject}

			testAttribute := apiAttributes.Attributes{}
			readFileToStruct(t, tt.attributeFile, &testAttribute)

			err := ValidateAndReplaceForProjects(testAttribute, testProjectArr)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error from test but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else if err == nil {
				expectedProject := v1alpha2.Project{}
				readFileToStruct(t, tt.outputFile, &expectedProject)
				expectedProjectArr := []v1alpha2.Project{expectedProject}
				assert.Equal(t, expectedProjectArr, testProjectArr, "The two values should be the same.")
			}
		})
	}
}

func TestValidateAndReplaceStarterProjects(t *testing.T) {

	tests := []struct {
		name          string
		testFile      string
		outputFile    string
		attributeFile string
		wantErr       bool
	}{
		{
			name:          "Good Substitution",
			testFile:      "test-fixtures/projects/starterproject.yaml",
			outputFile:    "test-fixtures/projects/starterproject-output.yaml",
			attributeFile: "test-fixtures/attributes/attributes-referenced.yaml",
			wantErr:       false,
		},
		{
			name:          "Invalid Reference",
			testFile:      "test-fixtures/projects/starterproject.yaml",
			attributeFile: "test-fixtures/attributes/attributes-notreferenced.yaml",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStarterProject := v1alpha2.StarterProject{}
			readFileToStruct(t, tt.testFile, &testStarterProject)
			testStarterProjectArr := []v1alpha2.StarterProject{testStarterProject}

			testAttribute := apiAttributes.Attributes{}
			readFileToStruct(t, tt.attributeFile, &testAttribute)

			err := ValidateAndReplaceForStarterProjects(testAttribute, testStarterProjectArr)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error from test but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else if err == nil {
				expectedStarterProject := v1alpha2.StarterProject{}
				readFileToStruct(t, tt.outputFile, &expectedStarterProject)
				expectedStarterProjectArr := []v1alpha2.StarterProject{expectedStarterProject}
				assert.Equal(t, expectedStarterProjectArr, testStarterProjectArr, "The two values should be the same.")
			}
		})
	}
}

func TestValidateAndReplaceProjectSrc(t *testing.T) {

	tests := []struct {
		name          string
		testFile      string
		outputFile    string
		attributeFile string
		wantErr       bool
	}{
		{
			name:          "Good Git Substitution",
			testFile:      "test-fixtures/projects/git.yaml",
			outputFile:    "test-fixtures/projects/git-output.yaml",
			attributeFile: "test-fixtures/attributes/attributes-referenced.yaml",
			wantErr:       false,
		},
		{
			name:          "Good Zip Substitution",
			testFile:      "test-fixtures/projects/zip.yaml",
			outputFile:    "test-fixtures/projects/zip-output.yaml",
			attributeFile: "test-fixtures/attributes/attributes-referenced.yaml",
			wantErr:       false,
		},
		{
			name:          "Invalid Git Reference",
			testFile:      "test-fixtures/projects/git.yaml",
			attributeFile: "test-fixtures/attributes/attributes-notreferenced.yaml",
			wantErr:       true,
		},
		{
			name:          "Invalid Zip Reference",
			testFile:      "test-fixtures/projects/zip.yaml",
			attributeFile: "test-fixtures/attributes/attributes-notreferenced.yaml",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testProjectSrc := v1alpha2.ProjectSource{}
			readFileToStruct(t, tt.testFile, &testProjectSrc)

			testAttribute := apiAttributes.Attributes{}
			readFileToStruct(t, tt.attributeFile, &testAttribute)

			err := validateandReplaceForProjectSource(testAttribute, &testProjectSrc)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error from test but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else if err == nil {
				expectedProjectSrc := v1alpha2.ProjectSource{}
				readFileToStruct(t, tt.outputFile, &expectedProjectSrc)
				assert.Equal(t, expectedProjectSrc, testProjectSrc, "The two values should be the same.")
			}
		})
	}
}
