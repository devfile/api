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

package devfile

import (
	"testing"

	attributes "github.com/devfile/api/v2/pkg/attributes"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestDecodeAttribute(t *testing.T) {
	devfileWithEnhancedHeader := `
schemaVersion: 2.0.0
metadata:
  name: "theName"
  version: "1.0.0"
  attributes:
    stringAttribute: stringValue
    boolAttribute: true
    numberAttribute: 9.9
    objectAttribute:
      stringField: stringFieldValue
      numberField: 8.8
    arrayAttribute:
      - number1
      - 7.7
`
	header := DevfileHeader{}
	err := yaml.Unmarshal([]byte(devfileWithEnhancedHeader), &header)

	assert.NoError(t, err)

	assert.Equal(t, "theName", header.Metadata.Name)
	assert.Equal(t, "1.0.0", header.Metadata.Version)
	assert.Equal(t, "stringValue", attributes.Attributes(header.Metadata.Attributes).GetString("stringAttribute", nil))
	assert.Equal(t, true, attributes.Attributes(header.Metadata.Attributes).GetBoolean("boolAttribute", nil))
	assert.Equal(t, 9.9, attributes.Attributes(header.Metadata.Attributes).GetNumber("numberAttribute", nil))
	assert.Equal(t, map[string]interface{}{
		"stringField": "stringFieldValue",
		"numberField": 8.8,
	}, attributes.Attributes(header.Metadata.Attributes).Get("objectAttribute", nil))
	assert.Equal(t, []interface{}{
		"number1",
		7.7,
	}, attributes.Attributes(header.Metadata.Attributes).Get("arrayAttribute", nil))
}
