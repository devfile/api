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

// Component that allows the developer to declare and configure a volume into his workspace
type VolumeComponent struct {
	BaseComponent `json:",inline"`
	Volume        `json:",inline"`
}

// Volume that should be mounted to a component container
type Volume struct {
	// +optional
	// Size of the volume
	Size string `json:"size,omitempty"`
}
