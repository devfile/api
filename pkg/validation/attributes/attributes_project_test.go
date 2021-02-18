package attributes

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
	"github.com/stretchr/testify/assert"
)

func TestValidateProjects(t *testing.T) {

	tests := []struct {
		name       string
		testFile   string
		expected   []v1alpha2.Project
		attributes apiAttributes.Attributes
		wantErr    bool
	}{
		{
			name:     "Good Substitution",
			testFile: "test-fixtures/projects/project.yaml",
			expected: []v1alpha2.Project{
				{
					Name:      "project1",
					ClonePath: "/FOO",
					SparseCheckoutDirs: []string{
						"/FOO",
						"/BAR",
					},
					ProjectSource: v1alpha2.ProjectSource{
						Git: &v1alpha2.GitProjectSource{
							GitLikeProjectSource: v1alpha2.GitLikeProjectSource{
								CheckoutFrom: &v1alpha2.CheckoutFrom{
									Revision: "FOO",
								},
								Remotes: map[string]string{
									"foo":    "BAR",
									"FOOBAR": "BARFOO",
								},
							},
						},
					},
				},
			},
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"bar": "BAR",
				"foo": "FOO",
			}, nil),
			wantErr: false,
		},
		{
			name:     "Invalid Reference",
			testFile: "test-fixtures/projects/project.yaml",
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"foo": "FOO",
			}, nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testProject := v1alpha2.Project{}

			readFileToStruct(t, tt.testFile, &testProject)

			testProjectArr := []v1alpha2.Project{testProject}

			err := ValidateProjects(tt.attributes, &testProjectArr)
			if tt.wantErr == (err == nil) {
				t.Errorf("error: %v", err)
				return
			} else if err == nil {
				assert.Equal(t, tt.expected, testProjectArr, "The two values should be the same.")
			}
		})
	}
}

func TestValidateStarterProjects(t *testing.T) {

	tests := []struct {
		name       string
		testFile   string
		expected   []v1alpha2.StarterProject
		attributes apiAttributes.Attributes
		wantErr    bool
	}{
		{
			name:     "Good Substitution",
			testFile: "test-fixtures/projects/starterproject.yaml",
			expected: []v1alpha2.StarterProject{
				{
					Name:        "starterproject1",
					Description: "FOOBAR is not BARFOO",
					SubDir:      "/FOO",
					ProjectSource: v1alpha2.ProjectSource{
						Zip: &v1alpha2.ZipProjectSource{
							Location: "/FOO",
						},
					},
				},
			},
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"bar": "BAR",
				"foo": "FOO",
			}, nil),
			wantErr: false,
		},
		{
			name:     "Invalid Reference",
			testFile: "test-fixtures/projects/starterproject.yaml",
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"foo": "FOO",
			}, nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStarterProject := v1alpha2.StarterProject{}

			readFileToStruct(t, tt.testFile, &testStarterProject)

			testStarterProjectArr := []v1alpha2.StarterProject{testStarterProject}

			err := ValidateStarterProjects(tt.attributes, &testStarterProjectArr)
			if tt.wantErr == (err == nil) {
				t.Errorf("error: %v", err)
				return
			} else if err == nil {
				assert.Equal(t, tt.expected, testStarterProjectArr, "The two values should be the same.")
			}
		})
	}
}

func TestValidateProjectSrc(t *testing.T) {

	tests := []struct {
		name       string
		testFile   string
		expected   v1alpha2.ProjectSource
		attributes apiAttributes.Attributes
		wantErr    bool
	}{
		{
			name:     "Good Git Substitution",
			testFile: "test-fixtures/projects/git.yaml",
			expected: v1alpha2.ProjectSource{
				Git: &v1alpha2.GitProjectSource{
					GitLikeProjectSource: v1alpha2.GitLikeProjectSource{
						CheckoutFrom: &v1alpha2.CheckoutFrom{
							Revision: "FOO",
							Remote:   "BAR",
						},
						Remotes: map[string]string{
							"foo":    "BAR",
							"FOOBAR": "BARFOO",
						},
					},
				},
			},
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"bar": "BAR",
				"foo": "FOO",
			}, nil),
			wantErr: false,
		},
		{
			name:     "Good Zip Substitution",
			testFile: "test-fixtures/projects/zip.yaml",
			expected: v1alpha2.ProjectSource{
				Zip: &v1alpha2.ZipProjectSource{
					Location: "/FOO",
				},
			},
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"foo": "FOO",
			}, nil),
			wantErr: false,
		},
		{
			name:     "Invalid Git Reference",
			testFile: "test-fixtures/projects/git.yaml",
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"foo": "FOO",
			}, nil),
			wantErr: true,
		},
		{
			name:     "Invalid Zip Reference",
			testFile: "test-fixtures/projects/zip.yaml",
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"bar": "BAR",
			}, nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testProjectSrc := v1alpha2.ProjectSource{}

			readFileToStruct(t, tt.testFile, &testProjectSrc)

			err := validateProjectSource(tt.attributes, &testProjectSrc)
			if tt.wantErr == (err == nil) {
				t.Errorf("error: %v", err)
				return
			} else if err == nil {
				assert.Equal(t, tt.expected, testProjectSrc, "The two values should be the same.")
			}
		})
	}
}
