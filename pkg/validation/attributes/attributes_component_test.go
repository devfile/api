package attributes

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
	"github.com/stretchr/testify/assert"
)

func TestValidateAndReplaceContainerComponent(t *testing.T) {

	tests := []struct {
		name          string
		testFile      string
		outputFile    string
		attributeFile string
		wantErr       bool
	}{
		{
			name:          "Good Substitution",
			testFile:      "test-fixtures/components/container.yaml",
			outputFile:    "test-fixtures/components/container-output.yaml",
			attributeFile: "test-fixtures/attributes/attributes-referenced.yaml",
			wantErr:       false,
		},
		{
			name:          "Invalid Reference",
			testFile:      "test-fixtures/components/container.yaml",
			attributeFile: "test-fixtures/attributes/attributes-notreferenced.yaml",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testContainerComponent := v1alpha2.ContainerComponent{}
			readFileToStruct(t, tt.testFile, &testContainerComponent)

			testAttribute := apiAttributes.Attributes{}
			readFileToStruct(t, tt.attributeFile, &testAttribute)

			err := validateAndReplaceForContainerComponent(testAttribute, &testContainerComponent)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error from test but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else if err == nil {
				expectedContainerComponent := v1alpha2.ContainerComponent{}
				readFileToStruct(t, tt.outputFile, &expectedContainerComponent)
				assert.Equal(t, expectedContainerComponent, testContainerComponent, "The two values should be the same.")
			}
		})
	}
}

func TestValidateAndReplaceOpenShiftKubernetesComponent(t *testing.T) {

	tests := []struct {
		name          string
		testFile      string
		outputFile    string
		attributeFile string
		wantErr       bool
	}{
		{
			name:          "Good Substitution",
			testFile:      "test-fixtures/components/openshift-kubernetes.yaml",
			outputFile:    "test-fixtures/components/openshift-kubernetes-output.yaml",
			attributeFile: "test-fixtures/attributes/attributes-referenced.yaml",
			wantErr:       false,
		},
		{
			name:          "Invalid Reference",
			testFile:      "test-fixtures/components/openshift-kubernetes.yaml",
			attributeFile: "test-fixtures/attributes/attributes-notreferenced.yaml",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOpenshiftComponent := v1alpha2.OpenshiftComponent{}
			testKubernetesComponent := v1alpha2.KubernetesComponent{}

			readFileToStruct(t, tt.testFile, &testOpenshiftComponent)
			readFileToStruct(t, tt.testFile, &testKubernetesComponent)

			testAttribute := apiAttributes.Attributes{}
			readFileToStruct(t, tt.attributeFile, &testAttribute)

			err := validateAndReplaceForOpenShiftComponent(testAttribute, &testOpenshiftComponent)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error from test but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else if err == nil {
				expectedOpenshiftComponent := v1alpha2.OpenshiftComponent{}
				readFileToStruct(t, tt.outputFile, &expectedOpenshiftComponent)
				assert.Equal(t, expectedOpenshiftComponent, testOpenshiftComponent, "The two values should be the same.")
			}

			err = validateAndReplaceForKubernetesComponent(testAttribute, &testKubernetesComponent)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error from test but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else if err == nil {
				expectedKubernetesComponent := v1alpha2.KubernetesComponent{}
				readFileToStruct(t, tt.outputFile, &expectedKubernetesComponent)
				assert.Equal(t, expectedKubernetesComponent, testKubernetesComponent, "The two values should be the same.")
			}
		})
	}
}

func TestValidateAndReplaceVolumeComponent(t *testing.T) {

	tests := []struct {
		name          string
		testFile      string
		outputFile    string
		attributeFile string
		wantErr       bool
	}{
		{
			name:          "Good Substitution",
			testFile:      "test-fixtures/components/volume.yaml",
			outputFile:    "test-fixtures/components/volume-output.yaml",
			attributeFile: "test-fixtures/attributes/attributes-referenced.yaml",
			wantErr:       false,
		},
		{
			name:          "Invalid Reference",
			testFile:      "test-fixtures/components/volume.yaml",
			attributeFile: "test-fixtures/attributes/attributes-notreferenced.yaml",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testVolumeComponent := v1alpha2.VolumeComponent{}
			readFileToStruct(t, tt.testFile, &testVolumeComponent)

			testAttribute := apiAttributes.Attributes{}
			readFileToStruct(t, tt.attributeFile, &testAttribute)

			err := validateAndReplaceForVolumeComponent(testAttribute, &testVolumeComponent)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error from test but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else if err == nil {
				expectedVolumeComponent := v1alpha2.VolumeComponent{}
				readFileToStruct(t, tt.outputFile, &expectedVolumeComponent)
				assert.Equal(t, expectedVolumeComponent, testVolumeComponent, "The two values should be the same.")
			}
		})
	}
}

func TestValidateAndReplaceEnv(t *testing.T) {

	tests := []struct {
		name          string
		testFile      string
		outputFile    string
		attributeFile string
		wantErr       bool
	}{
		{
			name:          "Good Substitution",
			testFile:      "test-fixtures/components/env.yaml",
			outputFile:    "test-fixtures/components/env-output.yaml",
			attributeFile: "test-fixtures/attributes/attributes-referenced.yaml",
			wantErr:       false,
		},
		{
			name:          "Invalid Reference",
			testFile:      "test-fixtures/components/env.yaml",
			attributeFile: "test-fixtures/attributes/attributes-notreferenced.yaml",
			wantErr:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testEnv := v1alpha2.EnvVar{}
			readFileToStruct(t, tt.testFile, &testEnv)
			testEnvArr := []v1alpha2.EnvVar{testEnv}

			testAttribute := apiAttributes.Attributes{}
			readFileToStruct(t, tt.attributeFile, &testAttribute)

			err := validateAndReplaceForEnv(testAttribute, testEnvArr)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error from test but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else if err == nil {
				expectedEnv := v1alpha2.EnvVar{}
				readFileToStruct(t, tt.outputFile, &expectedEnv)
				expectedEnvArr := []v1alpha2.EnvVar{expectedEnv}
				assert.Equal(t, expectedEnvArr, testEnvArr, "The two values should be the same.")
			}
		})
	}
}
