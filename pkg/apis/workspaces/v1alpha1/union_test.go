package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizingUnion_SetDiscriminator(t *testing.T) {
	original := ProjectSource{
		Git: &GitProjectSource{},
	}
	expected := ProjectSource{
		Git:        &GitProjectSource{},
		SourceType: "Git",
	}

	err := original.Normalize()
	assert.NoError(t, err)

	assert.Equal(t,
		expected,
		original,
		"The two values should be the same.")
}

func TestNormalizingUnion_CleanupOldValue(t *testing.T) {
	original := ProjectSource{
		Git:        &GitProjectSource{},
		Zip:        &ZipProjectSource{},
		SourceType: "Git",
	}
	expected := ProjectSource{
		Git:        &GitProjectSource{},
		SourceType: "Git",
	}

	err := original.Normalize()
	assert.NoError(t, err)

	assert.Equal(t,
		expected,
		original,
		"The two values should be the same.")
}

func TestSimplifyingUnion(t *testing.T) {
	original := ProjectSource{
		Git:        &GitProjectSource{},
		SourceType: "Git",
	}
	expected := ProjectSource{
		Git: &GitProjectSource{},
	}

	original.Simplify()

	assert.Equal(t,
		expected,
		original,
		"The two values should be the same.")
}
