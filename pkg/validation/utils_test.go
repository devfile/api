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
