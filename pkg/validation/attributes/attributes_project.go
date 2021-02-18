package attributes

import (
	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
)

// ValidateProjects validates the projects data for a global attribute
func ValidateProjects(attributes apiAttributes.Attributes, projects *[]v1alpha2.Project) error {

	if projects != nil {
		for i := range *projects {
			var err error

			// Validate project clonepath
			if (*projects)[i].ClonePath, err = validateAndReplaceDataWithAttribute((*projects)[i].ClonePath, attributes); err != nil {
				return err
			}

			// Validate project sparse checkout dir
			for j := range (*projects)[i].SparseCheckoutDirs {
				if (*projects)[i].SparseCheckoutDirs[j], err = validateAndReplaceDataWithAttribute((*projects)[i].SparseCheckoutDirs[j], attributes); err != nil {
					return err
				}
			}

			// Validate project source
			if err = validateProjectSource(attributes, &(*projects)[i].ProjectSource); err != nil {
				return err
			}
		}
	}

	return nil
}

// ValidateStarterProjects validates the starter projects data for a global attribute
func ValidateStarterProjects(attributes apiAttributes.Attributes, starterProjects *[]v1alpha2.StarterProject) error {

	if starterProjects != nil {
		for i := range *starterProjects {
			var err error

			// Validate starter project description
			if (*starterProjects)[i].Description, err = validateAndReplaceDataWithAttribute((*starterProjects)[i].Description, attributes); err != nil {
				return err
			}

			// Validate starter project sub dir
			if (*starterProjects)[i].SubDir, err = validateAndReplaceDataWithAttribute((*starterProjects)[i].SubDir, attributes); err != nil {
				return err
			}

			// Validate starter project source
			if err = validateProjectSource(attributes, &(*starterProjects)[i].ProjectSource); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateProjectSource validates a project source location for a global attribute
func validateProjectSource(attributes apiAttributes.Attributes, projectSource *v1alpha2.ProjectSource) error {

	var err error

	if projectSource != nil {
		switch {
		case projectSource.Zip != nil:
			if projectSource.Zip.Location, err = validateAndReplaceDataWithAttribute(projectSource.Zip.Location, attributes); err != nil {
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
				if gitProject.CheckoutFrom.Revision, err = validateAndReplaceDataWithAttribute(gitProject.CheckoutFrom.Revision, attributes); err != nil {
					return err
				}

				// // validate git checkout remote
				if gitProject.CheckoutFrom.Remote, err = validateAndReplaceDataWithAttribute(gitProject.CheckoutFrom.Remote, attributes); err != nil {
					return err
				}
			}

			// validate git remotes
			for k := range gitProject.Remotes {
				// update map value
				if gitProject.Remotes[k], err = validateAndReplaceDataWithAttribute(gitProject.Remotes[k], attributes); err != nil {
					return err
				}

				// update map key
				var updatedKey string
				if updatedKey, err = validateAndReplaceDataWithAttribute(k, attributes); err != nil {
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
