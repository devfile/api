package variables

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestValidateGlobalVariableBasic(t *testing.T) {

	tests := []struct {
		name        string
		testFile    string
		outputFile  string
		wantWarning VariableWarning
	}{
		{
			name:        "Successful global variable substitution",
			testFile:    "test-fixtures/all/devfile-good.yaml",
			outputFile:  "test-fixtures/all/devfile-good-output.yaml",
			wantWarning: VariableWarning{},
		},
		{
			name:       "Invalid Reference",
			testFile:   "test-fixtures/all/devfile-bad.yaml",
			outputFile: "test-fixtures/all/devfile-bad-output.yaml",
			wantWarning: VariableWarning{
				Commands: map[string][]string{
					"command1": {"tag", "BAR"},
					"command2": {"abc"},
					"command3": {"abc"},
				},
				Components: map[string][]string{
					"component1": {"a", "b", "c", "bar"},
					"component2": {"foo", "x", "bar"},
					"component3": {"xyz"},
					"component4": {"foo"},
				},
				Projects: map[string][]string{
					"project1": {"tag", "version1", "path", "dir", "version"},
					"project2": {"tag"},
				},
				StarterProjects: map[string][]string{
					"starterproject1": {"tag", "desc", "dir"},
					"starterproject2": {"tag"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDWT := v1alpha2.DevWorkspaceTemplateSpec{}
			readFileToStruct(t, tt.testFile, &testDWT)

			warning := ValidateAndReplaceGlobalVariable(&testDWT)

			expectedDWT := v1alpha2.DevWorkspaceTemplateSpec{}
			readFileToStruct(t, tt.outputFile, &expectedDWT)
			assert.Equal(t, expectedDWT, testDWT, "The two values should be the same.")

			// match the warning
			if !reflect.DeepEqual(tt.wantWarning, VariableWarning{}) {
				// commands
				for gotCommand, gotInvalidVars := range warning.Commands {
					if wantInvalidVars, ok := tt.wantWarning.Commands[gotCommand]; !ok {
						t.Errorf("unexpected command %s found in the warning", gotCommand)
					} else {
						if isEqual := testStringArrElements(wantInvalidVars, gotInvalidVars); !isEqual {
							t.Errorf("wantInvalidVars %+v not equal as gotInvalidVars %+v", wantInvalidVars, gotInvalidVars)
						}
					}
				}

				// components
				for gotComponent, gotInvalidVars := range warning.Components {
					if wantInvalidVars, ok := tt.wantWarning.Components[gotComponent]; !ok {
						t.Errorf("unexpected component %s found in the warning", gotComponent)
					} else {
						if isEqual := testStringArrElements(wantInvalidVars, gotInvalidVars); !isEqual {
							t.Errorf("wantInvalidVars %+v not equal as gotInvalidVars %+v", wantInvalidVars, gotInvalidVars)
						}
					}
				}

				// projects
				for gotProject, gotInvalidVars := range warning.Projects {
					if wantInvalidVars, ok := tt.wantWarning.Projects[gotProject]; !ok {
						t.Errorf("unexpected project %s found in the warning", gotProject)
					} else {
						if isEqual := testStringArrElements(wantInvalidVars, gotInvalidVars); !isEqual {
							t.Errorf("wantInvalidVars %+v not equal as gotInvalidVars %+v", wantInvalidVars, gotInvalidVars)
						}
					}
				}

				// starter projects
				for gotStarterProject, gotInvalidVars := range warning.StarterProjects {
					if wantInvalidVars, ok := tt.wantWarning.StarterProjects[gotStarterProject]; !ok {
						t.Errorf("unexpected starter project %s found in the warning", gotStarterProject)
					} else {
						if isEqual := testStringArrElements(wantInvalidVars, gotInvalidVars); !isEqual {
							t.Errorf("wantInvalidVars %+v not equal as gotInvalidVars %+v", wantInvalidVars, gotInvalidVars)
						}
					}
				}
			}
		})
	}
}

func TestValidateAndReplaceDataWithVariable(t *testing.T) {

	invalidVariableErr := ".*invalid variable references.*"

	tests := []struct {
		name       string
		testString string
		variables  map[string]string
		wantValue  string
		wantErr    *string
	}{
		{
			name:       "Valid variable reference",
			testString: "image-{{version}}:{{tag}}{{development}}-14",
			variables: map[string]string{
				"version":     "1.x.x",
				"tag":         "dev",
				"development": "sandbox",
			},
			wantValue: "image-1.x.x:devsandbox-14",
			wantErr:   nil,
		},
		{
			name:       "Invalid variable reference",
			testString: "image-{{version}}:{{tag}}{{invalid}}-14{{invalid}}",
			variables: map[string]string{
				"version": "1.x.x",
				"tag":     "dev",
			},
			wantValue: "image-1.x.x:dev{{invalid}}-14{{invald}}",
			wantErr:   &invalidVariableErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, err := validateAndReplaceDataWithVariable(tt.testString, tt.variables)
			if tt.wantErr != nil && assert.Error(t, err) {
				assert.Regexp(t, *tt.wantErr, err.Error(), "Error message should match")
			} else {
				assert.NoError(t, err, "Expected error to be nil")
				if gotValue != tt.wantValue {
					assert.Equal(t, tt.wantValue, gotValue, "Return value should match")
				}
			}
		})
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
