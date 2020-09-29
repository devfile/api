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

type Events struct {
	WorkspaceEvents `json:",inline"`
}

type WorkspaceEvents struct {
	// Names of commands that should be executed before the workspace start.
	// Kubernetes-wise, these commands would typically be executed in init containers of the workspace POD.
	// +optional
	PreStart []string `json:"preStart,omitempty"`

	// Names of commands that should be executed after the workspace is completely started.
	// In the case of Che-Theia, these commands should be executed after all plugins and extensions have started, including project cloning.
	// This means that those commands are not triggered until the user opens the IDE in his browser.
	// +optional
	PostStart []string `json:"postStart,omitempty"`

	// +optional
	// Names of commands that should be executed before stopping the workspace.
	PreStop []string `json:"preStop,omitempty"`

	// +optional
	// Names of commands that should be executed after stopping the workspace.
	PostStop []string `json:"postStop,omitempty"`
}
