//
//
// Copyright Red Hat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	if src == nil {
		return nil
	}

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
