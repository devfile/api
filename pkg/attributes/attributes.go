package attributes

import (
	"encoding/json"
	"errors"

	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// Attributes provides a way to add a map of arbitrary YAML/JSON
// objects.
// +kubebuilder:validation:Type=object
// +kubebuilder:validation:XPreserveUnknownFields
type Attributes map[string]apiext.JSON

// MarshalJSON implements custom JSON marshaling
// to support free-form attributes
func (attributes Attributes) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]apiext.JSON(attributes))
}

// UnmarshalJSON implements custom JSON unmarshaling
// to support free-form attributes
func (attributes *Attributes) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, (*map[string]apiext.JSON)(attributes))
}

// Exists returns `true` if the attribute with the given key
// exists in the attributes map.
func (attributes Attributes) Exists(key string) bool {
	_, exists := attributes[key]
	return exists
}

// GetString allows returning the attribute with the given key
// as a string. If the attribute JSON/YAML content is
// not a JSON string, then the result will be the empty string
// and an error is raised.
//
// An optional error holder can be passed as an argument
// to receive any error that might have be raised during the attribute
// decoding
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

// GetNumber allows returning the attribute with the given key
// as a float64. If the attribute JSON/YAML content is
// not a JSON number, then the result will be the zero value
// and an error is raised.
//
// An optional error holder can be passed as an argument
// to receive any error that might have be raised during the attribute
// decoding
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

// GetBoolean allows returning the attribute with the given key
// as a bool. If the attribute JSON/YAML content is
// not a JSON boolean, then the result will be the `false` zero value
// and an error is raised.
//
// An optional error holder can be passed as an argument
// to receive any error that might have be raised during the attribute
// decoding
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

// Get allows returning the attribute with the given key
// as an interface. The underlying type of the returned interface
// depends on the JSON/YAML content of the attribute. It can be either a simple type
// like a string, a float64 or a bool, either a structured type like
// a map of interfaces or an array of interfaces.
//
// An optional error holder can be passed as an argument
// to receive any error that might have occurred during the attribute
// decoding
func (attributes Attributes) Get(key string, errorHolder ...*error) interface{} {
	if attribute, exists := attributes[key]; exists {
		container := &[]interface{}{}
		err := json.Unmarshal([]byte("[ "+string(attribute.Raw)+" ]"), container)
		if err != nil && len(errorHolder) > 0 && errorHolder != nil {
			*errorHolder[0] = err
		}
		if len(*container) > 0 {
			return (*container)[0]
		}
	}
	return nil
}

// GetInto allows decoding the attribute with the given key
// into a given interface. The provided interface should be a pointer
// to a struct, to an array, or to any simple type.
//
// An error is returned if the provided interface type is not compatible
// with the attribute content
func (attributes Attributes) GetInto(key string, into interface{}) error {
	var err error
	if attribute, exists := attributes[key]; exists {
		err = json.Unmarshal(attribute.Raw, into)
	} else {
		err = errors.New("Key '" + key + "' doesn't exist")
	}
	return err
}

// Strings allows returning only the attributes whose content
// is a JSON string.
//
// An optional error holder can be passed as an argument
// to receive any error that might have be raised during the attribute
// decoding
func (attributes Attributes) Strings(errorHolder ...*error) map[string]string {
	result := map[string]string{}
	for key := range attributes {
		if value, isRightType := attributes.Get(key, errorHolder...).(string); isRightType {
			result[key] = value
		}
	}
	return result
}

// Numbers allows returning only the attributes whose content
// is a JSON number.
//
// An optional error holder can be passed as an argument
// to receive any error that might have be raised during the attribute
// decoding
func (attributes Attributes) Numbers(errorHolder ...*error) map[string]float64 {
	result := map[string]float64{}
	for key := range attributes {
		if value, isRightType := attributes.Get(key, errorHolder...).(float64); isRightType {
			result[key] = value
		}
	}
	return result
}

// Booleans allows returning only the attributes whose content
// is a JSON boolean.
//
// An optional error holder can be passed as an argument
// to receive any error that might have be raised during the attribute
// decoding
func (attributes Attributes) Booleans(errorHolder ...*error) map[string]bool {
	result := map[string]bool{}
	for key := range attributes {
		if value, isRightType := attributes.Get(key, errorHolder...).(bool); isRightType {
			result[key] = value
		}
	}
	return result
}

// Into allows decoding the whole attributes map
// into a given interface. The provided interface should be either a pointer
// to a struct, or to a map.
//
// An error is returned if the provided interface type is not compatible
// with the structure of the attributes
func (attributes Attributes) Into(into interface{}) error {
	rawJSON, err := json.Marshal(attributes)
	if err != nil {
		return err
	}

	err = json.Unmarshal(rawJSON, into)
	return err
}

