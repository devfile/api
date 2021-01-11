package validation

import (
	"fmt"
	"strings"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

// ValidateEvents validates all the devfile events
func ValidateEvents(events v1alpha2.Events, commands []v1alpha2.Command) error {

	eventErrors := ""

	commandMap := getCommandsMap(commands)

	switch {
	case len(events.PreStart) > 0:
		if preStartErr := isEventValid(events.PreStart, "preStart", commandMap); preStartErr != nil {
			eventErrors += fmt.Sprintf("\n%s", preStartErr.Error())
		}
		fallthrough
	case len(events.PostStart) > 0:
		if postStartErr := isEventValid(events.PostStart, "postStart", commandMap); postStartErr != nil {
			eventErrors += fmt.Sprintf("\n%s", postStartErr.Error())
		}
		fallthrough
	case len(events.PreStop) > 0:
		if preStopErr := isEventValid(events.PreStop, "preStop", commandMap); preStopErr != nil {
			eventErrors += fmt.Sprintf("\n%s", preStopErr.Error())
		}
		fallthrough
	case len(events.PostStop) > 0:
		if postStopErr := isEventValid(events.PostStop, "postStop", commandMap); postStopErr != nil {
			eventErrors += fmt.Sprintf("\n%s", postStopErr.Error())
		}
	}

	// if there is any validation error, return it
	if len(eventErrors) > 0 {
		return fmt.Errorf("devfile events validation error: %s", eventErrors)
	}

	return nil
}

// isEventValid checks if events belonging to a specific event type are valid ie; event should map to a valid devfile command
func isEventValid(eventNames []string, eventType string, commandMap map[string]v1alpha2.Command) error {
	var invalidEvents []string

	for _, eventName := range eventNames {
		if _, ok := commandMap[strings.ToLower(eventName)]; !ok {
			invalidEvents = append(invalidEvents, eventName)
		}
	}

	if len(invalidEvents) > 0 {
		eventErrors := fmt.Sprintf("\n%s does not map to a valid devfile command", strings.Join(invalidEvents, ", "))
		return &InvalidEventError{eventType: eventType, errorMsg: eventErrors}
	}

	return nil
}
