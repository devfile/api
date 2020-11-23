package v1alpha1

import (
	"encoding/json"
	"github.com/devfile/api/pkg/attributes"

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
	var srcAttributes map[string]string
	switch {
	case src.Exec != nil:
		srcAttributes = src.Exec.Attributes
	case src.Apply != nil:
		srcAttributes = src.Apply.Attributes
	case src.VscodeTask != nil:
		srcAttributes = src.VscodeTask.Attributes
	case src.VscodeLaunch != nil:
		srcAttributes = src.VscodeLaunch.Attributes
	case src.Composite != nil:
		srcAttributes = src.Composite.Attributes
	case src.Custom != nil:
		srcAttributes = src.Custom.Attributes
	}
	if srcAttributes != nil {
		dest.Attributes = attributes.Attributes{}
		convertAttributesTo_v1alpha2(srcAttributes, &dest.Attributes)
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
	if err != nil {
		return err
	}
	var destAttributes map[string]string
	if src.Attributes != nil {
		destAttributes = make(map[string]string)
		err = convertAttributesFrom_v1alpha2(&src.Attributes, destAttributes)
		if err != nil {
			return err
		}
	}

	switch {
	case dest.Apply != nil:
		dest.Apply.Attributes = destAttributes
		dest.Apply.Id = id
	case dest.Composite != nil:
		dest.Composite.Attributes = destAttributes
		dest.Composite.Id = id
	case dest.Custom != nil:
		dest.Custom.Attributes = destAttributes
		dest.Custom.Id = id
	case dest.Exec != nil:
		dest.Exec.Attributes = destAttributes
		dest.Exec.Id = id
	case dest.VscodeLaunch != nil:
		dest.VscodeLaunch.Attributes = destAttributes
		dest.VscodeLaunch.Id = id
	case dest.VscodeTask != nil:
		dest.VscodeTask.Attributes = destAttributes
		dest.VscodeTask.Id = id
	}
	return err
}
