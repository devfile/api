package v1alpha1

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
)

func TestEventsConversion_v1alpha1(t *testing.T) {
	f := fuzz.New().NilChance(fuzzNilChance)
	for i := 0; i < fuzzIterations; i++ {
		original := &Events{}
		intermediate := &v1alpha2.Events{}
		output := &Events{}
		f.Fuzz(original)
		input := original.DeepCopy()
		err := convertEventsTo_v1alpha2(input, intermediate)
		if !assert.NoError(t, err, "Should not return error when converting to v1alpha2") {
			return
		}
		err = convertEventsFrom_v1alpha2(intermediate, output)
		if !assert.NoError(t, err, "Should not return error when converting from v1alpha2") {
			return
		}
		assert.Equal(t, original, output, "Events should not be changed when converting between v1alpha1 and v1alpha2")
	}
}
