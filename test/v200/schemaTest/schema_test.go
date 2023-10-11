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

package schemaTest

import (
	"bytes"
	"encoding/json"
	"github.com/ghodss/yaml"
	"github.com/santhosh-tekuri/jsonschema"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Structure of json files providing schema to test and tests to run
type TestJson struct {
	SchemaFile    string   `json:"SchemaFile"`
	SchemaVersion string   `json:"SchemaVersion"`
	Tests         []string `json:"Tests"`
}

// Structure of a set of tests to run
type TestsToRun struct {
	Tests []TestToRun `json:"Tests"`
}

// Structre for an individual test
type TestToRun struct {
	FileName      string   `json:"FileName"`
	ExpectOutcome string   `json:"ExpectOutcome"`
	Disabled      bool     `json:"Disabled"`
	Files         []string `json:"Files"`
}

const testDir = "../"
const schemaDir = "../../../"
const jsonDir = "../json/"
const tempDir = "./tmp/"

func Test_Schema(t *testing.T) {

	// Clear the temp directory if it exists
	if _, err := os.Stat(tempDir); !os.IsNotExist(err) {
		os.RemoveAll(tempDir)
	}
	os.Mkdir(tempDir, 0755)

	// Read the content of the jso directory to find test files
	files, err := ioutil.ReadDir(jsonDir)
	if err != nil {
		t.Fatalf("Error finding test json files in : %s :  %v", jsonDir, err)
	}

	combinedTests := 0
	combinedPasses := 0
	for _, testJsonFile := range files {

		// if the file begins with test- and ends .json it can be processed
		if strings.HasPrefix(testJsonFile.Name(), "test-") && strings.HasSuffix(testJsonFile.Name(), ".json") {

			// Open the json file which defines the tests to run
			testJson, err := os.Open(filepath.Join(jsonDir, testJsonFile.Name()))
			if err != nil {
				t.Errorf("  FAIL : Failed to open %s : %s", testJsonFile.Name(), err)
				continue
			}

			// Read contents of the json file which defines the tests to run
			byteValue, err := ioutil.ReadAll(testJson)
			if err != nil {
				t.Errorf("FAIL : failed to read : %s : %v", testJsonFile.Name(), err)
			}

			var testJsonContent TestJson

			// Unmarshall the contents of the json file which defines the tests to run for each test
			err = json.Unmarshal(byteValue, &testJsonContent)
			if err != nil {
				t.Errorf("FAIL : failed to unmarshal : %s : %v", testJsonFile.Name(), err)
				continue
			}

			testJson.Close()

			t.Logf("INFO : File %s : SchemaFile : %s , SchemaVersion : %s", testJsonFile.Name(), testJsonContent.SchemaFile, testJsonContent.SchemaVersion)

			// Prepare the schema file
			compiler := jsonschema.NewCompiler()
			compiler.Draft = jsonschema.Draft7
			schema, err := compiler.Compile(filepath.Join(schemaDir, testJsonContent.SchemaFile))
			if err != nil {
				t.Fatalf("  FAIL : Schema compile failed : %s: %v", testJsonContent.SchemaFile, err)
			}

			// create the tempdirectoy to hold the generated yaml files
			testTempDir := tempDir
			testTempDir += strings.Split(testJsonFile.Name(), ".")[0]
			os.Mkdir(testTempDir, 0755)

			passTests := 0
			totalTests := 0

			// for each of the test files specified
			for m := 0; m < len(testJsonContent.Tests); m++ {

				// Open the json file which defines the tests to run
				testsToRunJson, err := os.Open(filepath.Join(jsonDir, testJsonContent.Tests[m]))
				if err != nil {
					t.Errorf("Failed to open tests %s : %s", testJsonContent.Tests[m], err)
					continue
				}

				// Read contents of the json file which defines the tests to run
				byteValue, err := ioutil.ReadAll(testsToRunJson)
				if err != nil {
					t.Fatalf("FAIL : failed to read : %s : %v", testJsonContent.Tests[m], err)
				}

				var testsToRunContent TestsToRun

				// Unmarshall the contents of the json file which defines the tests to run for each test
				err = json.Unmarshal(byteValue, &testsToRunContent)
				if err != nil {
					t.Fatalf("FAIL : failed to unmarshal : %s : %v", testJsonContent.Tests[m], err)
				}

				testsToRunJson.Close()

				// For each test defined in the test file
				for i := 0; i < len(testsToRunContent.Tests); i++ {

					if !testsToRunContent.Tests[i].Disabled {

						totalTests++

						// Open the file to containe the generated test yaml
						f, err := os.OpenFile(filepath.Join(testTempDir, testsToRunContent.Tests[i].FileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
						if err != nil {
							t.Errorf("FAIL : Failed to open %s : %v", filepath.Join(testTempDir, testsToRunContent.Tests[i].FileName), err)
							continue
						}

						// If test requires a schema write it to the yaml file
						if testJsonContent.SchemaVersion != "" {
							f.WriteString("schemaVersion: " + testJsonContent.SchemaVersion + "\n")
						}

						testYamlComplete := true
						// Now add each of the yaml sippets used the make the yaml file for test
						for j := 0; j < len(testsToRunContent.Tests[i].Files); j++ {
							// Read the snippet
							data, err := ioutil.ReadFile(filepath.Join(testDir, testsToRunContent.Tests[i].Files[j]))
							if err != nil {
								t.Errorf("FAIL: failed reading %s: %v", filepath.Join(testDir, testsToRunContent.Tests[i].Files[j]), err)
								testYamlComplete = false
								continue
							}
							if j > 0 {
								// Ensure appropriate line breaks
								f.WriteString("\n")
							}

							// Add snippet to yaml file
							f.Write(data)
						}

						if !testYamlComplete {
							f.Close()
							continue
						}

						// Read the created yaml file, ready for converison to json
						data, err := ioutil.ReadFile(filepath.Join(testTempDir, testsToRunContent.Tests[i].FileName))
						if err != nil {
							t.Errorf("  FAIL: unable to read %s: %v", testsToRunContent.Tests[i].FileName, err)
							f.Close()
							continue
						}

						f.Close()

						// Convert the yaml file to json
						yamldoc, err := yaml.YAMLToJSON(data)
						if err != nil {
							t.Errorf("  FAIL : %s : failed to convert to json : %v", testsToRunContent.Tests[i].FileName, err)
							continue
						}

						// validate the test yaml against the schema
						if err = schema.Validate(bytes.NewReader(yamldoc)); err != nil {
							if testsToRunContent.Tests[i].ExpectOutcome == "PASS" {
								t.Errorf("  FAIL : %s : %s : Validate failure : %s", testsToRunContent.Tests[i].FileName, testJsonContent.SchemaFile, err)
							} else if testsToRunContent.Tests[i].ExpectOutcome == "" {
								t.Errorf("  FAIL : %s : No expected ouctome was set : %s  got : %s", testsToRunContent.Tests[i].FileName, testsToRunContent.Tests[i].ExpectOutcome, err.Error())
							} else if !strings.Contains(err.Error(), testsToRunContent.Tests[i].ExpectOutcome) {
								t.Errorf("  FAIL : %s : %s : Did not fail as expected : %s  got : %s", testsToRunContent.Tests[i].FileName, testJsonContent.SchemaFile, testsToRunContent.Tests[i].ExpectOutcome, err.Error())
							} else {
								passTests++
								t.Logf("PASS : %s : %s: %s", testsToRunContent.Tests[i].FileName, testJsonContent.SchemaFile, testsToRunContent.Tests[i].ExpectOutcome)
							}
						} else if testsToRunContent.Tests[i].ExpectOutcome == "" {
							t.Errorf("  FAIL : %s : devfile was valid - No expected ouctome was set.", testsToRunContent.Tests[i].FileName)
						} else if testsToRunContent.Tests[i].ExpectOutcome != "PASS" {
							t.Errorf("  FAIL : %s : devfile was valid - Expected Error not found :  %s", testsToRunContent.Tests[i].FileName, testsToRunContent.Tests[i].ExpectOutcome)
						} else {
							passTests++
							t.Logf("PASS : %s : %s", testsToRunContent.Tests[i].FileName, testJsonContent.SchemaFile)
						}
						f.Close()
					}
				}
			}
			t.Logf("%s : %d of %d tests passed", testJsonFile.Name(), passTests, totalTests)
			t.Logf("")
			combinedTests += totalTests
			combinedPasses += passTests

		}
	}

	if combinedTests != combinedPasses {
		t.Errorf("OVERALL FAIL : %d of %d tests failed.", (combinedTests - combinedPasses), combinedTests)
	} else {
		t.Logf("OVERALL PASS : %d of %d tests passed.", combinedPasses, combinedTests)
	}

}
