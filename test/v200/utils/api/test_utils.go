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

package api

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	commonUtils "github.com/devfile/api/v2/test/v200/utils/common"
	"github.com/santhosh-tekuri/jsonschema"
	"sigs.k8s.io/yaml"
)

const (
	// numDevfiles : the number of devfiles to create for each test
	numDevfiles = 5

	schemaFileName = "../../../schemas/latest/ide-targeted/devfile.json"
)

var schemas = make(map[string]SchemaFile)

// SchemaFile - represents the schema stucture
type SchemaFile struct {
	Schema *jsonschema.Schema
}

// DevfileValidator struct for DevfileValidator interface defined in common utils.
type DevfileValidator struct{}

// WriteAndValidate implements Saved.DevfileValidator interface.
// writes to disk and validates the specified devfile
func (devfileValidator DevfileValidator) WriteAndValidate(devfile *commonUtils.TestDevfile) error {
	err := writeDevfile(devfile)
	if err != nil {
		commonUtils.LogErrorMessage(fmt.Sprintf("Error writing file : %s : %v", devfile.FileName, err))
	} else {
		err = validateDevfile(devfile)
		if err != nil {
			commonUtils.LogErrorMessage(fmt.Sprintf("Error vaidating file : %s : %v", devfile.FileName, err))
		}
	}
	return err
}

// checkWithSchema checks the validity of a devfile against the schema.
func (schemaFile *SchemaFile) checkWithSchema(devfile string, expectedMessage string) error {

	// Read the created yaml file, ready for converison to json
	devfileData, err := ioutil.ReadFile(devfile)
	if err != nil {
		commonUtils.LogErrorMessage(fmt.Sprintf("  FAIL: schema : unable to read %s: %v", devfile, err))
		return err
	}

	// Convert the yaml file to json
	devfileDataAsJSON, err := yaml.YAMLToJSON(devfileData)
	if err != nil {
		commonUtils.LogErrorMessage(fmt.Sprintf("  FAIL : %s : schema : failed to convert to json : %v", devfile, err))
		return err
	}

	validationErr := schemaFile.Schema.Validate(bytes.NewReader(devfileDataAsJSON))
	if validationErr != nil {
		if len(expectedMessage) > 0 {
			if !strings.Contains(validationErr.Error(), expectedMessage) {
				err = errors.New(commonUtils.LogErrorMessage(fmt.Sprintf("  FAIL : schema : %s : Did not fail as expected : %s  got : %v", devfile, expectedMessage, validationErr)))
			} else {
				commonUtils.LogInfoMessage(fmt.Sprintf("PASS: schema :  Expected Error received : %s", expectedMessage))
			}
		} else {
			err = errors.New(commonUtils.LogErrorMessage(fmt.Sprintf("  FAIL : schema : %s : Did not pass as expected, got : %v", devfile, validationErr)))
		}
	} else {
		if len(expectedMessage) > 0 {
			err = errors.New(commonUtils.LogErrorMessage(fmt.Sprintf("  FAIL : schema : %s :  was valid - Expected Error not found : %v", devfile, validationErr)))
		} else {
			commonUtils.LogInfoMessage(fmt.Sprintf("  PASS : schema : %s : devfile was valid.", devfile))
		}
	}
	return err
}

// getSchema downloads and saves a schema from the provided url
func getSchema(schemaFileName string) (SchemaFile, error) {

	var err error
	schemaFile, found := schemas[schemaFileName]
	if !found {

		schemaFile = SchemaFile{}

		// Prepare the schema file
		compiler := jsonschema.NewCompiler()
		// Use Draft 7, github.com/santhosh-tekuri/jsonschema provides 4,6 an 7 so use the latest
		compiler.Draft = jsonschema.Draft7
		schemaFile.Schema, err = compiler.Compile(schemaFileName)
		if err != nil {
			commonUtils.LogErrorMessage(fmt.Sprintf("FAIL : Failed to compile schema  %v", err))
		} else {
			commonUtils.LogInfoMessage(fmt.Sprintf("Schema compiled from file: %s)", schemaFileName))
			schemas[schemaFileName] = schemaFile
		}
	}
	return schemaFile, err
}

// writeDevfile creates a devfile on disk for use in a test.
func writeDevfile(devfile *commonUtils.TestDevfile) error {
	var err error

	fileName := devfile.FileName
	if !strings.HasSuffix(fileName, ".yaml") {
		fileName += ".yaml"
	}

	commonUtils.LogInfoMessage(fmt.Sprintf("Marshall and write devfile %s", devfile.FileName))

	c, marshallErr := yaml.Marshal(&(devfile.SchemaDevFile))

	if marshallErr != nil {
		err = errors.New(commonUtils.LogErrorMessage(fmt.Sprintf("Marshall devfile %s : %v", devfile.FileName, marshallErr)))
	} else {
		err = ioutil.WriteFile(fileName, c, 0644)
		if err != nil {
			commonUtils.LogErrorMessage(fmt.Sprintf("Write devfile %s : %v", devfile.FileName, err))
		}
	}
	return err
}

// validateDevfile check the provided defile against the schema
func validateDevfile(devfile *commonUtils.TestDevfile) error {

	var err error
	var schemaFile SchemaFile

	schemaFile, err = getSchema(schemaFileName)
	if err != nil {
		commonUtils.LogErrorMessage(fmt.Sprintf("Failed to get devfile schema : %v", err))
	} else {
		err = schemaFile.checkWithSchema(devfile.FileName, "")
		if err != nil {
			commonUtils.LogErrorMessage(fmt.Sprintf("Verification with devfile schema failed : %v", err))
		} else {
			commonUtils.LogInfoMessage(fmt.Sprintf("Devfile validated using JSONSchema schema : %s", devfile.FileName))
		}
	}

	return err
}

// RunTest : Runs a test to create and verify a devfile based on the content of the specified TestContent
func RunTest(testContent commonUtils.TestContent, t *testing.T) {

	commonUtils.LogMessage(fmt.Sprintf("Start test for %s", testContent.FileName))

	validator := DevfileValidator{}

	devfileName := testContent.FileName
	for i := 1; i <= numDevfiles; i++ {

		testContent.FileName = commonUtils.AddSuffixToFileName(devfileName, strconv.Itoa(i))

		testDevfile, err := commonUtils.GetDevfile(testContent.FileName, nil, validator)
		if err != nil {
			t.Fatalf(commonUtils.LogMessage(fmt.Sprintf("Error creating devfile : %v", err)))
		}

		testDevfile.RunTest(testContent, t)
	}
}
