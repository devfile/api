package validation

import (
	"fmt"
	"strings"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	// EnvProjectsSrc is the env defined for path to the project source in a component container
	EnvProjectsSrc = "PROJECT_SOURCE"

	// EnvProjectsRoot is the env defined for project mount in a component container when component's mountSources=true
	EnvProjectsRoot = "PROJECTS_ROOT"
)

// ValidateComponents validates that the components
// 1. makes sure the container components reference a valid volume component if it uses volume mounts
// 2. makes sure the volume components are unique
// 3. checks the URI specified in openshift components and kubernetest components are with valid format
// 4. makes sure the component name is not a numeric string
// 5. makes sure the component name is unique
func ValidateComponents(components []v1alpha2.Component) error {

	processedVolumes := make(map[string]bool)
	processedVolumeMounts := make(map[string]bool)
	processedEndPointName := make(map[string]bool)
	processedEndPointPort := make(map[int]bool)
	componentNameMap := make(map[string]bool)

	for _, component := range components {
		if isInt(component.Name) {
			return &InvalidNameOrIdError{name:component.Name, resourceType: "component"}
		}
		if _,exists:= componentNameMap[component.Name]; exists {
			return &InvalidComponentError{name: component.Name}
		}
		componentNameMap[component.Name] = true


		if component.Container != nil {
			// Process all the volume mounts in container components to validate them later
			for _, volumeMount := range component.Container.VolumeMounts {
				if _, ok := processedVolumeMounts[volumeMount.Name]; !ok {
					processedVolumeMounts[volumeMount.Name] = true
				}
			}

			// Check if any containers are customizing the reserved PROJECT_SOURCE or PROJECTS_ROOT env
			for _, env := range component.Container.Env {
				if env.Name == EnvProjectsSrc {
					return &ReservedEnvError{envName: EnvProjectsSrc, componentName: component.Name}
				} else if env.Name == EnvProjectsRoot {
					return &ReservedEnvError{envName: EnvProjectsRoot, componentName: component.Name}
				}
			}

			// Check if all the endpoint names are unique across components
			// and check if endpoint port are unique across component containers ie;
			// two component containers cannot have the same target port but two endpoints
			// in a single component container can have the same target port
			err := validateEndpoints(component.Container.Endpoints, processedEndPointPort,processedEndPointName )
			if err != nil {
				return err
			}
		} else if component.Volume != nil {
			if _, ok := processedVolumes[component.Name]; !ok {
				processedVolumes[component.Name] = true
				if len(component.Volume.Size) > 0 {
					// Only validate on Kubernetes since Docker volumes do not use sizes
					// We use the Kube API for validation because there are so many ways to
					// express storage in Kubernetes. For reference, you may check doc
					// https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
					if _, err := resource.ParseQuantity(component.Volume.Size); err != nil {
						return &InvalidVolumeError{name: component.Name, reason: fmt.Sprintf("size %s for volume component is invalid, %v. Example - 2Gi, 1024Mi", component.Volume.Size, err)}
					}
				}
			} else {
				return &InvalidVolumeError{name: component.Name, reason: "duplicate volume components present with the same name"}
			}
		} else if component.Openshift != nil {
			if component.Openshift.Uri != "" {
				return ValidateURI( component.Openshift.Uri )
			}

			err := validateEndpoints(component.Openshift.Endpoints, processedEndPointPort,processedEndPointName )
			if err != nil {
				return err
			}
		}else if component.Kubernetes != nil {
			if component.Kubernetes.Uri != "" {
				return ValidateURI( component.Kubernetes.Uri )
			}
			err := validateEndpoints(component.Kubernetes.Endpoints, processedEndPointPort,processedEndPointName )
			if err != nil {
				return err
			}
		}

	}

	// Check if the volume mounts mentioned in the containers are referenced by a volume component
	var invalidVolumeMounts []string
	for volumeMountName := range processedVolumeMounts {
		if _, ok := processedVolumes[volumeMountName]; !ok {
			invalidVolumeMounts = append(invalidVolumeMounts, volumeMountName)
		}
	}

	if len(invalidVolumeMounts) > 0 {
		return &MissingVolumeMountError{volumeName: strings.Join(invalidVolumeMounts, ",")}
	}

	// Successful
	return nil
}
