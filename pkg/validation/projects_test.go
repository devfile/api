package validation

import (
	"testing"

	"github.com/devfile/api/v2/pkg/attributes"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
)

func generateDummyGitStarterProject(name string, checkoutRemote *v1alpha2.CheckoutFrom, remotes map[string]string, projectAttribute attributes.Attributes) v1alpha2.StarterProject {
	return v1alpha2.StarterProject{
		Attributes: projectAttribute,
		Name:       name,
		ProjectSource: v1alpha2.ProjectSource{
			Git: &v1alpha2.GitProjectSource{
				GitLikeProjectSource: v1alpha2.GitLikeProjectSource{
					Remotes:      remotes,
					CheckoutFrom: checkoutRemote,
				},
			},
		},
	}
}

func generateDummyGitProject(name string, checkoutRemote *v1alpha2.CheckoutFrom, remotes map[string]string, projectAttribute attributes.Attributes) v1alpha2.Project {
	return v1alpha2.Project{
		Attributes: projectAttribute,
		Name:       name,
		ProjectSource: v1alpha2.ProjectSource{
			Git: &v1alpha2.GitProjectSource{
				GitLikeProjectSource: v1alpha2.GitLikeProjectSource{
					Remotes:      remotes,
					CheckoutFrom: checkoutRemote,
				},
			},
		},
	}
}

func TestValidateStarterProjects(t *testing.T) {

	oneRemoteErr := "starterProject .* should have one remote only"
	wrongCheckoutErr := "unable to find the checkout remote .* in the remotes for project.*"
	atleastOneRemoteErr := "starterProject .* should have at least one remote"

	parentOverridesFromMainDevfile := attributes.Attributes{}.PutString(ImportSourceAttribute,
		"uri: http://127.0.0.1:8080").PutString(ParentOverrideAttribute, "main devfile")
	wrongCheckoutErrWithImportAttributes := "unable to find the checkout remote .* in the remotes for project.*, imported from uri: http://127.0.0.1:8080, in parent overrides from main devfile"

	tests := []struct {
		name            string
		starterProjects []v1alpha2.StarterProject
		wantErr         []string
	}{
		{
			name: "Valid Starter Project",
			starterProjects: []v1alpha2.StarterProject{
				generateDummyGitStarterProject("project1", &v1alpha2.CheckoutFrom{Remote: "origin"}, map[string]string{"origin": "originremote"}, attributes.Attributes{}),
				generateDummyGitStarterProject("project2", &v1alpha2.CheckoutFrom{Remote: "origin"}, map[string]string{"origin": "originremote2"}, attributes.Attributes{}),
			},
		},
		{
			name: "Invalid Starter Project",
			starterProjects: []v1alpha2.StarterProject{
				generateDummyGitStarterProject("project1", &v1alpha2.CheckoutFrom{Remote: "origin"}, map[string]string{"origin": "originremote", "test": "testremote"}, attributes.Attributes{}),
			},
			wantErr: []string{oneRemoteErr},
		},
		{
			name: "Invalid Starter Project with wrong checkout",
			starterProjects: []v1alpha2.StarterProject{
				generateDummyGitStarterProject("project1", &v1alpha2.CheckoutFrom{Remote: "origin"}, map[string]string{"test": "testremote"}, attributes.Attributes{}),
			},
			wantErr: []string{wrongCheckoutErr},
		},
		{
			name: "Valid Starter Project with empty checkout remote",
			starterProjects: []v1alpha2.StarterProject{
				generateDummyGitStarterProject("project1", &v1alpha2.CheckoutFrom{Remote: ""}, map[string]string{"origin": "originremote"}, attributes.Attributes{}),
			},
		},
		{
			name: "Valid Starter Project with no checkout remote",
			starterProjects: []v1alpha2.StarterProject{
				generateDummyGitStarterProject("project1", nil, map[string]string{"origin": "originremote"}, attributes.Attributes{}),
			},
		},
		{
			name: "Multiple errors: Starter Project with empty remotes, Starter Project with multiple remotes",
			starterProjects: []v1alpha2.StarterProject{
				generateDummyGitStarterProject("project1", &v1alpha2.CheckoutFrom{Remote: "origin"}, map[string]string{}, attributes.Attributes{}),
				generateDummyGitStarterProject("project3", &v1alpha2.CheckoutFrom{Remote: "origin"}, map[string]string{"origin": "originremote", "test": "testremote"}, attributes.Attributes{}),
			},
			wantErr: []string{atleastOneRemoteErr, oneRemoteErr},
		},
		{
			name: "Invalid Starter Project due to wrong checkout with import source attributes",
			starterProjects: []v1alpha2.StarterProject{
				generateDummyGitStarterProject("project1", &v1alpha2.CheckoutFrom{Remote: "origin"}, map[string]string{"test": "testremote"}, parentOverridesFromMainDevfile),
			},
			wantErr: []string{wrongCheckoutErrWithImportAttributes},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStarterProjects(tt.starterProjects)

			if tt.wantErr != nil {
				assert.Equal(t, len(tt.wantErr), len(err), "Error list length should match")
				for i := 0; i < len(err); i++ {
					assert.Regexp(t, tt.wantErr[i], err[i].Error(), "Error message should match")
				}
			} else {
				assert.Equal(t, 0, len(err), "Error list should be empty")
			}
		})
	}
}

