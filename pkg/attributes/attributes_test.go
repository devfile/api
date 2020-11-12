package attributes

import (
	//	"encoding/json"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	//	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func runBuildWithValidJson(t *testing.T, lenientBuild bool, builder Builder, expectedResult interface{}) {
	var attributes Attributes
	var err error

	if lenientBuild {
		attributes = builder.LenientBuild()
	} else {
		attributes, err = builder.Build()
		assert.NoError(t, err)
	}
	
	expectedJson, _ := json.Marshal(expectedResult)
	actualJson, _ := json.Marshal(attributes)
	assert.Equal(t, string(expectedJson), string(actualJson))
}

func runBuildWithInvalidValidJson(t *testing.T, lenientBuild bool, builder Builder, expectedResult interface{}, expectedError string) {
	var attributes Attributes
	var err error

	if lenientBuild {
		attributes = builder.LenientBuild()
		expectedJson, _ := json.Marshal(expectedResult)
		actualJson, _ := json.Marshal(attributes)
		assert.Equal(t, string(expectedJson), string(actualJson))
	} else {
		attributes, err = builder.Build()
		assert.EqualError(t, err, expectedError)
	}
}

type buildAttributesTestCase struct {
	name string
	lenient bool
	builder Builder
	expectedResult interface{}
	expectedError string
}

var buildAttributesTestCases []buildAttributesTestCase = []buildAttributesTestCase {
	{
		name: "Valid Attributes Build / simple types",
		lenient: false,
		builder: Builder{
			"boolAttribute": FromInterface(true),
			"numberAttribute": FromInterface(12),
			"stringAttribute": FromInterface("stringValue"),
		},
		expectedResult : &struct {
			BoolAttribute bool `json:"boolAttribute"`
			NumberAttribute int `json:"numberAttribute"`
			StringAttribute string `json:"stringAttribute"`
		} { BoolAttribute: true, StringAttribute: "stringValue", NumberAttribute: 12},
	},
	{
		name: "Valid Attributes Lenient Build / simple types",
		lenient: true,
		builder: Builder{
			"boolAttribute": FromInterface(true),
			"numberAttribute": FromInterface(12),
			"stringAttribute": FromInterface("stringValue"),
		},
		expectedResult : &struct {
			BoolAttribute bool `json:"boolAttribute"`
			NumberAttribute int `json:"numberAttribute"`
			StringAttribute string `json:"stringAttribute"`
		} { BoolAttribute: true, StringAttribute: "stringValue", NumberAttribute: 12},
	},
	{
		name: "Valid Attributes Build / structured types",
		lenient: false,
		builder: Builder{
			"arrayAttribute": FromInterface([]string{"stringOne, stringTwo"}),
			"objectAttribute": FromInterface(struct { SubField int }{SubField: 10}),
		},
		expectedResult : &struct {
			ArrayAttribute []string `json:"arrayAttribute"`
			ObjectAttribute struct { SubField int } `json:"objectAttribute"`
		} { ObjectAttribute: struct { SubField int }{SubField: 10}, ArrayAttribute: []string{"stringOne, stringTwo"}},
	},
	{
		name: "Valid Attributes Lenient Build / structured types",
		lenient: true,
		builder: Builder{
			"arrayAttribute": FromInterface([]string{"stringOne, stringTwo"}),
			"objectAttribute": FromInterface(struct { SubField int }{SubField: 10}),
		},
		expectedResult : &struct {
			ArrayAttribute []string `json:"arrayAttribute"`
			ObjectAttribute struct { SubField int } `json:"objectAttribute"`
		} { ObjectAttribute: struct { SubField int }{SubField: 10}, ArrayAttribute: []string{"stringOne, stringTwo"}},
	},
	{
		name: "Invalid Attributes Build",
		lenient: false,
		builder: Builder{
			"boolAttribute": FromInterface(true),
			"channelAttribute": FromInterface(make(chan int)),
		},
		expectedError: "json: unsupported type: chan int",
		expectedResult: Attributes{},
	},
	{
		name: "Invalid Attributes Lenient Build",
		lenient: true,
		builder: Builder{
			"boolAttribute": FromInterface(true),
			"channelAttribute": FromInterface(make(chan int)),
		},
		expectedResult : &struct {
			BoolAttribute bool `json:"boolAttribute"`
			ChannelAttribute interface{} `json:"channelAttribute"`
		} { BoolAttribute: true, ChannelAttribute: nil},
	},
}

