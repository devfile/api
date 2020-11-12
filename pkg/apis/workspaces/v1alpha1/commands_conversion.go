package v1alpha1

import (
	"encoding/json"

	"github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
)

func convertCommandTo_v1alpha2(src *Command, dest *v1alpha2.Command) error {
	id, err := src.Key()
	if err != nil {
		return err
	}
	jsonCommand, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonCommand, dest)
	if err != nil {
		return err
	}
	dest.Id = id
	return nil
}

func convertCommandFrom_v1alpha2(src *v1alpha2.Command, dest *Command) error {
	id := src.Key()
	jsonCommand, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonCommand, dest)
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