func TestValidateProjects(t *testing.T) {

	wrongCheckoutErr := "unable to find the checkout remote .* in the remotes for project.*"
	atleastOneRemoteErr := "projects .* should have at least one remote"
	missingCheckOutFromRemoteErr := "project .* has more than one remote defined, but has no checkoutfrom remote defined"

	parentOverridesFromMainDevfile := attributes.Attributes{}.PutString(ImportSourceAttribute,
		"uri: http://127.0.0.1:8080").PutString(ParentOverrideAttribute, "main devfile")
	wrongCheckoutErrWithImportAttributes := "unable to find the checkout remote .* in the remotes for project.*, imported from uri: http://127.0.0.1:8080, in parent overrides from main devfile"

	tests := []struct {
		name     string
		projects []v1alpha2.Project
		wantErr  []string
	}{
		{
			name: "Valid Project",
			projects: []v1alpha2.Project{
				generateDummyGitProject("project1", &v1alpha2.CheckoutFrom{Remote: "origin"}, map[string]string{"origin": "originremote"}, attributes.Attributes{}),
				generateDummyGitProject("project2", &v1alpha2.CheckoutFrom{Remote: "origin"}, map[string]string{"origin": "originremote"}, attributes.Attributes{}),
			},
		},
		{
			name: "Invalid Project with multiple remotes but no checkoutfrom",
			projects: []v1alpha2.Project{
				generateDummyGitProject("project2", nil, map[string]string{"origin": "originremote", "test": "testremote"}, attributes.Attributes{}),
			},
			wantErr: []string{missingCheckOutFromRemoteErr},
		},
		{
			name: "Invalid Project with multiple remote and empty checkout remote",
			projects: []v1alpha2.Project{
				generateDummyGitProject("project2", &v1alpha2.CheckoutFrom{Remote: "origin"}, map[string]string{"origin": "originremote"}, attributes.Attributes{}),
				generateDummyGitProject("project1", &v1alpha2.CheckoutFrom{Remote: ""}, map[string]string{"origin": "originremote", "test": "testremote"}, attributes.Attributes{}),
			},
			wantErr: []string{missingCheckOutFromRemoteErr},
		},
		{
			name: "Invalid Project with wrong checkout",
			projects: []v1alpha2.Project{
				generateDummyGitProject("project1", &v1alpha2.CheckoutFrom{Remote: "origin1"}, map[string]string{"origin": "originremote", "test": "testremote"}, attributes.Attributes{}),
			},
			wantErr: []string{wrongCheckoutErr},
		},
		{
			name: "Valid Project with empty checkout remote",
			projects: []v1alpha2.Project{
				generateDummyGitProject("project1", &v1alpha2.CheckoutFrom{Remote: ""}, map[string]string{"origin": "originremote"}, attributes.Attributes{}),
			},
		},
		{
			name: "Multiple errors: invalid Project with empty remotes, invalid Project with wrong checkout",
			projects: []v1alpha2.Project{
				generateDummyGitProject("project1", &v1alpha2.CheckoutFrom{Remote: "origin"}, map[string]string{}, attributes.Attributes{}),
				generateDummyGitProject("project2", &v1alpha2.CheckoutFrom{Remote: "origins"}, map[string]string{"origin": "originremote", "test": "testremote"}, attributes.Attributes{}),
			},
			wantErr: []string{atleastOneRemoteErr, wrongCheckoutErr},
		},
		{
			name: "Invalid Project due to wrong checkout with import source attributes",
			projects: []v1alpha2.Project{
				generateDummyGitProject("project1", &v1alpha2.CheckoutFrom{Remote: "origin"}, map[string]string{"test": "testremote"}, parentOverridesFromMainDevfile),
			},
			wantErr: []string{wrongCheckoutErrWithImportAttributes},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjects(tt.projects)

			if tt.wantErr != nil {
				assert.Equal(t, len(tt.wantErr), len(err), "Error list length should match")
				for i := 0; i < len(err); i++ {
					assert.Regexp(t, tt.wantErr[i], err[i].Error(), "Error message should match")
				}
			} else {
				assert.Equal(t, 0, len(err), "Error list should be empty")
			}
		})
	}
}
