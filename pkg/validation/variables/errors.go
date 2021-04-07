package variables

import (
	"fmt"
	"strings"
)

// InvalidKeysError returns an error for the invalid keys
type InvalidKeysError struct {
	Keys []string
}

func (e *InvalidKeysError) Error() string {
	return fmt.Sprintf("invalid variable references - %s", strings.Join(e.Keys, ","))
}
