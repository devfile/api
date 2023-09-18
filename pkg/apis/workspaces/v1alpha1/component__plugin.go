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

type PluginComponent struct {
	BaseComponent   `json:",inline"`
	ImportReference `json:",inline"`
	PluginOverrides `json:",inline"`

	// +optional
	// Optional name that allows referencing the component
	// in commands, or inside a parent
	// If omitted it will be infered from the location (uri or registryEntry)
	Name string `json:"name,omitempty"`
}
