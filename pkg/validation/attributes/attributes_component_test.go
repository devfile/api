package attributes

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
	"github.com/stretchr/testify/assert"
)

func TestValidateContainerComponent(t *testing.T) {

	tests := []struct {
		name       string
		testFile   string
		expected   v1alpha2.ContainerComponent
		attributes apiAttributes.Attributes
		wantErr    bool
	}{
		{
			name:     "Good Substitution",
			testFile: "test-fixtures/components/container.yaml",
			expected: v1alpha2.ContainerComponent{
				Container: v1alpha2.Container{
					Image:   "image-1",
					Command: []string{"tail", "-f", "/dev/null"},
					Args:    []string{"/dev/null"},
					Env: []v1alpha2.EnvVar{
						{
							Name:  "FOO",
							Value: "BAR",
						},
					},
					VolumeMounts: []v1alpha2.VolumeMount{
						{
							Name: "vol1",
							Path: "/FOO",
						},
					},
					MemoryLimit:   "FOO",
					MemoryRequest: "FOO",
					SourceMapping: "FOO",
				},
				Endpoints: []v1alpha2.Endpoint{
					{
						Name:     "endpoint1",
						Exposure: "public",
						Path:     "FOO",
					},
				},
			},
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"version": "1",
				"devnull": "/dev/null",
				"bar":     "BAR",
				"foo":     "FOO",
			}, nil),
			wantErr: false,
		},
		{
			name:     "Invalid Reference",
			testFile: "test-fixtures/components/container.yaml",
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"foo": "FOO",
			}, nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testContainerComponent := v1alpha2.ContainerComponent{}

			readFileToStruct(t, tt.testFile, &testContainerComponent)

			err := validateContainerComponent(tt.attributes, &testContainerComponent)
			if tt.wantErr == (err == nil) {
				t.Errorf("error: %v", err)
				return
			} else if err == nil {
				assert.Equal(t, tt.expected, testContainerComponent, "The two values should be the same.")
			}
		})
	}
}

func TestValidateOpenShiftAndKubernetesComponent(t *testing.T) {

	expectedKubeLikeComp := v1alpha2.K8sLikeComponent{
		K8sLikeComponentLocation: v1alpha2.K8sLikeComponentLocation{
			Uri:     "uri",
			Inlined: "inlined",
		},
		Endpoints: []v1alpha2.Endpoint{
			{
				Name:     "endpoint1",
				Exposure: "public",
				Path:     "FOO",
			},
		},
	}

	tests := []struct {
		name               string
		testFile           string
		expectedOpenShift  v1alpha2.OpenshiftComponent
		expectedKubernetes v1alpha2.KubernetesComponent
		attributes         apiAttributes.Attributes
		wantErr            bool
	}{
		{
			name:     "Good Substitution",
			testFile: "test-fixtures/components/openshift-kubernetes.yaml",
			expectedOpenShift: v1alpha2.OpenshiftComponent{
				K8sLikeComponent: expectedKubeLikeComp,
			},
			expectedKubernetes: v1alpha2.KubernetesComponent{
				K8sLikeComponent: expectedKubeLikeComp,
			},
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"uri":     "uri",
				"inlined": "inlined",
				"foo":     "FOO",
			}, nil),
			wantErr: false,
		},
		{
			name:     "Invalid Reference",
			testFile: "test-fixtures/components/openshift-kubernetes.yaml",
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"foo": "FOO",
			}, nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testOpenshiftComponent := v1alpha2.OpenshiftComponent{}
			testKubernetesComponent := v1alpha2.KubernetesComponent{}

			readFileToStruct(t, tt.testFile, &testOpenshiftComponent)
			readFileToStruct(t, tt.testFile, &testKubernetesComponent)

			err := validateOpenShiftComponent(tt.attributes, &testOpenshiftComponent)
			if tt.wantErr == (err == nil) {
				t.Errorf("error: %v", err)
				return
			} else if err == nil {
				assert.Equal(t, tt.expectedOpenShift, testOpenshiftComponent, "The two values should be the same.")
			}

			err = validateKubernetesComponent(tt.attributes, &testKubernetesComponent)
			if tt.wantErr == (err == nil) {
				t.Errorf("error: %v", err)
				return
			} else if err == nil {
				assert.Equal(t, tt.expectedKubernetes, testKubernetesComponent, "The two values should be the same.")
			}
		})
	}
}

func TestValidateVolumeComponent(t *testing.T) {

	tests := []struct {
		name       string
		testFile   string
		expected   v1alpha2.VolumeComponent
		attributes apiAttributes.Attributes
		wantErr    bool
	}{
		{
			name:     "Good Substitution",
			testFile: "test-fixtures/components/volume.yaml",
			expected: v1alpha2.VolumeComponent{
				Volume: v1alpha2.Volume{
					Size: "1Gi",
				},
			},
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"size": "1Gi",
			}, nil),
			wantErr: false,
		},
		{
			name:     "Invalid Reference",
			testFile: "test-fixtures/components/volume.yaml",
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"foo": "FOO",
			}, nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testVolumeComponent := v1alpha2.VolumeComponent{}

			readFileToStruct(t, tt.testFile, &testVolumeComponent)

			err := validateVolumeComponent(tt.attributes, &testVolumeComponent)
			if tt.wantErr == (err == nil) {
				t.Errorf("error: %v", err)
				return
			} else if err == nil {
				assert.Equal(t, tt.expected, testVolumeComponent, "The two values should be the same.")
			}
		})
	}
}

func TestValidateEnv(t *testing.T) {

	tests := []struct {
		name       string
		testFile   string
		expected   []v1alpha2.EnvVar
		attributes apiAttributes.Attributes
		wantErr    bool
	}{
		{
			name:     "Good Substitution",
			testFile: "test-fixtures/components/env.yaml",
			expected: []v1alpha2.EnvVar{
				{
					Name:  "FOO",
					Value: "BAR",
				},
			},
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"foo": "FOO",
				"bar": "BAR",
			}, nil),
			wantErr: false,
		},
		{
			name:     "Invalid Reference",
			testFile: "test-fixtures/components/env.yaml",
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"foo": "FOO",
			}, nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testEnv := v1alpha2.EnvVar{}

			readFileToStruct(t, tt.testFile, &testEnv)

			testEnvArr := []v1alpha2.EnvVar{testEnv}

			err := validateEnv(tt.attributes, &testEnvArr)
			if tt.wantErr == (err == nil) {
				t.Errorf("error: %v", err)
				return
			} else if err == nil {
				assert.Equal(t, tt.expected, testEnvArr, "The two values should be the same.")
			}
		})
	}
}
