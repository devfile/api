package unions

import (
	"testing"

	workspaces "github.com/devfile/kubernetes-api/pkg/apis/workspaces/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestNormalizingUnion_SetDiscriminator(t *testing.T) {
	original := workspaces.DevWorkspaceTemplateSpecContent{
		Projects: []workspaces.Project {
			{
				Name: "MyProject",
				ProjectSource: workspaces.ProjectSource {
					Git: &workspaces.GitProjectSource{},
				},
			},
		},
	}

	err := Normalize(original)
	assert.Equal(t,
		nil,
		err,
		"The two values should be the same.")
	
	assert.Equal(t,
		workspaces.DevWorkspaceTemplateSpecContent{
			Projects: []workspaces.Project {
				{
					Name: "MyProject",
					ProjectSource: workspaces.ProjectSource {
						Git: &workspaces.GitProjectSource{},
						SourceType: "Git",
					},
				},
			},
		},
	original,
	"The two values should be the same.")
}

func TestNormalizingUnion_CleanupOldValue(t *testing.T) {
	original := workspaces.DevWorkspaceTemplateSpecContent{
		Projects: []workspaces.Project {
			{
				Name: "MyProject",
				ProjectSource: workspaces.ProjectSource {
					Git: &workspaces.GitProjectSource{},
					Zip: &workspaces.ZipProjectSource{},
					SourceType: "Git",
				},
			},
		},
	}

	err := Normalize(original)
	assert.Equal(t,
		nil,
		err,
		"The two values should be the same.")
	
	assert.Equal(t,
		workspaces.DevWorkspaceTemplateSpecContent{
			Projects: []workspaces.Project {
				{
					Name: "MyProject",
					ProjectSource: workspaces.ProjectSource {
						Git: &workspaces.GitProjectSource{},
						SourceType: "Git",
					},
				},
			},
		},
	original,
	"The two values should be the same.")
}

func TestSimplifyingUnion(t *testing.T) {
	original := workspaces.DevWorkspaceTemplateSpecContent{
		Projects: []workspaces.Project {
			{
				Name: "MyProject",
				ProjectSource: workspaces.ProjectSource {
					Git: &workspaces.GitProjectSource{},
					SourceType: "Git",
				},
			},
		},
	}

	Simplify(original)

	assert.Equal(t,
		workspaces.DevWorkspaceTemplateSpecContent{
			Projects: []workspaces.Project {
				{
					Name: "MyProject",
					ProjectSource: workspaces.ProjectSource {
					Git: &workspaces.GitProjectSource{},
				},
			},
		},
	},
	original,
	"The two values should be the same.")
}
