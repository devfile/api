package v1alpha1

import (
	"encoding/json"

	"github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
)

func convertProjectTo_v1alpha2(src *Project, dest *v1alpha2.Project) error {
	jsonProject, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonProject, dest)
	if err != nil {
		return err
	}
	var sparseCheckoutDir string
	switch {
	case src.Git != nil:
		sparseCheckoutDir = src.Git.SparseCheckoutDir
	case src.Github != nil:
		sparseCheckoutDir = src.Github.SparseCheckoutDir
	case src.Zip != nil:
		sparseCheckoutDir = src.Zip.SparseCheckoutDir
	}
	if sparseCheckoutDir != "" {
		dest.SparseCheckoutDirs = []string{sparseCheckoutDir}
	}
	return nil
}

func convertProjectFrom_v1alpha2(src *v1alpha2.Project, dest *Project) error {
	jsonProject, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonProject, dest)
	if err != nil {
		return err
	}
	// **Note**: These aren't technically compatible:
	// - v1alpha2 allows us to specify multiple sparse checkout dirs; v1alpha1 only supports one
	//   -> we ignore all but the first sparseCheckoutDir
	// - v1alpha2 doesn't forbid sparse checkout dir for a custom project source
	//   -> we ignore all sparseCheckoutDirs when project source is Custom
	if len(src.SparseCheckoutDirs) > 0 {
		sparseCheckoutDir := src.SparseCheckoutDirs[0]
		switch {
		case src.Git != nil:
			dest.Git.SparseCheckoutDir = sparseCheckoutDir
		case src.Github != nil:
			dest.Github.SparseCheckoutDir = sparseCheckoutDir
		case src.Zip != nil:
			dest.Zip.SparseCheckoutDir = sparseCheckoutDir
		}
	}
	return nil
}

func convertStarterProjectTo_v1alpha2(src *StarterProject, dest *v1alpha2.StarterProject) error {
	jsonProject, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonProject, dest)
	if err != nil {
		return err
	}
	// **Note**: There are API differences for starter projects between v1alpha1 and v1alpha2:
	// - ClonePath is removed from starter projects in v1alpha2; we drop it on conversion
	// - SparseCheckoutDir is removed and SubDir is added in its place. For conversion purposes, we make these fields
	//   equivalent.
	switch {
	case src.Git != nil:
		dest.SubDir = src.Git.SparseCheckoutDir
	case src.Github != nil:
		dest.SubDir = src.Github.SparseCheckoutDir
	case src.Zip != nil:
		dest.SubDir = src.Zip.SparseCheckoutDir
	}

	return nil
}

func convertStarterProjectFrom_v1alpha2(src *v1alpha2.StarterProject, dest *StarterProject) error {
	jsonProject, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonProject, dest)
	if err != nil {
		return err
	}

	if src.SubDir != "" {
		switch {
		case src.Git != nil:
			dest.Git.SparseCheckoutDir = src.SubDir
		case src.Github != nil:
			dest.Github.SparseCheckoutDir = src.SubDir
		case src.Zip != nil:
			dest.Zip.SparseCheckoutDir = src.SubDir
		}
	}

	return nil
}
