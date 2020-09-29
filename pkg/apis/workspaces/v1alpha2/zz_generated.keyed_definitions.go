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
