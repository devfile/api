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
)

func convertParentTo_v1alpha2(src *Parent, dest *v1alpha2.Parent) error {
	dest.Id = src.Id
	dest.Uri = src.Uri
	dest.ImportReferenceType = v1alpha2.ImportReferenceType(src.ImportReferenceType)
	dest.RegistryUrl = src.RegistryUrl
	if src.Kubernetes != nil {
		kube := v1alpha2.KubernetesCustomResourceImportReference(*src.Kubernetes)
		dest.Kubernetes = &kube
	}

	for _, srcCommand := range src.Commands {
		srcCommand := srcCommand
		if srcCommand.Custom != nil {
			// v1alpha2 does not support Parent Custom commands, so we have to drop them here
			continue
		}
		destCommand := v1alpha2.CommandParentOverride{}
		err := convertParentCommandTo_v1alpha2(&srcCommand, &destCommand)
		if err != nil {
			return err
		}
		dest.Commands = append(dest.Commands, destCommand)
	}

	for _, srcComponent := range src.Components {
		srcComponent := srcComponent
		if srcComponent.Custom != nil {
			// v1alpha2 does not support Parent Custom Components, so we have to drop them here
			continue
		}
		destComponent := v1alpha2.ComponentParentOverride{}
		err := convertParentComponentTo_v1alpha2(&srcComponent, &destComponent)
		if err != nil {
			return err
		}
		dest.Components = append(dest.Components, destComponent)
	}

	for _, srcProject := range src.Projects {
		srcProject := srcProject
		destProject := v1alpha2.Project{}
		err := convertProjectTo_v1alpha2(&srcProject, &destProject)
		if err != nil {
			return err
		}
		destParentProject := v1alpha2.ProjectParentOverride{}
		jsonProject, err := json.Marshal(destProject)
		if err != nil {
			return err
		}
		err = json.Unmarshal(jsonProject, &destParentProject)
		if err != nil {
			return err
		}
		dest.Projects = append(dest.Projects, destParentProject)
	}

	for _, srcProject := range src.StarterProjects {
		srcProject := srcProject
		destProject := v1alpha2.StarterProject{}
		err := convertStarterProjectTo_v1alpha2(&srcProject, &destProject)
		if err != nil {
			return err
		}
		destParentProject := v1alpha2.StarterProjectParentOverride{}
		jsonProject, err := json.Marshal(destProject)
		if err != nil {
			return err
		}
		err = json.Unmarshal(jsonProject, &destParentProject)
		if err != nil {
			return err
		}
		dest.StarterProjects = append(dest.StarterProjects, destParentProject)
	}

	return nil
}

func convertParentCommandTo_v1alpha2(src *Command, dest *v1alpha2.CommandParentOverride) error {
	srcId, err := src.Key()
	if err != nil {
		return err
	}
	jsonCommand, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonCommand, &dest)
	if err != nil {
		return err
	}
	dest.Id = srcId
	return nil
}

func convertParentComponentTo_v1alpha2(src *Component, dest *v1alpha2.ComponentParentOverride) error {
	srcName, err := src.Key()
	if err != nil {
		return err
	}

	if src.Plugin != nil {
		destPluginComponent := &v1alpha2.PluginComponentParentOverride{}
		pluginComponent := v1alpha2.Component{}
		err := convertPluginComponentTo_v1alpha2(src, &pluginComponent)
		if err != nil {
			return err
		}
		// Though identical in json representation, we can't assign between PluginComponentParentOverride and Plugin
		jsonPlugin, err := json.Marshal(pluginComponent)
		if err != nil {
			return err
		}
		err = json.Unmarshal(jsonPlugin, destPluginComponent)
		if err != nil {
			return err
		}
		dest.Plugin = destPluginComponent
	} else {
		jsonComponent, err := json.Marshal(src)
		if err != nil {
			return err
		}
		err = json.Unmarshal(jsonComponent, &dest)
		if err != nil {
			return err
		}
	}

	dest.Name = srcName
	return nil
}

