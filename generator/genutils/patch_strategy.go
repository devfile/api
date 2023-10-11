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

package genutils

import (
	"strings"

	"sigs.k8s.io/controller-tools/pkg/markers"
)

const (
	patchStrategyTagKey = "patchStrategy"
	patchMergeKeyTagKey = "patchMergeKey"
	// MergePatchStrategy is the name of the Merge patch strategy
	MergePatchStrategy = "merge"
	// ReplacePatchStrategy is the name of the Replace patch strategy
	ReplacePatchStrategy = "replace"
)

// ContainsPatchStrategy reads the field tags to check whether the given patch strategy is defined
func ContainsPatchStrategy(field *markers.FieldInfo, strategy string) bool {
	patchStrategy := field.Tag.Get(patchStrategyTagKey)
	if patchStrategy == "" {
		return false
	}

	for _, s := range strings.Split(patchStrategy, ",") {
		if s == strategy {
			return true
		}
	}
	return false
}

// GetPatchMergeKey reads the field tags to retrieve the patch merge key. It returns nil if no such key is defined
func GetPatchMergeKey(field *markers.FieldInfo) string {
	return field.Tag.Get(patchMergeKeyTagKey)
}
