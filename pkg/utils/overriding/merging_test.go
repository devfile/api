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

package overriding

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	dw "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
)

func mergingPatchTest(main, parent []byte, expected dw.DevWorkspaceTemplateSpecContent, expectedError string, plugins ...[]byte) func(t *testing.T) {
	return func(t *testing.T) {
		actual, err := MergeDevWorkspaceTemplateSpecBytes(main, parent, plugins...)
		if err != nil {
			compareErrorMessages(t, expectedError, err.Error(), "wrong error")
			return
		}
		if expectedError != "" {
			t.Error("Expected error but did not get one")
			return
		}

		assert.Equal(t, &expected, actual, "The two values should be the same")
	}
}

func TestMerging(t *testing.T) {
	filepath.Walk("test-fixtures/merges", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() == "main.yaml" {
			if err != nil {
				t.Error(err)
				return nil
			}
			main, err := ioutil.ReadFile(path)
			if err != nil {
				t.Error(err)
				return nil
			}
			dirPath := filepath.Dir(path)
			parent := []byte{}
			parentFile := filepath.Join(dirPath, "parent.yaml")
			if _, err = os.Stat(parentFile); err == nil {
				parent, err = ioutil.ReadFile(parentFile)
				if err != nil {
					t.Error(err)
					return nil
				}
			}

			plugins := [][]byte{}
			pluginFile := filepath.Join(dirPath, "plugin.yaml")
			if _, err = os.Stat(pluginFile); err == nil {
				plugin, err := ioutil.ReadFile(filepath.Join(dirPath, "plugin.yaml"))
				if err != nil {
					t.Error(err)
					return nil
				}
				plugins = append(plugins, plugin)
			}
			var resultTemplate dw.DevWorkspaceTemplateSpecContent
			resultError := ""
			errorFile := filepath.Join(dirPath, "result-error.txt")
			if _, err = os.Stat(errorFile); err == nil {
				resultErrorBytes, err := ioutil.ReadFile(errorFile)
				if err != nil {
					t.Error(err)
					return nil
				}
				resultError = string(resultErrorBytes)
			} else {
				readFileToStruct(t, filepath.Join(dirPath, "result.yaml"), &resultTemplate)
			}
			testName := filepath.Base(dirPath)

			t.Run(testName, mergingPatchTest(main, parent, resultTemplate, resultError, plugins...))
		}
		return nil
	})
}

func TestMergingOnlyPlugins(t *testing.T) {
	baseFile := "test-fixtures/merges/no-parent/main.yaml"
	pluginFile := "test-fixtures/merges/no-parent/plugin.yaml"
	resultFile := "test-fixtures/merges/no-parent/result.yaml"

	baseDWT := dw.DevWorkspaceTemplateSpecContent{}
	pluginDWT := dw.DevWorkspaceTemplateSpecContent{}
	expectedDWT := dw.DevWorkspaceTemplateSpecContent{}

	readFileToStruct(t, baseFile, &baseDWT)
	readFileToStruct(t, pluginFile, &pluginDWT)
	readFileToStruct(t, resultFile, &expectedDWT)

	gotDWT, err := MergeDevWorkspaceTemplateSpec(&baseDWT, nil, &pluginDWT)
	if assert.NoError(t, err) {
		assert.Equal(t, &expectedDWT, gotDWT)
	}
}

func TestMergingOnlyParent(t *testing.T) {
	// Reuse only plugin case since it's compatible
	baseFile := "test-fixtures/merges/no-parent/main.yaml"
	parentFile := "test-fixtures/merges/no-parent/plugin.yaml"
	resultFile := "test-fixtures/merges/no-parent/result.yaml"

	baseDWT := dw.DevWorkspaceTemplateSpecContent{}
	parentDWT := dw.DevWorkspaceTemplateSpecContent{}
	expectedDWT := dw.DevWorkspaceTemplateSpecContent{}

	readFileToStruct(t, baseFile, &baseDWT)
	readFileToStruct(t, parentFile, &parentDWT)
	readFileToStruct(t, resultFile, &expectedDWT)

	gotDWT, err := MergeDevWorkspaceTemplateSpec(&baseDWT, &parentDWT)
	if assert.NoError(t, err) {
		assert.Equal(t, &expectedDWT, gotDWT)
	}
}
