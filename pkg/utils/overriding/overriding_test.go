package overriding

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	workspaces "github.com/devfile/api/pkg/apis/workspaces/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/json"
	yamlMachinery "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

func TestBasicToplevelOverriding(t *testing.T) {
	original := workspaces.DevWorkspaceTemplateSpecContent{
		Commands: []workspaces.Command{
			{
				Exec: &workspaces.ExecCommand{
					LabeledCommand: workspaces.LabeledCommand{
						BaseCommand: workspaces.BaseCommand{
							Id: "commandWithTypeChanged",
						},
					},
				},
			},
			{
				Exec: &workspaces.ExecCommand{
					LabeledCommand: workspaces.LabeledCommand{
						BaseCommand: workspaces.BaseCommand{
							Id: "commandToReplace",
						},
					},
					Env: []workspaces.EnvVar{
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
			{
				Exec: &workspaces.ExecCommand{
					LabeledCommand: workspaces.LabeledCommand{
						BaseCommand: workspaces.BaseCommand{
							Id: "commandNotChanged",
						},
					},
				},
			},
		},
	}

	patch := workspaces.Overrides{
		OverridesBase: workspaces.OverridesBase{
			Commands: []workspaces.Command{
				{
					Apply: &workspaces.ApplyCommand{
						LabeledCommand: workspaces.LabeledCommand{
							BaseCommand: workspaces.BaseCommand{
								Id: "commandWithTypeChanged",
							},
						},
						Component: "mycomponent",
					},
				},
				{
					Exec: &workspaces.ExecCommand{
						LabeledCommand: workspaces.LabeledCommand{
							BaseCommand: workspaces.BaseCommand{
								Id: "commandToReplace",
							},
						},
						Env: []workspaces.EnvVar{
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
	}

	expected := &workspaces.DevWorkspaceTemplateSpecContent{
		Commands: []workspaces.Command{
			{
				Apply: &workspaces.ApplyCommand{
					LabeledCommand: workspaces.LabeledCommand{
						BaseCommand: workspaces.BaseCommand{
							Id: "commandWithTypeChanged",
						},
					},
					Component: "mycomponent",
				},
			},
			{
				Exec: &workspaces.ExecCommand{
					LabeledCommand: workspaces.LabeledCommand{
						BaseCommand: workspaces.BaseCommand{
							Id: "commandToReplace",
						},
					},
					Env: []workspaces.EnvVar{
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
			{
				Exec: &workspaces.ExecCommand{
					LabeledCommand: workspaces.LabeledCommand{
						BaseCommand: workspaces.BaseCommand{
							Id: "commandNotChanged",
						},
					},
				},
			},
		},
	}

	result, err := OverrideDevWorkspaceTemplateSpec(&original, patch)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, expected, result, "The two values should be the same.")
}

func overridingPatchTest(original, patch, expected []byte, expectedError string) func(t *testing.T) {
	return func(t *testing.T) {
		result, err := OverrideDevWorkspaceTemplateSpecBytes(original, patch)
		if err != nil {
			assert.Equal(t, strings.TrimSpace(expectedError), strings.TrimSpace(err.Error()), "Wrong error")
			return
		}

		resultJson, err := json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
		resultYaml, err := yaml.JSONToYAML(resultJson)
		if err != nil {
			t.Error(err)
		}

		expectedJson, err := yamlMachinery.ToJSON(expected)
		if err != nil {
			t.Error(err)
		}
		expectedYaml, err := yaml.JSONToYAML(expectedJson)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, string(expectedYaml), string(resultYaml), "The two values should be the same.")
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
			result := []byte{}
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
				result, err = ioutil.ReadFile(filepath.Join(dirPath, "result.yaml"))
				if err != nil {
					t.Error(err)
					return nil
				}
			}
			testName := filepath.Base(dirPath)

			t.Run(testName, overridingPatchTest(original, patch, result, resultError))
		}
		return nil
	})
}

func mergingPatchTest(main, parent, expected []byte, expectedError string, plugins ...[]byte) func(t *testing.T) {
	return func(t *testing.T) {
		result, err := MergeDevWorkspaceTemplateSpecBytes(main, parent, plugins...)
		if err != nil {
			assert.Equal(t, strings.TrimSpace(expectedError), strings.TrimSpace(err.Error()), "Wrong error")
			return
		}

		resultJson, err := json.Marshal(result)
		if err != nil {
			t.Error(err)
		}
		resultYaml, err := yaml.JSONToYAML(resultJson)
		if err != nil {
			t.Error(err)
		}

		expectedJson, err := yamlMachinery.ToJSON(expected)
		if err != nil {
			t.Error(err)
		}
		expectedYaml, err := yaml.JSONToYAML(expectedJson)
		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, string(expectedYaml), string(resultYaml), "The two values should be the same.")
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
			result := []byte{}
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
				result, err = ioutil.ReadFile(filepath.Join(dirPath, "result.yaml"))
				if err != nil {
					t.Error(err)
					return nil
				}
			}
			testName := filepath.Base(dirPath)

			t.Run(testName, mergingPatchTest(main, parent, result, resultError, plugins...))
		}
		return nil
	})
}
