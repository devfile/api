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

package attributes

import (
	//	"encoding/json"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type buildAttributesTestCase struct {
	name           string
	builder        func(errorHolder *error) Attributes
	expectedResult interface{}
	expectedError  string
}

var buildAttributesTestCases []buildAttributesTestCase = []buildAttributesTestCase{
	{
		name: "FromStringMap",
		builder: func(*error) Attributes {
			return Attributes{}.FromStringMap(map[string]string{
				"field1": "value1",
				"field2": "value2",
			})
		},
		expectedResult: &struct {
			Field1 string `json:"field1"`
			Field2 string `json:"field2"`
		}{Field1: "value1", Field2: "value2"},
	},
	{
		name: "FromBooleanMap",
		builder: func(*error) Attributes {
			return Attributes{}.FromBooleanMap(map[string]bool{
				"field1": true,
				"field2": false,
			})
		},
		expectedResult: &struct {
			Field1 bool `json:"field1"`
			Field2 bool `json:"field2"`
		}{Field1: true, Field2: false},
	},
	{
		name: "FromFloatMap",
		builder: func(*error) Attributes {
			return Attributes{}.FromFloatMap(map[string]float64{
				"field1": 9.9,
				"field2": 10,
			})
		},
		expectedResult: &struct {
			Field1 float64 `json:"field1"`
			Field2 float64 `json:"field2"`
		}{Field1: 9.9, Field2: 10},
	},
	{
		name: "FromIntegerMap",
		builder: func(*error) Attributes {
			return Attributes{}.FromIntegerMap(map[string]int{
				"field1": 9,
				"field2": 10,
			})
		},
		expectedResult: &struct {
			Field1 float64 `json:"field1"`
			Field2 float64 `json:"field2"`
		}{Field1: 9, Field2: 10},
	},
	{
		name: "FromMap / simple types / No errors",
		builder: func(errorHolder *error) Attributes {
			return Attributes{}.FromMap(map[string]interface{}{
				"boolAttribute":   true,
				"numberAttribute": 12,
				"stringAttribute": "stringValue",
			}, errorHolder)
		},
		expectedResult: &struct {
			BoolAttribute   bool   `json:"boolAttribute"`
			NumberAttribute int    `json:"numberAttribute"`
			StringAttribute string `json:"stringAttribute"`
		}{BoolAttribute: true, StringAttribute: "stringValue", NumberAttribute: 12},
	},
	{
		name: "FromMap / invalid attributes",
		builder: func(errorHolder *error) Attributes {
			return Attributes{}.FromMap(map[string]interface{}{
				"boolAttribute":   true,
				"numberAttribute": 12,
				"wrongAttribute":  make(chan int),
			}, errorHolder)
		},
		expectedResult: &struct {
			BoolAttribute   bool        `json:"boolAttribute"`
			NumberAttribute int         `json:"numberAttribute"`
			WrongAttribute  interface{} `json:"wrongAttribute"`
		}{BoolAttribute: true, NumberAttribute: 12, WrongAttribute: nil},
		expectedError: "json: unsupported type: chan int",
	},
	{
		name: "FromMap / structured types",
		builder: func(errorHolder *error) Attributes {
			return Attributes{}.FromMap(map[string]interface{}{
				"arrayAttribute":  []string{"stringOne, stringTwo"},
				"objectAttribute": struct{ SubField int }{SubField: 10},
			}, errorHolder)
		},
		expectedResult: &struct {
			ArrayAttribute  []string               `json:"arrayAttribute"`
			ObjectAttribute struct{ SubField int } `json:"objectAttribute"`
		}{ObjectAttribute: struct{ SubField int }{SubField: 10}, ArrayAttribute: []string{"stringOne, stringTwo"}},
	},
	{
		name: "Put",
		builder: func(*error) Attributes {
			return Attributes{}.FromStringMap(map[string]string{
				"field1": "value1",
			}).Put("field2", map[string]string{"subfield1": "subvalue1", "subfield2": "subvalue2"}, nil)
		},
		expectedResult: &struct {
			Field1 string                 `json:"field1"`
			Field2 map[string]interface{} `json:"field2"`
		}{Field1: "value1", Field2: map[string]interface{}{"subfield1": "subvalue1", "subfield2": "subvalue2"}},
	},
	{
		name: "PutString",
		builder: func(*error) Attributes {
			return Attributes{}.FromStringMap(map[string]string{
				"field1": "value1",
			}).PutString("field2", "value2")
		},
		expectedResult: &struct {
			Field1 string `json:"field1"`
			Field2 string `json:"field2"`
		}{Field1: "value1", Field2: "value2"},
	},
	{
		name: "PutBoolean",
		builder: func(*error) Attributes {
			return Attributes{}.FromStringMap(map[string]string{
				"field1": "value1",
			}).PutBoolean("field2", true)
		},
		expectedResult: &struct {
			Field1 string `json:"field1"`
			Field2 bool   `json:"field2"`
		}{Field1: "value1", Field2: true},
	},
	{
		name: "PutFloat",
		builder: func(*error) Attributes {
			return Attributes{}.FromStringMap(map[string]string{
				"field1": "value1",
			}).PutFloat("field2", 9.9)
		},
		expectedResult: &struct {
			Field1 string  `json:"field1"`
			Field2 float64 `json:"field2"`
		}{Field1: "value1", Field2: 9.9},
	},
	{
		name: "PutInteger",
		builder: func(*error) Attributes {
			return Attributes{}.FromStringMap(map[string]string{
				"field1": "value1",
			}).PutInteger("field2", 9)
		},
		expectedResult: &struct {
			Field1 string  `json:"field1"`
			Field2 float64 `json:"field2"`
		}{Field1: "value1", Field2: 9},
	},
	{
		name: "FromInterface / Struct",
		builder: func(*error) Attributes {
			return Attributes{}.FromInterface(struct {
				Field1 string  `json:"field1"`
				Field2 float64 `json:"field2"`
			}{Field1: "value1", Field2: 9.9}, nil)
		},
		expectedResult: &struct {
			Field1 string  `json:"field1"`
			Field2 float64 `json:"field2"`
		}{Field1: "value1", Field2: 9.9},
	},
	{
		name: "FromInterface / Struct Pointer",
		builder: func(*error) Attributes {
			return Attributes{}.FromInterface(&struct {
				Field1 string  `json:"field1"`
				Field2 float64 `json:"field2"`
			}{Field1: "value1", Field2: 9.9}, nil)
		},
		expectedResult: &struct {
			Field1 string  `json:"field1"`
			Field2 float64 `json:"field2"`
		}{Field1: "value1", Field2: 9.9},
	},
}

