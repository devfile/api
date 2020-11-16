package attributes

import (
	"encoding/json"
	"errors"

	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// Attributes provides a way to add a map of arbitrary YAML/JSON
// objects.
type Attributes map[string]apiext.JSON

func (attributes Attributes) GetDecodedInto(key string, into interface{}) error {
	var err error
	if attribute, exists := attributes[key]; exists {
		err = json.Unmarshal(attribute.Raw, into)
	} else {
		err = errors.New("Key '" + key + "' doesn't exist")
	}
	return err
}

func (attributes Attributes) Get(key string, errorHolder ...*error) interface{} {
	if attribute, exists := attributes[key]; exists {
		container := &[]interface{}{}
		err := json.Unmarshal([]byte("[ "+string(attribute.Raw)+" ]"), container)
		if err != nil && len(errorHolder) > 0 && errorHolder != nil {
			*errorHolder[0] = err
			return nil
		}
		return (*container)[0]
	}
	return nil
}

func (attributes Attributes) GetString(key string, errorHolder ...*error) string {
	if attribute, exists := attributes[key]; exists {
		result := new(string)
		err := json.Unmarshal(attribute.Raw, result)
		if err == nil {
			return *result
		}
		if len(errorHolder) > 0 && errorHolder != nil {
			*errorHolder[0] = err
		}
	}
	return ""
}

func (attributes Attributes) GetNumber(key string, errorHolder ...*error) float64 {
	if attribute, exists := attributes[key]; exists {
		result := new(float64)
		err := json.Unmarshal(attribute.Raw, result)
		if err == nil {
			return *result
		}
		if len(errorHolder) > 0 && errorHolder != nil {
			*errorHolder[0] = err
		}
	}
	return 0
}

func (attributes Attributes) GetBoolean(key string, errorHolder ...*error) bool {
	if attribute, exists := attributes[key]; exists {
		result := new(bool)
		err := json.Unmarshal(attribute.Raw, result)
		if err == nil {
			return *result
		}
		if len(errorHolder) > 0 && errorHolder != nil {
			*errorHolder[0] = err
		}
	}
	return false
}

func (attributes Attributes) Exists(key string) bool {
	_, exists := attributes[key]
	return exists
}

func (attributes Attributes) Strings(errorHolder ...*error) map[string]string {
	result := map[string]string{}
	for key := range attributes {
		if value, isRightType := attributes.Get(key, errorHolder...).(string); isRightType {
			result[key] = value
		}
	}
	return result
}

func (attributes Attributes) Numbers(errorHolder ...*error) map[string]float64 {
	result := map[string]float64{}
	for key := range attributes {
		if value, isRightType := attributes.Get(key, errorHolder...).(float64); isRightType {
			result[key] = value
		}
	}
	return result
}

func (attributes Attributes) Booleans(errorHolder ...*error) map[string]bool {
	result := map[string]bool{}
	for key := range attributes {
		if value, isRightType := attributes.Get(key, errorHolder...).(bool); isRightType {
			result[key] = value
		}
	}
	return result
}

func (attributes Attributes) DecodeInto(into interface{}) error {
	rawJson, err := json.Marshal(attributes)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rawJson, into)
	return err
}

func (attributes Attributes) Interface(errorHolder ...*error) interface{} {
	rawJson, err := json.Marshal(attributes)
	if err != nil && len(errorHolder) > 0 && errorHolder != nil {
		*errorHolder[0] = err
		return nil
	}

	container := &[]interface{}{}
	err = json.Unmarshal([]byte("[ "+string(rawJson)+" ]"), container)
	if err != nil && len(errorHolder) > 0 && errorHolder != nil {
		*errorHolder[0] = err
		return nil
	}

	return (*container)[0]
}

func (attributes Attributes) PutString(key string, value string) Attributes {
	rawJSON, _ := json.Marshal(value)
	attributes[key] = apiext.JSON{
		Raw: rawJSON,
	}
	return attributes
}

func (attributes Attributes) FromStringMap(strings map[string]string) Attributes {
	for key, value := range strings {
		attributes.PutString(key, value)
	}
	return attributes
}

func (attributes Attributes) PutFloat(key string, value float64) Attributes {
	rawJSON, _ := json.Marshal(value)
	attributes[key] = apiext.JSON{
		Raw: rawJSON,
	}
	return attributes
}

func (attributes Attributes) FromFloatMap(strings map[string]float64) Attributes {
	for key, value := range strings {
		attributes.PutFloat(key, value)
	}
	return attributes
}

func (attributes Attributes) PutInteger(key string, value int) Attributes {
	rawJSON, _ := json.Marshal(value)
	attributes[key] = apiext.JSON{
		Raw: rawJSON,
	}
	return attributes
}

func (attributes Attributes) FromIntegerMap(strings map[string]int) Attributes {
	for key, value := range strings {
		rawJSON, _ := json.Marshal(value)
		attributes[key] = apiext.JSON{
			Raw: rawJSON,
		}
	}
	return attributes
}

func (attributes Attributes) PutBoolean(key string, value bool) Attributes {
	rawJSON, _ := json.Marshal(value)
	attributes[key] = apiext.JSON{
		Raw: rawJSON,
	}
	return attributes
}

func (attributes Attributes) FromBooleanMap(strings map[string]bool) Attributes {
	for key, value := range strings {
		rawJSON, _ := json.Marshal(value)
		attributes[key] = apiext.JSON{
			Raw: rawJSON,
		}
	}
	return attributes
}

func (attributes Attributes) Put(key string, value interface{}, errorHolder ...*error) Attributes {
	rawJSON, err := json.Marshal(value)
	if err != nil && len(errorHolder) > 0 && errorHolder != nil {
		*errorHolder[0] = err
	}

	attributes[key] = apiext.JSON{
		Raw: rawJSON,
	}
	return attributes
}

func (attributes Attributes) FromMap(strings map[string]interface{}, errorHolder ...*error) Attributes {
	for key, value := range strings {
		attributes.Put(key, value, errorHolder...)
	}
	return attributes
}

func (attributes Attributes) FromInterface(structure interface{}, errorHolder ...*error) Attributes {
	completeJSON, err := json.Marshal(structure)
	if err != nil && len(errorHolder) > 0 && errorHolder != nil {
		*errorHolder[0] = err
	}

	err = json.Unmarshal(completeJSON, &attributes)
	return attributes
}
