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
			name:       "Numeric string",
			arg:        "1234",
			wantResult: true,
		},
		{
			name:       "Alphanumeric string",
			arg:        "pp1234abc-1223",
			wantResult: false,
		},
		{
			name:       "String with numbers and character",
			arg:        "12-34",
			wantResult: false,
		},
		{
			name:       "Hexadecimal string",
			arg:        "0xff",
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
