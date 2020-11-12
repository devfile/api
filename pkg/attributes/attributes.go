package attributes

import (
	"encoding/json"

	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type Attribute apiext.JSON

func (a Attribute) MarshalJSON() ([]byte, error) {
	return apiext.JSON(a).MarshalJSON()
}

func (a *Attribute) UnmarshalJSON(data []byte) error {
	return (*apiext.JSON)(a).UnmarshalJSON(data)
}

func (attr Attribute) DeepCopy() *Attribute {
	resultBytes := make([]byte, len(attr.Raw))
	copy(resultBytes, attr.Raw)
	return &Attribute {
		Raw: resultBytes,
	}
}

func (_ Attribute) OpenAPISchemaType() []string {
	return nil
}

func (_ Attribute) OpenAPISchemaFormat() string { return "" }


// Attributes provides a way to add a map of arbitrary YAML/JSON
// objects.
type Attributes map[string]apiext.JSON

func (attributes Attributes) Get(key string) Attribute {
	return Attribute(attributes[key])
}

func (attributes Attributes) Exists(key string) bool {
	_, exists := attributes[key]
	return exists
}

func (attribute Attribute) Interface() interface{} {
	container := &[]interface{}{}
	err := json.Unmarshal([]byte("[ " + string(attribute.Raw) + " ]"), container)
	print(err)
	return (*container)[0]
}

func (attribute Attribute) DecodeInto(into interface{}) error {
	err := json.Unmarshal(attribute.Raw, into)
	return err
}

type attributeBuilder func() (Attribute, error)

func FromInterface(value interface{}) attributeBuilder {
	return func() (Attribute, error) {
		attribute := Attribute{}
		rawJson, err := json.Marshal(value)
		if err == nil {
			attribute.Raw = rawJson
		}
		return attribute, err
	}
}

type Builder map[string]attributeBuilder

func (builder Builder) Build() (Attributes, error) {
	result := Attributes{}
	for key, buildAttribute := range builder {
		attribute, err := buildAttribute()
		if err != nil {
			return Attributes{}, err
		}
		result[key] = apiext.JSON(attribute)
	}
  return result, nil
}

func (builder Builder) LenientBuild() Attributes {
	result := Attributes{}
	for key, buildAttribute := range builder {
		attribute, err := buildAttribute()
		if err != nil {
			result[key] = apiext.JSON(Attribute{})
		} else {
			result[key] = apiext.JSON(attribute)
		}
	}
  return result
}
