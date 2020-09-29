//
// Copyright (c) 2019-2020 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//

package v1alpha2

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
