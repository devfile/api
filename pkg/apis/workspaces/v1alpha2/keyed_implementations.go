//
// Copyright (c) 2019-2020 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//

package v1alpha2

import (
	"reflect"
)

func extractKeys(keyedList interface{}) []Keyed {
	value := reflect.ValueOf(keyedList)
	keys := make([]Keyed, 0, value.Len())
	for i := 0; i < value.Len(); i++ {
		elem := value.Index(i)
		if elem.CanInterface() {
			i := elem.Interface()
			if keyed, ok := i.(Keyed); ok {
				keys = append(keys, keyed)
			}
		}
	}
	return keys
}
