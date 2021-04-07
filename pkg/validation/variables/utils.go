package variables

// checkForInvalidError checks for InvalidKeysError and stores the key in the map
func checkForInvalidError(invalidKeys map[string]bool, err error) {
	if verr, ok := err.(*InvalidKeysError); ok {
		for _, key := range verr.Keys {
			invalidKeys[key] = true
		}
	}
}

// processInvalidKeys processes the invalid keys and return InvalidKeysError if present
func processInvalidKeys(invalidKeys map[string]bool) error {
	var invalidKeysArr []string
	for key := range invalidKeys {
		invalidKeysArr = append(invalidKeysArr, key)
	}

	if len(invalidKeysArr) > 0 {
		return &InvalidKeysError{Keys: invalidKeysArr}
	}

	return nil
}
