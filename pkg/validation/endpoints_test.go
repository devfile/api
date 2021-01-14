package validation

import (
	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"testing"
)

func TestValidateEndpoints(t *testing.T) {

	tests := []struct {
		name        string
		endpoints  []v1alpha2.Endpoint
		processedEndpointName map[string]bool
		processedEndpointPort map[int]bool
		wantErr     bool
	}{
		{
			name: "Case 1: Duplicate endpoint name",
			endpoints: []v1alpha2.Endpoint{
				generateDummyEndpoint("url1", 8080),
				generateDummyEndpoint("url1", 8081),
			},
			processedEndpointName: map[string]bool{},
			processedEndpointPort: map[int]bool{},
			wantErr: true,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEndpoints(tt.endpoints, tt.processedEndpointPort, tt.processedEndpointName)

			if tt.wantErr && err == nil {
				t.Errorf("TestValidateEndpoints error - expected an err but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("TestValidateEndpoints error - unexpected err %v", err)
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