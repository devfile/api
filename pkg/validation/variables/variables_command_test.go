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
			outputFile:   "test-fixtures/commands/exec.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      true,
		},
		{
			name:         "Not an exec command",
			testFile:     "test-fixtures/components/volume.yaml",
			outputFile:   "test-fixtures/components/volume.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testExecCommand := v1alpha2.ExecCommand{}
			readFileToStruct(t, tt.testFile, &testExecCommand)

			testVariable := make(map[string]string)
			readFileToStruct(t, tt.variableFile, &testVariable)

			var err error
			if reflect.DeepEqual(testExecCommand, v1alpha2.ExecCommand{}) {
				err = validateAndReplaceForExecCommand(testVariable, nil)
			} else {
				err = validateAndReplaceForExecCommand(testVariable, &testExecCommand)
			}
			_, ok := err.(*InvalidKeysError)
			if tt.wantErr && !ok {
				t.Errorf("Expected InvalidKeysError error from test but got %+v", err)
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else {
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
			outputFile:   "test-fixtures/commands/composite.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      true,
		},
		{
			name:         "Not a composite command",
			testFile:     "test-fixtures/components/volume.yaml",
			outputFile:   "test-fixtures/components/volume.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCompositeCommand := v1alpha2.CompositeCommand{}
			readFileToStruct(t, tt.testFile, &testCompositeCommand)

			testVariable := make(map[string]string)
			readFileToStruct(t, tt.variableFile, &testVariable)

			var err error
			if reflect.DeepEqual(testCompositeCommand, v1alpha2.CompositeCommand{}) {
				err = validateAndReplaceForCompositeCommand(testVariable, nil)
			} else {
				err = validateAndReplaceForCompositeCommand(testVariable, &testCompositeCommand)
			}
			_, ok := err.(*InvalidKeysError)
			if tt.wantErr && !ok {
				t.Errorf("Expected InvalidKeysError error from test but got %+v", err)
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else {
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
			outputFile:   "test-fixtures/commands/apply.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      true,
		},
		{
			name:         "Not an apply command",
			testFile:     "test-fixtures/components/volume.yaml",
			outputFile:   "test-fixtures/components/volume.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testApplyCommand := v1alpha2.ApplyCommand{}
			readFileToStruct(t, tt.testFile, &testApplyCommand)

			testVariable := make(map[string]string)
			readFileToStruct(t, tt.variableFile, &testVariable)

			var err error
			if reflect.DeepEqual(testApplyCommand, v1alpha2.ApplyCommand{}) {
				err = validateAndReplaceForApplyCommand(testVariable, nil)
			} else {
				err = validateAndReplaceForApplyCommand(testVariable, &testApplyCommand)
			}
			_, ok := err.(*InvalidKeysError)
			if tt.wantErr && !ok {
				t.Errorf("Expected InvalidKeysError error from test but got %+v", err)
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else {
				expectedApplyCommand := v1alpha2.ApplyCommand{}
				readFileToStruct(t, tt.outputFile, &expectedApplyCommand)
				assert.Equal(t, expectedApplyCommand, testApplyCommand, "The two values should be the same.")
			}
		})
	}
}
