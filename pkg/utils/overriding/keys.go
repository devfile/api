package overriding

import (
	workspaces "github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	"github.com/hashicorp/go-multierror"
	"k8s.io/apimachinery/pkg/util/sets"
)

type checkFn func(elementType string, keysSets []sets.String) []error

// checkKeys provides a generic way to apply some validation on the content of each type of top-level list
// contained in the `toplevelListContainers` passed in argument.
//
// For each type of top-level list, the `keysSets` argument that will be passed to the `doCheck` function
// contains the the key sets that correspond to the `toplevelListContainers` passed to this method,
// in the same order.
func checkKeys(doCheck checkFn, toplevelListContainers ...workspaces.TopLevelListContainer) error {
	var errors *multierror.Error

	// Retrieve the top-level lists (Commands, Projects, Components) for each TopLevelListContainer
	// (DevWorkspaceTemplateSpec, PluginOverrides, ...)
	topLevelListsByContainer := []workspaces.TopLevelLists{}
	for _, toplevelListContainer := range toplevelListContainers {
		topLevelListsByContainer = append(topLevelListsByContainer, toplevelListContainer.GetToplevelLists())
	}

	// Gather all types of top-level lists: Commands, Projects, etc ...
	listTypes := sets.String{}
	for _, toplevelLists := range topLevelListsByContainer {
		for listType := range toplevelLists {
			listTypes.Insert(listType)
		}
	}

	// For each type of top-level list (Commands, Projects), etc ...
	for _, listType := range listTypes.List() {
		// For each toplevel-list container, extract the set of keys of the given type of top-level list
		keySets := []sets.String{}
		for _, toplevelLists := range topLevelListsByContainer {
			keys := sets.String{}
			for _, keyed := range toplevelLists[listType] {
				keys.Insert(keyed.Key())
			}
			keySets = append(keySets, keys)
		}
		// Perform the check on the sets of keys for the given type of top-level list
		checkErrors := doCheck(listType, keySets)
		errors = multierror.Append(errors, checkErrors...)
	}

	return errors.ErrorOrNil()
}
