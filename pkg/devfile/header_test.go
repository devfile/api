package devfile

import (
	"testing"

	attributes "github.com/devfile/api/pkg/attributes"
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
	assert.Equal(t, "stringValue", attributes.Attributes(header.Metadata.Attributes).GetString("stringAttribute"))
	assert.Equal(t, true, attributes.Attributes(header.Metadata.Attributes).GetBoolean("boolAttribute"))
	assert.Equal(t, 9.9, attributes.Attributes(header.Metadata.Attributes).GetNumber("numberAttribute"))
	assert.Equal(t, map[string]interface{}{
		"stringField": "stringFieldValue",
		"numberField": 8.8,
	}, attributes.Attributes(header.Metadata.Attributes).Get("objectAttribute"))
	assert.Equal(t, []interface{}{
		"number1",
		7.7,
	}, attributes.Attributes(header.Metadata.Attributes).Get("arrayAttribute"))
}
