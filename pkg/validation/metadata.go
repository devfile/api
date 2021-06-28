package validation

import (
	"fmt"
	"github.com/devfile/api/v2/pkg/devfile"
	"strings"
)

var validArchitectures = map[string]bool{
	"amd64":   true,
	"arm64":   true,
	"ppc64le": true,
	"s390x":   true,
}

// ValidateMetadata validates the devfile metadata
func ValidateMetadata(metadata devfile.DevfileMetadata) (err error) {

	if len(metadata.Architectures) > 0 {
		err = validateArchitectures(metadata.Architectures)
	}

	return err
}

// validateArchitectures validates the architectures property to ensure that the architectures
// mentioned conform to the architecture convention for container images
func validateArchitectures(architectures []string) error {

	var err error
	var invalidArchitectures []string
	for _, arch := range architectures {
		if ok := validArchitectures[arch]; !ok {
			invalidArchitectures = append(invalidArchitectures, arch)
		}
	}

	if len(invalidArchitectures) > 0 {
		err = fmt.Errorf("architecture: %s not valid. Please ensure that the architecture list conforms to specification", strings.Join(invalidArchitectures, ","))
	}

	return err
}
