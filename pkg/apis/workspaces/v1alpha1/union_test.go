package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizingUnion_SetDiscriminator(t *testing.T) {
	original := ProjectSource {
		Git: &GitProjectSource{},
	}

	err := original.Normalize()
	assert.Equal(t,
		nil,
		err,
		"The two values should be the same.")
	
	assert.Equal(t,
	ProjectSource {
		Git: &GitProjectSource{},
		SourceType: "Git",
	},
	original,
	"The two values should be the same.")
}

func TestNormalizingUnion_CleanupOldValue(t *testing.T) {
	original := ProjectSource {
		Git: &GitProjectSource{},
		Zip: &ZipProjectSource{},
		SourceType: "Git",
	}

	err := original.Normalize()
	assert.Equal(t,
		nil,
		err,
		"The two values should be the same.")
	
	assert.Equal(t,
	ProjectSource {
		Git: &GitProjectSource{},
		SourceType: "Git",
	},
	original,
	"The two values should be the same.")
}
