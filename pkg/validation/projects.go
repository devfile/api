package validation

import (
	"fmt"
	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

// ValidateStarterProjects checks if starter project has only one remote configured
// and if the checkout remote matches the renote configured
func ValidateStarterProjects(starterProjects []v1alpha2.StarterProject) (errList []error) {

	for _, starterProject := range starterProjects {
		var gitSource v1alpha2.GitLikeProjectSource
		if starterProject.Git != nil {
			gitSource = starterProject.Git.GitLikeProjectSource
		} else {
			continue
		}

		switch len(gitSource.Remotes) {
		case 0:
			starterProjectErr := fmt.Errorf("starterProject %s should have at least one remote", starterProject.Name)
			newErr := resolveErrorMessageWithImportAttributes(starterProjectErr, starterProject.Attributes)
			errList = append(errList, newErr)
		case 1:
			if gitSource.CheckoutFrom != nil && gitSource.CheckoutFrom.Remote != "" {
				err := validateRemoteMap(gitSource.Remotes, gitSource.CheckoutFrom.Remote, starterProject.Name)
				if err != nil {
					newErr := resolveErrorMessageWithImportAttributes(err, starterProject.Attributes)
					errList = append(errList, newErr)
				}
			}
		default: // len(gitSource.Remotes) >= 2
			starterProjectErr := fmt.Errorf("starterProject %s should have one remote only", starterProject.Name)
			newErr := resolveErrorMessageWithImportAttributes(starterProjectErr, starterProject.Attributes)
			errList = append(errList, newErr)
		}
	}

	return errList
}

// ValidateProjects checks if the project has more than one remote configured then a checkout
// remote is mandatory and if the checkout remote matches the renote configured
func ValidateProjects(projects []v1alpha2.Project) (errList []error)  {

	for _, project := range projects {
		var gitSource v1alpha2.GitLikeProjectSource
		if project.Git != nil {
			gitSource = project.Git.GitLikeProjectSource
		} else {
			continue
		}
		switch len(gitSource.Remotes) {
		case 0:
			projectErr := fmt.Errorf("projects %s should have at least one remote", project.Name)
			newErr := resolveErrorMessageWithImportAttributes(projectErr, project.Attributes)
			errList = append(errList, newErr)
		case 1:
			if gitSource.CheckoutFrom != nil && gitSource.CheckoutFrom.Remote != "" {
				if err := validateRemoteMap(gitSource.Remotes, gitSource.CheckoutFrom.Remote, project.Name); err != nil {
					newErr := resolveErrorMessageWithImportAttributes(err, project.Attributes)
					errList = append(errList, newErr)
				}
			}
		default: // len(gitSource.Remotes) >= 2
			if gitSource.CheckoutFrom == nil || gitSource.CheckoutFrom.Remote == "" {
				projectErr := fmt.Errorf("project %s has more than one remote defined, but has no checkoutfrom remote defined", project.Name)
				newErr := resolveErrorMessageWithImportAttributes(projectErr, project.Attributes)
				errList = append(errList, newErr)
				continue
			}
			if err := validateRemoteMap(gitSource.Remotes, gitSource.CheckoutFrom.Remote, project.Name); err != nil {
				newErr := resolveErrorMessageWithImportAttributes(err, project.Attributes)
				errList = append(errList, newErr)
			}
		}
	}

	return errList
}

// validateRemoteMap checks if the checkout remote is present in the project remote map
func validateRemoteMap(remotes map[string]string, checkoutRemote, projectName string) error {

	if _, ok := remotes[checkoutRemote]; !ok {
		return fmt.Errorf("unable to find the checkout remote %s in the remotes for project %s", checkoutRemote, projectName)
	}

	return nil
}
