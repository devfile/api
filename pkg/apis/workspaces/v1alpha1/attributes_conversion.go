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
	"fmt"

	"github.com/devfile/api/v2/pkg/attributes"
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
	case command.Composite != nil:
		return command.Composite.Attributes
	case command.Custom != nil:
		return command.Custom.Attributes
	}
	return nil
}