func TestBuildAttributes(t *testing.T) {
	for _, test := range buildAttributesTestCases {
		t.Run(test.name, func(t *testing.T) {
			var attributes Attributes
			var err error

			attributes = test.builder(&err)
			if test.expectedError != "" {
				assert.EqualError(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}

			expectedJson, _ := json.Marshal(test.expectedResult)
			actualJson, _ := json.Marshal(attributes)
			assert.Equal(t, string(expectedJson), string(actualJson))
		})
	}
}

type decodeAttributeTestCase struct {
	name                    string
	attributeKey            string
	attributeJson           string
	expectedNumber          float64
	expectedNumberError     string
	expectedString          string
	expectedStringError     string
	expectedBool            bool
	expectedBoolError       string
	expectedInterface       interface{}
	expectedInterfaceError  string
	decodeInto              interface{}
	decodeIntoExpectedValue interface{}
	decodeIntoError         string
}

var invalidKey = "randomKey"

var keyNotFoundErr error = &KeyNotFoundError{Key: invalidKey}

var decodeAttributeTestCases []decodeAttributeTestCase = []decodeAttributeTestCase{
	{
		name:                    "DecodeSimpleString",
		attributeKey:            "test",
		attributeJson:           `"simpleString"`,
		expectedInterface:       "simpleString",
		expectedString:          "simpleString",
		expectedBoolError:       "json: cannot unmarshal string into Go value of type bool",
		expectedNumberError:     "json: cannot unmarshal string into Go value of type float64",
		decodeInto:              new(string),
		decodeIntoExpectedValue: "simpleString",
	},
	{
		name:                    "DecodeSimpleString / true",
		attributeKey:            "test",
		attributeJson:           `"true"`,
		expectedInterface:       "true",
		expectedString:          "true",
		expectedBool:            true,
		expectedNumberError:     "json: cannot unmarshal string into Go value of type float64",
		decodeInto:              new(string),
		decodeIntoExpectedValue: "true",
	},
	{
		name:                    "DecodeSimpleString / false",
		attributeKey:            "test",
		attributeJson:           `"false"`,
		expectedInterface:       "false",
		expectedString:          "false",
		expectedBool:            false,
		expectedNumberError:     "json: cannot unmarshal string into Go value of type float64",
		decodeInto:              new(string),
		decodeIntoExpectedValue: "false",
	},
	{
		name:                    "DecodeSimpleString / True",
		attributeKey:            "test",
		attributeJson:           `"True"`,
		expectedInterface:       "True",
		expectedString:          "True",
		expectedBool:            true,
		expectedNumberError:     "json: cannot unmarshal string into Go value of type float64",
		decodeInto:              new(string),
		decodeIntoExpectedValue: "True",
	},
	{
		name:                    "DecodeSimpleString / False",
		attributeKey:            "test",
		attributeJson:           `"False"`,
		expectedInterface:       "False",
		expectedString:          "False",
		expectedBool:            false,
		expectedNumberError:     "json: cannot unmarshal string into Go value of type float64",
		decodeInto:              new(string),
		decodeIntoExpectedValue: "False",
	},
	{
		name:                    "DecodeSimpleString / TRUE",
		attributeKey:            "test",
		attributeJson:           `"TRUE"`,
		expectedInterface:       "TRUE",
		expectedString:          "TRUE",
		expectedBool:            true,
		expectedNumberError:     "json: cannot unmarshal string into Go value of type float64",
		decodeInto:              new(string),
		decodeIntoExpectedValue: "TRUE",
	},
	{
		name:                    "DecodeSimpleString / FALSE",
		attributeKey:            "test",
		attributeJson:           `"FALSE"`,
		expectedInterface:       "FALSE",
		expectedString:          "FALSE",
		expectedBool:            false,
		expectedNumberError:     "json: cannot unmarshal string into Go value of type float64",
		decodeInto:              new(string),
		decodeIntoExpectedValue: "FALSE",
	},
	{
		name:                    "DecodeSimpleString / t",
		attributeKey:            "test",
		attributeJson:           `"t"`,
		expectedInterface:       "t",
		expectedString:          "t",
		expectedBool:            true,
		expectedNumberError:     "json: cannot unmarshal string into Go value of type float64",
		decodeInto:              new(string),
		decodeIntoExpectedValue: "t",
	},
	{
		name:                    "DecodeSimpleString / f",
		attributeKey:            "test",
		attributeJson:           `"f"`,
		expectedInterface:       "f",
		expectedString:          "f",
		expectedBool:            false,
		expectedNumberError:     "json: cannot unmarshal string into Go value of type float64",
		decodeInto:              new(string),
		decodeIntoExpectedValue: "f",
	},
	{
		name:                    "DecodeSimpleString / T",
		attributeKey:            "test",
		attributeJson:           `"T"`,
		expectedInterface:       "T",
		expectedString:          "T",
		expectedBool:            true,
		expectedNumberError:     "json: cannot unmarshal string into Go value of type float64",
		decodeInto:              new(string),
		decodeIntoExpectedValue: "T",
	},
	{
		name:                    "DecodeSimpleString / F",
		attributeKey:            "test",
		attributeJson:           `"F"`,
		expectedInterface:       "F",
		expectedString:          "F",
		expectedBool:            false,
		expectedNumberError:     "json: cannot unmarshal string into Go value of type float64",
		decodeInto:              new(string),
		decodeIntoExpectedValue: "F",
	},
	{
		name:                    "DecodeSimpleString / 1",
		attributeKey:            "test",
		attributeJson:           `"1"`,
		expectedInterface:       "1",
		expectedString:          "1",
		expectedBool:            true,
		expectedNumber:          1.0,
		decodeInto:              new(string),
		decodeIntoExpectedValue: "1",
	},
	{
		name:                    "DecodeSimpleString / 0",
		attributeKey:            "test",
		attributeJson:           `"0"`,
		expectedInterface:       "0",
		expectedString:          "0",
		expectedBool:            false,
		expectedNumber:          0.0,
		decodeInto:              new(string),
		decodeIntoExpectedValue: "0",
	},
	{
		name:                    "DecodeSimpleString / Number",
		attributeKey:            "test",
		attributeJson:           `"9.9"`,
		expectedInterface:       "9.9",
		expectedString:          "9.9",
		expectedNumber:          9.9,
		expectedBoolError:       "json: cannot unmarshal string into Go value of type bool",
		decodeInto:              new(string),
		decodeIntoExpectedValue: "9.9",
	},
	{
		name:                    "DecodeSimpleInt",
		attributeKey:            "test",
		attributeJson:           `9`,
		expectedInterface:       float64(9),
		expectedNumber:          float64(9),
		expectedBoolError:       "json: cannot unmarshal number into Go value of type bool",
		expectedString:          "9",
		decodeInto:              new(int),
		decodeIntoExpectedValue: 9,
	},
	{
		name:                    "DecodeSimpleFloat",
		attributeKey:            "test",
		attributeJson:           `9.9`,
		expectedInterface:       9.9,
		expectedNumber:          float64(9.9),
		expectedBoolError:       "json: cannot unmarshal number into Go value of type bool",
		expectedString:          "9.9",
		decodeInto:              new(float64),
		decodeIntoExpectedValue: 9.9,
	},
	{
		name:                    "DecodeSimpleBool",
		attributeKey:            "test",
		attributeJson:           `true`,
		expectedInterface:       true,
		expectedBool:            true,
		expectedNumberError:     "json: cannot unmarshal bool into Go value of type float64",
		expectedString:          "true",
		decodeInto:              new(bool),
		decodeIntoExpectedValue: true,
	},
	{
		name:                    "DecodeArray",
		attributeKey:            "test",
		attributeJson:           `[ 1, 2 ]`,
		expectedInterface:       []interface{}{float64(1), float64(2)},
		expectedBoolError:       "json: cannot unmarshal array into Go value of type bool",
		expectedNumberError:     "json: cannot unmarshal array into Go value of type float64",
		expectedStringError:     "json: cannot unmarshal array into Go value of type string",
		decodeInto:              &[]int{},
		decodeIntoExpectedValue: []int{1, 2},
	},
	{
		name:          "DecodeObject",
		attributeKey:  "test",
		attributeJson: `{ "Field1": "value1", "Field2": 9 }`,
		expectedInterface: map[string]interface{}{
			"Field1": "value1",
			"Field2": float64(9)},
		expectedBoolError:   "json: cannot unmarshal object into Go value of type bool",
		expectedNumberError: "json: cannot unmarshal object into Go value of type float64",
		expectedStringError: "json: cannot unmarshal object into Go value of type string",
		decodeInto: &struct {
			Field1 string
			Field2 int
		}{},
		decodeIntoExpectedValue: struct {
			Field1 string
			Field2 int
		}{
			Field1: "value1",
			Field2: 9,
		},
	},
	{
		name:          "DecodeObjectIntoIncompleteStruct",
		attributeKey:  "test",
		attributeJson: `{ "Field1": "value1", "Field2": 9 }`,
		expectedInterface: map[string]interface{}{
			"Field1": "value1",
			"Field2": float64(9)},
		expectedBoolError:   "json: cannot unmarshal object into Go value of type bool",
		expectedNumberError: "json: cannot unmarshal object into Go value of type float64",
		expectedStringError: "json: cannot unmarshal object into Go value of type string",
		decodeInto:          &struct{ Field1 string }{},
		decodeIntoExpectedValue: struct{ Field1 string }{
			Field1: "value1",
		},
	},
	{
		name:          "DecodeObjectIntoStringFails",
		attributeKey:  "test",
		attributeJson: `{ "Field1": "value1", "Field2": 9 }`,
		expectedInterface: map[string]interface{}{
			"Field1": "value1",
			"Field2": float64(9)},
		expectedBoolError:       "json: cannot unmarshal object into Go value of type bool",
		expectedNumberError:     "json: cannot unmarshal object into Go value of type float64",
		expectedStringError:     "json: cannot unmarshal object into Go value of type string",
		decodeInto:              new(string),
		decodeIntoExpectedValue: "",
		decodeIntoError:         "json: cannot unmarshal object into Go value of type string",
	},
	{
		name:                    "DecodeInvalidObject",
		attributeKey:            "test",
		attributeJson:           `{ invalidObject }`,
		expectedInterface:       nil,
		expectedInterfaceError:  "invalid character 'i' looking for beginning of object key string",
		expectedBoolError:       "invalid character 'i' looking for beginning of object key string",
		expectedNumberError:     "invalid character 'i' looking for beginning of object key string",
		expectedStringError:     "invalid character 'i' looking for beginning of object key string",
		decodeInto:              &map[string]interface{}{},
		decodeIntoExpectedValue: map[string]interface{}{},
		decodeIntoError:         "invalid character 'i' looking for beginning of object key string",
	},
	{
		name:                    "GetInvalidKey",
		attributeKey:            invalidKey,
		attributeJson:           `9`,
		expectedInterface:       nil,
		expectedNumber:          float64(0),
		expectedInterfaceError:  keyNotFoundErr.Error(),
		expectedBoolError:       keyNotFoundErr.Error(),
		expectedNumberError:     keyNotFoundErr.Error(),
		expectedStringError:     keyNotFoundErr.Error(),
		decodeInto:              &map[string]interface{}{},
		decodeIntoExpectedValue: map[string]interface{}{},
		decodeIntoError:         keyNotFoundErr.Error(),
	},
}

func checkError(t *testing.T, err error, expectedError string) {
	if expectedError != "" {
		assert.EqualError(t, err, expectedError)
	} else {
		assert.NoError(t, err)
	}
}

func TestDecodeAttribute(t *testing.T) {
	for _, test := range decodeAttributeTestCases {
		t.Run(test.name, func(t *testing.T) {
			json := apiext.JSON{}
			if err := json.UnmarshalJSON([]byte(test.attributeJson)); err != nil {
				// This should never happen
				panic(err)
			}

			attributes := Attributes{
				"test": json,
			}

			var err error = nil
			assert.Equal(t, test.expectedBool, attributes.GetBoolean(test.attributeKey, &err))
			checkError(t, err, test.expectedBoolError)

			err = nil
			assert.Equal(t, test.expectedString, attributes.GetString(test.attributeKey, &err))
			checkError(t, err, test.expectedStringError)

			err = nil
			assert.Equal(t, test.expectedNumber, attributes.GetNumber(test.attributeKey, &err))
			checkError(t, err, test.expectedNumberError)

			err = nil
			assert.Equal(t, test.expectedInterface, attributes.Get(test.attributeKey, &err))
			checkError(t, err, test.expectedInterfaceError)

			err = attributes.GetInto(test.attributeKey, test.decodeInto)
			checkError(t, err, test.decodeIntoError)

			decodedValue := reflect.ValueOf(test.decodeInto)
			if decodedValue.Kind() == reflect.Ptr {
				decodedValue = decodedValue.Elem()
			}
			assert.Equal(t, test.decodeIntoExpectedValue, decodedValue.Interface())
		})
	}
}

type decodeAttributesTestCase struct {
	name                    string
	attributes              Attributes
	expectedBooleans        map[string]bool
	expectedBooleansError   string
	expectedStrings         map[string]string
	expectedStringsError    string
	expectedNumbers         map[string]float64
	expectedNumbersError    string
	expectedInterface       interface{}
	expectedInterfaceError  string
	decodeInto              interface{}
	decodeIntoExpectedValue interface{}
	decodeIntoError         string
}

var decodeAttributesTestCases []decodeAttributesTestCase = []decodeAttributesTestCase{
	{
		name: "DecodeSimpleStringMap",
		attributes: Attributes{}.FromStringMap(map[string]string{
			"firstString":  "firstStringValue",
			"secondString": "secondStringValue",
		}),
		expectedInterface: map[string]interface{}{
			"firstString":  "firstStringValue",
			"secondString": "secondStringValue",
		},
		decodeInto: &map[string]string{},
		decodeIntoExpectedValue: map[string]string{
			"firstString":  "firstStringValue",
			"secondString": "secondStringValue",
		},
	},
	{
		name: "DecodeStruct",
		attributes: Attributes{}.FromMap(map[string]interface{}{
			"attribute1": "value1",
			"attribute2": 9.9,
			"attribute3": true,
		}, nil),
		expectedInterface: map[string]interface{}{
			"attribute1": "value1",
			"attribute2": 9.9,
			"attribute3": true,
		},
		expectedBooleans: map[string]bool{
			"attribute3": true,
		},
		expectedNumbers: map[string]float64{
			"attribute2": 9.9,
		},
		expectedStrings: map[string]string{
			"attribute1": "value1",
		},
		decodeInto: &struct {
			Attribute1 string  `json:"attribute1"`
			Attribute2 float64 `json:"attribute2"`
			Attribute3 bool    `json:"attribute3"`
		}{},
		decodeIntoExpectedValue: struct {
			Attribute1 string  `json:"attribute1"`
			Attribute2 float64 `json:"attribute2"`
			Attribute3 bool    `json:"attribute3"`
		}{
			Attribute1: "value1",
			Attribute2: 9.9,
			Attribute3: true,
		},
	},
	{
		name: "DecodeStruct / missing attribute",
		attributes: Attributes{}.FromMap(map[string]interface{}{
			"attribute1": "value1",
		}, nil),
		expectedInterface: map[string]interface{}{
			"attribute1": "value1",
		},
		decodeInto: &struct {
			Attribute1 string  `json:"attribute1"`
			Attribute2 float64 `json:"attribute2"`
		}{},
		decodeIntoExpectedValue: struct {
			Attribute1 string  `json:"attribute1"`
			Attribute2 float64 `json:"attribute2"`
		}{
			Attribute1: "value1",
			Attribute2: 0.0,
		},
	},
	{
		name: "DecodeStruct / into wrong type",
		attributes: Attributes{}.FromMap(map[string]interface{}{
			"attribute1": "value1",
			"attribute2": 9.9,
		}, nil),
		expectedInterface: map[string]interface{}{
			"attribute1": "value1",
			"attribute2": 9.9,
		},
		decodeInto:              &[]string{},
		decodeIntoExpectedValue: []string{},
		decodeIntoError:         "json: cannot unmarshal object into Go value of type []string",
	},
	{
		name: "DecodeStruct / invalid type / With Get Error",
		attributes: Attributes{
			"attributes1": apiext.JSON{
				Raw: []byte("{ invalidObject }"),
			},
		},
		expectedInterface:       nil,
		expectedInterfaceError:  "json: error calling MarshalJSON for type attributes.Attributes: json: error calling MarshalJSON for type v1.JSON: invalid character 'i' looking for beginning of object key string",
		decodeInto:              &[]string{},
		decodeIntoExpectedValue: []string{},
		decodeIntoError:         "json: error calling MarshalJSON for type attributes.Attributes: json: error calling MarshalJSON for type v1.JSON: invalid character 'i' looking for beginning of object key string",
	},
	{
		name:                    "DecodeMap / nil / Should not change the provided map",
		attributes:              nil,
		expectedInterface:       nil,
		decodeInto:              &map[string]interface{}{},
		decodeIntoExpectedValue: map[string]interface{}{},
	},
}

func TestDecodeAttributes(t *testing.T) {
	for _, test := range decodeAttributesTestCases {
		t.Run(test.name, func(t *testing.T) {
			var err error

			if test.expectedBooleans != nil {
				err = nil
				assert.Equal(t, test.expectedBooleans, test.attributes.Booleans(&err))
				checkError(t, err, test.expectedBooleansError)
			}
			if test.expectedStrings != nil {
				err = nil
				assert.Equal(t, test.expectedStrings, test.attributes.Strings(&err))
				checkError(t, err, test.expectedStringsError)
			}
			if test.expectedNumbers != nil {
				err = nil
				assert.Equal(t, test.expectedNumbers, test.attributes.Numbers(&err))
				checkError(t, err, test.expectedNumbersError)
			}

			err = nil
			assert.Equal(t, test.expectedInterface, test.attributes.AsInterface(&err))
			checkError(t, err, test.expectedInterfaceError)

			err = test.attributes.Into(test.decodeInto)
			checkError(t, err, test.decodeIntoError)

			decodedValue := reflect.ValueOf(test.decodeInto)
			if decodedValue.Kind() == reflect.Ptr {
				decodedValue = decodedValue.Elem()
			}
			assert.Equal(t, test.decodeIntoExpectedValue, decodedValue.Interface())
		})
	}
}
