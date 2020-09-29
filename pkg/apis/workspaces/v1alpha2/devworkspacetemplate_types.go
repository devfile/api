//
// Copyright (c) 2019-2020 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevWorkspaceTemplate is the Schema for the devworkspacetemplates API
// +kubebuilder:resource:path=devworkspacetemplates,scope=Namespaced
// +devfile:jsonschema:generate
type DevWorkspaceTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DevWorkspaceTemplateSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DevWorkspaceTemplateList contains a list of DevWorkspaceTemplate
type DevWorkspaceTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DevWorkspaceTemplate `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DevWorkspaceTemplate{}, &DevWorkspaceTemplateList{})
}
