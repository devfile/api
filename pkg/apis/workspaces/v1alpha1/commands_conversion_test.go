package v1alpha1

import (
	"github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	"github.com/google/go-cmp/cmp"
	"github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommandConversion_v1alpha1(t *testing.T) {
	f := fuzz.New().NilChance(fuzzNilChance).MaxDepth(100).Funcs(
		CommandFuzzFunc,
		RawExtFuzzFunc,
	)
	for i := 0; i < fuzzIterations; i++ {
		original := &Command{}
		intermediate := &v1alpha2.Command{}
		output := &Command{}
		f.Fuzz(original)
		input := original.DeepCopy()
		err := convertCommandTo_v1alpha2(input, intermediate)
		if !assert.NoError(t, err, "Should not return error when converting to v1alpha2") {
			return
		}
		err = convertCommandFrom_v1alpha2(intermediate, output)
		if !assert.NoError(t, err, "Should not return error when converting from v1alpha2") {
			return
		}
		if !assert.True(t, cmp.Equal(original, output), "Command should not be changed when converting between v1alpha1 and v1alpha2") {
			t.Logf("Diff: \n%s\n", cmp.Diff(original, output))
		}
	}
}
