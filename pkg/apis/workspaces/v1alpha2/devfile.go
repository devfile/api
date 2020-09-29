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

// Devfile describes the structure of a cloud-native workspace and development environment.
// +devfile:jsonschema:generate:omitCustomUnionMembers=true
type Devfile struct {
	// Devfile schema version
	// +kubebuilder:validation:Pattern=^([2-9][0-9]*)\.([0-9]+)\.([0-9]+)(\-[0-9a-z-]+(\.[0-9a-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$
	SchemaVersion string `json:"schemaVersion"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	// Optional metadata
	Metadata DevfileMetadata `json:"metadata,omitempty"`

	DevWorkspaceTemplateSpec `json:",inline"`
}

type DevfileMetadata struct {
	// Optional devfile name
	Name string `json:"name,omitempty"`

	// Optional semver-compatible version
	// +kubebuilder:validation:Pattern=^([0-9]+)\.([0-9]+)\.([0-9]+)(\-[0-9a-z-]+(\.[0-9a-z-]+)*)?(\+[0-9A-Za-z-]+(\.[0-9A-Za-z-]+)*)?$
	Version string `json:"version,omitempty"`
}
