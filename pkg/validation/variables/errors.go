package variables

import (
	"fmt"
	"sort"
	"strings"
)

// InvalidKeysError returns an error for the invalid keys
type InvalidKeysError struct {
	Keys []string
}

func (e *InvalidKeysError) Error() string {
	return fmt.Sprintf("invalid variable references - %s", strings.Join(e.Keys, ","))
}

// processInvalidKeys processes the invalid keys and return InvalidKeysError if present
func processInvalidKeys(invalidKeys map[string]bool) error {
	var invalidKeysArr []string
	for key := range invalidKeys {
		invalidKeysArr = append(invalidKeysArr, key)
	}

	if len(invalidKeysArr) > 0 {
		sort.Strings(invalidKeysArr)
		return &InvalidKeysError{Keys: invalidKeysArr}
	}

	return nil
}
