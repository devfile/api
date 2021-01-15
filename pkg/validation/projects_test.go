package validation

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

func generateDummyGitStarterProject(name, checkoutRemote string, remotes map[string]string) v1alpha2.StarterProject {
	return v1alpha2.StarterProject{
		Name: name,
		ProjectSource: v1alpha2.ProjectSource{
			Git: &v1alpha2.GitProjectSource{
				GitLikeProjectSource: v1alpha2.GitLikeProjectSource{
					Remotes: remotes,
					CheckoutFrom: &v1alpha2.CheckoutFrom{
						Remote: checkoutRemote,
					},
				},
			},
		},
	}
}

func generateDummyGithubStarterProject(name, checkoutRemote string, remotes map[string]string) v1alpha2.StarterProject {
	return v1alpha2.StarterProject{
		Name: name,
		ProjectSource: v1alpha2.ProjectSource{
			Github: &v1alpha2.GithubProjectSource{
				GitLikeProjectSource: v1alpha2.GitLikeProjectSource{
					Remotes: remotes,
					CheckoutFrom: &v1alpha2.CheckoutFrom{
						Remote: checkoutRemote,
					},
				},
			},
		},
	}
}

func generateDummyGitProject(name, checkoutRemote string, remotes map[string]string) v1alpha2.Project {
	return v1alpha2.Project{
		Name: name,
		ProjectSource: v1alpha2.ProjectSource{
			Git: &v1alpha2.GitProjectSource{
				GitLikeProjectSource: v1alpha2.GitLikeProjectSource{
					Remotes: remotes,
					CheckoutFrom: &v1alpha2.CheckoutFrom{
						Remote: checkoutRemote,
					},
				},
			},
		},
	}
}

func generateDummyGithubProject(name, checkoutRemote string, remotes map[string]string) v1alpha2.Project {
	return v1alpha2.Project{
		Name: name,
		ProjectSource: v1alpha2.ProjectSource{
			Github: &v1alpha2.GithubProjectSource{
				GitLikeProjectSource: v1alpha2.GitLikeProjectSource{
					Remotes: remotes,
					CheckoutFrom: &v1alpha2.CheckoutFrom{
						Remote: checkoutRemote,
					},
				},
			},
		},
	}
}

func TestValidateStarterProjects(t *testing.T) {

	tests := []struct {
		name           string
		starterProject v1alpha2.StarterProject
		wantErr        bool
	}{
		{
			name:           "Case 1: Valid Starter Project",
			starterProject: generateDummyGitStarterProject("project1", "origin", map[string]string{"origin": "originremote"}),
			wantErr:        false,
		},
		{
			name:           "Case 2: Invalid Starter Project",
			starterProject: generateDummyGithubStarterProject("project1", "origin", map[string]string{"origin": "originremote", "test": "testremote"}),
			wantErr:        true,
		},
		{
			name:           "Case 3: Invalid Starter Project with wrong checkout",
			starterProject: generateDummyGithubStarterProject("project1", "origin", map[string]string{"test": "testremote"}),
			wantErr:        true,
		},
		{
			name:           "Case 4: Valid Starter Project with empty checkout remote",
			starterProject: generateDummyGitStarterProject("project1", "", map[string]string{"origin": "originremote"}),
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStarterProjects(tt.starterProject)
			if !tt.wantErr == (err != nil) {
				t.Errorf("TestValidateStarterProjects unexpected error: %v", err)
			}
		})
	}
}

func TestValidateProjects(t *testing.T) {

	tests := []struct {
		name    string
		project v1alpha2.Project
		wantErr bool
	}{
		{
			name:    "Case 1: Valid Project",
			project: generateDummyGitProject("project1", "origin", map[string]string{"origin": "originremote"}),
			wantErr: false,
		},
		{
			name:    "Case 2: Invalid Project with multiple remote and empty checkout remote",
			project: generateDummyGithubProject("project1", "", map[string]string{"origin": "originremote", "test": "testremote"}),
			wantErr: true,
		},
		{
			name:    "Case 3: Invalid Project with wrong checkout",
			project: generateDummyGithubProject("project1", "invalidorigin", map[string]string{"origin": "originremote", "test": "testremote"}),
			wantErr: true,
		},
		{
			name:    "Case 4: Valid Project with empty checkout remote",
			project: generateDummyGitProject("project1", "", map[string]string{"origin": "originremote"}),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjects(tt.project)
			if !tt.wantErr == (err != nil) {
				t.Errorf("TestValidateProjects unexpected error: %v", err)
			}
		})
	}
}
