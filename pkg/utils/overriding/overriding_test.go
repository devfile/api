package overriding

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	workspaces "github.com/devfile/kubernetes-api/pkg/apis/workspaces/v1alpha1"
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
				{
					Exec: &workspaces.ExecCommand{
						LabeledCommand: workspaces.LabeledCommand{
							BaseCommand: workspaces.BaseCommand{
								Id: "commandToAdd",
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
							Id: "commandToAdd",
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

func overridingPatchTest(original, patch, expected []byte) func(t *testing.T) {
	return func(t *testing.T) {
		result, err := OverrideDevWorkspaceTemplateSpecBytes(original, patch)
		if err != nil {
			t.Error(err)
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
	filepath.Walk("test-fixtures", func(path string, info os.FileInfo, err error) error {
		if ! info.IsDir() && info.Name() == "original.yaml" {
			if err != nil {
				return err
			}
			original, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			dirPath := filepath.Dir(path)
			patch, err := ioutil.ReadFile(filepath.Join(dirPath, "patch.yaml"))
			if err != nil {
				return err
			}
			result, err := ioutil.ReadFile(filepath.Join(dirPath, "result.yaml"))
			if err != nil {
				return err
			}
			testName := filepath.Base(dirPath)

			t.Run(testName, overridingPatchTest(original, patch, result))
		}
		return nil
	})
}


