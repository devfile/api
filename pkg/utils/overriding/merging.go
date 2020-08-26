package overriding

import (
	"fmt"
	"strings"

	workspaces "github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func ensureNoConflictWithParent(mainContent *workspaces.DevWorkspaceTemplateSpecContent, parentflattenedContent *workspaces.DevWorkspaceTemplateSpecContent) error {
	return checkKeys(func(elementType string, keysSets []sets.String) []error {
		mainKeys := keysSets[0]
		parentOrPluginKeys := keysSets[1]
		overriddenElementsInMainContent := mainKeys.Intersection(parentOrPluginKeys)
		if overriddenElementsInMainContent.Len() > 0 {
			return []error{fmt.Errorf("Some %s are already defined in parent: %s. "+
				"If you want to override them, you should do it in the parent scope.",
				elementType,
				strings.Join(overriddenElementsInMainContent.List(), ", "))}
		}
		return []error{}
	},
		mainContent, parentflattenedContent)
}

func ensureNoConflictsWithPlugins(mainContent *workspaces.DevWorkspaceTemplateSpecContent, pluginFlattenedContents ...*workspaces.DevWorkspaceTemplateSpecContent) error {
	getPluginKey := func(pluginIndex int) string {
		index := 0
		for _, comp := range mainContent.Components {
			if comp.Plugin != nil {
				if pluginIndex == index {
					return comp.Name
				}
				index++
			}
		}
		return "unknown"
	}

	allSpecs := []workspaces.TopLevelListContainer{mainContent}
	for _, plugipluginFlattenedContent := range pluginFlattenedContents {
		allSpecs = append(allSpecs, plugipluginFlattenedContent)
	}
	return checkKeys(func(elementType string, keysSets []sets.String) []error {
		mainKeys := keysSets[0]
		pluginKeysSets := keysSets[1:]
		errs := []error{}
		for pluginNumber, pluginKeys := range pluginKeysSets {
			overriddenElementsInMainContent := mainKeys.Intersection(pluginKeys)

			if overriddenElementsInMainContent.Len() > 0 {
				errs = append(errs, fmt.Errorf("Some %s are already defined in plugin '%s': %s. "+
					"If you want to override them, you should do it in the plugin scope.",
					elementType,
					getPluginKey(pluginNumber),
					strings.Join(overriddenElementsInMainContent.List(), ", ")))
			}
		}
		return errs
	},
		allSpecs...)
}

