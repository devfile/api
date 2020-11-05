package v1alpha1

import (
	"github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	"sigs.k8s.io/yaml"
)

func convertCommandTo_v1alpha2(src *Command, dest *v1alpha2.Command) error {
	id, err := src.Key()
	if err != nil {
		return err
	}
	yamlCommand, err := yaml.Marshal(src)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlCommand, dest)
	if err != nil {
		return err
	}
	dest.Id = id
	return nil
}

func convertCommandFrom_v1alpha2(src *v1alpha2.Command, dest *Command) error {
	id := src.Key()
	yamlCommand, err := yaml.Marshal(src)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlCommand, dest)
	switch {
	case dest.Apply != nil:
		dest.Apply.Id = id
	case dest.Composite != nil:
		dest.Composite.Id = id
	case dest.Custom != nil:
		dest.Custom.Id = id
	case dest.Exec != nil:
		dest.Exec.Id = id
	case dest.VscodeLaunch != nil:
		dest.VscodeLaunch.Id = id
	case dest.VscodeTask != nil:
		dest.VscodeTask.Id = id
	}
	return nil
}
