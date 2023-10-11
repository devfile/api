//
//
// Copyright Red Hat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
)

func TestCommandConversion_v1alpha1(t *testing.T) {
	f := fuzz.New().NilChance(fuzzNilChance).MaxDepth(100).Funcs(
		commandFuzzFunc,
		rawExtFuzzFunc,
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

func TestCommandConversionFrom_v1alpha2(t *testing.T) {

	src := &v1alpha2.Command{
		Id: "test1",
		CommandUnion: v1alpha2.CommandUnion{
			Exec: &v1alpha2.ExecCommand{
				LabeledCommand: v1alpha2.LabeledCommand{
					BaseCommand: v1alpha2.BaseCommand{
						Group: &v1alpha2.CommandGroup{
							Kind: v1alpha2.DeployCommandGroupKind,
						},
					},
				},
			},
		},
	}

	output := &Command{}

	err := convertCommandFrom_v1alpha2(src, output)
	if !assert.NoError(t, err, "Should not return error when converting from v1alpha2") {
		return
	}

	assert.Equal(t, &Command{}, output, "Conversion from v1alpha2 should be skipped for deploy kind command")
}
