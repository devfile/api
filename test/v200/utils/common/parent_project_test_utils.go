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

package common

import (
	"fmt"

	schema "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

// parentProjectAdded adds a new project to the test schema data
func (testDevfile *TestDevfile) parentProjectAdded(project schema.ProjectParentOverride) {
	LogInfoMessage(fmt.Sprintf("Parent project added Name: %s", project.Name))
	testDevfile.SchemaDevFile.Parent.ParentOverrides.Projects = append(testDevfile.SchemaDevFile.Parent.ParentOverrides.Projects, project)
}

// parentProjectUpdated updates a project in the test schema data
func (testDevfile *TestDevfile) parentProjectUpdated(project schema.ProjectParentOverride) {
	LogInfoMessage(fmt.Sprintf("Parent project updated Name: %s", project.Name))
	testDevfile.replaceParentSchemaProject(project)
}

// parentStarterProjectAdded adds a new starter project to the test schema data
func (testDevfile *TestDevfile) parentStarterProjectAdded(starterProject schema.StarterProjectParentOverride) {
	LogInfoMessage(fmt.Sprintf("Parent starter project added Name: %s", starterProject.Name))
	testDevfile.SchemaDevFile.Parent.ParentOverrides.StarterProjects = append(testDevfile.SchemaDevFile.Parent.StarterProjects, starterProject)
}

// parentStarterProjectUpdated updates a starterproject in the test schema data
func (testDevfile *TestDevfile) parentStarterProjectUpdated(starterProject schema.StarterProjectParentOverride) {
	LogInfoMessage(fmt.Sprintf("Parent starter project updated Name: %s", starterProject.Name))
	testDevfile.replaceParentSchemaStarterProject(starterProject)
}

// replaceParentSchemaProject replaces a Project in the saved devfile schema structure
func (testDevfile *TestDevfile) replaceParentSchemaProject(project schema.ProjectParentOverride) {
	for i := 0; i < len(testDevfile.SchemaDevFile.Projects); i++ {
		if testDevfile.SchemaDevFile.Parent.Projects[i].Name == project.Name {
			testDevfile.SchemaDevFile.Parent.Projects[i] = project
			break
		}
	}
}

// replaceParentSchemaStarterProject replaces a Starter Project in the saved devfile schema structure
func (testDevfile *TestDevfile) replaceParentSchemaStarterProject(starterProject schema.StarterProjectParentOverride) {
	for i := 0; i < len(testDevfile.SchemaDevFile.StarterProjects); i++ {
		if testDevfile.SchemaDevFile.Parent.StarterProjects[i].Name == starterProject.Name {
			testDevfile.SchemaDevFile.Parent.StarterProjects[i] = starterProject
			break
		}
	}
}

// AddParentProject adds a project of the specified type, with random attributes, to the devfile schema
func (testDevfile *TestDevfile) AddParentProject(projectType schema.ProjectSourceType) string {
	project := testDevfile.createParentProject(projectType)
	testDevfile.SetParentProjectValues(&project)
	return project.Name
}

// AddParentStarterProject adds a starter project of the specified type, with random attributes, to the devfile schema
func (testDevfile *TestDevfile) AddParentStarterProject(projectType schema.ProjectSourceType) string {
	starterProject := testDevfile.createParentStarterProject(projectType)
	testDevfile.SetParentStarterProjectValues(&starterProject)
	return starterProject.Name
}

// createProject creates a project of a specified type with only required attributes set
func (testDevfile *TestDevfile) createParentProject(projectType schema.ProjectSourceType) schema.ProjectParentOverride {
	project := schema.ProjectParentOverride{}
	project.Name = GetRandomUniqueString(GetRandomNumber(8, 63), true)
	LogInfoMessage(fmt.Sprintf("Create Parent Project Name: %s", project.Name))

	if projectType == schema.GitProjectSourceType {
		project.Git = createParentGitProject(GetRandomNumber(1, 5))
	} else if projectType == schema.ZipProjectSourceType {
		project.Zip = createParentZipProject()
	}
	testDevfile.parentProjectAdded(project)
	return project
}

// createParentStarterProject creates a starter project of a specified type with only required attributes set
func (testDevfile *TestDevfile) createParentStarterProject(projectType schema.ProjectSourceType) schema.StarterProjectParentOverride {
	starterProject := schema.StarterProjectParentOverride{}
	starterProject.Name = GetRandomUniqueString(GetRandomNumber(8, 63), true)
	LogInfoMessage(fmt.Sprintf("Create Parent StarterProject Name: %s", starterProject.Name))

	if projectType == schema.GitProjectSourceType {
		//there can only be one remote for a starter project
		starterProject.Git = createParentGitProject(1)
	} else if projectType == schema.ZipProjectSourceType {
		starterProject.Zip = createParentZipProject()
	}
	testDevfile.parentStarterProjectAdded(starterProject)
	return starterProject

}

// createParentGitProject creates a git project structure with mandatory attributes set
func createParentGitProject(numRemotes int) *schema.GitProjectSourceParentOverride {
	project := schema.GitProjectSourceParentOverride{}
	project.Remotes = getRemotes(numRemotes)
	return &project
}

// createParentZipProject creates a zip project structure
func createParentZipProject() *schema.ZipProjectSourceParentOverride {
	project := schema.ZipProjectSourceParentOverride{}
	return &project
}

// SetParentProjectValues sets project attributes, common to all projects, to random values.
func (testDevfile *TestDevfile) SetParentProjectValues(project *schema.ProjectParentOverride) {

	if GetBinaryDecision() {
		project.ClonePath = "./" + GetRandomString(GetRandomNumber(4, 12), false)
		LogInfoMessage(fmt.Sprintf("Set ClonePath : %s", project.ClonePath))
	}

	if project.Git != nil {
		setParentGitProjectValues(project.Git)
	} else if project.Zip != nil {
		setParentZipProjectValues(project.Zip)
	}

	testDevfile.parentProjectUpdated(*project)
}

// SetParentStarterProjectValues sets starter project attributes, common to all starter projects, to random values.
func (testDevfile *TestDevfile) SetParentStarterProjectValues(starterProject *schema.StarterProjectParentOverride) {

	if GetBinaryDecision() {
		numWords := GetRandomNumber(2, 6)
		for i := 0; i < numWords; i++ {
			if i > 0 {
				starterProject.Description += " "
			}
			starterProject.Description += GetRandomString(8, false)
		}
		LogInfoMessage(fmt.Sprintf("Set Description : %s", starterProject.Description))
	}

	if GetBinaryDecision() {
		starterProject.SubDir = GetRandomString(12, false)
		LogInfoMessage(fmt.Sprintf("Set SubDir : %s", starterProject.SubDir))
	}

	if starterProject.Git != nil {
		setParentGitProjectValues(starterProject.Git)
	} else if starterProject.Zip != nil {
		setParentZipProjectValues(starterProject.Zip)
	}

	testDevfile.parentStarterProjectUpdated(*starterProject)

}

// setParentGitProjectValues randomly sets attributes for a Git project
func setParentGitProjectValues(gitProject *schema.GitProjectSourceParentOverride) {

	if len(gitProject.Remotes) > 1 {
		numKey := GetRandomNumber(1, len(gitProject.Remotes))
		for key, _ := range gitProject.Remotes {
			numKey--
			if numKey <= 0 {
				gitProject.CheckoutFrom = &schema.CheckoutFromParentOverride{}
				gitProject.CheckoutFrom.Remote = key
				gitProject.CheckoutFrom.Revision = GetRandomString(8, false)
				LogInfoMessage(fmt.Sprintf("set CheckoutFrom remote = %s, and revision = %s", gitProject.CheckoutFrom.Remote, gitProject.CheckoutFrom.Revision))
				break
			}
		}
	}
}

// setParentZipProjectValues randomly sets attributes for a Zip Project
func setParentZipProjectValues(zipProject *schema.ZipProjectSourceParentOverride) {
	zipProject.Location = GetRandomString(GetRandomNumber(8, 16), false)
}
