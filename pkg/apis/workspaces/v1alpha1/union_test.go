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
