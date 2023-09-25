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

// ImportReferenceType describes the type of location
// from where the referenced template structure should be retrieved.
// Only one of the following parent locations may be specified.
// +kubebuilder:validation:Enum=Uri;Id;Kubernetes
type ImportReferenceType string

const (
	UriImportReferenceType        ImportReferenceType = "Uri"
	IdImportReferenceType         ImportReferenceType = "Id"
	KubernetesImportReferenceType ImportReferenceType = "Kubernetes"
)

// Location from where the an import reference is retrieved
// +k8s:openapi-gen=true
// +union
type ImportReferenceUnion struct {
	// type of location from where the referenced template structure should be retrieved
	// +
	// +unionDiscriminator
	// +optional
	ImportReferenceType ImportReferenceType `json:"importReferenceType,omitempty"`

	// Uri of a Devfile yaml file
	// +optional
	Uri string `json:"uri,omitempty"`

	// Id in a registry that contains a Devfile yaml file
	// +optional
	Id string `json:"id,omitempty"`

	// Reference to a Kubernetes CRD of type DevWorkspaceTemplate
	// +optional
	Kubernetes *KubernetesCustomResourceImportReference `json:"kubernetes,omitempty"`
}

type KubernetesCustomResourceImportReference struct {
	Name string `json:"name"`

	// +optional
	Namespace string `json:"namespace,omitempty"`
}

type ImportReference struct {
	ImportReferenceUnion `json:",inline"`
	// +optional
	RegistryUrl string `json:"registryUrl,omitempty"`
}