func convertParentFrom_v1alpha2(src *v1alpha2.Parent, dest *Parent) error {
	dest.Id = src.Id
	dest.Uri = src.Uri
	dest.ImportReferenceType = ImportReferenceType(src.ImportReferenceType)
	dest.RegistryUrl = src.RegistryUrl
	if src.Kubernetes != nil {
		kube := KubernetesCustomResourceImportReference(*src.Kubernetes)
		dest.Kubernetes = &kube
	}
	for _, srcCommand := range src.Commands {
		srcCommand := srcCommand
		destCommand := Command{}
		err := convertParentCommandFrom_v1alpha2(&srcCommand, &destCommand)
		if err != nil {
			return err
		}
		dest.Commands = append(dest.Commands, destCommand)
	}

	for _, srcComponent := range src.Components {
		srcComponent := srcComponent
		destComponent := Component{}
		err := convertParentComponentFrom_v1alpha2(&srcComponent, &destComponent)
		if err != nil {
			return err
		}
		dest.Components = append(dest.Components, destComponent)
	}

	for _, srcParentProject := range src.Projects {
		destProject := Project{}
		srcProject := v1alpha2.Project{}
		jsonProject, err := json.Marshal(srcParentProject)
		if err != nil {
			return err
		}
		err = json.Unmarshal(jsonProject, &srcProject)
		if err != nil {
			return err
		}
		err = convertProjectFrom_v1alpha2(&srcProject, &destProject)
		if err != nil {
			return err
		}
		dest.Projects = append(dest.Projects, destProject)
	}

	for _, srcParentProject := range src.StarterProjects {
		destProject := StarterProject{}
		srcProject := v1alpha2.StarterProject{}
		jsonProject, err := json.Marshal(srcParentProject)
		if err != nil {
			return err
		}
		err = json.Unmarshal(jsonProject, &srcProject)
		if err != nil {
			return err
		}
		err = convertStarterProjectFrom_v1alpha2(&srcProject, &destProject)
		if err != nil {
			return err
		}
		dest.StarterProjects = append(dest.StarterProjects, destProject)
	}

	return nil
}

func convertParentCommandFrom_v1alpha2(src *v1alpha2.CommandParentOverride, dest *Command) error {
	srcId := src.Key()
	jsonCommand, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonCommand, &dest)
	if err != nil {
		return err
	}
	switch {
	case src.Apply != nil:
		dest.Apply.Id = srcId
	case src.Composite != nil:
		dest.Composite.Id = srcId
	case src.Exec != nil:
		dest.Exec.Id = srcId
	}
	return nil
}

func convertParentComponentFrom_v1alpha2(src *v1alpha2.ComponentParentOverride, dest *Component) error {
	srcName := src.Key()

	if src.Plugin != nil {
		// If the parent component is a Plugin, we need to first convert it to v1alpha2.Component, then to a v1alpha1.Component
		// Through the json representation of v1alpha2 Plugin and PluginComponentParentOverride is identical they're not assignable in go
		// so we convert with a json intermediary
		srcPluginComponent := &v1alpha2.PluginComponent{}
		v1alpha2Component := v1alpha2.Component{}
		jsonPlugin, err := json.Marshal(src.Plugin)
		if err != nil {
			return err
		}
		err = json.Unmarshal(jsonPlugin, srcPluginComponent)
		if err != nil {
			return err
		}
		v1alpha2Component.Plugin = srcPluginComponent

		err = convertPluginComponentFrom_v1alpha2(&v1alpha2Component, dest)
		if err != nil {
			return err
		}
	} else {
		jsonComponent, err := json.Marshal(src)
		if err != nil {
			return err
		}
		err = json.Unmarshal(jsonComponent, &dest)
		if err != nil {
			return err
		}
	}

	switch {
	case src.Container != nil:
		dest.Container.Name = srcName
	case src.Plugin != nil:
		dest.Plugin.Name = srcName
	case src.Volume != nil:
		dest.Volume.Name = srcName
	case src.Openshift != nil:
		dest.Openshift.Name = srcName
	case src.Kubernetes != nil:
		dest.Kubernetes.Name = srcName
	}
	return nil
}
