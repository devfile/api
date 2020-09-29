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

func (container DevWorkspaceTemplateSpecContent) GetToplevelLists() TopLevelLists {
	return TopLevelLists{
		"Components":      extractKeys(container.Components),
		"Projects":        extractKeys(container.Projects),
		"StarterProjects": extractKeys(container.StarterProjects),
		"Commands":        extractKeys(container.Commands),
	}
}

func (container ParentOverrides) GetToplevelLists() TopLevelLists {
	return TopLevelLists{
		"Components":      extractKeys(container.Components),
		"Projects":        extractKeys(container.Projects),
		"StarterProjects": extractKeys(container.StarterProjects),
		"Commands":        extractKeys(container.Commands),
	}
}

func (container PluginOverridesParentOverride) GetToplevelLists() TopLevelLists {
	return TopLevelLists{
		"Components": extractKeys(container.Components),
		"Commands":   extractKeys(container.Commands),
	}
}

func (container PluginOverrides) GetToplevelLists() TopLevelLists {
	return TopLevelLists{
		"Components": extractKeys(container.Components),
		"Commands":   extractKeys(container.Commands),
	}
}
