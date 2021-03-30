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
