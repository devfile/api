package validation

import (
	"testing"
)

func TestIsInt(t *testing.T) {

	tests := []struct {
		name       string
		arg        string
		wantResult bool
	}{
		{
			name:       "Case 1: numeric string",
			arg:        "1234",
			wantResult: true,
		},
		{
			name:       "Case 2: alphanumeric string",
			arg:        "1234abc",
			wantResult: false,
		},
		{
			name:       "Case 3: string with numbers and character",
			arg:        "12_34",
			wantResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isInt(tt.arg)
			if result != tt.wantResult {
				t.Errorf("TestIsInt result: %v, wantResult: %v", result, tt.wantResult)
			}
		})
	}
}

func TestValidateURI(t *testing.T) {

	tests := []struct {
		name    string
		uri     string
		wantErr bool
	}{
		{
			name:    "Case 1: valid uri format starts with http",
			uri:     "http://devfile.yaml",
			wantErr: false,
		},
		{
			name:    "Case 2: invalid uri format starts with http",
			uri:     "http//devfile.yaml",
			wantErr: true,
		},
		{
			name:    "Case 3: invalid uri format does not start with http",
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
