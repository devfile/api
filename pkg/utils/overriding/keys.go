package overriding

import (
	"reflect"

	dw "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	attributesPkg "github.com/devfile/api/v2/pkg/attributes"
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
func checkKeys(doCheck checkFn, toplevelListContainers ...dw.TopLevelListContainer) error {
	var errors *multierror.Error

	// intermediate storage for the conversion []map[string]KeyedList -> map[string][]sets.String
	listTypeToKeys := map[string][]sets.String{}

	// Flatten []map[string]KeyedList -> map[string][]KeyedList based on map keys and convert each KeyedList
	// into a sets.String
	for _, topLevelListContainer := range toplevelListContainers {
		topLevelList := topLevelListContainer.GetToplevelLists()
		for listType, listElem := range topLevelList {
			listTypeToKeys[listType] = append(listTypeToKeys[listType], sets.NewString(listElem.GetKeys()...))
		}

		value := reflect.ValueOf(topLevelListContainer)

		var attributeValue reflect.Value
		// toplevelListContainers can contain either a pointer or a struct and needs to be safeguarded when using reflect
		if value.Kind() == reflect.Ptr {
			attributeValue = value.Elem().FieldByName("Attributes")
		} else {
			attributeValue = value.FieldByName("Attributes")
		}

		if attributeValue.CanInterface() {
			attributes := attributeValue.Interface().(attributesPkg.Attributes)
			var attributeKeys []string
			for k := range attributes {
				attributeKeys = append(attributeKeys, k)
			}
			listTypeToKeys["Attributes"] = append(listTypeToKeys["Attributes"], sets.NewString(attributeKeys...))
		}
	}

	for listType, keySets := range listTypeToKeys {
		errors = multierror.Append(errors, doCheck(listType, keySets)...)
	}
	return errors.ErrorOrNil()
}
