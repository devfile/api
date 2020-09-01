package overriding

import (
	workspaces "github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	"github.com/hashicorp/go-multierror"
	"k8s.io/apimachinery/pkg/util/sets"
)

type checkFn func(elementType string, keysSets []sets.String) []error

func checkKeys(doCheck checkFn, toplevelListContainers ...workspaces.TopLevelListContainer) error {
	var errors *multierror.Error = nil

	listOfTopLevelLists := []workspaces.TopLevelLists{}
	for _, toplevelListContainer := range toplevelListContainers {
		listOfTopLevelLists = append(listOfTopLevelLists, toplevelListContainer.GetToplevelLists())
	}

	elementTypes := sets.String{}
	for _, toplevelLists := range listOfTopLevelLists {
		for elementType := range toplevelLists {
			elementTypes.Insert(elementType)
		}
	}

	for _, elementType := range elementTypes.List() {
		keySets := []sets.String{}
		for _, toplevelLists := range listOfTopLevelLists {
			keys := sets.String{}
			for _, keyed := range toplevelLists[elementType] {
				keys.Insert(keyed.Key())
			}
			keySets = append(keySets, keys)
		}
		checkErrors := doCheck(elementType, keySets)
		errors = multierror.Append(errors, checkErrors...)
	}

	return errors.ErrorOrNil()
}
