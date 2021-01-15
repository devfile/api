package validation

import (
	"fmt"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

// ValidateStarterProjects validates the starter projects
func ValidateStarterProjects(starterProject v1alpha2.StarterProject) (err error) {

	var gitSource v1alpha2.GitLikeProjectSource
	if starterProject.Git != nil {
		gitSource = starterProject.Git.GitLikeProjectSource
	} else if starterProject.Github != nil {
		gitSource = starterProject.Github.GitLikeProjectSource
	} else {
		return
	}

	if len(gitSource.Remotes) != 1 {
		return fmt.Errorf("starterProject can have only one remote")
	} else if gitSource.CheckoutFrom.Remote != "" {
		return validateRemoteMap(gitSource.Remotes, gitSource.CheckoutFrom.Remote)
	}

	return
}

// ValidateProjects validates the projects
func ValidateProjects(project v1alpha2.Project) (err error) {

	var gitSource v1alpha2.GitLikeProjectSource
	if project.Git != nil {
		gitSource = project.Git.GitLikeProjectSource
	} else if project.Github != nil {
		gitSource = project.Github.GitLikeProjectSource
	} else {
		return
	}

	if len(gitSource.Remotes) > 1 || gitSource.CheckoutFrom.Remote != "" {
		return validateRemoteMap(gitSource.Remotes, gitSource.CheckoutFrom.Remote)
	}

	return
}

// validateRemoteMap checks if the checkout remote is present in the project remote map
func validateRemoteMap(remotes map[string]string, checkoutRemote string) (err error) {

	if _, ok := remotes[checkoutRemote]; !ok {
		return fmt.Errorf("unable to find the checkout remote in project remotes map")
	}

	return
}
