package v1alpha1

import (
	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

func convertEventsTo_v1alpha2(src *Events, dest *v1alpha2.Events) error {
	if src != nil {
		dest.DevWorkspaceEvents = v1alpha2.DevWorkspaceEvents(src.WorkspaceEvents)
	}
	return nil
}

func convertEventsFrom_v1alpha2(src *v1alpha2.Events, dest *Events) error {
	if src != nil {
		dest.WorkspaceEvents = WorkspaceEvents(src.DevWorkspaceEvents)
	}
	return nil
}
