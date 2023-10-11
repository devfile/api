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

package variables

import (
	"reflect"
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
)

func TestValidateAndReplaceProjects(t *testing.T) {

	tests := []struct {
		name         string
		testFile     string
		outputFile   string
		variableFile string
		wantErr      bool
	}{
		{
			name:         "Good Substitution",
			testFile:     "test-fixtures/projects/project.yaml",
			outputFile:   "test-fixtures/projects/project-output.yaml",
			variableFile: "test-fixtures/variables/variables-referenced.yaml",
			wantErr:      false,
		},
		{
			name:         "Invalid Reference",
			testFile:     "test-fixtures/projects/project.yaml",
			outputFile:   "test-fixtures/projects/project.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testProject := v1alpha2.Project{}
			readFileToStruct(t, tt.testFile, &testProject)
			testProjectArr := []v1alpha2.Project{testProject}

			testVariable := make(map[string]string)
			readFileToStruct(t, tt.variableFile, &testVariable)

			ValidateAndReplaceForProjects(testVariable, testProjectArr)
			expectedProject := v1alpha2.Project{}
			readFileToStruct(t, tt.outputFile, &expectedProject)
			expectedProjectArr := []v1alpha2.Project{expectedProject}
			assert.Equal(t, expectedProjectArr, testProjectArr, "The two values should be the same.")
		})
	}
}

func TestValidateAndReplaceStarterProjects(t *testing.T) {

	tests := []struct {
		name         string
		testFile     string
		outputFile   string
		variableFile string
	}{
		{
			name:         "Good Substitution",
			testFile:     "test-fixtures/projects/starterproject.yaml",
			outputFile:   "test-fixtures/projects/starterproject-output.yaml",
			variableFile: "test-fixtures/variables/variables-referenced.yaml",
		},
		{
			name:         "Invalid Reference",
			testFile:     "test-fixtures/projects/starterproject.yaml",
			outputFile:   "test-fixtures/projects/starterproject.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStarterProject := v1alpha2.StarterProject{}
			readFileToStruct(t, tt.testFile, &testStarterProject)
			testStarterProjectArr := []v1alpha2.StarterProject{testStarterProject}

			testVariable := make(map[string]string)
			readFileToStruct(t, tt.variableFile, &testVariable)

			ValidateAndReplaceForStarterProjects(testVariable, testStarterProjectArr)
			expectedStarterProject := v1alpha2.StarterProject{}
			readFileToStruct(t, tt.outputFile, &expectedStarterProject)
			expectedStarterProjectArr := []v1alpha2.StarterProject{expectedStarterProject}
			assert.Equal(t, expectedStarterProjectArr, testStarterProjectArr, "The two values should be the same.")
		})
	}
}

func TestValidateAndReplaceProjectSrc(t *testing.T) {

	tests := []struct {
		name         string
		testFile     string
		outputFile   string
		variableFile string
		wantErr      bool
	}{
		{
			name:         "Good Git Substitution",
			testFile:     "test-fixtures/projects/git.yaml",
			outputFile:   "test-fixtures/projects/git-output.yaml",
			variableFile: "test-fixtures/variables/variables-referenced.yaml",
			wantErr:      false,
		},
		{
			name:         "Good Zip Substitution",
			testFile:     "test-fixtures/projects/zip.yaml",
			outputFile:   "test-fixtures/projects/zip-output.yaml",
			variableFile: "test-fixtures/variables/variables-referenced.yaml",
			wantErr:      false,
		},
		{
			name:         "Invalid Git Reference",
			testFile:     "test-fixtures/projects/git.yaml",
			outputFile:   "test-fixtures/projects/git.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      true,
		},
		{
			name:         "Invalid Zip Reference",
			testFile:     "test-fixtures/projects/zip.yaml",
			outputFile:   "test-fixtures/projects/zip.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      true,
		},
		{
			name:         "Not a project source",
			testFile:     "test-fixtures/components/volume.yaml",
			outputFile:   "test-fixtures/components/volume.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testProjectSrc := v1alpha2.ProjectSource{}
			readFileToStruct(t, tt.testFile, &testProjectSrc)

			testVariable := make(map[string]string)
			readFileToStruct(t, tt.variableFile, &testVariable)

			var err error
			if reflect.DeepEqual(testProjectSrc, v1alpha2.ProjectSource{}) {
				err = validateandReplaceForProjectSource(testVariable, nil)
			} else {
				err = validateandReplaceForProjectSource(testVariable, &testProjectSrc)
			}
			_, ok := err.(*InvalidKeysError)
			if tt.wantErr && !ok {
				t.Errorf("Expected InvalidKeysError error from test but got %+v", err)
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else {
				expectedProjectSrc := v1alpha2.ProjectSource{}
				readFileToStruct(t, tt.outputFile, &expectedProjectSrc)
				assert.Equal(t, expectedProjectSrc, testProjectSrc, "The two values should be the same.")
			}
		})
	}
}

func TestValidateAndReplaceGitProjectSrc(t *testing.T) {

	tests := []struct {
		name         string
		testFile     string
		outputFile   string
		variableFile string
		wantErr      bool
	}{
		{
			name:         "Good Git Substitution",
			testFile:     "test-fixtures/projects/git.yaml",
			outputFile:   "test-fixtures/projects/git-output.yaml",
			variableFile: "test-fixtures/variables/variables-referenced.yaml",
			wantErr:      false,
		},
		{
			name:         "Not a git roject source",
			testFile:     "test-fixtures/projects/zip.yaml",
			outputFile:   "test-fixtures/projects/zip.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testProjectGitSrc := v1alpha2.GitLikeProjectSource{}
			readFileToStruct(t, tt.testFile, &testProjectGitSrc)

			testVariable := make(map[string]string)
			readFileToStruct(t, tt.variableFile, &testVariable)

			var err error
			if reflect.DeepEqual(testProjectGitSrc, v1alpha2.GitLikeProjectSource{}) {
				err = validateAndReplaceForGitProjectSource(testVariable, nil)
			} else {
				err = validateAndReplaceForGitProjectSource(testVariable, &testProjectGitSrc)
			}
			_, ok := err.(*InvalidKeysError)
			if tt.wantErr && !ok {
				t.Errorf("Expected InvalidKeysError error from test but got %+v", err)
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else {
				expectedProjectSrc := v1alpha2.GitLikeProjectSource{}
				readFileToStruct(t, tt.outputFile, &expectedProjectSrc)
				assert.Equal(t, expectedProjectSrc, testProjectGitSrc, "The two values should be the same.")
			}
		})
	}
}
