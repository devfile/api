package validation

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

func TestValidateEndpoints(t *testing.T) {

	tests := []struct {
		name                  string
		endpoints             []v1alpha2.Endpoint
		processedEndpointName map[string]bool
		processedEndpointPort map[int]bool
		wantErr               bool
	}{
		{
			name: "Case 1: Duplicate endpoint name",
			endpoints: []v1alpha2.Endpoint{
				generateDummyEndpoint("url1", 8080),
				generateDummyEndpoint("url1", 8081),
			},
			processedEndpointName: map[string]bool{},
			processedEndpointPort: map[int]bool{},
			wantErr:               true,
		},
		{
			name: "Case 2:  Duplicate endpoint name across components",
			endpoints: []v1alpha2.Endpoint{
				generateDummyEndpoint("url1", 8080),
			},
			processedEndpointName: map[string]bool{
				"url1": true,
			},
			processedEndpointPort: map[int]bool{},
			wantErr:               true,
		},
		{
			name: "Case 3:  Duplicate endpoint port within same component",
			endpoints: []v1alpha2.Endpoint{
				generateDummyEndpoint("url1", 8080),
				generateDummyEndpoint("url2", 8080),
			},
			processedEndpointName: map[string]bool{},
			processedEndpointPort: map[int]bool{},
			wantErr:               false,
		},
		{
			name: "Case 4:  Duplicate endpoint port across components",
			endpoints: []v1alpha2.Endpoint{
				generateDummyEndpoint("url1", 8080),
				generateDummyEndpoint("url2", 8081),
			},
			processedEndpointName: map[string]bool{},
			processedEndpointPort: map[int]bool{
				8080: true,
			},
			wantErr: true,
		},
		{
			name: "Case 5: numeric endpoint name",
			endpoints: []v1alpha2.Endpoint{
				generateDummyEndpoint("123", 8080),
			},
			processedEndpointName: map[string]bool{},
			processedEndpointPort: map[int]bool{},
			wantErr:               true,
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
