package overriding

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	dw "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
	yamlMachinery "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

func mergingPatchTest(main, parent, expected []byte, expectedError string, plugins ...[]byte) func(t *testing.T) {
	return func(t *testing.T) {
		result, err := MergeDevWorkspaceTemplateSpecBytes(main, parent, plugins...)
		if err != nil {
			compareErrorMessages(t, expectedError, err.Error(), "wrong error")
			return
		}
		if expectedError != "" {
			t.Error("Expected error but did not get one")
			return
		}

		resultYaml, err := yaml.Marshal(result)
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
