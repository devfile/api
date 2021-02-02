package validation

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
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

	oneRemoteErr := "starterProject .* should have one remote only"
	wrongCheckoutErr := "unable to find the checkout remote .* in the remotes for project.*"
	atleastOneRemoteErr := "starterProject .* should have at least one remote"

	tests := []struct {
		name            string
		starterProjects []v1alpha2.StarterProject
		wantErr         *string
	}{
		{
			name: "Valid Starter Project",
			starterProjects: []v1alpha2.StarterProject{
				generateDummyGitStarterProject("project1", "origin", map[string]string{"origin": "originremote"}),
				generateDummyGitStarterProject("project2", "origin", map[string]string{"origin": "originremote2"}),
			},
		},
		{
			name: "Invalid Starter Project",
			starterProjects: []v1alpha2.StarterProject{
				generateDummyGithubStarterProject("project1", "origin", map[string]string{"origin": "originremote", "test": "testremote"}),
			},
			wantErr: &oneRemoteErr,
		},
		{
			name: "Invalid Starter Project with wrong checkout",
			starterProjects: []v1alpha2.StarterProject{
				generateDummyGithubStarterProject("project1", "origin", map[string]string{"test": "testremote"}),
			},
			wantErr: &wrongCheckoutErr,
		},
		{
			name: "Valid Starter Project with empty checkout remote",
			starterProjects: []v1alpha2.StarterProject{
				generateDummyGitStarterProject("project1", "", map[string]string{"origin": "originremote"}),
			},
		},
		{
			name: "Invalid Starter Project with empty remotes",
			starterProjects: []v1alpha2.StarterProject{
				generateDummyGithubStarterProject("project1", "origin", map[string]string{}),
				generateDummyGithubStarterProject("project3", "origin", map[string]string{"origin": "originremote", "test": "testremote"}),
			},
			wantErr: &atleastOneRemoteErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStarterProjects(tt.starterProjects)

			if tt.wantErr != nil && assert.Error(t, err) {
				assert.Regexp(t, *tt.wantErr, err.Error(), "Error message should match")
			} else {
				assert.NoError(t, err, "Expected error to be nil")
			}
		})
	}
}

func TestValidateProjects(t *testing.T) {

	wrongCheckoutErr := "unable to find the checkout remote .* in the remotes for project.*"
	atleastOneRemoteErr := "projects .* should have at least one remote"

	tests := []struct {
		name     string
		projects []v1alpha2.Project
		wantErr  *string
	}{
		{
			name: "Valid Project",
			projects: []v1alpha2.Project{
				generateDummyGitProject("project1", "origin", map[string]string{"origin": "originremote"}),
				generateDummyGithubProject("project2", "origin", map[string]string{"origin": "originremote"}),
			},
		},
		{
			name: "Invalid Project with multiple remote and empty checkout remote",
			projects: []v1alpha2.Project{
				generateDummyGitProject("project2", "origin", map[string]string{"origin": "originremote"}),
				generateDummyGithubProject("project1", "", map[string]string{"origin": "originremote", "test": "testremote"}),
			},
			wantErr: &wrongCheckoutErr,
		},
		{
			name: "Invalid Project with wrong checkout",
			projects: []v1alpha2.Project{
				generateDummyGithubProject("project1", "origin", map[string]string{"origin": "originremote", "test": "testremote"}),
				generateDummyGitProject("project2", "origin1", map[string]string{"origin2": "originremote2"}),
			},
			wantErr: &wrongCheckoutErr,
		},
		{
			name: "Valid Project with empty checkout remote",
			projects: []v1alpha2.Project{
				generateDummyGitProject("project1", "", map[string]string{"origin": "originremote"}),
			},
		},
		{
			name: "Invalid Project with empty remotes",
			projects: []v1alpha2.Project{
				generateDummyGitProject("project1", "origin", map[string]string{}),
				generateDummyGithubProject("project2", "origins", map[string]string{"origin": "originremote", "test": "testremote"}),
			},
			wantErr: &atleastOneRemoteErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjects(tt.projects)

			if tt.wantErr != nil && assert.Error(t, err) {
				assert.Regexp(t, *tt.wantErr, err.Error(), "Error message should match")
			} else {
				assert.NoError(t, err, "Expected error to be nil")
			}
		})
	}
}
