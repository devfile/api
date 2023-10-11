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
	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// Spokes for conversion have to satisfy the Convertible interface.
var _ conversion.Convertible = (*DevWorkspaceTemplate)(nil)

func (src *DevWorkspaceTemplate) ConvertTo(destRaw conversion.Hub) error {
	dest := destRaw.(*v1alpha2.DevWorkspaceTemplate)
	return convertDevWorkspaceTemplateTo_v1alpha2(src, dest)

}

func (dest *DevWorkspaceTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.DevWorkspaceTemplate)
	return convertDevWorkspaceTemplateFrom_v1alpha2(src, dest)
}
