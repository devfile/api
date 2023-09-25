//
//
// Copyright Red Hat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package validation

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
)

func TestValidateEndpoints(t *testing.T) {

	duplicateNameErr := "multiple endpoint entries with same name"
	duplicatePortErr := "devfile contains multiple containers with same endpoint targetPort"

	tests := []struct {
		name                  string
		endpoints             []v1alpha2.Endpoint
		processedEndpointName map[string]bool
		processedEndpointPort map[int]bool
		wantErr               []string
	}{
		{
			name: "Duplicate endpoint name",
			endpoints: []v1alpha2.Endpoint{
				generateDummyEndpoint("url1", 8080),
				generateDummyEndpoint("url1", 8081),
			},
			processedEndpointName: map[string]bool{},
			processedEndpointPort: map[int]bool{},
			wantErr:               []string{duplicateNameErr},
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
			wantErr:               []string{duplicateNameErr},
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
			wantErr: []string{duplicatePortErr},
		},
		{
			name: "Multiple errors: Duplicate endpoint name, duplicate endpoint port",
			endpoints: []v1alpha2.Endpoint{
				generateDummyEndpoint("url1", 8080),
				generateDummyEndpoint("url2", 8081),
			},
			processedEndpointName: map[string]bool{
				"url1": true,
			},
			processedEndpointPort: map[int]bool{
				8080: true,
				8081: true,
			},
			wantErr: []string{duplicateNameErr, duplicatePortErr, duplicatePortErr},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEndpoints(tt.endpoints, tt.processedEndpointPort, tt.processedEndpointName)

			if tt.wantErr != nil {
				if assert.Equal(t, len(tt.wantErr), len(err), "Error list length should match") {
					for i := 0; i < len(err); i++ {
						assert.Regexp(t, tt.wantErr[i], err[i].Error(), "Error message should match")
					}
				}
			} else {
				assert.Equal(t, 0, len(err), "Error list should be empty")
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
