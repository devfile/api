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

// Component that allows the developer to declare and configure a volume into his workspace
type VolumeComponent struct {
	BaseComponent `json:",inline"`
	Volume        `json:",inline"`
}

// Volume that should be mounted to a component container
type Volume struct {
	// Mandatory name that allows referencing the Volume component
	// in Container volume mounts or inside a parent
	Name string `json:"name"`

	// +optional
	// Size of the volume
	Size string `json:"size,omitempty"`

	// +optional
	// Ephemeral volumes are not stored persistently across restarts. Defaults
	// to false
	Ephemeral bool `json:"ephemeral,omitempty"`
}
