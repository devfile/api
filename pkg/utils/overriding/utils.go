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

package overriding

import (
	"k8s.io/apimachinery/pkg/util/json"
)

func handleUnmarshal(j []byte) (map[string]interface{}, error) {
	if j == nil {
		j = []byte("{}")
	}

	m := map[string]interface{}{}
	err := json.Unmarshal(j, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// UnionStrings returns the union of two string slices, preserving the order
// of first occurrence and removing duplicates.
//
// Elements from the first slice `a` are added in order. Elements from the second
// slice `b` are appended only if they do not already exist in `a`.
//
// The returned slice is guaranteed to be non-nil, even if both input slices are nil.
//
// Example:
//
//	UnionStrings([]string{"a", "b"}, []string{"b", "c"})
//	// Returns: []string{"a", "b", "c"}
func UnionStrings(a, b []string) []string {
	alreadyExistsInList := make(map[string]bool)
	result := make([]string, 0)

	for _, s := range a {
		if !alreadyExistsInList[s] {
			alreadyExistsInList[s] = true
			result = append(result, s)
		}
	}
	for _, s := range b {
		if !alreadyExistsInList[s] {
			alreadyExistsInList[s] = true
			result = append(result, s)
		}
	}
	return result
}
