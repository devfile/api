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
