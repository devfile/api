package v1alpha1

import (
	"github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
)

func convertEventsTo_v1alpha2(src *Events, dest *v1alpha2.Events) error {
	if src != nil {
		dest.WorkspaceEvents = v1alpha2.WorkspaceEvents(src.WorkspaceEvents)
	}
	return nil
}

func convertEventsFrom_v1alpha2(src *v1alpha2.Events, dest *Events) error {
	if src != nil {
		dest.WorkspaceEvents = WorkspaceEvents(src.WorkspaceEvents)
	}
	return nil
}
