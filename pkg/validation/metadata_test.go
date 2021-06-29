package validation

import (
	"github.com/devfile/api/v2/pkg/devfile"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateMetadata(t *testing.T) {

	invalidArchErr := "architecture:.* not valid. Please ensure that the architecture list conforms.*"

	tests := []struct {
		name     string
		metadata devfile.DevfileMetadata
		wantErr  *string
	}{
		{
			name: "Architecture Present - Valid",
			metadata: devfile.DevfileMetadata{
				Architectures: []string{"amd64", "s390x"},
			},
		},
		{
			name: "Architecture Present - Invalid",
			metadata: devfile.DevfileMetadata{
				Architectures: []string{"amd64", "386", "arm", "s390x"},
			},
			wantErr: &invalidArchErr,
		},
		{
			name:     "Architecture Absent",
			metadata: devfile.DevfileMetadata{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMetadata(tt.metadata)
			if tt.wantErr != nil && assert.Error(t, err) {
				assert.Regexp(t, *tt.wantErr, err.Error(), "Error message should match")
			} else {
				assert.NoError(t, err, "Expected error to be nil")
			}
		})
	}
}
