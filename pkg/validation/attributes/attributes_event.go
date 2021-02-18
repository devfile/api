package attributes

import (
	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
)

// ValidateEvents validates the events data for a global attribute
func ValidateEvents(attributes apiAttributes.Attributes, events *v1alpha2.Events) error {
	var err error

	if events != nil {
		switch {
		case len(events.PreStart) > 0:
			for i := range events.PreStart {
				if events.PreStart[i], err = validateAndReplaceDataWithAttribute(events.PreStart[i], attributes); err != nil {
					return err
				}
			}
			fallthrough
		case len(events.PostStart) > 0:
			for i := range events.PostStart {
				if events.PostStart[i], err = validateAndReplaceDataWithAttribute(events.PostStart[i], attributes); err != nil {
					return err
				}
			}
			fallthrough
		case len(events.PreStop) > 0:
			for i := range events.PreStop {
				if events.PreStop[i], err = validateAndReplaceDataWithAttribute(events.PreStop[i], attributes); err != nil {
					return err
				}
			}
			fallthrough
		case len(events.PostStop) > 0:
			for i := range events.PostStop {
				if events.PostStop[i], err = validateAndReplaceDataWithAttribute(events.PostStop[i], attributes); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
