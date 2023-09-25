//
//
// Copyright Red Hat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"encoding/json"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/api/v2/pkg/attributes"
)

const (
	GitHubConversionFromAttributeValue = "GitHub"
)

func convertProjectTo_v1alpha2(src *Project, dest *v1alpha2.Project) error {
	// Convert Github type projects in v1alpha1 to Git-type projects in v1alpha2, since Github was dropped
	if src.Github != nil {
		src.Git = &GitProjectSource{
			GitLikeProjectSource: src.Github.GitLikeProjectSource,
		}
		if dest.Attributes == nil {
			dest.Attributes = attributes.Attributes{}
		}
		dest.Attributes.PutString(ConvertedFromAttribute, GitHubConversionFromAttributeValue)
	}

	jsonProject, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonProject, dest)
	if err != nil {
		return err
	}

	// Make sure we didn't modify underlying src struct
	if src.Github != nil {
		src.Git = nil
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

	// Check if a Git-type project was originally a Github-type project in v1alpha1
	if src.Git != nil && src.Attributes != nil {
		convertedFrom := src.Attributes.GetString(ConvertedFromAttribute, nil)
		if convertedFrom == GitHubConversionFromAttributeValue {
			dest.Github = &GithubProjectSource{
				GitLikeProjectSource: dest.Git.GitLikeProjectSource,
			}
			dest.Git = nil
		}
	}

	return nil
}

func convertStarterProjectTo_v1alpha2(src *StarterProject, dest *v1alpha2.StarterProject) error {
	// Convert Github type projects in v1alpha1 to Git-type projects in v1alpha2, since Github was dropped
	if src.Github != nil {
		src.Git = &GitProjectSource{
			GitLikeProjectSource: src.Github.GitLikeProjectSource,
		}
		if dest.Attributes == nil {
			dest.Attributes = attributes.Attributes{}
		}
		dest.Attributes.PutString(ConvertedFromAttribute, GitHubConversionFromAttributeValue)
	}

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
	case src.Zip != nil:
		dest.SubDir = src.Zip.SparseCheckoutDir
	}

	// Make sure we didn't modify underlying src struct
	if src.Github != nil {
		src.Git = nil
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
		case src.Zip != nil:
			dest.Zip.SparseCheckoutDir = src.SubDir
		}
	}

	// Check if a Git-type project was originally a Github-type project in v1alpha1
	if src.Git != nil && src.Attributes != nil {
		convertedFrom := src.Attributes.GetString(ConvertedFromAttribute, nil)
		if convertedFrom == GitHubConversionFromAttributeValue {
			dest.Github = &GithubProjectSource{
				GitLikeProjectSource: dest.Git.GitLikeProjectSource,
			}
			dest.Git = nil
		}
	}

	return nil
}
