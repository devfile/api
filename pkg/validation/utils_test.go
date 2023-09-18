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
)

func TestValidateURI(t *testing.T) {

	tests := []struct {
		name    string
		uri     string
		wantErr bool
	}{
		{
			name:    "Valid URI format starts with http",
			uri:     "http://devfile.yaml",
			wantErr: false,
		},
		{
			name:    "Invalid URI format starts with http",
			uri:     "http//devfile.yaml",
			wantErr: true,
		},
		{
			name:    "Valid URI format does not start with http",
			uri:     "./devfile.yaml",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURI(tt.uri)
			if err != nil && !tt.wantErr {
				t.Errorf("TestValidateURI error: %v", err)
			}
		})
	}
}
