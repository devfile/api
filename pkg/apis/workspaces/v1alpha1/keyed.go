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

// +k8s:deepcopy-gen=false
type Keyed interface {
	Key() (string, error)
}

func (union Component) Key() (string, error) {
	key := ""
	err := union.Visit(ComponentVisitor{
		Container: func(container *ContainerComponent) error {
			key = container.Name
			return nil
		},
		Plugin: func(plugin *PluginComponent) error {
			if plugin.Name != "" {
				key = plugin.Name
				return nil
			}
			return plugin.ImportReference.ImportReferenceUnion.Visit(ImportReferenceUnionVisitor{
				Uri: func(uri string) error {
					key = uri
					return nil
				},
				Id: func(id string) error {
					key = plugin.Id
					if plugin.RegistryUrl != "" {
						key = plugin.RegistryUrl + "/" + key
					}
					return nil
				},
				Kubernetes: func(cr *KubernetesCustomResourceImportReference) error {
					key = cr.Name
					if cr.Namespace != "" {
						key = cr.Namespace + "/" + key
					}
					return nil
				},
			})
		},
		Kubernetes: func(k8s *KubernetesComponent) error {
			key = k8s.Name
			return nil
		},
		Openshift: func(os *OpenshiftComponent) error {
			key = os.Name
			return nil
		},
		Volume: func(vol *VolumeComponent) error {
			key = vol.Name
			return nil
		},
		Custom: func(custom *CustomComponent) error {
			key = custom.Name
			return nil
		},
	})
	return key, err
}

func (union PluginComponentsOverride) Key() (string, error) {
	key := ""
	err := union.Visit(PluginComponentsOverrideVisitor{
		Container: func(container *ContainerComponent) error {
			key = container.Name
			return nil
		},
		Kubernetes: func(k8s *KubernetesComponent) error {
			key = k8s.Name
			return nil
		},
		Openshift: func(os *OpenshiftComponent) error {
			key = os.Name
			return nil
		},
		Volume: func(vol *VolumeComponent) error {
			key = vol.Name
			return nil
		},
	})
	return key, err
}

func (keyed Command) Key() (string, error) {
	key := ""
	err := keyed.Visit(CommandVisitor{
		Apply: func(command *ApplyCommand) error {
			key = command.Id
			return nil
		},
		Exec: func(command *ExecCommand) error {
			key = command.Id
			return nil
		},
		Composite: func(command *CompositeCommand) error {
			key = command.Id
			return nil
		},
		Custom: func(command *CustomCommand) error {
			key = command.Id
			return nil
		},
		VscodeLaunch: func(command *VscodeConfigurationCommand) error {
			key = command.Id
			return nil
		},
		VscodeTask: func(command *VscodeConfigurationCommand) error {
			key = command.Id
			return nil
		},
	})
	return key, err
}
