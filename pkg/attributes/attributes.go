package v1alpha2

import (
	"encoding/json"

	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type Attribute apiext.JSON

// Attributes provides a way to add a map of arbitrary YAML/JSON
// objects.
type Attributes map[string]Attribute

func (attribute Attribute) Decode() interface{} {
	var result interface{} = nil

	json.Unmarshal(attribute.Raw, result)
	return result
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
		if err != nil {
			attribute.Raw = rawJson
		}
		return attribute, err
	}
}

type Builder map[string]attributeBuilder

func (builder Builder) Build() (Attributes, error) {

}

func (builder Builder) LenientBuild() Attributes {

}
