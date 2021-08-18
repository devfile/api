package variables

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
)

func TestValidateAndReplaceContainerComponent(t *testing.T) {

	tests := []struct {
		name         string
		testFile     string
		outputFile   string
		variableFile string
		wantErr      bool
	}{
		{
			name:         "Good Substitution",
			testFile:     "test-fixtures/components/container.yaml",
			outputFile:   "test-fixtures/components/container-output.yaml",
			variableFile: "test-fixtures/variables/variables-referenced.yaml",
			wantErr:      false,
		},
		{
			name:         "Invalid Reference",
			testFile:     "test-fixtures/components/container.yaml",
			outputFile:   "test-fixtures/components/container.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testContainerComponent := v1alpha2.ContainerComponent{}
			readFileToStruct(t, tt.testFile, &testContainerComponent)

			testVariable := make(map[string]string)
			readFileToStruct(t, tt.variableFile, &testVariable)

			err := validateAndReplaceForContainerComponent(testVariable, &testContainerComponent)
			_, ok := err.(*InvalidKeysError)
			if tt.wantErr && !ok {
				t.Errorf("Expected InvalidKeysError error from test but got %+v", err)
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else {
				expectedContainerComponent := v1alpha2.ContainerComponent{}
				readFileToStruct(t, tt.outputFile, &expectedContainerComponent)
				assert.Equal(t, expectedContainerComponent, testContainerComponent, "The two values should be the same.")
			}
		})
	}
}

func TestValidateAndReplaceOpenShiftKubernetesComponent(t *testing.T) {

	tests := []struct {
		name         string
		testFile     string
		outputFile   string
		variableFile string
		wantErr      bool
	}{
		{
			name:         "Good Substitution",
			testFile:     "test-fixtures/components/openshift-kubernetes.yaml",
			outputFile:   "test-fixtures/components/openshift-kubernetes-output.yaml",
			variableFile: "test-fixtures/variables/variables-referenced.yaml",
			wantErr:      false,
		},
		{
			name:         "Invalid Reference",
			testFile:     "test-fixtures/components/openshift-kubernetes.yaml",
			outputFile:   "test-fixtures/components/openshift-kubernetes.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOpenshiftComponent := v1alpha2.OpenshiftComponent{}
			testKubernetesComponent := v1alpha2.KubernetesComponent{}

			readFileToStruct(t, tt.testFile, &testOpenshiftComponent)
			readFileToStruct(t, tt.testFile, &testKubernetesComponent)

			testVariable := make(map[string]string)
			readFileToStruct(t, tt.variableFile, &testVariable)

			err := validateAndReplaceForOpenShiftComponent(testVariable, &testOpenshiftComponent)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error from test but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else if err == nil {
				expectedOpenshiftComponent := v1alpha2.OpenshiftComponent{}
				readFileToStruct(t, tt.outputFile, &expectedOpenshiftComponent)
				assert.Equal(t, expectedOpenshiftComponent, testOpenshiftComponent, "The two values should be the same.")
			}

			err = validateAndReplaceForKubernetesComponent(testVariable, &testKubernetesComponent)
			_, ok := err.(*InvalidKeysError)
			if tt.wantErr && !ok {
				t.Errorf("Expected InvalidKeysError error from test but got %+v", err)
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else {
				expectedKubernetesComponent := v1alpha2.KubernetesComponent{}
				readFileToStruct(t, tt.outputFile, &expectedKubernetesComponent)
				assert.Equal(t, expectedKubernetesComponent, testKubernetesComponent, "The two values should be the same.")
			}
		})
	}
}

