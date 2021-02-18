package overriding

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	workspaces "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	attributesPkg "github.com/devfile/api/v2/pkg/attributes"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/json"
	yamlMachinery "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

func TestBasicMerging(t *testing.T) {

	tests := []struct {
		name                    string
		mainContent             *workspaces.DevWorkspaceTemplateSpecContent
		parentFlattenedContent  *workspaces.DevWorkspaceTemplateSpecContent
		pluginFlattenedContents []*workspaces.DevWorkspaceTemplateSpecContent
		expected                *workspaces.DevWorkspaceTemplateSpecContent
		wantErr                 *string
	}{
		{
			name: "Basic Merging",
			mainContent: &workspaces.DevWorkspaceTemplateSpecContent{
				Attributes: attributesPkg.Attributes{}.FromMap(map[string]interface{}{
					"main": true,
				}, nil),
				Commands: []workspaces.Command{
					{
						Id: "mainCommand",
						CommandUnion: workspaces.CommandUnion{
							Exec: &workspaces.ExecCommand{
								WorkingDir: "dir",
							},
						},
					},
				},
				Components: []workspaces.Component{
					{
						Name: "mainComponent",
						ComponentUnion: workspaces.ComponentUnion{
							Container: &workspaces.ContainerComponent{
								Container: workspaces.Container{
									Image: "image",
								},
							},
						},
					},
					{
						Name: "mainPluginComponent",
						ComponentUnion: workspaces.ComponentUnion{
							Plugin: &workspaces.PluginComponent{
								ImportReference: workspaces.ImportReference{
									ImportReferenceUnion: workspaces.ImportReferenceUnion{
										Uri: "uri",
									},
								},
							},
						},
					},
				},
				Events: &workspaces.Events{
					WorkspaceEvents: workspaces.WorkspaceEvents{
						PostStop: []string{"post-stop-main"},
					},
				},
			},
			pluginFlattenedContents: []*workspaces.DevWorkspaceTemplateSpecContent{
				{
					Attributes: attributesPkg.Attributes{}.FromMap(map[string]interface{}{
						"version2": "plugin",
					}, nil),
					Commands: []workspaces.Command{
						{
							Id: "pluginCommand",
							CommandUnion: workspaces.CommandUnion{
								Exec: &workspaces.ExecCommand{
									WorkingDir: "dir",
								},
							},
						},
					},
					Components: []workspaces.Component{
						{
							Name: "pluginComponent",
							ComponentUnion: workspaces.ComponentUnion{
								Container: &workspaces.ContainerComponent{
									Container: workspaces.Container{
										Image: "image",
									},
								},
							},
						},
					},
					Events: &workspaces.Events{
						WorkspaceEvents: workspaces.WorkspaceEvents{
							PostStop: []string{"post-stop-plugin"},
						},
					},
				},
			},
			parentFlattenedContent: &workspaces.DevWorkspaceTemplateSpecContent{
				Attributes: attributesPkg.Attributes{}.FromMap(map[string]interface{}{
					"version": "parent",
				}, nil),
				Commands: []workspaces.Command{
					{
						Id: "parentCommand",
						CommandUnion: workspaces.CommandUnion{
							Exec: &workspaces.ExecCommand{
								WorkingDir: "dir",
							},
						},
					},
				},
				Components: []workspaces.Component{
					{
						Name: "parentComponent",
						ComponentUnion: workspaces.ComponentUnion{
							Container: &workspaces.ContainerComponent{
								Container: workspaces.Container{
									Image: "image",
								},
							},
						},
					},
				},
				Events: &workspaces.Events{
					WorkspaceEvents: workspaces.WorkspaceEvents{
						PostStop:  []string{"post-stop-parent"},
						PostStart: []string{"post-start-parent"},
					},
				},
			},
			expected: &workspaces.DevWorkspaceTemplateSpecContent{
				Attributes: attributesPkg.Attributes{}.FromMap(map[string]interface{}{
					"version":  "parent",
					"version2": "plugin",
					"main":     true,
				}, nil),
				Commands: []workspaces.Command{
					{
						Id: "parentCommand",
						CommandUnion: workspaces.CommandUnion{
							Exec: &workspaces.ExecCommand{
								WorkingDir: "dir",
							},
						},
					},
					{
						Id: "pluginCommand",
						CommandUnion: workspaces.CommandUnion{
							Exec: &workspaces.ExecCommand{
								WorkingDir: "dir",
							},
						},
					},
					{
						Id: "mainCommand",
						CommandUnion: workspaces.CommandUnion{
							Exec: &workspaces.ExecCommand{
								WorkingDir: "dir",
							},
						},
					},
				},
				Components: []workspaces.Component{
					{
						Name: "parentComponent",
						ComponentUnion: workspaces.ComponentUnion{
							Container: &workspaces.ContainerComponent{
								Container: workspaces.Container{
									Image: "image",
								},
							},
						},
					},
					{
						Name: "pluginComponent",
						ComponentUnion: workspaces.ComponentUnion{
							Container: &workspaces.ContainerComponent{
								Container: workspaces.Container{
									Image: "image",
								},
							},
						},
					},
					{
						Name: "mainComponent",
						ComponentUnion: workspaces.ComponentUnion{
							Container: &workspaces.ContainerComponent{
								Container: workspaces.Container{
									Image: "image",
								},
							},
						},
					},
				},
				Events: &workspaces.Events{
					WorkspaceEvents: workspaces.WorkspaceEvents{
						PreStart:  []string{},
						PostStart: []string{"post-start-parent"},
						PreStop:   []string{},
						PostStop:  []string{"post-stop-main", "post-stop-parent", "post-stop-plugin"},
					},
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mergedContent, err := MergeDevWorkspaceTemplateSpec(tt.mainContent, tt.parentFlattenedContent, tt.pluginFlattenedContents...)
			if err != nil {
				t.Error(err)
				return
			}

			assert.Equal(t, tt.expected, mergedContent, "The two values should be the same.")
		})
	}
}

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

func TestMergingOnlyPlugins(t *testing.T) {
	baseFile := "test-fixtures/merges/no-parent/main.yaml"
	pluginFile := "test-fixtures/merges/no-parent/plugin.yaml"
	resultFile := "test-fixtures/merges/no-parent/result.yaml"

	baseDWT := workspaces.DevWorkspaceTemplateSpecContent{}
	pluginDWT := workspaces.DevWorkspaceTemplateSpecContent{}
	expectedDWT := workspaces.DevWorkspaceTemplateSpecContent{}

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

	baseDWT := workspaces.DevWorkspaceTemplateSpecContent{}
	parentDWT := workspaces.DevWorkspaceTemplateSpecContent{}
	expectedDWT := workspaces.DevWorkspaceTemplateSpecContent{}

	readFileToStruct(t, baseFile, &baseDWT)
	readFileToStruct(t, parentFile, &parentDWT)
	readFileToStruct(t, resultFile, &expectedDWT)

	gotDWT, err := MergeDevWorkspaceTemplateSpec(&baseDWT, &parentDWT)
	if assert.NoError(t, err) {
		assert.Equal(t, &expectedDWT, gotDWT)
	}
}
