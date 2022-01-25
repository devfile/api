package v1alpha1

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
)

func TestProjectConversion_v1alpha1(t *testing.T) {
	f := fuzz.New().NilChance(fuzzNilChance).Funcs(
		rawExtFuzzFunc,
		projectFuzzFunc,
	)
	for i := 0; i < fuzzIterations; i++ {
		original := &Project{}
		intermediate := &v1alpha2.Project{}
		output := &Project{}
		f.Fuzz(original)
		input := original.DeepCopy()
		err := convertProjectTo_v1alpha2(input, intermediate)
		if !assert.NoError(t, err, "Should not return error when converting to v1alpha2") {
			return
		}
		err = convertProjectFrom_v1alpha2(intermediate, output)
		if !assert.NoError(t, err, "Should not return error when converting from v1alpha2") {
			return
		}
		assert.Equal(t, original, output, "Projects should not be changed when converting between v1alpha1 and v1alpha2")
	}
}

func TestStarterProjectConversion_v1alpha1(t *testing.T) {
	f := fuzz.New().NilChance(fuzzNilChance).Funcs(
		rawExtFuzzFunc,
		projectFuzzFunc,
	)
	for i := 0; i < fuzzIterations; i++ {
		original := &StarterProject{}
		intermediate := &v1alpha2.StarterProject{}
		output := &StarterProject{}
		f.Fuzz(original)
		input := original.DeepCopy()
		err := convertStarterProjectTo_v1alpha2(input, intermediate)
		if !assert.NoError(t, err, "Should not return error when converting to v1alpha2") {
			return
		}
		err = convertStarterProjectFrom_v1alpha2(intermediate, output)
		if !assert.NoError(t, err, "Should not return error when converting from v1alpha2") {
			return
		}
		assert.Equal(t, original, output, "StarterProjects should not be changed when converting between v1alpha1 and v1alpha2")
	}
}