func TestValidateAndReplaceImageComponent(t *testing.T) {

	tests := []struct {
		name         string
		testFile     string
		outputFile   string
		variableFile string
		wantErr      bool
	}{
		{
			name:         "Good Substitution - dockerfile uri src",
			testFile:     "test-fixtures/components/image-dockerfile-uri.yaml",
			outputFile:   "test-fixtures/components/image-dockerfile-uri-output.yaml",
			variableFile: "test-fixtures/variables/variables-referenced.yaml",
			wantErr:      false,
		},
		{
			name:         "Good Substitution - dockerfile git src",
			testFile:     "test-fixtures/components/image-dockerfile-git.yaml",
			outputFile:   "test-fixtures/components/image-dockerfile-git-output.yaml",
			variableFile: "test-fixtures/variables/variables-referenced.yaml",
			wantErr:      false,
		},
		{
			name:         "Good Substitution - dockerfile registry src",
			testFile:     "test-fixtures/components/image-dockerfile-registry.yaml",
			outputFile:   "test-fixtures/components/image-dockerfile-registry-output.yaml",
			variableFile: "test-fixtures/variables/variables-referenced.yaml",
			wantErr:      false,
		},
		{
			name:         "Invalid Reference - dockerfile uri src",
			testFile:     "test-fixtures/components/image-dockerfile-uri.yaml",
			outputFile:   "test-fixtures/components/image-dockerfile-uri.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      true,
		},
		{
			name:         "Invalid Reference - dockerfile git src",
			testFile:     "test-fixtures/components/image-dockerfile-git.yaml",
			outputFile:   "test-fixtures/components/image-dockerfile-git.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      true,
		},
		{
			name:         "Invalid Reference - dockerfile registry src",
			testFile:     "test-fixtures/components/image-dockerfile-registry.yaml",
			outputFile:   "test-fixtures/components/image-dockerfile-registry.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testImageComponent := v1alpha2.ImageComponent{}
			readFileToStruct(t, tt.testFile, &testImageComponent)

			testVariable := make(map[string]string)
			readFileToStruct(t, tt.variableFile, &testVariable)

			err := validateAndReplaceForImageComponent(testVariable, &testImageComponent)
			_, ok := err.(*InvalidKeysError)
			if tt.wantErr && !ok {
				t.Errorf("Expected InvalidKeysError error from test but got %+v", err)
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else {
				expectedImageComponent := v1alpha2.ImageComponent{}
				readFileToStruct(t, tt.outputFile, &expectedImageComponent)
				assert.Equal(t, expectedImageComponent, testImageComponent, "The two values should be the same.")
			}
		})
	}
}

func TestValidateAndReplaceVolumeComponent(t *testing.T) {

	tests := []struct {
		name         string
		testFile     string
		outputFile   string
		variableFile string
		wantErr      bool
	}{
		{
			name:         "Good Substitution",
			testFile:     "test-fixtures/components/volume.yaml",
			outputFile:   "test-fixtures/components/volume-output.yaml",
			variableFile: "test-fixtures/variables/variables-referenced.yaml",
			wantErr:      false,
		},
		{
			name:         "Invalid Reference",
			testFile:     "test-fixtures/components/volume.yaml",
			outputFile:   "test-fixtures/components/volume.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testVolumeComponent := v1alpha2.VolumeComponent{}
			readFileToStruct(t, tt.testFile, &testVolumeComponent)

			testVariable := make(map[string]string)
			readFileToStruct(t, tt.variableFile, &testVariable)

			err := validateAndReplaceForVolumeComponent(testVariable, &testVolumeComponent)
			_, ok := err.(*InvalidKeysError)
			if tt.wantErr && !ok {
				t.Errorf("Expected InvalidKeysError error from test but got %+v", err)
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else {
				expectedVolumeComponent := v1alpha2.VolumeComponent{}
				readFileToStruct(t, tt.outputFile, &expectedVolumeComponent)
				assert.Equal(t, expectedVolumeComponent, testVolumeComponent, "The two values should be the same.")
			}
		})
	}
}

func TestValidateAndReplaceEnv(t *testing.T) {

	tests := []struct {
		name         string
		testFile     string
		outputFile   string
		variableFile string
		wantErr      bool
	}{
		{
			name:         "Good Substitution",
			testFile:     "test-fixtures/components/env.yaml",
			outputFile:   "test-fixtures/components/env-output.yaml",
			variableFile: "test-fixtures/variables/variables-referenced.yaml",
			wantErr:      false,
		},
		{
			name:         "Invalid Reference",
			testFile:     "test-fixtures/components/env.yaml",
			outputFile:   "test-fixtures/components/env.yaml",
			variableFile: "test-fixtures/variables/variables-notreferenced.yaml",
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testEnv := v1alpha2.EnvVar{}
			readFileToStruct(t, tt.testFile, &testEnv)
			testEnvArr := []v1alpha2.EnvVar{testEnv}

			testVariable := make(map[string]string)
			readFileToStruct(t, tt.variableFile, &testVariable)

			err := validateAndReplaceForEnv(testVariable, testEnvArr)
			_, ok := err.(*InvalidKeysError)
			if tt.wantErr && !ok {
				t.Errorf("Expected InvalidKeysError error from test but got %+v", err)
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			} else {
				expectedEnv := v1alpha2.EnvVar{}
				readFileToStruct(t, tt.outputFile, &expectedEnv)
				expectedEnvArr := []v1alpha2.EnvVar{expectedEnv}
				assert.Equal(t, expectedEnvArr, testEnvArr, "The two values should be the same.")
			}
		})
	}
}
