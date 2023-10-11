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

package unions

import (
	"testing"

	dw "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
)

func TestNormalizingUnion_SetDiscriminator(t *testing.T) {
	original := dw.DevWorkspaceTemplateSpecContent{
		Projects: []dw.Project{
			{
				Name: "MyProject",
				ProjectSource: dw.ProjectSource{
					Git: &dw.GitProjectSource{},
				},
			},
		},
	}
	expected := dw.DevWorkspaceTemplateSpecContent{
		Projects: []dw.Project{
			{
				Name: "MyProject",
				ProjectSource: dw.ProjectSource{
					Git:        &dw.GitProjectSource{},
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
	original := dw.DevWorkspaceTemplateSpecContent{
		Projects: []dw.Project{
			{
				Name: "MyProject",
				ProjectSource: dw.ProjectSource{
					Git:        &dw.GitProjectSource{},
					Zip:        &dw.ZipProjectSource{},
					SourceType: "Git",
				},
			},
		},
	}
	expected := dw.DevWorkspaceTemplateSpecContent{
		Projects: []dw.Project{
			{
				Name: "MyProject",
				ProjectSource: dw.ProjectSource{
					Git:        &dw.GitProjectSource{},
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
	original := dw.DevWorkspaceTemplateSpecContent{
		Projects: []dw.Project{
			{
				Name: "MyProject",
				ProjectSource: dw.ProjectSource{
					Git:        &dw.GitProjectSource{},
					SourceType: "Git",
				},
			},
		},
	}
	expected := dw.DevWorkspaceTemplateSpecContent{
		Projects: []dw.Project{
			{
				Name: "MyProject",
				ProjectSource: dw.ProjectSource{
					Git: &dw.GitProjectSource{},
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
