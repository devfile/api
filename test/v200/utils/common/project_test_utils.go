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

// projectAdded adds a new project to the test schema data and notifies the follower
func (testDevfile *TestDevfile) projectAdded(project schema.Project) {
	LogInfoMessage(fmt.Sprintf("project added Name: %s", project.Name))
	testDevfile.SchemaDevFile.Projects = append(testDevfile.SchemaDevFile.Projects, project)
	if testDevfile.Follower != nil {
		testDevfile.Follower.AddProject(project)
	}
}

// projectUpdated notifies the follower of the project which has been updated
func (testDevfile *TestDevfile) projectUpdated(project schema.Project) {
	LogInfoMessage(fmt.Sprintf("project updated Name: %s", project.Name))
	testDevfile.replaceSchemaProject(project)
	if testDevfile.Follower != nil {
		testDevfile.Follower.UpdateProject(project)
	}
}

// starterProjectAdded adds a new starter project to the test schema data and notifies the follower
func (testDevfile *TestDevfile) starterProjectAdded(starterProject schema.StarterProject) {
	LogInfoMessage(fmt.Sprintf("starter project added Name: %s", starterProject.Name))
	testDevfile.SchemaDevFile.StarterProjects = append(testDevfile.SchemaDevFile.StarterProjects, starterProject)
	if testDevfile.Follower != nil {
		testDevfile.Follower.AddStarterProject(starterProject)
	}
}

// starterProjectUpdated notifies the follower of the starter project which has been updated
func (testDevfile *TestDevfile) starterProjectUpdated(starterProject schema.StarterProject) {
	LogInfoMessage(fmt.Sprintf("starter project updated Name: %s", starterProject.Name))
	testDevfile.replaceSchemaStarterProject(starterProject)
	if testDevfile.Follower != nil {
		testDevfile.Follower.UpdateStarterProject(starterProject)
	}
}

// replaceSchemaProject replaces a Project in the saved devfile schema structure
func (testDevfile *TestDevfile) replaceSchemaProject(project schema.Project) {
	for i := 0; i < len(testDevfile.SchemaDevFile.Projects); i++ {
		if testDevfile.SchemaDevFile.Projects[i].Name == project.Name {
			testDevfile.SchemaDevFile.Projects[i] = project
			break
		}
	}
}

// replaceSchemaStarterProject replaces a Starter Project in the saved devfile schema structure
func (testDevfile *TestDevfile) replaceSchemaStarterProject(starterProject schema.StarterProject) {
	for i := 0; i < len(testDevfile.SchemaDevFile.StarterProjects); i++ {
		if testDevfile.SchemaDevFile.StarterProjects[i].Name == starterProject.Name {
			testDevfile.SchemaDevFile.StarterProjects[i] = starterProject
			break
		}
	}
}

// getRemotes creates and returns a map of remotes
func getRemotes(numRemotes int) map[string]string {
	remotes := make(map[string]string)
	for i := 0; i < numRemotes; i++ {
		key := GetRandomUniqueString(GetRandomNumber(6, 12), false)
		remotes[key] = GetRandomUniqueString(GetRandomNumber(6, 12), false)
		LogInfoMessage(fmt.Sprintf("Set remote key= %s, value= %s", key, remotes[key]))
	}
	return remotes
}

// AddProject adds a project of the specified type, with random attributes, to the devfile schema
func (testDevfile *TestDevfile) AddProject(projectType schema.ProjectSourceType) string {
	project := testDevfile.createProject(projectType)
	testDevfile.SetProjectValues(&project)
	return project.Name
}

// AddStarterProject adds a starter project of the specified type, with random attributes, to the devfile schema
func (testDevfile *TestDevfile) AddStarterProject(projectType schema.ProjectSourceType) string {
	starterProject := testDevfile.createStarterProject(projectType)
	testDevfile.SetStarterProjectValues(&starterProject)
	return starterProject.Name
}

// createProject creates a project of a specified type with only required attributes set
func (testDevfile *TestDevfile) createProject(projectType schema.ProjectSourceType) schema.Project {
	project := schema.Project{}
	project.Name = GetRandomUniqueString(GetRandomNumber(8, 63), true)
	LogInfoMessage(fmt.Sprintf("Create Project Name: %s", project.Name))

	if projectType == schema.GitProjectSourceType {
		project.Git = createGitProject(GetRandomNumber(1, 5))
	} else if projectType == schema.ZipProjectSourceType {
		project.Zip = createZipProject()
	}
	testDevfile.projectAdded(project)
	return project
}

// createStarterProject creates a starter project of a specified type with only required attributes set
func (testDevfile *TestDevfile) createStarterProject(projectType schema.ProjectSourceType) schema.StarterProject {
	starterProject := schema.StarterProject{}
	starterProject.Name = GetRandomUniqueString(GetRandomNumber(8, 63), true)
	LogInfoMessage(fmt.Sprintf("Create StarterProject Name: %s", starterProject.Name))

	if projectType == schema.GitProjectSourceType {
		starterProject.Git = createGitProject(1)
	} else if projectType == schema.ZipProjectSourceType {
		starterProject.Zip = createZipProject()
	}
	testDevfile.starterProjectAdded(starterProject)
	return starterProject

}

// createGitProject creates a git project structure with mandatory attributes set
func createGitProject(numRemotes int) *schema.GitProjectSource {
	project := schema.GitProjectSource{}
	project.Remotes = getRemotes(numRemotes)
	return &project
}

// createZipProject creates a zip project structure
func createZipProject() *schema.ZipProjectSource {
	project := schema.ZipProjectSource{}
	return &project
}

// SetProjectValues sets project attributes, common to all projects, to random values.
func (testDevfile *TestDevfile) SetProjectValues(project *schema.Project) {

	if GetBinaryDecision() {
		project.ClonePath = "./" + GetRandomString(GetRandomNumber(4, 12), false)
		LogInfoMessage(fmt.Sprintf("Set ClonePath : %s", project.ClonePath))
	}

	if project.Git != nil {
		setGitProjectValues(project.Git)
	} else if project.Zip != nil {
		setZipProjectValues(project.Zip)
	}

	testDevfile.projectUpdated(*project)
}

// SetStarterProjectValues sets starter project attributes, common to all starter projects, to random values.
func (testDevfile *TestDevfile) SetStarterProjectValues(starterProject *schema.StarterProject) {

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
		setGitProjectValues(starterProject.Git)
	} else if starterProject.Zip != nil {
		setZipProjectValues(starterProject.Zip)
	}

	testDevfile.starterProjectUpdated(*starterProject)

}

// setGitProjectValues randomly sets attributes for a Git project
func setGitProjectValues(gitProject *schema.GitProjectSource) {

	if len(gitProject.Remotes) > 1 {
		numKey := GetRandomNumber(1, len(gitProject.Remotes))
		for key, _ := range gitProject.Remotes {
			numKey--
			if numKey <= 0 {
				gitProject.CheckoutFrom = &schema.CheckoutFrom{}
				gitProject.CheckoutFrom.Remote = key
				gitProject.CheckoutFrom.Revision = GetRandomString(8, false)
				LogInfoMessage(fmt.Sprintf("set CheckoutFrom remote = %s, and revision = %s", gitProject.CheckoutFrom.Remote, gitProject.CheckoutFrom.Revision))
				break
			}
		}
	}
}

// setZipProjectValues randomly sets attributes for a Zip Project
func setZipProjectValues(zipProject *schema.ZipProjectSource) {
	zipProject.Location = GetRandomString(GetRandomNumber(8, 16), false)
}
