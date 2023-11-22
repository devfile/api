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
	"strconv"

	schema "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

// parentComponentAdded adds a new component to the test schema data and to the parser data
func (devfile *TestDevfile) parentComponentAdded(component schema.ComponentParentOverride) {
	LogInfoMessage(fmt.Sprintf("parent component added Name: %s", component.Name))
	devfile.SchemaDevFile.Parent.Components = append(devfile.SchemaDevFile.Parent.Components, component)
}

// addParentVolume returns volumeMounts in a schema structure based on a specified number of volumes
func (devfile *TestDevfile) addParentVolume(numVols int) []schema.VolumeMountParentOverride {
	commandVols := make([]schema.VolumeMountParentOverride, numVols)
	for i := 0; i < numVols; i++ {
		volumeComponent := devfile.AddParentComponent(schema.VolumeComponentType)
		commandVols[i].Name = volumeComponent.Name
		commandVols[i].Path = "/Path_" + GetRandomString(5, false)
		LogInfoMessage(fmt.Sprintf("....... Add Volume: %s", commandVols[i]))
	}
	return commandVols
}

// AddParentComponent adds a component of the specified type, with random attributes, to the devfile schema
func (devfile *TestDevfile) AddParentComponent(componentType schema.ComponentType) schema.ComponentParentOverride {
	LogInfoMessage(fmt.Sprintf("Create a %v parent component :", componentType))
	component := schema.ComponentParentOverride{}
	component.Name = GetRandomUniqueString(8, true)
	LogInfoMessage(fmt.Sprintf("....... Name: %s", component.Name))
	switch componentType {
	case schema.ContainerComponentType:
		component.Container = &schema.ContainerComponentParentOverride{}
		devfile.SetParentContainerComponentValues(&component)
	case schema.KubernetesComponentType:
		component.Kubernetes = &schema.KubernetesComponentParentOverride{}
		devfile.SetParentK8sComponentValues(&component)
	case schema.OpenshiftComponentType:
		component.Openshift = &schema.OpenshiftComponentParentOverride{}
		devfile.SetParentK8sComponentValues(&component)
	case schema.VolumeComponentType:
		component.Volume = &schema.VolumeComponentParentOverride{}
		devfile.SetParentVolumeComponentValues(&component)
	}

	devfile.parentComponentAdded(component)
	return component
}

// GetParentContainerName returns the name of an existing, or newly created, container.
func (devfile *TestDevfile) GetParentContainerName() string {

	componentName := ""
	for _, currentComponent := range devfile.SchemaDevFile.Parent.Components {
		if currentComponent.Container != nil {
			componentName = currentComponent.Name
			LogInfoMessage(fmt.Sprintf("return existing container from GetParentContainerName  : %s", componentName))
			break
		}
	}

	if componentName == "" {
		component := devfile.AddParentComponent(schema.ContainerComponentType)
		component.Container.Image = GetRandomUniqueString(GetRandomNumber(8, 18), false)
		componentName = component.Name
		LogInfoMessage(fmt.Sprintf("return new container from GetParentContainerName : %s", componentName))
	}

	return componentName
}

