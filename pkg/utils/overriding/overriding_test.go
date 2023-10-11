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
	"strings"
	"testing"

	dw "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	attributesPkg "github.com/devfile/api/v2/pkg/attributes"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestBasicToplevelOverriding(t *testing.T) {
	original := dw.DevWorkspaceTemplateSpecContent{
		Commands: []dw.Command{
			{
				Id: "command-with-type-changed",
				CommandUnion: dw.CommandUnion{
					Exec: &dw.ExecCommand{},
				},
			},
			{
				Id: "command-to-replace",
				CommandUnion: dw.CommandUnion{
					Exec: &dw.ExecCommand{
						Env: []dw.EnvVar{
							{
								Name:  "envVarToReplace",
								Value: "envVarToReplaceOriginalValue",
							},
							{
								Name:  "envVarNotChanged",
								Value: "envVarNotChangedOriginalValue",
							},
						},
					},
				},
			},
			{
				Id: "command-not-changed",
				CommandUnion: dw.CommandUnion{
					Exec: &dw.ExecCommand{
						LabeledCommand: dw.LabeledCommand{
							Label: "commandNotChangedLabel",
						},
					},
				},
			},
		},
		Variables: map[string]string{
			"version": "main",
			"xyz":     "xyz",
		},
		Attributes: attributesPkg.Attributes{}.FromMap(map[string]interface{}{
			"version": "main",
			"xyz":     "xyz",
		}, nil),
	}

	patch := dw.ParentOverrides{
		Commands: []dw.CommandParentOverride{
			{
				Id: "command-with-type-changed",
				CommandUnionParentOverride: dw.CommandUnionParentOverride{
					Apply: &dw.ApplyCommandParentOverride{
						Component: "mycomponent",
					},
				},
			},
			{
				Id: "command-to-replace",
				CommandUnionParentOverride: dw.CommandUnionParentOverride{
					Exec: &dw.ExecCommandParentOverride{
						Env: []dw.EnvVarParentOverride{
							{
								Name:  "envVarToReplace",
								Value: "envVarToReplaceNewValue",
							},
							{
								Name:  "endVarToAdd",
								Value: "endVarToAddValue",
							},
						},
					},
				},
			},
		},
		Variables: map[string]string{
			"version": "patch",
		},
		Attributes: attributesPkg.Attributes{}.FromMap(map[string]interface{}{
			"version": map[string]interface{}{
				"patch": true,
			},
		}, nil),
	}

	expected := &dw.DevWorkspaceTemplateSpecContent{
		Commands: []dw.Command{
			{
				Id: "command-with-type-changed",
				CommandUnion: dw.CommandUnion{
					Apply: &dw.ApplyCommand{
						Component: "mycomponent",
					},
				},
			},
			{
				Id: "command-to-replace",
				CommandUnion: dw.CommandUnion{
					Exec: &dw.ExecCommand{
						Env: []dw.EnvVar{
							{
								Name:  "envVarToReplace",
								Value: "envVarToReplaceNewValue",
							},
							{
								Name:  "endVarToAdd",
								Value: "endVarToAddValue",
							},
							{
								Name:  "envVarNotChanged",
								Value: "envVarNotChangedOriginalValue",
							},
						},
					},
				},
			},
			{
				Id: "command-not-changed",
				CommandUnion: dw.CommandUnion{
					Exec: &dw.ExecCommand{
						LabeledCommand: dw.LabeledCommand{
							Label: "commandNotChangedLabel",
						},
					},
				},
			},
		},
		Variables: map[string]string{
			"version": "patch",
			"xyz":     "xyz",
		},
		Attributes: attributesPkg.Attributes{}.FromMap(map[string]interface{}{
			"version": map[string]interface{}{
				"patch": true,
			},
			"xyz": "xyz",
		}, nil),
	}

	result, err := OverrideDevWorkspaceTemplateSpec(&original, &patch)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, expected, result, "The two values should be the same.")
}

func overridingPatchTest(original, patch []byte, expected dw.DevWorkspaceTemplateSpecContent, expectedError string) func(t *testing.T) {
	return func(t *testing.T) {
		actual, err := OverrideDevWorkspaceTemplateSpecBytes(original, patch)
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

func TestOverridingPatches(t *testing.T) {
	filepath.Walk("test-fixtures/patches", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && info.Name() == "original.yaml" {
			if err != nil {
				t.Error(err)
				return nil
			}
			original, err := ioutil.ReadFile(path)
			if err != nil {
				t.Error(err)
				return nil
			}
			dirPath := filepath.Dir(path)
			patch, err := ioutil.ReadFile(filepath.Join(dirPath, "patch.yaml"))
			if err != nil {
				t.Error(err)
				return nil
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

			t.Run(testName, overridingPatchTest(original, patch, resultTemplate, resultError))
		}
		return nil
	})
}

func TestPluginOverrides(t *testing.T) {
	originalFile := "test-fixtures/patches/override-just-plugin/original.yaml"
	patchFile := "test-fixtures/patches/override-just-plugin/patch.yaml"
	resultFile := "test-fixtures/patches/override-just-plugin/result.yaml"

	originalDWT := dw.DevWorkspaceTemplateSpecContent{}
	patch := dw.PluginOverrides{}
	expectedDWT := dw.DevWorkspaceTemplateSpecContent{}

	readFileToStruct(t, originalFile, &originalDWT)
	readFileToStruct(t, patchFile, &patch)
	readFileToStruct(t, resultFile, &expectedDWT)

	gotDWT, err := OverrideDevWorkspaceTemplateSpec(&originalDWT, patch)
	if assert.NoError(t, err) {
		assert.Equal(t, &expectedDWT, gotDWT)
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

// Since order of error message lines is not deterministic, it's necessary to compare
// in a weaker way than asserting string equality.
func compareErrorMessages(t *testing.T, expected, actual string, failReason string) {
	if expected == "" {
		t.Error("Received error but did not expect one: " + actual)
		return
	}
	expectedLines := strings.Split(strings.TrimSpace(expected), "\n")
	actualLines := strings.Split(strings.TrimSpace(actual), "\n")
	assert.ElementsMatch(t, expectedLines, actualLines, failReason)
}
