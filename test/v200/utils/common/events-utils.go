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

package common

import (
	"fmt"

	schema "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

// EventsAdded adds event to the test schema and notifies a registered follower
func (devfile *TestDevfile) EventsAdded(events *schema.Events) {
	LogInfoMessage(fmt.Sprintf("events added"))
	devfile.SchemaDevFile.Events = events
	if devfile.Follower != nil {
		devfile.Follower.AddEvent(*events)
	}
}

// EventsUpdated notifies a registered follower that the events have been updated
func (devfile *TestDevfile) EventsUpdated(events *schema.Events) {
	LogInfoMessage(fmt.Sprintf("events updated"))
	if devfile.Follower != nil {
		devfile.Follower.UpdateEvent(*events)
	}
}

// AddEvents adds events in the test schema structure and populates it with random attributes
func (devfile *TestDevfile) AddEvents() schema.Events {
	events := schema.Events{}
	devfile.EventsAdded(&events)
	devfile.SetEventsValues(&events)
	return events
}

// SetEventsValues randomly adds/modifies attributes of the supplied events
func (devfile *TestDevfile) SetEventsValues(events *schema.Events) {
	if GetRandomDecision(4, 1) {
		numPreStart := GetRandomNumber(1, 5)
		LogInfoMessage(fmt.Sprintf("   ....... add %d command(s) to PreStart event", numPreStart))
		for i := 0; i < numPreStart; i++ {
			if GetRandomDecision(4, 1) {
				events.PreStart = append(events.PreStart, devfile.AddCommand(schema.ApplyCommandType).Id)
			} else {
				compositeCommand := devfile.AddCommand(schema.CompositeCommandType)
				devfile.SetCompositeCommandCommands(&compositeCommand, schema.ApplyCommandType)
				events.PreStart = append(events.PreStart, compositeCommand.Id)
			}
		}
	}
	if GetRandomDecision(4, 1) {
		numPostStart := GetRandomNumber(1, 5)
		LogInfoMessage(fmt.Sprintf("   ....... add %d command(s) to PostStart event", numPostStart))
		for i := 0; i < numPostStart; i++ {
			if GetRandomDecision(4, 1) {
				events.PostStart = append(events.PostStart, devfile.AddCommand(schema.ExecCommandType).Id)
			} else {
				compositeCommand := devfile.AddCommand(schema.CompositeCommandType)
				devfile.SetCompositeCommandCommands(&compositeCommand, schema.ExecCommandType)
				events.PostStart = append(events.PostStart, compositeCommand.Id)
			}
		}
	}
	if GetRandomDecision(4, 1) {
		numPreStop := GetRandomNumber(1, 5)
		LogInfoMessage(fmt.Sprintf("   ....... add %d command(s) to PreStop event", numPreStop))
		for i := 0; i < numPreStop; i++ {
			if GetRandomDecision(4, 1) {
				events.PreStop = append(events.PreStop, devfile.AddCommand(schema.ExecCommandType).Id)
			} else {
				compositeCommand := devfile.AddCommand(schema.CompositeCommandType)
				devfile.SetCompositeCommandCommands(&compositeCommand, schema.ExecCommandType)
				events.PreStop = append(events.PreStop, compositeCommand.Id)
			}
		}
	}
	if GetRandomDecision(4, 1) {
		numPostStop := GetRandomNumber(1, 5)
		LogInfoMessage(fmt.Sprintf("   ....... add %d command(s) to PostStop event", numPostStop))
		for i := 0; i < numPostStop; i++ {
			if GetRandomDecision(4, 1) {
				events.PostStop = append(events.PostStop, devfile.AddCommand(schema.ApplyCommandType).Id)
			} else {
				compositeCommand := devfile.AddCommand(schema.CompositeCommandType)
				devfile.SetCompositeCommandCommands(&compositeCommand, schema.ApplyCommandType)
				events.PostStop = append(events.PostStop, compositeCommand.Id)
			}
		}
	}
	devfile.EventsUpdated(events)
}