// SetParentContainerComponentValues randomly sets/updates container component attributes to random values
func (devfile *TestDevfile) SetParentContainerComponentValues(component *schema.ComponentParentOverride) {

	containerComponent := component.Container.ContainerParentOverride

	containerComponent.Image = GetRandomUniqueString(GetRandomNumber(8, 18), false)

	if GetBinaryDecision() {
		numCommands := GetRandomNumber(1, 3)
		containerComponent.Command = make([]string, numCommands)
		for i := 0; i < numCommands; i++ {
			containerComponent.Command[i] = GetRandomString(GetRandomNumber(4, 16), false)
			LogInfoMessage(fmt.Sprintf("....... command %d of %d : %s", i, numCommands, containerComponent.Command[i]))
		}
	}

	if GetBinaryDecision() {
		numArgs := GetRandomNumber(1, 3)
		containerComponent.Args = make([]string, numArgs)
		for i := 0; i < numArgs; i++ {
			containerComponent.Args[i] = GetRandomString(GetRandomNumber(8, 18), false)
			LogInfoMessage(fmt.Sprintf("....... arg %d of %d : %s", i, numArgs, containerComponent.Args[i]))
		}
	}

	value := GetBinaryDecision()
	containerComponent.DedicatedPod = &value
	LogInfoMessage(fmt.Sprintf("....... DedicatedPod: %t", *(containerComponent.DedicatedPod)))

	if GetBinaryDecision() {
		containerComponent.MemoryLimit = strconv.Itoa(GetRandomNumber(4, 124)) + "M"
		LogInfoMessage(fmt.Sprintf("....... MemoryLimit: %s", containerComponent.MemoryLimit))
	}

	if GetBinaryDecision() {
		setMountSources := GetBinaryDecision()
		containerComponent.MountSources = &setMountSources
		LogInfoMessage(fmt.Sprintf("....... MountSources: %t", *containerComponent.MountSources))

		if setMountSources {
			containerComponent.SourceMapping = "/" + GetRandomString(8, false)
			LogInfoMessage(fmt.Sprintf("....... SourceMapping: %s", containerComponent.SourceMapping))
		}
	}

	if GetBinaryDecision() {
		containerComponent.Env = addParentEnv(GetRandomNumber(1, 4))
	} else {
		containerComponent.Env = nil
	}

	if len(containerComponent.VolumeMounts) == 0 {
		if GetBinaryDecision() {
			containerComponent.VolumeMounts = devfile.addParentVolume(GetRandomNumber(1, 4))
		}
	}

	if GetBinaryDecision() {
		component.Container.Endpoints = devfile.CreateParentEndpoints()
	}

	LogInfoMessage(fmt.Sprintf("component updated Name: %s", component.Name))

}

// SetParentK8sComponentValues randomly sets the required properties of a Kubernetes or Openshift component
func (devfile *TestDevfile) SetParentK8sComponentValues(component *schema.ComponentParentOverride) {
	var k8type *schema.K8sLikeComponentParentOverride = &schema.K8sLikeComponentParentOverride{}

	if component.Kubernetes != nil {
		k8type = &component.Kubernetes.K8sLikeComponentParentOverride
	} else if component.Openshift != nil {
		k8type = &component.Openshift.K8sLikeComponentParentOverride
	}

	if k8type.Inlined != "" {
		k8type.Inlined = GetRandomString(GetRandomNumber(8, 18), false)
		LogInfoMessage(fmt.Sprintf("....... updated k8type.Inlined: %s", k8type.Inlined))
	} else if k8type.Uri != "" {
		k8type.Uri = GetRandomString(GetRandomNumber(8, 18), false)
		LogInfoMessage(fmt.Sprintf("....... updated k8type.Uri: %s", k8type.Uri))
	} else {
		//This is the component creation scenario when no inlined or uri property is set
		if GetBinaryDecision() {
			k8type.Inlined = GetRandomString(GetRandomNumber(8, 18), false)
			LogInfoMessage(fmt.Sprintf("....... created Inlined: %s", k8type.Inlined))
		} else {
			k8type.Uri = GetRandomString(GetRandomNumber(8, 18), false)
			LogInfoMessage(fmt.Sprintf("....... created Uri: %s", k8type.Uri))
		}
	}

	if GetBinaryDecision() {
		k8type.Endpoints = devfile.CreateParentEndpoints()
	}

	LogInfoMessage(fmt.Sprintf("component updated Name: %s", component.Name))
}

// SetParentVolumeComponentValues randomly sets/updates volume component attributes to random values
func (devfile *TestDevfile) SetParentVolumeComponentValues(component *schema.ComponentParentOverride) {

	component.Volume.Size = strconv.Itoa(4+GetRandomNumber(64, 256)) + "G"
	LogInfoMessage(fmt.Sprintf("....... volumeComponent.Size: %s", component.Volume.Size))
	LogInfoMessage(fmt.Sprintf("component updated Name: %s", component.Name))

	value := GetBinaryDecision()
	component.Volume.Ephemeral = &value
	LogInfoMessage(fmt.Sprintf("....... volumeComponent.Ephemeral: %t", *component.Volume.Ephemeral))

}
