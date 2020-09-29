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

package unions

import (
	"testing"

	workspaces "github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
)

func TestNormalizingUnion_SetDiscriminator(t *testing.T) {
	original := workspaces.DevWorkspaceTemplateSpecContent{
		Projects: []workspaces.Project{
			{
				Name: "MyProject",
				ProjectSource: workspaces.ProjectSource{
					Git: &workspaces.GitProjectSource{},
				},
			},
		},
	}
	expected := workspaces.DevWorkspaceTemplateSpecContent{
		Projects: []workspaces.Project{
			{
				Name: "MyProject",
				ProjectSource: workspaces.ProjectSource{
					Git:        &workspaces.GitProjectSource{},
					SourceType: "Git",
				},
			},
		},
	}

	err := Normalize(original)
	assert.NoError(t, err)

	assert.Equal(t,
		expected,
		original,
		"The two values should be the same.")
}

func TestNormalizingUnion_CleanupOldValue(t *testing.T) {
	original := workspaces.DevWorkspaceTemplateSpecContent{
		Projects: []workspaces.Project{
			{
				Name: "MyProject",
				ProjectSource: workspaces.ProjectSource{
					Git:        &workspaces.GitProjectSource{},
					Zip:        &workspaces.ZipProjectSource{},
					SourceType: "Git",
				},
			},
		},
	}
	expected := workspaces.DevWorkspaceTemplateSpecContent{
		Projects: []workspaces.Project{
			{
				Name: "MyProject",
				ProjectSource: workspaces.ProjectSource{
					Git:        &workspaces.GitProjectSource{},
					SourceType: "Git",
				},
			},
		},
	}

	err := Normalize(original)
	assert.NoError(t, err)

	assert.Equal(t,
		expected,
		original,
		"The two values should be the same.")
}

func TestSimplifyingUnion(t *testing.T) {
	original := workspaces.DevWorkspaceTemplateSpecContent{
		Projects: []workspaces.Project{
			{
				Name: "MyProject",
				ProjectSource: workspaces.ProjectSource{
					Git:        &workspaces.GitProjectSource{},
					SourceType: "Git",
				},
			},
		},
	}
	expected := workspaces.DevWorkspaceTemplateSpecContent{
		Projects: []workspaces.Project{
			{
				Name: "MyProject",
				ProjectSource: workspaces.ProjectSource{
					Git: &workspaces.GitProjectSource{},
				},
			},
		},
	}

	Simplify(original)

	assert.Equal(t,
		expected,
		original,
		"The two values should be the same.")
}
