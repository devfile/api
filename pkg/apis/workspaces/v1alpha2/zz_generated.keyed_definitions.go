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

package v1alpha2

func (keyed Component) Key() string {
	return keyed.Name
}

func (keyed Project) Key() string {
	return keyed.Name
}

func (keyed StarterProject) Key() string {
	return keyed.Name
}

func (keyed Command) Key() string {
	return keyed.Id
}

func (keyed ComponentParentOverride) Key() string {
	return keyed.Name
}

func (keyed ProjectParentOverride) Key() string {
	return keyed.Name
}

func (keyed StarterProjectParentOverride) Key() string {
	return keyed.Name
}

func (keyed CommandParentOverride) Key() string {
	return keyed.Id
}

func (keyed ComponentPluginOverrideParentOverride) Key() string {
	return keyed.Name
}

func (keyed CommandPluginOverrideParentOverride) Key() string {
	return keyed.Id
}

func (keyed ComponentPluginOverride) Key() string {
	return keyed.Name
}

func (keyed CommandPluginOverride) Key() string {
	return keyed.Id
}
