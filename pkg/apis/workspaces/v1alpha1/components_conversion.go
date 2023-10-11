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

func convertComponentTo_v1alpha2(src *Component, dest *v1alpha2.Component) error {
	if src.Plugin != nil {
		// Need to handle plugin components separately.
		return convertPluginComponentTo_v1alpha2(src, dest)
	}
	name, err := src.Key()
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
	dest.Name = name
	return nil
}

func convertComponentFrom_v1alpha2(src *v1alpha2.Component, dest *Component) error {
	if src.Plugin != nil {
		// Need to handle plugin components separately.
		return convertPluginComponentFrom_v1alpha2(src, dest)
	} else if src.Image != nil {
		// Skip converting an Image component since v1alpha1 does not have an Image component
		return nil
	}
	name := src.Key()
	jsonComponent, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonComponent, dest)
	switch {
	case dest.Container != nil:
		dest.Container.Name = name
	case dest.Plugin != nil:
		dest.Plugin.Name = name
	case dest.Volume != nil:
		dest.Volume.Name = name
	case dest.Openshift != nil:
		dest.Openshift.Name = name
	case dest.Kubernetes != nil:
		dest.Kubernetes.Name = name
	case dest.Custom != nil:
		dest.Custom.Name = name
	}
	return nil
}
