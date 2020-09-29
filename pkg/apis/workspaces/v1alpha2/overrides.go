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

// +k8s:deepcopy-gen=false
type Overrides interface {
	TopLevelListContainer
	isOverride()
}

// OverridesBase is used in the Overrides generator in order to provide a common base for the generated Overrides
// So please be careful when renaming
type OverridesBase struct{}
