package v1alpha1

import (
	"encoding/json"

	"github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
)

func convertComponentTo_v1alpha2(src *Component, dest *v1alpha2.Component) error {
	if src.Plugin != nil {
		// Need to handle plugin components separately.
		return convertPluginComponentTo_v1alpha2(src, dest)
	}
	name, err := src.Key()
	if err != nil {
		return err
	}
	jsonComponent, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonComponent, dest)
	if err != nil {
		return err
	}
	dest.Name = name
	return nil
}

func convertComponentFrom_v1alpha2(src *v1alpha2.Component, dest *Component) error {
	if src.Plugin != nil {
		// Need to handle plugin components separately.
		return convertPluginComponentFrom_v1alpha2(src, dest)
	}
	name := src.Key()
	jsonComponent, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonComponent, dest)
	switch {
	case dest.Container != nil:
		dest.Container.Name = name
	case dest.Plugin != nil:
		dest.Plugin.Name = name
	case dest.Volume != nil:
		dest.Volume.Name = name
	case dest.Openshift != nil:
		dest.Openshift.Name = name
	case dest.Kubernetes != nil:
		dest.Kubernetes.Name = name
	case dest.Custom != nil:
		dest.Custom.Name = name
	}
	return nil
}