func TestBuildAttributes(t *testing.T) {
	for _, test := range buildAttributesTestCases {
		t.Run(test.name, func(t *testing.T) {
			var attributes Attributes
			var err error
		
			if test.lenient {
				attributes = test.builder.LenientBuild()
			} else {
				attributes, err = test.builder.Build()
				if test.expectedError != "" {
					assert.EqualError(t, err, test.expectedError)
				} else {
					assert.NoError(t, err)
				}
			}
			
			expectedJson, _ := json.Marshal(test.expectedResult)
			actualJson, _ := json.Marshal(attributes)
			assert.Equal(t, string(expectedJson), string(actualJson))
		})
	}
}	

type decodeAttributeTestCase struct {
	name string
	attributeJson string
	expectedInterface interface{}
	decodeInto interface{}
	decodeIntoExpectedValue interface{}
	decodeIntoError string
}

var decodeAttributeTestCases []decodeAttributeTestCase = []decodeAttributeTestCase {
	{
		name: "DecodeSimpleString",
		attributeJson: `"simpleString"`,
		expectedInterface: "simpleString",
		decodeInto: new(string),
		decodeIntoExpectedValue : "simpleString",
	},
	{
		name: "DecodeSimpleInt",
		attributeJson: `9`,
		expectedInterface: float64(9),
		decodeInto: new(int),
		decodeIntoExpectedValue : 9,
	},
	{
		name: "DecodeSimpleFloat",
		attributeJson: `9.4`,
		expectedInterface: 9.4,
		decodeInto: new(float64),
		decodeIntoExpectedValue : 9.4,
	},
	{
		name: "DecodeSimpleBool",
		attributeJson: `true`,
		expectedInterface: true,
		decodeInto: new(bool),
		decodeIntoExpectedValue : true,
	},
	{
		name: "DecodeArray",
		attributeJson: `[ 1, 2 ]`,
		expectedInterface: []interface{}{float64(1), float64(2)},
		decodeInto: &[]int{},
		decodeIntoExpectedValue : []int{1, 2},
	},
	{
		name: "DecodeObject",
		attributeJson: `{ "Field1": "value1", "Field2": 9 }`,
		expectedInterface: map[string]interface{}{
			"Field1": "value1",
			"Field2": float64(9)},
		decodeInto: &struct{ Field1 string; Field2 int }{},
		decodeIntoExpectedValue : struct{ Field1 string; Field2 int }{
			Field1: "value1",
			Field2: 9,
		},
	},
	{
		name: "DecodeObjectIntoIncompleteStruct",
		attributeJson: `{ "Field1": "value1", "Field2": 9 }`,
		expectedInterface: map[string]interface{}{
			"Field1": "value1",
			"Field2": float64(9)},
		decodeInto: &struct{ Field1 string }{},
		decodeIntoExpectedValue : struct{ Field1 string }{
			Field1: "value1",
		},
	},
	{
		name: "DecodeObjectIntoStringFails",
		attributeJson: `{ "Field1": "value1", "Field2": 9 }`,
		expectedInterface: map[string]interface{}{
			"Field1": "value1",
			"Field2": float64(9)},
		decodeInto: new(string),
		decodeIntoError: "json: cannot unmarshal object into Go value of type string",
	},
}

func TestDecodeAttributes(t *testing.T) {
	for _, test := range decodeAttributeTestCases {
		t.Run(test.name, func(t *testing.T) {
			attribute := Attribute{}
			if err := attribute.UnmarshalJSON([]byte(test.attributeJson)); err != nil {
				// This should never happen 
				panic(err)
			}
			
			assert.Equal(t, test.expectedInterface, attribute.Interface())

			err := attribute.DecodeInto(test.decodeInto)
			if test.decodeIntoError != "" {
				assert.EqualError(t, err, test.decodeIntoError)
				return
			} else {
				assert.NoError(t, err)
			}

			decodedValue := reflect.ValueOf(test.decodeInto)
			if decodedValue.Kind() == reflect.Ptr {
				decodedValue = decodedValue.Elem()
			}
			assert.Equal(t, test.decodeIntoExpectedValue, decodedValue.Interface())
		})
	}
}	
