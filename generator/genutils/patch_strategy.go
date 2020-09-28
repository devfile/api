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
