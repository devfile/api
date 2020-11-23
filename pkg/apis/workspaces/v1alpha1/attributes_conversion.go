package v1alpha1

import (
	"fmt"
	"github.com/devfile/api/pkg/attributes"
)

func convertAttributesTo_v1alpha2(src map[string]string, dest *attributes.Attributes) {
	dest.FromStringMap(src)
}

func convertAttributesFrom_v1alpha2(src *attributes.Attributes, dest map[string]string) error {
	if dest == nil {
		return fmt.Errorf("trying to insert into a nil map")
	}
	var err error
	stringAttributes := src.Strings(&err)
	if err != nil {
		return err
	}
	for k, v := range stringAttributes {
		dest[k] = v
	}
	return nil
}

func getCommandAttributes(command *Command) map[string]string {
	switch {
	case command.Exec != nil:
		return command.Exec.Attributes
	case command.Apply != nil:
		return command.Apply.Attributes
	case command.VscodeTask != nil:
		return command.VscodeTask.Attributes
	case command.VscodeLaunch != nil:
		return command.VscodeLaunch.Attributes
	case command.Composite != nil:
		return command.Composite.Attributes
	case command.Custom != nil:
		return command.Custom.Attributes
	}
	return nil
}
