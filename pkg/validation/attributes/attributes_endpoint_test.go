package attributes

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
	"github.com/stretchr/testify/assert"
)

func TestValidateEndpoint(t *testing.T) {

	tests := []struct {
		name       string
		testFile   string
		expected   []v1alpha2.Endpoint
		attributes apiAttributes.Attributes
		wantErr    bool
	}{
		{
			name:     "Good Substitution",
			testFile: "test-fixtures/components/endpoint.yaml",
			expected: []v1alpha2.Endpoint{
				{
					Name:       "endpoint1",
					TargetPort: 9999,
					Exposure:   "public",
					Protocol:   "https",
					Path:       "/FOO",
				},
			},
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"foo": "FOO",
			}, nil),
			wantErr: false,
		},
		{
			name:     "Invalid Reference",
			testFile: "test-fixtures/components/endpoint.yaml",
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"notfoo": "FOO",
			}, nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testEndpoint := v1alpha2.Endpoint{}

			readFileToStruct(t, tt.testFile, &testEndpoint)

			testEndpointArr := []v1alpha2.Endpoint{testEndpoint}

			err := validateEndpoint(tt.attributes, &testEndpointArr)
			if tt.wantErr == (err == nil) {
				t.Errorf("error: %v", err)
				return
			} else if err == nil {
				assert.Equal(t, tt.expected, testEndpointArr, "The two values should be the same.")
			}
		})
	}
}
