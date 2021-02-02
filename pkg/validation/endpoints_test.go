package validation

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
)

func TestValidateEndpoints(t *testing.T) {

	duplicateNameErr := "multiple endpoint entries with same name"
	duplicatePortErr := "devfile contains multiple containers with same TargetPort"

	tests := []struct {
		name                  string
		endpoints             []v1alpha2.Endpoint
		processedEndpointName map[string]bool
		processedEndpointPort map[int]bool
		wantErr               *string
	}{
		{
			name: "Duplicate endpoint name",
			endpoints: []v1alpha2.Endpoint{
				generateDummyEndpoint("url1", 8080),
				generateDummyEndpoint("url1", 8081),
			},
			processedEndpointName: map[string]bool{},
			processedEndpointPort: map[int]bool{},
			wantErr:               &duplicateNameErr,
		},
		{
			name: "Duplicate endpoint name across components",
			endpoints: []v1alpha2.Endpoint{
				generateDummyEndpoint("url1", 8080),
			},
			processedEndpointName: map[string]bool{
				"url1": true,
			},
			processedEndpointPort: map[int]bool{},
			wantErr:               &duplicateNameErr,
		},
		{
			name: "Duplicate endpoint port within same component",
			endpoints: []v1alpha2.Endpoint{
				generateDummyEndpoint("url1", 8080),
				generateDummyEndpoint("url2", 8080),
			},
			processedEndpointName: map[string]bool{},
			processedEndpointPort: map[int]bool{},
		},
		{
			name: "Duplicate endpoint port across components",
			endpoints: []v1alpha2.Endpoint{
				generateDummyEndpoint("url1", 8080),
				generateDummyEndpoint("url2", 8081),
			},
			processedEndpointName: map[string]bool{},
			processedEndpointPort: map[int]bool{
				8080: true,
			},
			wantErr: &duplicatePortErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEndpoints(tt.endpoints, tt.processedEndpointPort, tt.processedEndpointName)

			if tt.wantErr != nil && assert.Error(t, err) {
				assert.Regexp(t, *tt.wantErr, err.Error(), "Error message should match")
			} else {
				assert.NoError(t, err, "Expected error to be nil")
			}
		})
	}

}

func generateDummyEndpoint(name string, port int) v1alpha2.Endpoint {
	return v1alpha2.Endpoint{
		Name:       name,
		TargetPort: port,
	}
}
