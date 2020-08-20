package v1alpha1

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalizingUnion_SetDiscriminator(t *testing.T) {
	original := ProjectSource{
		Git: &GitProjectSource{},
	}

	err := original.Normalize()
	assert.Equal(t,
		nil,
		err,
		"The two values should be the same.")

	assert.Equal(t,
		ProjectSource{
			Git:        &GitProjectSource{},
			SourceType: "Git",
		},
		original,
		"The two values should be the same.")
}

func TestNormalizingUnion_CleanupOldValue(t *testing.T) {
	original := ProjectSource{
		Git:        &GitProjectSource{},
		Zip:        &ZipProjectSource{},
		SourceType: "Git",
	}

	err := original.Normalize()
	assert.Equal(t,
		nil,
		err,
		"The two values should be the same.")

	assert.Equal(t,
		ProjectSource{
			Git:        &GitProjectSource{},
			SourceType: "Git",
		},
		original,
		"The two values should be the same.")
}

func TestSimplifyingUnion(t *testing.T) {
	original := ProjectSource{
		Git:        &GitProjectSource{},
		SourceType: "Git",
	}

	original.Simplify()

	assert.Equal(t,
		ProjectSource{
			Git: &GitProjectSource{},
		},
		original,
		"The two values should be the same.")
}
