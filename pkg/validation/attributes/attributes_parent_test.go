package attributes

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
	"github.com/stretchr/testify/assert"
)

func TestValidateParent(t *testing.T) {

	tests := []struct {
		name       string
		testFile   string
		expected   v1alpha2.Parent
		attributes apiAttributes.Attributes
		wantErr    bool
	}{
		{
			name:     "Good Uri Substitution",
			testFile: "test-fixtures/parent/parent-uri.yaml",
			expected: v1alpha2.Parent{
				ImportReference: v1alpha2.ImportReference{
					ImportReferenceUnion: v1alpha2.ImportReferenceUnion{
						Uri: "FOO",
					},
					RegistryUrl: "FOO",
				},
			},
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"foo": "FOO",
			}, nil),
			wantErr: false,
		},
		{
			name:     "Good Id Substitution",
			testFile: "test-fixtures/parent/parent-id.yaml",
			expected: v1alpha2.Parent{
				ImportReference: v1alpha2.ImportReference{
					ImportReferenceUnion: v1alpha2.ImportReferenceUnion{
						Id: "FOO",
					},
					RegistryUrl: "FOO",
				},
			},
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"foo": "FOO",
			}, nil),
			wantErr: false,
		},
		{
			name:     "Good Kube Substitution",
			testFile: "test-fixtures/parent/parent-kubernetes.yaml",
			expected: v1alpha2.Parent{
				ImportReference: v1alpha2.ImportReference{
					ImportReferenceUnion: v1alpha2.ImportReferenceUnion{
						Kubernetes: &v1alpha2.KubernetesCustomResourceImportReference{
							Name:      "FOO",
							Namespace: "FOO",
						},
					},
					RegistryUrl: "FOO",
				},
			},
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"foo": "FOO",
			}, nil),
			wantErr: false,
		},
		{
			name:     "Invalid Reference",
			testFile: "test-fixtures/parent/parent-id.yaml",
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"bar": "BAR",
			}, nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testParent := v1alpha2.Parent{}

			readFileToStruct(t, tt.testFile, &testParent)

			err := ValidateParent(tt.attributes, &testParent)
			if tt.wantErr == (err == nil) {
				t.Errorf("error: %v", err)
				return
			} else if err == nil {
				assert.Equal(t, tt.expected, testParent, "The two values should be the same.")
			}
		})
	}
}
