package attributes

import (
	"io/ioutil"
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestValidateGlobalAttributeBasic(t *testing.T) {

	tests := []struct {
		name     string
		testFile string
		expected v1alpha2.DevWorkspaceTemplateSpec
		wantErr  bool
	}{
		{
			name:     "Successful global attribute substitution",
			testFile: "test-fixtures/all/devfile-good.yaml",
			expected: v1alpha2.DevWorkspaceTemplateSpec{
				DevWorkspaceTemplateSpecContent: v1alpha2.DevWorkspaceTemplateSpecContent{
					Attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
						"tag":     "xyz",
						"version": "1",
						"foo":     "FOO",
						"devnull": "/dev/null",
					}, nil),
					Components: []v1alpha2.Component{
						{
							Name: "component1",
							ComponentUnion: v1alpha2.ComponentUnion{
								Container: &v1alpha2.ContainerComponent{
									Container: v1alpha2.Container{
										Image:   "image",
										Command: []string{"tail", "-f", "/dev/null"},
										Env: []v1alpha2.EnvVar{
											{
												Name:  "BAR",
												Value: "FOO",
											},
											{
												Name:  "FOO",
												Value: "BAR",
											},
										},
									},
								},
							},
						},
						{
							Name: "component2",
							ComponentUnion: v1alpha2.ComponentUnion{
								Kubernetes: &v1alpha2.KubernetesComponent{
									K8sLikeComponent: v1alpha2.K8sLikeComponent{
										K8sLikeComponentLocation: v1alpha2.K8sLikeComponentLocation{
											Inlined: "FOO",
										},
										Endpoints: []v1alpha2.Endpoint{
											{
												Name:       "endpoint1",
												Exposure:   "public",
												TargetPort: 9999,
											},
										},
									},
								},
							},
						},
					},
					Commands: []v1alpha2.Command{
						{
							Id: "command1",
							CommandUnion: v1alpha2.CommandUnion{
								Exec: &v1alpha2.ExecCommand{
									CommandLine: "test-xyz",
									Env: []v1alpha2.EnvVar{
										{
											Name:  "tag",
											Value: "xyz",
										},
										{
											Name:  "FOO",
											Value: "BAR",
										},
									},
								},
							},
						},
						{
							Id: "command2",
							CommandUnion: v1alpha2.CommandUnion{
								Composite: &v1alpha2.CompositeCommand{
									Commands: []string{
										"xyz",
										"command1",
									},
								},
							},
						},
					},
					Events: &v1alpha2.Events{
						WorkspaceEvents: v1alpha2.WorkspaceEvents{
							PreStart: []string{
								"xyz",
								"test",
							},
							PreStop: []string{
								"1",
							},
						},
					},
					Projects: []v1alpha2.Project{
						{
							Name: "project1",
							ProjectSource: v1alpha2.ProjectSource{
								Git: &v1alpha2.GitProjectSource{
									GitLikeProjectSource: v1alpha2.GitLikeProjectSource{
										CheckoutFrom: &v1alpha2.CheckoutFrom{
											Revision: "xyz",
										},
										Remotes: map[string]string{
											"xyz": "/dev/null",
											"1":   "test",
										},
									},
								},
							},
						},
						{
							Name: "project2",
							ProjectSource: v1alpha2.ProjectSource{
								Zip: &v1alpha2.ZipProjectSource{
									Location: "xyz",
								},
							},
						},
					},
					StarterProjects: []v1alpha2.StarterProject{
						{
							Name: "starterproject1",
							ProjectSource: v1alpha2.ProjectSource{
								Git: &v1alpha2.GitProjectSource{
									GitLikeProjectSource: v1alpha2.GitLikeProjectSource{
										CheckoutFrom: &v1alpha2.CheckoutFrom{
											Revision: "xyz",
										},
										Remotes: map[string]string{
											"xyz": "/dev/null",
											"1":   "test",
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:     "Invalid Reference",
			testFile: "test-fixtures/all/devfile-bad.yaml",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDWT := v1alpha2.DevWorkspaceTemplateSpec{}

			readFileToStruct(t, tt.testFile, &testDWT)

			err := ValidateGlobalAttribute(&testDWT)
			if tt.wantErr == (err == nil) {
				t.Errorf("error: %v", err)
				return
			} else if err == nil {
				assert.Equal(t, tt.expected, testDWT, "The two values should be the same.")
			}
		})
	}
}

func TestValidateAndReplaceDataWithAttribute(t *testing.T) {

	invalidAttributeErr := ".*Attribute with key .* does not exist.*"
	wrongAttributeTypeErr := ".*cannot unmarshal object into Go value of type string.*"

	tests := []struct {
		name       string
		testString string
		attributes apiAttributes.Attributes
		wantValue  string
		wantErr    *string
	}{
		{
			name:       "Valid attribute reference",
			testString: "image-{{version}}:{{tag}}-14",
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"version": "1.x.x",
				"tag":     "dev",
				"import": map[string]interface{}{
					"strategy": "Dockerfile",
				},
			}, nil),
			wantValue: "image-1.x.x:dev-14",
			wantErr:   nil,
		},
		{
			name:       "Invalid attribute reference",
			testString: "image-{{version}}:{{invalid}}-14",
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"version": "1.x.x",
				"tag":     "dev",
			}, nil),
			wantErr: &invalidAttributeErr,
		},
		{
			name:       "Attribute reference with non-string type value",
			testString: "image-{{version}}:{{invalid}}-14",
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"version": "1.x.x",
				"invalid": map[string]interface{}{
					"key": "value",
				},
			}, nil),
			wantErr: &wrongAttributeTypeErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotValue, err := validateAndReplaceDataWithAttribute(tt.testString, tt.attributes)
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