// MergeDevWorkspaceTemplateSpec implements the merging logic of a main devfile content with flattened, already-overriden parent devfiles or plugins.
// On an `main` `DevfileWorkspaceTemplateSpec` (which is the core part of a devfile, without the `apiVersion` and `metadata`),
// it allows adding all the new overriden elements provided by flattened parent and plugins
//
// It is not allowed for to have duplicate (== with same key) commands, components or projects between the main content and the parent or plugins.
// An error would be thrown
//
// The result is a transformed `DevfileWorkspaceTemplateSpec` object, that does not contain any `plugin` component
// (since they are expected to be provided as flattened overriden devfiles in teh arguments)
func MergeDevWorkspaceTemplateSpec(
	mainContent *workspaces.DevWorkspaceTemplateSpecContent,
	parentFlattenedContent *workspaces.DevWorkspaceTemplateSpecContent,
	pluginFlattenedContents ...*workspaces.DevWorkspaceTemplateSpecContent) (*workspaces.DevWorkspaceTemplateSpecContent, error) {

	allContents := []*workspaces.DevWorkspaceTemplateSpecContent{parentFlattenedContent}
	allContents = append(allContents, pluginFlattenedContents...)
	allContents = append(allContents, mainContent)

	if err := ensureNoConflictWithParent(mainContent, parentFlattenedContent); err != nil {
		return nil, err
	}

	if err := ensureNoConflictsWithPlugins(mainContent, pluginFlattenedContents...); err != nil {
		return nil, err
	}

	result := workspaces.DevWorkspaceTemplateSpecContent{}

	preStartCommands := sets.String{}
	postStartCommands := sets.String{}
	preStopCommands := sets.String{}
	postStopCommands := sets.String{}

	for _, content := range allContents {
		if content.Commands != nil && len(content.Commands) > 0 {
			if result.Commands == nil {
				result.Commands = []workspaces.Command{}
			}
			for _, command := range content.Commands {
				result.Commands = append(result.Commands, command)
			}
		}

		if content.Components != nil && len(content.Components) > 0 {
			for _, component := range content.Components {
				// We skip the plugins of the main content, since they have been provided by flattened plugin content.
				if content == mainContent && component.Plugin != nil {
					continue
				}
				if result.Components == nil {
					result.Components = []workspaces.Component{}
				}
				result.Components = append(result.Components, component)
			}
		}

		if content.Projects != nil && len(content.Projects) > 0 {
			if result.Projects == nil {
				result.Projects = []workspaces.Project{}
			}
			for _, project := range content.Projects {
				result.Projects = append(result.Projects, project)
			}
		}

		if content.Events != nil {
			if result.Events == nil {
				result.Events = &workspaces.Events{}
			}
			preStartCommands = preStartCommands.Union(sets.NewString(content.Events.PreStart...))
			postStartCommands = postStartCommands.Union(sets.NewString(content.Events.PostStart...))
			preStopCommands = preStopCommands.Union(sets.NewString(content.Events.PreStop...))
			postStopCommands = postStopCommands.Union(sets.NewString(content.Events.PostStop...))
		}
	}

	if result.Events != nil {
		result.Events.PreStart = preStartCommands.List()
		result.Events.PostStart = postStartCommands.List()
		result.Events.PreStop = preStopCommands.List()
		result.Events.PostStop = postStopCommands.List()
	}

	return &result, nil
}

// MergeDevWorkspaceTemplateSpecBytes implements the merging logic of a main devfile content with flattened, already-overridden parent devfiles or plugins.
// On an json or yaml document that contains the core content of the devfile (which is the core part of a devfile, without the `apiVersion` and `metadata`),
// it allows adding all the new overridden elements provided by flattened parent and plugins (also provided as json or yaml documents)
//
// It is not allowed for to have duplicate (== with same key) commands, components or projects between the main content and the parent or plugins.
// An error would be thrown
//
// The result is a transformed `DevfileWorkspaceTemplateSpec` object, that does not contain any `plugin` component
// (since they are expected to be provided as flattened overridden devfiles in the arguments)
func MergeDevWorkspaceTemplateSpecBytes(originalBytes []byte, flattenedParentBytes []byte, flattenPluginsBytes ...[]byte) (*workspaces.DevWorkspaceTemplateSpecContent, error) {
	originalJson, err := yaml.ToJSON(originalBytes)
	if err != nil {
		return nil, err
	}

	original := workspaces.DevWorkspaceTemplateSpecContent{}
	err = json.Unmarshal(originalJson, &original)
	if err != nil {
		return nil, err
	}

	flattenedParentJson, err := yaml.ToJSON(flattenedParentBytes)
	if err != nil {
		return nil, err
	}

	flattenedParent := workspaces.DevWorkspaceTemplateSpecContent{}
	err = json.Unmarshal(flattenedParentJson, &flattenedParent)
	if err != nil {
		return nil, err
	}

	flattenedPlugins := []*workspaces.DevWorkspaceTemplateSpecContent{}
	for _, flattenedPluginBytes := range flattenPluginsBytes {
		flattenedPluginJson, err := yaml.ToJSON(flattenedPluginBytes)
		if err != nil {
			return nil, err
		}

		flattenedPlugin := workspaces.DevWorkspaceTemplateSpecContent{}
		err = json.Unmarshal(flattenedPluginJson, &flattenedPlugin)
		if err != nil {
			return nil, err
		}
		flattenedPlugins = append(flattenedPlugins, &flattenedPlugin)
	}

	return MergeDevWorkspaceTemplateSpec(&original, &flattenedParent, flattenedPlugins...)
}
