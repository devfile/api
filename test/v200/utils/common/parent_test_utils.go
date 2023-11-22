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

// AddParent adds properties in the test schema structure and populates it with random attributes
func (devfile *TestDevfile) AddParent() schema.Parent {
	parent := schema.Parent{}
	devfile.ParentAdded(&parent)
	devfile.setParentValues(&parent)
	return parent
}

// ParentAdded adds the parent obj to the test schema
func (devfile *TestDevfile) ParentAdded(parent *schema.Parent) {
	LogInfoMessage("Parent added")
	devfile.SchemaDevFile.Parent = parent
}

// setParentValues randomly adds/modifies object properties.
func (devfile *TestDevfile) setParentValues(parent *schema.Parent) {

	//we can use generated random values for api tests but for library parser tests,  values are dependent on what's in the
	//parent devfile so we will need to use pre-existing test artifacts

	//Set the mandatory importRefTypes (oneof Id, Kubernetes, or Uri)
	switch schema.ImportReferenceType(GetRandomValue(ImportReferenceTypes).String()) {
	case schema.IdImportReferenceType:
		parent.Id = GetRandomString(8, false)
		LogInfoMessage(fmt.Sprintf("   ....... parent.Id %s", parent.Id))
		parent.RegistryUrl = "https://" + GetRandomString(8, false)
		LogInfoMessage(fmt.Sprintf("   ....... parent.RegistryUrl %s", parent.RegistryUrl))
	case schema.KubernetesImportReferenceType:
		parent.Kubernetes = &schema.KubernetesCustomResourceImportReference{}
		parent.Kubernetes.Name = GetRandomString(8, false)
		LogInfoMessage(fmt.Sprintf("   ....... parent.Kubernetes.Name %s", parent.Kubernetes.Name))
		if GetBinaryDecision() {
			parent.Kubernetes.Namespace = GetRandomString(8, false)
			LogInfoMessage(fmt.Sprintf("   ....... parent.Kubernetes.Namespace %s", parent.Kubernetes.Namespace))
		}
	case schema.UriImportReferenceType:
		parent.Uri = GetRandomString(8, false)
		LogInfoMessage(fmt.Sprintf("   ....... parent.Uri %s", parent.Uri))
	}

	numCommands := GetRandomNumber(1, maxCommands)
	for i := 0; i < numCommands; i++ {
		devfile.AddParentCommand(schema.CommandType(GetRandomValue(CommandTypes).String()))
	}

	numComponents := GetRandomNumber(1, maxComponents)
	for i := 0; i < numComponents; i++ {
		devfile.AddParentComponent(schema.ComponentType(GetRandomValue(ComponentTypes).String()))
	}

	numProjects := GetRandomNumber(1, maxProjects)
	for i := 0; i < numProjects; i++ {
		devfile.AddParentProject(schema.ProjectSourceType(GetRandomValue(ProjectSourceTypes).String()))

	}

	numStarterProjects := GetRandomNumber(1, maxStarterProjects)
	for i := 0; i < numStarterProjects; i++ {
		devfile.AddParentStarterProject(schema.ProjectSourceType(GetRandomValue(ProjectSourceTypes).String()))
	}

	LogInfoMessage("Parent updated")
}
