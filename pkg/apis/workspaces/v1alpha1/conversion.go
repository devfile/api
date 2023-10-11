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
)

const (
	// ConvertedFromAttribute marks a devfile element as being converted from a different underlying field. For example,
	// since v1alpha2 does not support
	ConvertedFromAttribute = "conversion.api.devfile.io/converted-from"
)

func convertDevWorkspaceTo_v1alpha2(src *DevWorkspace, dest *v1alpha2.DevWorkspace) error {
	dest.ObjectMeta = src.ObjectMeta
	dest.Status.DevWorkspaceId = src.Status.WorkspaceId
	dest.Status.MainUrl = src.Status.IdeUrl
	dest.Status.Phase = v1alpha2.DevWorkspacePhase(src.Status.Phase)
	dest.Status.Message = src.Status.Message
	convertConditionsTo_v1alpha2(src, dest)
	dest.Spec.RoutingClass = src.Spec.RoutingClass
	dest.Spec.Started = src.Spec.Started

	return convertDevWorkspaceTemplateSpecTo_v1alpha2(&src.Spec.Template, &dest.Spec.Template)
}

func convertDevWorkspaceFrom_v1alpha2(src *v1alpha2.DevWorkspace, dest *DevWorkspace) error {
	dest.ObjectMeta = src.ObjectMeta
	dest.Status.WorkspaceId = src.Status.DevWorkspaceId
	dest.Status.IdeUrl = src.Status.MainUrl
	dest.Status.Phase = WorkspacePhase(src.Status.Phase)
	dest.Status.Message = src.Status.Message
	convertConditionsFrom_v1alpha2(src, dest)
	dest.Spec.RoutingClass = src.Spec.RoutingClass
	dest.Spec.Started = src.Spec.Started

	return convertDevWorkspaceTemplateSpecFrom_v1alpha2(&src.Spec.Template, &dest.Spec.Template)
}

func convertDevWorkspaceTemplateTo_v1alpha2(src *DevWorkspaceTemplate, dest *v1alpha2.DevWorkspaceTemplate) error {
	dest.ObjectMeta = src.ObjectMeta
	return convertDevWorkspaceTemplateSpecTo_v1alpha2(&src.Spec, &dest.Spec)
}

func convertDevWorkspaceTemplateFrom_v1alpha2(src *v1alpha2.DevWorkspaceTemplate, dest *DevWorkspaceTemplate) error {
	dest.ObjectMeta = src.ObjectMeta
	return convertDevWorkspaceTemplateSpecFrom_v1alpha2(&src.Spec, &dest.Spec)
}

func convertDevWorkspaceTemplateSpecTo_v1alpha2(src *DevWorkspaceTemplateSpec, dest *v1alpha2.DevWorkspaceTemplateSpec) error {
	if src.Parent != nil {
		dest.Parent = &v1alpha2.Parent{}
		err := convertParentTo_v1alpha2(src.Parent, dest.Parent)
		if err != nil {
			return err
		}
	}
	for _, srcComponent := range src.Components {
		srcComponent := srcComponent
		destComponent := v1alpha2.Component{}
		err := convertComponentTo_v1alpha2(&srcComponent, &destComponent)
		if err != nil {
			return err
		}
		dest.Components = append(dest.Components, destComponent)
	}
	for _, srcProject := range src.Projects {
		srcProject := srcProject
		destProject := v1alpha2.Project{}
		err := convertProjectTo_v1alpha2(&srcProject, &destProject)
		if err != nil {
			return err
		}
		dest.Projects = append(dest.Projects, destProject)
	}
	for _, srcStarterProject := range src.StarterProjects {
		srcStarterProject := srcStarterProject
		destStarterProject := v1alpha2.StarterProject{}
		err := convertStarterProjectTo_v1alpha2(&srcStarterProject, &destStarterProject)
		if err != nil {
			return err
		}
		dest.StarterProjects = append(dest.StarterProjects, destStarterProject)
	}
	for _, srcCommand := range src.Commands {
		srcCommand := srcCommand
		destCommand := v1alpha2.Command{}
		err := convertCommandTo_v1alpha2(&srcCommand, &destCommand)
		if err != nil {
			return err
		}
		dest.Commands = append(dest.Commands, destCommand)
	}
	if src.Events != nil {
		dest.Events = &v1alpha2.Events{}
		err := convertEventsTo_v1alpha2(src.Events, dest.Events)
		if err != nil {
			return err
		}
	}
	return nil
}

func convertDevWorkspaceTemplateSpecFrom_v1alpha2(src *v1alpha2.DevWorkspaceTemplateSpec, dest *DevWorkspaceTemplateSpec) error {
	if src.Parent != nil {
		dest.Parent = &Parent{}
		err := convertParentFrom_v1alpha2(src.Parent, dest.Parent)
		if err != nil {
			return err
		}
	}
	for _, srcComponent := range src.Components {
		srcComponent := srcComponent
		destComponent := Component{}
		err := convertComponentFrom_v1alpha2(&srcComponent, &destComponent)
		if err != nil {
			return err
		}
		dest.Components = append(dest.Components, destComponent)
	}
	for _, srcProject := range src.Projects {
		srcProject := srcProject
		destProject := Project{}
		err := convertProjectFrom_v1alpha2(&srcProject, &destProject)
		if err != nil {
			return err
		}
		dest.Projects = append(dest.Projects, destProject)
	}
	for _, srcStarterProject := range src.StarterProjects {
		srcStarterProject := srcStarterProject
		destStarterProject := StarterProject{}
		err := convertStarterProjectFrom_v1alpha2(&srcStarterProject, &destStarterProject)
		if err != nil {
			return err
		}
		dest.StarterProjects = append(dest.StarterProjects, destStarterProject)
	}
	for _, srcCommand := range src.Commands {
		srcCommand := srcCommand
		destCommand := Command{}
		err := convertCommandFrom_v1alpha2(&srcCommand, &destCommand)
		if err != nil {
			return err
		}
		dest.Commands = append(dest.Commands, destCommand)
	}
	if src.Events != nil {
		dest.Events = &Events{}
		err := convertEventsFrom_v1alpha2(src.Events, dest.Events)
		if err != nil {
			return err
		}
	}
	return nil
}

func convertConditionsTo_v1alpha2(src *DevWorkspace, dest *v1alpha2.DevWorkspace) {
	for _, srcCondition := range src.Status.Conditions {
		dstCondition := v1alpha2.DevWorkspaceCondition{}
		dstCondition.Status = srcCondition.Status
		dstCondition.LastTransitionTime = srcCondition.LastTransitionTime
		dstCondition.Message = srcCondition.Message
		dstCondition.Reason = srcCondition.Reason
		dstCondition.Type = v1alpha2.DevWorkspaceConditionType(srcCondition.Type)
		dest.Status.Conditions = append(dest.Status.Conditions, dstCondition)
	}
}

func convertConditionsFrom_v1alpha2(src *v1alpha2.DevWorkspace, dest *DevWorkspace) {
	for _, srcCondition := range src.Status.Conditions {
		dstCondition := WorkspaceCondition{}
		dstCondition.Status = srcCondition.Status
		dstCondition.LastTransitionTime = srcCondition.LastTransitionTime
		dstCondition.Message = srcCondition.Message
		dstCondition.Reason = srcCondition.Reason
		dstCondition.Type = WorkspaceConditionType(srcCondition.Type)
		dest.Status.Conditions = append(dest.Status.Conditions, dstCondition)
	}
}
