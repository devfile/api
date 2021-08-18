package v1alpha1

import (
	"encoding/json"

	"github.com/devfile/api/v2/pkg/attributes"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
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

// getGroup returns the group the command belongs to
func getGroup(dc v1alpha2.Command) *v1alpha2.CommandGroup {
	switch {
	case dc.Composite != nil:
		return dc.Composite.Group
	case dc.Exec != nil:
		return dc.Exec.Group
	case dc.Apply != nil:
		return dc.Apply.Group
	case dc.Custom != nil:
		return dc.Custom.Group

	default:
		return nil
	}
}

func convertCommandFrom_v1alpha2(src *v1alpha2.Command, dest *Command) error {
	if src != nil {
		id := src.Key()

		srcCmdGroup := getGroup(*src)
		if srcCmdGroup != nil && srcCmdGroup.Kind == v1alpha2.DeployCommandGroupKind {
			// skip converting deploy kind commands as deploy kind commands are not supported in v1alpha1
			return nil
		}

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
		}
		return err
	}
	return nil
}
