package api

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	common "github.com/devfile/api/v2/test/v200/utils/common"
	"github.com/santhosh-tekuri/jsonschema"
	"sigs.k8s.io/yaml"
)

const (
	// numDevfiles : the number of devfiles to create for each test
	numDevfiles = 5

	schemaFileName = "../../../schemas/latest/ide-targeted/devfile.json"
)

var schemas = make(map[string]SchemaFile)

type SchemaFile struct {
	Schema *jsonschema.Schema
}

type DevfileValidator struct{}

// WriteAndVerify implements Saved.DevfileValidator interface.
// writes to disk and validates the specified devfile
func (devfileValidator DevfileValidator) WriteAndValidate(devfile *common.TestDevfile) error {
	err := writeDevfile(devfile)
	if err != nil {
		common.LogErrorMessage(fmt.Sprintf("Error writing file : %s : %v", devfile.FileName, err))
	} else {
		err = validateDevfile(devfile)
		if err != nil {
			common.LogErrorMessage(fmt.Sprintf("Error vaidating file : %s : %v", devfile.FileName, err))
		}
	}
	return err
}

// CheckWithSchema checks the validity of aa devfile against the schema.
func (schemaFile *SchemaFile) CheckWithSchema(devfile string, expectedMessage string) error {

	// Read the created yaml file, ready for converison to json
	devfileData, err := ioutil.ReadFile(devfile)
	if err != nil {
		common.LogErrorMessage(fmt.Sprintf("  FAIL: schema : unable to read %s: %v", devfile, err))
		return err
	}

	// Convert the yaml file to json
	devfileDataAsJSON, err := yaml.YAMLToJSON(devfileData)
	if err != nil {
		common.LogErrorMessage(fmt.Sprintf("  FAIL : %s : schema : failed to convert to json : %v", devfile, err))
		return err
	}

	validationErr := schemaFile.Schema.Validate(bytes.NewReader(devfileDataAsJSON))
	if validationErr != nil {
		if len(expectedMessage) > 0 {
			if !strings.Contains(validationErr.Error(), expectedMessage) {
				err = errors.New(common.LogErrorMessage(fmt.Sprintf("  FAIL : schema : %s : Did not fail as expected : %s  got : %v", devfile, expectedMessage, validationErr)))
			} else {
				common.LogInfoMessage(fmt.Sprintf("PASS: schema :  Expected Error received : %s", expectedMessage))
			}
		} else {
			err = errors.New(common.LogErrorMessage(fmt.Sprintf("  FAIL : schema : %s : Did not pass as expected, got : %v", devfile, validationErr)))
		}
	} else {
		if len(expectedMessage) > 0 {
			err = errors.New(common.LogErrorMessage(fmt.Sprintf("  FAIL : schema : %s :  was valid - Expected Error not found : %v", devfile, validationErr)))
		} else {
			common.LogInfoMessage(fmt.Sprintf("  PASS : schema : %s : devfile was valid.", devfile))
		}
	}
	return err
}

// GetSchema downloads and saves a schema from the provided url
func GetSchema(schemafile string) (SchemaFile, error) {

	var err error
	schemaFile, found := schemas[schemafile]
	if !found {

		schemaFile = SchemaFile{}

		// Prepare the schema file
		compiler := jsonschema.NewCompiler()
		compiler.Draft = jsonschema.Draft7
		schemaFile.Schema, err = compiler.Compile(schemafile)
		if err != nil {
			//t.Fatalf("  FAIL : Schema compile failed : %s: %v", testJsonContent.SchemaFile, err)
			common.LogErrorMessage(fmt.Sprintf("FAIL : Failed to compile schema  %v", err))
		} else {
			common.LogInfoMessage(fmt.Sprintf("Schema compiled from file: %s)", schemafile))
			schemas[schemafile] = schemaFile
		}
	}
	return schemaFile, err
}

// WriteDevfile creates a devfile on disk for use in a test.
func writeDevfile(devfile *common.TestDevfile) error {
	var err error

	fileName := devfile.FileName
	if !strings.HasSuffix(fileName, ".yaml") {
		fileName += ".yaml"
	}

	common.LogInfoMessage(fmt.Sprintf("Marshall and write devfile %s", devfile.FileName))

	c, marshallErr := yaml.Marshal(&(devfile.SchemaDevFile))

	if marshallErr != nil {
		err = errors.New(common.LogErrorMessage(fmt.Sprintf("Marshall devfile %s : %v", devfile.FileName, marshallErr)))
	} else {
		err = ioutil.WriteFile(fileName, c, 0644)
		if err != nil {
			common.LogErrorMessage(fmt.Sprintf("Write devfile %s : %v", devfile.FileName, err))
		}
	}
	return err
}

// validateDevfile check the provided defile against the schema
func validateDevfile(devfile *common.TestDevfile) error {

	var err error
	var schemaFile SchemaFile

	schemaFile, err = GetSchema(schemaFileName)
	if err != nil {
		common.LogErrorMessage(fmt.Sprintf("Failed to get devfile schema : %v", err))
	} else {
		err = schemaFile.CheckWithSchema(devfile.FileName, "")
		if err != nil {
			common.LogErrorMessage(fmt.Sprintf("Verification with devfile schema failed : %v", err))
		} else {
			common.LogInfoMessage(fmt.Sprintf("Devfile validated using JSONSchema schema : %s", devfile.FileName))
		}
	}

	return err
}

// RunTest : Runs a test to create and verify a devfile based on the content of the specified TestContent
func RunTest(testContent common.TestContent, t *testing.T) {

	common.LogMessage(fmt.Sprintf("Start test for %s", testContent.FileName))

	validator := DevfileValidator{}

	devfileName := testContent.FileName
	for i := 1; i <= numDevfiles; i++ {

		testContent.FileName = common.AddSuffixToFileName(devfileName, strconv.Itoa(i))

		testDevfile, err := common.GetDevfile(testContent.FileName, nil, validator)
		if err != nil {
			t.Fatalf(common.LogMessage(fmt.Sprintf("Error creating devfile : %v", err)))
		}

		testDevfile.RunTest(testContent, t)
	}
}
