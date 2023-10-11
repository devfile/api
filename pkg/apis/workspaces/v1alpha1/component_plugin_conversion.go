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

func convertPluginComponentTo_v1alpha2(srcComponent *Component, destComponent *v1alpha2.Component) error {
	src := srcComponent.Plugin
	if destComponent.Plugin == nil {
		destComponent.Plugin = &v1alpha2.PluginComponent{}
	}
	dest := destComponent.Plugin
	dest.Id = src.Id
	dest.RegistryUrl = src.RegistryUrl
	dest.Uri = src.Uri
	dest.ImportReferenceType = v1alpha2.ImportReferenceType(src.ImportReferenceType)
	if src.Kubernetes != nil {
		kube := v1alpha2.KubernetesCustomResourceImportReference(*src.Kubernetes)
		dest.Kubernetes = &kube
	}
	pluginKey, err := srcComponent.Key()
	if err != nil {
		return err
	}
	destComponent.Name = pluginKey

	for _, srcCommand := range src.Commands {
		srcCommand := srcCommand
		if srcCommand.Custom != nil {
			// v1alpha2 does not support Plugin Custom commands, so we have to drop them here
			continue
		}
		destCommand := v1alpha2.CommandPluginOverride{}
		err := convertPluginComponentCommandTo_v1alpha2(&srcCommand, &destCommand)
		if err != nil {
			return err
		}
		dest.Commands = append(dest.Commands, destCommand)
	}

	for _, srcComponent := range src.Components {
		srcComponent := srcComponent
		destComponent := v1alpha2.ComponentPluginOverride{}
		err := convertPluginComponentSubComponentTo_v1alpha2(&srcComponent, &destComponent)
		if err != nil {
			return err
		}
		dest.Components = append(dest.Components, destComponent)
	}
	return nil
}

func convertPluginComponentCommandTo_v1alpha2(src *Command, dest *v1alpha2.CommandPluginOverride) error {
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

func convertPluginComponentSubComponentTo_v1alpha2(src *PluginComponentsOverride, dest *v1alpha2.ComponentPluginOverride) error {
	srcName, err := src.Key()
	if err != nil {
		return err
	}
	jsonComponent, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonComponent, dest)
	if err != nil {
		return err
	}
	dest.Name = srcName
	return nil
}

func convertPluginComponentFrom_v1alpha2(srcComponent *v1alpha2.Component, destComponent *Component) error {
	src := srcComponent.Plugin
	if destComponent.Plugin == nil {
		destComponent.Plugin = &PluginComponent{}
	}
	dest := destComponent.Plugin
	dest.Id = src.Id
	dest.RegistryUrl = src.RegistryUrl
	dest.Uri = src.Uri
	dest.ImportReferenceType = ImportReferenceType(src.ImportReferenceType)
	if src.Kubernetes != nil {
		kube := KubernetesCustomResourceImportReference(*src.Kubernetes)
		dest.Kubernetes = &kube
	}
	destComponent.Plugin.Name = srcComponent.Name

	for _, srcCommand := range src.Commands {
		srcCommand := srcCommand
		destCommand := Command{}
		err := convertPluginComponentCommandFrom_v1alpha2(&srcCommand, &destCommand)
		if err != nil {
			return err
		}
		dest.Commands = append(dest.Commands, destCommand)
	}

	for _, srcComponent := range src.Components {
		srcComponent := srcComponent
		destComponent := PluginComponentsOverride{}
		err := convertPluginComponentSubComponentFrom_v1alpha2(&srcComponent, &destComponent)
		if err != nil {
			return err
		}
		dest.Components = append(dest.Components, destComponent)
	}

	return nil
}

func convertPluginComponentCommandFrom_v1alpha2(src *v1alpha2.CommandPluginOverride, dest *Command) error {
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

func convertPluginComponentSubComponentFrom_v1alpha2(src *v1alpha2.ComponentPluginOverride, dest *PluginComponentsOverride) error {
	srcName := src.Key()
	jsonComponent, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonComponent, &dest)
	if err != nil {
		return err
	}
	switch {
	case src.Container != nil:
		dest.Container.Name = srcName
	case src.Volume != nil:
		dest.Volume.Name = srcName
	case src.Openshift != nil:
		dest.Openshift.Name = srcName
	case src.Kubernetes != nil:
		dest.Kubernetes.Name = srcName
	}
	return nil
}
