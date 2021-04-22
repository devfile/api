package validation

import (
	"fmt"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

// ValidateStarterProjects checks if starter project has only one remote configured
// and if the checkout remote matches the renote configured
func ValidateStarterProjects(starterProjects []v1alpha2.StarterProject) error {

	var errString string
	for _, starterProject := range starterProjects {
		var gitSource v1alpha2.GitLikeProjectSource
		if starterProject.Git != nil {
			gitSource = starterProject.Git.GitLikeProjectSource
		} else if starterProject.Github != nil {
			gitSource = starterProject.Github.GitLikeProjectSource
		} else {
			continue
		}

		switch len(gitSource.Remotes) {
		case 0:
			starterProjectErr := fmt.Errorf("\nstarterProject %s should have at least one remote", starterProject.Name)
			newErr := resolveErrorMessageWithImportArrtibutes(starterProjectErr, starterProject.Attributes)
			errString += newErr.Error()
		case 1:
			if gitSource.CheckoutFrom != nil && gitSource.CheckoutFrom.Remote != "" {
				err := validateRemoteMap(gitSource.Remotes, gitSource.CheckoutFrom.Remote, starterProject.Name)
				if err != nil {
					newErr := resolveErrorMessageWithImportArrtibutes(err, starterProject.Attributes)
					errString += newErr.Error()
					errString += fmt.Sprintf("\n%s", newErr.Error())
				}
			}
		default: // len(gitSource.Remotes) >= 2
			starterProjectErr := fmt.Errorf("\nstarterProject %s should have one remote only", starterProject.Name)
			newErr := resolveErrorMessageWithImportArrtibutes(starterProjectErr, starterProject.Attributes)
			errString += newErr.Error()
		}
	}

	var err error
	if len(errString) > 0 {
		err = fmt.Errorf("error validating starter projects:%s", errString)
	}

	return err
}

// ValidateProjects checks if the project has more than one remote configured then a checkout
// remote is mandatory and if the checkout remote matches the renote configured
func ValidateProjects(projects []v1alpha2.Project) error {

	var errString string
	for _, project := range projects {
		var gitSource v1alpha2.GitLikeProjectSource
		if project.Git != nil {
			gitSource = project.Git.GitLikeProjectSource
		} else if project.Github != nil {
			gitSource = project.Github.GitLikeProjectSource
		} else {
			continue
		}

		switch len(gitSource.Remotes) {
		case 0:
			projectErr := fmt.Errorf("\nprojects %s should have at least one remote", project.Name)
			newErr := resolveErrorMessageWithImportArrtibutes(projectErr, project.Attributes)
			errString += newErr.Error()
		case 1:
			if gitSource.CheckoutFrom != nil && gitSource.CheckoutFrom.Remote != "" {
				if err := validateRemoteMap(gitSource.Remotes, gitSource.CheckoutFrom.Remote, project.Name); err != nil {
					newErr := resolveErrorMessageWithImportArrtibutes(err, project.Attributes)
					errString += fmt.Sprintf("\n%s", newErr.Error())
				}
			}
		default: // len(gitSource.Remotes) >= 2
			if gitSource.CheckoutFrom == nil || gitSource.CheckoutFrom.Remote == "" {
				projectErr := fmt.Errorf("\nproject %s has more than one remote defined, but has no checkoutfrom remote defined", project.Name)
				newErr := resolveErrorMessageWithImportArrtibutes(projectErr, project.Attributes)
				errString += newErr.Error()
				continue
			}
			if err := validateRemoteMap(gitSource.Remotes, gitSource.CheckoutFrom.Remote, project.Name); err != nil {
				newErr := resolveErrorMessageWithImportArrtibutes(err, project.Attributes)
				errString += fmt.Sprintf("\n%s", newErr.Error())
			}
		}
	}

	var err error
	if len(errString) > 0 {
		err = fmt.Errorf("error validating projects:%s", errString)
	}

	return err
}

// validateRemoteMap checks if the checkout remote is present in the project remote map
func validateRemoteMap(remotes map[string]string, checkoutRemote, projectName string) error {

	if _, ok := remotes[checkoutRemote]; !ok {
		return fmt.Errorf("unable to find the checkout remote %s in the remotes for project %s", checkoutRemote, projectName)
	}

	return nil
}