// AsInterface allows returning the whole attributes map
// as an interface. When the attributes are not empty,
// the returned interface will be a map
// of interfaces.
//
// An optional error holder can be passed as an argument
// to receive any error that might have occured during the attributes
// decoding
func (attributes Attributes) AsInterface(errorHolder ...*error) interface{} {
	rawJSON, err := json.Marshal(attributes)
	if err != nil && len(errorHolder) > 0 && errorHolder != nil {
		*errorHolder[0] = err
		return nil
	}

	container := &[]interface{}{}
	err = json.Unmarshal([]byte("[ "+string(rawJSON)+" ]"), container)
	if err != nil && len(errorHolder) > 0 && errorHolder != nil {
		*errorHolder[0] = err
		return nil
	}

	return (*container)[0]
}

// PutString allows adding a string attribute to the
// current map of attributes
func (attributes Attributes) PutString(key string, value string) Attributes {
	rawJSON, _ := json.Marshal(value)
	attributes[key] = apiext.JSON{
		Raw: rawJSON,
	}
	return attributes
}

// FromStringMap allows adding into the current map of attributes all
// the attributes contained in the given string map
func (attributes Attributes) FromStringMap(strings map[string]string) Attributes {
	for key, value := range strings {
		attributes.PutString(key, value)
	}
	return attributes
}

// PutFloat allows adding a float attribute to the
// current map of attributes
func (attributes Attributes) PutFloat(key string, value float64) Attributes {
	rawJSON, _ := json.Marshal(value)
	attributes[key] = apiext.JSON{
		Raw: rawJSON,
	}
	return attributes
}

// FromFloatMap allows adding into the current map of attributes all
// the attributes contained in the given map of floats
func (attributes Attributes) FromFloatMap(strings map[string]float64) Attributes {
	for key, value := range strings {
		attributes.PutFloat(key, value)
	}
	return attributes
}

// PutInteger allows adding an integer attribute to the
// current map of attributes
func (attributes Attributes) PutInteger(key string, value int) Attributes {
	rawJSON, _ := json.Marshal(value)
	attributes[key] = apiext.JSON{
		Raw: rawJSON,
	}
	return attributes
}

// FromIntegerMap allows adding into the current map of attributes all
// the attributes contained in the given map of integers
func (attributes Attributes) FromIntegerMap(strings map[string]int) Attributes {
	for key, value := range strings {
		rawJSON, _ := json.Marshal(value)
		attributes[key] = apiext.JSON{
			Raw: rawJSON,
		}
	}
	return attributes
}

// PutBoolean allows adding a boolean attribute to the
// current map of attributes
func (attributes Attributes) PutBoolean(key string, value bool) Attributes {
	rawJSON, _ := json.Marshal(value)
	attributes[key] = apiext.JSON{
		Raw: rawJSON,
	}
	return attributes
}

// FromBooleanMap allows adding into the current map of attributes all
// the attributes contained in the given map of booleans
func (attributes Attributes) FromBooleanMap(strings map[string]bool) Attributes {
	for key, value := range strings {
		rawJSON, _ := json.Marshal(value)
		attributes[key] = apiext.JSON{
			Raw: rawJSON,
		}
	}
	return attributes
}

// Put allows adding an attribute to the
// current map of attributes.
// The attribute is provided as an interface, and can be any value
// that supports Json Marshaling.
//
// An optional error holder can be passed as an argument
// to receive any error that might have occured during the attributes
// decoding
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

// FromMap allows adding into the current map of attributes all
// the attributes contained in the given map of interfaces
// each attribute of the given map is provided as an interface, and can be any value
// that supports Json Marshaling.
//
// An optional error holder can be passed as an argument
// to receive any error that might have occured during the attributes
// decoding
func (attributes Attributes) FromMap(strings map[string]interface{}, errorHolder ...*error) Attributes {
	for key, value := range strings {
		attributes.Put(key, value, errorHolder...)
	}
	return attributes
}

// FromInterface allows completing the map of attributes from the given interface.
// The given interface, and can be any value
// that supports Json Marshaling and will be marshalled as a JSON object.
//
// This is especially useful to create attributes from well-known, but
// implementation- dependent Go structures.
//
// An optional error holder can be passed as an argument
// to receive any error that might have occured during the attributes
// decoding
func (attributes Attributes) FromInterface(structure interface{}, errorHolder ...*error) Attributes {
	newAttributes := Attributes{}
	completeJSON, err := json.Marshal(structure)
	if err != nil && len(errorHolder) > 0 && errorHolder != nil {
		*errorHolder[0] = err
	}

	err = json.Unmarshal(completeJSON, &newAttributes)
	for key, value := range newAttributes {
		attributes[key] = value
	}
	return attributes
}
