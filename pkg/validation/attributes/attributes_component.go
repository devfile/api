package attributes

import (
	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
)

// ValidateAndReplaceForComponents validates the components data for global attribute references and replaces them with the attribute value
func ValidateAndReplaceForComponents(attributes apiAttributes.Attributes, components []v1alpha2.Component) error {

	for i := range components {
		var err error

		// Validate various component types
		switch {
		case components[i].Container != nil:
			if err = validateAndReplaceForContainerComponent(attributes, components[i].Container); err != nil {
				return err
			}
		case components[i].Kubernetes != nil:
			if err = validateAndReplaceForKubernetesComponent(attributes, components[i].Kubernetes); err != nil {
				return err
			}
		case components[i].Openshift != nil:
			if err = validateAndReplaceForOpenShiftComponent(attributes, components[i].Openshift); err != nil {
				return err
			}
		case components[i].Volume != nil:
			if err = validateAndReplaceForVolumeComponent(attributes, components[i].Volume); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateAndReplaceForContainerComponent validates the container component data for global attribute references and replaces them with the attribute value
func validateAndReplaceForContainerComponent(attributes apiAttributes.Attributes, container *v1alpha2.ContainerComponent) error {
	var err error

	if container != nil {
		// Validate container image
		if container.Image, err = validateAndReplaceDataWithAttribute(container.Image, attributes); err != nil {
			return err
		}

		// Validate container commands
		for i := range container.Command {
			if container.Command[i], err = validateAndReplaceDataWithAttribute(container.Command[i], attributes); err != nil {
				return err
			}
		}

		// Validate container args
		for i := range container.Args {
			if container.Args[i], err = validateAndReplaceDataWithAttribute(container.Args[i], attributes); err != nil {
				return err
			}
		}

		// Validate memory limit
		if container.MemoryLimit, err = validateAndReplaceDataWithAttribute(container.MemoryLimit, attributes); err != nil {
			return err
		}

		// Validate memory limit
		if container.MemoryRequest, err = validateAndReplaceDataWithAttribute(container.MemoryRequest, attributes); err != nil {
			return err
		}

		// Validate source mapping
		if container.SourceMapping, err = validateAndReplaceDataWithAttribute(container.SourceMapping, attributes); err != nil {
			return err
		}

		// Validate container env
		if len(container.Env) > 0 {
			if err = validateAndReplaceForEnv(attributes, container.Env); err != nil {
				return err
			}
		}

		// Validate container volume mounts
		for i := range container.VolumeMounts {
			if container.VolumeMounts[i].Name, err = validateAndReplaceDataWithAttribute(container.VolumeMounts[i].Name, attributes); err != nil {
				return err
			}
			if container.VolumeMounts[i].Path, err = validateAndReplaceDataWithAttribute(container.VolumeMounts[i].Path, attributes); err != nil {
				return err
			}
		}

		// Validate container endpoints
		if len(container.Endpoints) > 0 {
			if err = validateAndReplaceForEndpoint(attributes, container.Endpoints); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateAndReplaceForEnv validates the env data for global attribute references and replaces them with the attribute value
func validateAndReplaceForEnv(attributes apiAttributes.Attributes, env []v1alpha2.EnvVar) error {

	for i := range env {
		var err error

		// Validate env name
		if env[i].Name, err = validateAndReplaceDataWithAttribute(env[i].Name, attributes); err != nil {
			return err
		}

		// Validate env value
		if env[i].Value, err = validateAndReplaceDataWithAttribute(env[i].Value, attributes); err != nil {
			return err
		}
	}

	return nil
}

// validateAndReplaceForKubernetesComponent validates the kubernetes component data for global attribute references and replaces them with the attribute value
func validateAndReplaceForKubernetesComponent(attributes apiAttributes.Attributes, kubernetes *v1alpha2.KubernetesComponent) error {
	var err error

	if kubernetes != nil {
		// Validate kubernetes uri
		if kubernetes.Uri, err = validateAndReplaceDataWithAttribute(kubernetes.Uri, attributes); err != nil {
			return err
		}

		// Validate kubernetes inlined
		if kubernetes.Inlined, err = validateAndReplaceDataWithAttribute(kubernetes.Inlined, attributes); err != nil {
			return err
		}

		// Validate kubernetes endpoints
		if len(kubernetes.Endpoints) > 0 {
			if err = validateAndReplaceForEndpoint(attributes, kubernetes.Endpoints); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateAndReplaceForOpenShiftComponent validates the openshift component data for global attribute references and replaces them with the attribute value
func validateAndReplaceForOpenShiftComponent(attributes apiAttributes.Attributes, openshift *v1alpha2.OpenshiftComponent) error {
	var err error

	if openshift != nil {
		// Validate openshift uri
		if openshift.Uri, err = validateAndReplaceDataWithAttribute(openshift.Uri, attributes); err != nil {
			return err
		}

		// Validate openshift inlined
		if openshift.Inlined, err = validateAndReplaceDataWithAttribute(openshift.Inlined, attributes); err != nil {
			return err
		}

		// Validate openshift endpoints
		if len(openshift.Endpoints) > 0 {
			if err = validateAndReplaceForEndpoint(attributes, openshift.Endpoints); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateAndReplaceForVolumeComponent validates the volume component data for global attribute references and replaces them with the attribute value
func validateAndReplaceForVolumeComponent(attributes apiAttributes.Attributes, volume *v1alpha2.VolumeComponent) error {
	var err error

	if volume != nil {
		// Validate volume size
		if volume.Size, err = validateAndReplaceDataWithAttribute(volume.Size, attributes); err != nil {
			return err
		}
	}

	return nil
}
