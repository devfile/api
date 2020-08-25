package overriding

import (
	workspaces "github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	"github.com/hashicorp/go-multierror"
	"k8s.io/apimachinery/pkg/util/sets"
)

type keyedCommand struct {
	Id string `json:"id,omitempty"`

	workspaces.Command `json:",inline"`
}

type keyedComponent struct {
	Name string `json:"name,omitempty"`

	workspaces.Component `json:",inline"`
}

type keyedDevWorkspaceTemplateSpecContent struct {
	// Projects worked on in the workspace, containing names and sources locations
	// +optional
	Projects []workspaces.Project `json:"projects,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// List of the workspace components, such as editor and plugins,
	// user-provided containers, or other types of components.
	// +optional
	Components []keyedComponent `json:"components,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// Predefined, ready-to-use, workspace-related commands.
	// +optional
	Commands []keyedCommand `json:"commands,omitempty" patchStrategy:"merge" patchMergeKey:"id"`

	// Bindings of commands to events.
	// Each command is referred-to by its name.
	// +optional
	Events *workspaces.Events `json:"events,omitempty"`
}

func addKeys(keyedSpec *keyedDevWorkspaceTemplateSpecContent) error {
	for idx := range keyedSpec.Commands {
		key, err := keyedSpec.Commands[idx].Key()
		if err != nil {
			return err
		}
		keyedSpec.Commands[idx].Id = key
	}
	for idx := range keyedSpec.Components {
		key, err := keyedSpec.Components[idx].Key()
		if err != nil {
			return err
		}
		keyedSpec.Components[idx].Name = key
	}
	return nil
}

func removeKeys(keyedSpec *keyedDevWorkspaceTemplateSpecContent) *workspaces.DevWorkspaceTemplateSpecContent {
	result := workspaces.DevWorkspaceTemplateSpecContent{
		Projects: keyedSpec.Projects,
		Events:   keyedSpec.Events,
	}

	if keyedSpec.Commands != nil && len(keyedSpec.Commands) > 0 {
		result.Commands = []workspaces.Command{}
	}

	if keyedSpec.Components != nil && len(keyedSpec.Components) > 0 {
		result.Components = []workspaces.Component{}
	}

	for _, keyedCommand := range keyedSpec.Commands {
		result.Commands = append(result.Commands, keyedCommand.Command)
	}
	for _, keyedComponent := range keyedSpec.Components {
		result.Components = append(result.Components, keyedComponent.Component)
	}
	return &result
}

func getCommandKeySets(keyedSpecs ...*keyedDevWorkspaceTemplateSpecContent) []sets.String {
	keySets := []sets.String{}
	for _, keyedSpec := range keyedSpecs {
		keys := sets.String{}
		for _, item := range keyedSpec.Commands {
			keys.Insert(item.Id)
		}
		keySets = append(keySets, keys)
	}
	return keySets
}

func getComponentKeySets(keyedSpecs ...*keyedDevWorkspaceTemplateSpecContent) []sets.String {
	keySets := []sets.String{}
	for _, keyedSpec := range keyedSpecs {
		keys := sets.String{}
		for _, item := range keyedSpec.Components {
			keys.Insert(item.Name)
		}
		keySets = append(keySets, keys)
	}
	return keySets
}

func getProjectKeySets(keyedSpecs ...*keyedDevWorkspaceTemplateSpecContent) []sets.String {
	keySets := []sets.String{}
	for _, keyedSpec := range keyedSpecs {
		keys := sets.String{}
		for _, item := range keyedSpec.Projects {
			keys.Insert(item.Name)
		}
		keySets = append(keySets, keys)
	}
	return keySets
}

type checkFn func(elementType string, keysSets []sets.String) []error

func checkKeys(doCheck checkFn, keyedSpecs ...*keyedDevWorkspaceTemplateSpecContent) error {
	var errors *multierror.Error = nil

	for _, test := range []struct {
		elementType      string
		keySetsRetriever func(...*keyedDevWorkspaceTemplateSpecContent) []sets.String
	}{
		{
			"Command",
			getCommandKeySets,
		},
		{
			"Component",
			getComponentKeySets,
		},
		{
			"Project",
			getProjectKeySets,
		},
	} {
		for _, err := range doCheck(test.elementType, test.keySetsRetriever(keyedSpecs...)) {
			errors = multierror.Append(errors, err)
		}
	}

	return errors.ErrorOrNil()
}
