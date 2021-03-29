package variables

import (
	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

// ValidateAndReplaceForProjects validates the projects data for global variable references and replaces them with the variable value
func ValidateAndReplaceForProjects(variables map[string]string, projects []v1alpha2.Project) error {

	for i := range projects {
		var err error

		// Validate project clonepath
		if projects[i].ClonePath, err = validateAndReplaceDataWithVariable(projects[i].ClonePath, variables); err != nil {
			return err
		}

		// Validate project sparse checkout dir
		for j := range projects[i].SparseCheckoutDirs {
			if projects[i].SparseCheckoutDirs[j], err = validateAndReplaceDataWithVariable(projects[i].SparseCheckoutDirs[j], variables); err != nil {
				return err
			}
		}

		// Validate project source
		if err = validateandReplaceForProjectSource(variables, &projects[i].ProjectSource); err != nil {
			return err
		}
	}

	return nil
}

// ValidateAndReplaceForStarterProjects validates the starter projects data for global variable references and replaces them with the variable value
func ValidateAndReplaceForStarterProjects(variables map[string]string, starterProjects []v1alpha2.StarterProject) error {

	for i := range starterProjects {
		var err error

		// Validate starter project description
		if starterProjects[i].Description, err = validateAndReplaceDataWithVariable(starterProjects[i].Description, variables); err != nil {
			return err
		}

		// Validate starter project sub dir
		if starterProjects[i].SubDir, err = validateAndReplaceDataWithVariable(starterProjects[i].SubDir, variables); err != nil {
			return err
		}

		// Validate starter project source
		if err = validateandReplaceForProjectSource(variables, &starterProjects[i].ProjectSource); err != nil {
			return err
		}
	}

	return nil
}

// validateandReplaceForProjectSource validates a project source location for global variable references and replaces them with the variable value
func validateandReplaceForProjectSource(variables map[string]string, projectSource *v1alpha2.ProjectSource) error {

	var err error

	if projectSource != nil {
		switch {
		case projectSource.Zip != nil:
			if projectSource.Zip.Location, err = validateAndReplaceDataWithVariable(projectSource.Zip.Location, variables); err != nil {
				return err
			}
		case projectSource.Git != nil || projectSource.Github != nil:
			var gitProject *v1alpha2.GitLikeProjectSource
			if projectSource.Git != nil {
				gitProject = &projectSource.Git.GitLikeProjectSource
			} else if projectSource.Github != nil {
				gitProject = &projectSource.Github.GitLikeProjectSource
			}

			if gitProject.CheckoutFrom != nil {
				// validate git checkout revision
				if gitProject.CheckoutFrom.Revision, err = validateAndReplaceDataWithVariable(gitProject.CheckoutFrom.Revision, variables); err != nil {
					return err
				}

				// // validate git checkout remote
				if gitProject.CheckoutFrom.Remote, err = validateAndReplaceDataWithVariable(gitProject.CheckoutFrom.Remote, variables); err != nil {
					return err
				}
			}

			// validate git remotes
			for k := range gitProject.Remotes {
				// update map value
				if gitProject.Remotes[k], err = validateAndReplaceDataWithVariable(gitProject.Remotes[k], variables); err != nil {
					return err
				}

				// update map key
				var updatedKey string
				if updatedKey, err = validateAndReplaceDataWithVariable(k, variables); err != nil {
					return err
				} else if updatedKey != k {
					gitProject.Remotes[updatedKey] = gitProject.Remotes[k]
					delete(gitProject.Remotes, k)
				}
			}
		}
	}

	return nil
}
