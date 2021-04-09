package variables

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckForInvalidError(t *testing.T) {

	tests := []struct {
		name            string
		wantInvalidKeys map[string]bool
		err             error
	}{
		{
			name:            "No error",
			wantInvalidKeys: make(map[string]bool),
			err:             nil,
		},
		{
			name:            "Different error",
			wantInvalidKeys: make(map[string]bool),
			err:             fmt.Errorf("an error"),
		},
		{
			name: "InvalidKeysError error",
			wantInvalidKeys: map[string]bool{
				"key1": true,
				"key2": true,
			},
			err: &InvalidKeysError{Keys: []string{"key1", "key2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testMap := make(map[string]bool)
			checkForInvalidError(testMap, tt.err)
			assert.Equal(t, tt.wantInvalidKeys, testMap, "the key map should be the same")
		})
	}
}

func TestProcessInvalidKeys(t *testing.T) {

	tests := []struct {
		name        string
		invalidKeys map[string]bool
		wantErr     bool
	}{
		{
			name:        "No invalid keys",
			invalidKeys: make(map[string]bool),
			wantErr:     false,
		},
		{
			name: "InvalidKeysError error",
			invalidKeys: map[string]bool{
				"key1": true,
				"key2": true,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := processInvalidKeys(tt.invalidKeys)
			if tt.wantErr && err == nil {
				t.Errorf("Expected error from test but got nil")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Got unexpected error: %s", err)
			}
		})
	}
}
