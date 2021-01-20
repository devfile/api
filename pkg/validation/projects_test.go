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
		name            string
		starterProjects []v1alpha2.StarterProject
		wantErr         bool
	}{
		{
			name: "Case 1: Valid Starter Project",
			starterProjects: []v1alpha2.StarterProject{
				generateDummyGitStarterProject("project1", "origin", map[string]string{"origin": "originremote"}),
				generateDummyGitStarterProject("project2", "origin", map[string]string{"origin": "originremote2"}),
			},
			wantErr: false,
		},
		{
			name: "Case 2: Invalid Starter Project",
			starterProjects: []v1alpha2.StarterProject{
				generateDummyGithubStarterProject("project1", "origin", map[string]string{"origin": "originremote", "test": "testremote"}),
				generateDummyGithubStarterProject("project2", "origin", map[string]string{"origin": "originremote", "test": "testremote"}),
			},
			wantErr: true,
		},
		{
			name: "Case 3: Invalid Starter Project with wrong checkout",
			starterProjects: []v1alpha2.StarterProject{
				generateDummyGithubStarterProject("project1", "origin", map[string]string{"test": "testremote"}),
				generateDummyGitStarterProject("project2", "origin", map[string]string{"origin": "originremote2"}),
				generateDummyGithubStarterProject("project3", "origin", map[string]string{"origin": "originremote", "test": "testremote"}),
			},
			wantErr: true,
		},
		{
			name: "Case 4: Valid Starter Project with empty checkout remote",
			starterProjects: []v1alpha2.StarterProject{
				generateDummyGitStarterProject("project1", "", map[string]string{"origin": "originremote"}),
			},
			wantErr: false,
		},
		{
			name: "Case 5: Invalid Starter Project with empty remotes",
			starterProjects: []v1alpha2.StarterProject{
				generateDummyGithubStarterProject("project1", "origin", map[string]string{}),
				generateDummyGithubStarterProject("project3", "origin", map[string]string{"origin": "originremote", "test": "testremote"}),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStarterProjects(tt.starterProjects)
			if !tt.wantErr == (err != nil) {
				t.Errorf("TestValidateStarterProjects unexpected error: %v", err)
			}
		})
	}
}

func TestValidateProjects(t *testing.T) {

	tests := []struct {
		name     string
		projects []v1alpha2.Project
		wantErr  bool
	}{
		{
			name: "Case 1: Valid Project",
			projects: []v1alpha2.Project{
				generateDummyGitProject("project1", "origin", map[string]string{"origin": "originremote"}),
				generateDummyGithubProject("project2", "origin", map[string]string{"origin": "originremote"}),
			},
			wantErr: false,
		},
		{
			name: "Case 2: Invalid Project with multiple remote and empty checkout remote",
			projects: []v1alpha2.Project{
				generateDummyGitProject("project2", "origin", map[string]string{"origin": "originremote"}),
				generateDummyGithubProject("project1", "", map[string]string{"origin": "originremote", "test": "testremote"}),
			},
			wantErr: true,
		},
		{
			name: "Case 3: Invalid Project with wrong checkout",
			projects: []v1alpha2.Project{
				generateDummyGithubProject("project1", "origin", map[string]string{"origin": "originremote", "test": "testremote"}),
				generateDummyGitProject("project2", "origin1", map[string]string{"origin2": "originremote2"}),
				generateDummyGitProject("project3", "origin3", map[string]string{"origin2": "originremote2"}),
			},
			wantErr: true,
		},
		{
			name: "Case 4: Valid Project with empty checkout remote",
			projects: []v1alpha2.Project{
				generateDummyGitProject("project1", "", map[string]string{"origin": "originremote"}),
			},
			wantErr: false,
		},
		{
			name: "Case 5: Invalid Project with empty remotes",
			projects: []v1alpha2.Project{
				generateDummyGitProject("project1", "origin", map[string]string{}),
				generateDummyGithubProject("project2", "origins", map[string]string{"origin": "originremote", "test": "testremote"}),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjects(tt.projects)
			if !tt.wantErr == (err != nil) {
				t.Errorf("TestValidateProjects unexpected error: %v", err)
			}
		})
	}
}
