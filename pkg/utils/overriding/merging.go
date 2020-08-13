package overriding

import (
	//	"errors"
	"errors"
	"strings"

	workspaces "github.com/devfile/kubernetes-api/pkg/apis/workspaces/v1alpha1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func ensureNoConflicWithParent(mainContent *keyedDevWorkspaceTemplateSpecContent, parentflattenedContent *keyedDevWorkspaceTemplateSpecContent) error {
	return checkKeys(func(elementType string, keysSets []sets.String) []error {
		mainKeys := keysSets[0]
		parentOrPluginKeys := keysSets[1]
		overridenElementsInMainContent := mainKeys.Intersection(parentOrPluginKeys)
		if overridenElementsInMainContent.Len() > 0 {
			return []error{errors.New("Some " +
				elementType +
				" elements are already defined in parent: " +
				strings.Join(overridenElementsInMainContent.List(), ", ") +
				". If you want to override them, you should do it in the parent scope.")}
		}
		return []error{}
	},
		mainContent, parentflattenedContent)
}

func ensureNoConflicsWithPlugins(mainContent *keyedDevWorkspaceTemplateSpecContent, pluginFlattenedContents ...*keyedDevWorkspaceTemplateSpecContent) error {
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

	allSpecs := append([]*keyedDevWorkspaceTemplateSpecContent{mainContent}, pluginFlattenedContents...)
	return checkKeys(func(elementType string, keysSets []sets.String) []error {
		mainKeys := keysSets[0]
		pluginKeysSets := keysSets[1:]
		errs := []error{}
		for pluginNumber, pluginKeys := range pluginKeysSets {
			overridenElementsInMainContent := mainKeys.Intersection(pluginKeys)

			if overridenElementsInMainContent.Len() > 0 {
				errs = append(errs, errors.New("Some "+
					elementType+
					" elements are already defined in plugin '"+
					getPluginKey(pluginNumber)+
					"': "+
					strings.Join(overridenElementsInMainContent.List(), ", ")+
					". If you want to override them, you should do it in the plugin scope."))
			}
		}
		return errs
	},
		allSpecs...)
}

func mergeKeyedDevWorkspaceTemplateSpec(mainContent *keyedDevWorkspaceTemplateSpecContent, parentFlattenedContent *keyedDevWorkspaceTemplateSpecContent, pluginFlattenedContents ...*keyedDevWorkspaceTemplateSpecContent) (*keyedDevWorkspaceTemplateSpecContent, error) {
	allContents := []*keyedDevWorkspaceTemplateSpecContent { parentFlattenedContent }
	allContents = append(allContents, pluginFlattenedContents...)
	allContents = append(allContents, mainContent)
	for _, content := range allContents {
		addKeys(content)
	}
	
	if err := ensureNoConflicWithParent(mainContent, parentFlattenedContent); err != nil {
		return nil, err
	}

	if err := ensureNoConflicsWithPlugins(mainContent, pluginFlattenedContents...); err != nil {
		return nil, err
	}

	result := keyedDevWorkspaceTemplateSpecContent{
	}

	preStartCommands := sets.String{}
	postStartCommands := sets.String{}
	preStopCommands := sets.String{}
	postStopCommands := sets.String{}

	for _, content := range allContents {
		if content.Commands != nil && len(content.Commands) > 0 {
			if result.Commands == nil {
				result.Commands = []keyedCommand{}
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
					result.Components = []keyedComponent{}
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

// MergeDevWorkspaceTemplateSpecBytes implements the merging logic of a main devfile content with flattened, already-overriden parent devfiles or plugins.
// On an Json or Yaml document that contains the core content of the devfile (which is the core part of a devfile, without the `apiVersion` and `metadata`),
// it allows adding all the new overriden elements provided by flattened parent and plugins (also provided as Json or Yaml documents)
//
// It is not allowed for to have duplicate (== with same key) commands, components or projects between the main content and the parent or plugins.
// An error would be thrown  
// 
// The result is a transformed `DevfileWorkspaceTemplateSpec` object, that does not contain any `plugin` component
// (since they are expected to be provided as flattened overriden devfiles in teh arguments)
func MergeDevWorkspaceTemplateSpecBytes(originalBytes []byte, flattenedParentBytes []byte, flattenPluginsBytes ...[]byte) (*workspaces.DevWorkspaceTemplateSpecContent, error) {
	originalJson, err := yaml.ToJSON(originalBytes)
	if err != nil {
		return nil, err
	}

	keyedOriginal := keyedDevWorkspaceTemplateSpecContent{}
	err = json.Unmarshal(originalJson, &keyedOriginal)
	if err != nil {
		return nil, err
	}

	flattenedParentJson, err := yaml.ToJSON(flattenedParentBytes)
	if err != nil {
		return nil, err
	}

	keyedFlattenedParent := keyedDevWorkspaceTemplateSpecContent{}
	err = json.Unmarshal(flattenedParentJson, &keyedFlattenedParent)
	if err != nil {
		return nil, err
	}

	keyedFlattenedPlugins := []*keyedDevWorkspaceTemplateSpecContent{}
	for _, flattenedPluginBytes := range flattenPluginsBytes {
		flattenedPluginJson, err := yaml.ToJSON(flattenedPluginBytes)
		if err != nil {
			return nil, err
		}

		keyedFlattenedPlugin := keyedDevWorkspaceTemplateSpecContent{}
		err = json.Unmarshal(flattenedPluginJson, &keyedFlattenedPlugin)
		if err != nil {
			return nil, err
		}
		keyedFlattenedPlugins = append(keyedFlattenedPlugins, &keyedFlattenedPlugin)
	}

	result, err := mergeKeyedDevWorkspaceTemplateSpec(&keyedOriginal, &keyedFlattenedParent, keyedFlattenedPlugins...)
	if err != nil {
		return nil, err
	}
	return removeKeys(result), nil
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
func MergeDevWorkspaceTemplateSpec(main *workspaces.DevWorkspaceTemplateSpecContent, flattenedParent *workspaces.DevWorkspaceTemplateSpecContent, flattenedPlugins ...*workspaces.DevWorkspaceTemplateSpecContent) (*workspaces.DevWorkspaceTemplateSpecContent, error) {
	mainBytes, err := json.Marshal(main)
	if err != nil {
		return nil, err
	}

	flattenedParentBytes, err := json.Marshal(flattenedParent)
	if err != nil {
		return nil, err
	}

	flattenedPluginsBytes := [][]byte{}
	for _, flattenedPlugin := range flattenedPlugins {
		flattenedPluginBytes, err := json.Marshal(flattenedPlugin)
		if err != nil {
			return nil, err
		}
		flattenedPluginsBytes = append(flattenedPluginsBytes, flattenedPluginBytes)
	}

	return MergeDevWorkspaceTemplateSpecBytes(mainBytes, flattenedParentBytes, flattenedPluginsBytes...)
}
