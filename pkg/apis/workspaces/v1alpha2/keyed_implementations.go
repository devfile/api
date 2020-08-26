package v1alpha1

import (
	"reflect"
)

func extractKeys(keyedList interface{}) []string {
	value := reflect.ValueOf(keyedList)
	keys := make([]string, 0, value.Len())
	for i := 0; i<value.Len(); i++ {
		elem := value.Index(i)
		if elem.CanInterface() {
			i := elem.Interface()
			if keyed, ok := i.(Keyed); ok {
				keys = append(keys, keyed.Key())
			}
		}
	}
	return keys
}
