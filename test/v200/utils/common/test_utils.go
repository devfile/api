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

package common

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	schema "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	header "github.com/devfile/api/v2/pkg/devfile"
)

const (
	// maxCommands : The maximum number of commands to include in a generated devfile
	maxCommands = 10
	// maxComponents : The maximum number of components to include in a generated devfile
	maxComponents = 10
	// maxProjects : The maximum number of projects to include in a generated devfile
	maxProjects = 10
	// maxStarterProjects : the number of starterProjects to create for each test
	maxStarterProjects = 10

	defaultTempDir = "./tmp/"
	logFileName    = "test.log"
	// logToFileOnly - If set to false the log output will also be output to the console
	logToFileOnly = true
)

var (
	// tmpDir temporary directory in use
	tmpDir     string
	testLogger *log.Logger
)

// DevfileFollower interface implemented by the parser tests for updating the parser data
type DevfileFollower interface {
	AddCommand(schema.Command) error
	UpdateCommand(schema.Command)
	AddComponent(schema.Component) error
	UpdateComponent(schema.Component)
	AddProject(schema.Project) error
	UpdateProject(schema.Project)
	AddStarterProject(schema.StarterProject) error
	UpdateStarterProject(schema.StarterProject)
	AddEvent(schema.Events) error
	UpdateEvent(schema.Events)
	SetParent(schema.Parent) error
	UpdateParent(schema.Parent)

	SetMetaData(header.DevfileMetadata) error
	UpdateMetaData(header.DevfileMetadata)
	SetSchemaVersion(string)
}

// DevfileValidator interface implemented by the parser and api tests for verifying generated devfiles
type DevfileValidator interface {
	WriteAndValidate(*TestDevfile) error
}

// TestContent - structure used by a test to configure the tests to run
type TestContent struct {
	CommandTypes        []schema.CommandType
	ComponentTypes      []schema.ComponentType
	ProjectTypes        []schema.ProjectSourceType
	StarterProjectTypes []schema.ProjectSourceType
	AddParent           bool
	AddMetaData         bool
	AddEvents           bool
	FileName            string
	EditContent         bool
}

// init creates:
//   - the temporary directory used by the test to store logs and generated devfiles.
//   - the log file
func init() {
	tmpDir = defaultTempDir
	if _, err := os.Stat(tmpDir); !os.IsNotExist(err) {
		os.RemoveAll(tmpDir)
	}
	if err := os.Mkdir(tmpDir, 0755); err != nil {
		fmt.Printf("Failed to create temp directory, will use current directory : %v ", err)
		tmpDir = "./"
	}
	f, err := os.OpenFile(filepath.Join(tmpDir, logFileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error creating Log file : %v", err)
	} else {
		if logToFileOnly {
			testLogger = log.New(f, "", log.LstdFlags|log.Lmicroseconds)
		} else {
			writer := io.MultiWriter(f, os.Stdout)
			testLogger = log.New(writer, "", log.LstdFlags|log.Lmicroseconds)
		}
		testLogger.Println("Test Starting:")
	}

	rand.Seed(time.Now().UnixNano())

}

// CreateTempDir creates a specified sub directory under the temp directory if it does not exist.
// Returns the name of the created directory.
func CreateTempDir(subdir string) string {
	tempDir := tmpDir + subdir + "/"
	var err error
	if _, err = os.Stat(tempDir); os.IsNotExist(err) {
		err = os.Mkdir(tempDir, 0755)
	}
	if err != nil {
		// if cannot create subdirectory just use the base tmp directory
		LogErrorMessage(fmt.Sprintf("Failed to create temp directory %s will use %s : %v", tempDir, tmpDir, err))
		tempDir = tmpDir
	}
	return tempDir
}

// GetDevFileName returns a qualified name of a devfile for use in a test.
// The devfile will be in a temporary directory and is named using the calling function's name.
func GetDevFileName() string {
	pc, fn, _, ok := runtime.Caller(1)
	if !ok {
		return tmpDir + "DefaultDevfile"
	}

	testFile := filepath.Base(fn)
	testFileExtension := filepath.Ext(testFile)
	subdir := testFile[0 : len(testFile)-len(testFileExtension)]
	destDir := CreateTempDir(subdir)
	callerName := runtime.FuncForPC(pc).Name()
	pos1 := strings.LastIndex(callerName, "Test_")
	devfileName := destDir + callerName[pos1:len(callerName)] + ".yaml"

	LogInfoMessage(fmt.Sprintf("GetDevFileName : %s", devfileName))

	return devfileName
}

// AddSuffixToFileName adds a specified suffix to the name of a specified file.
// For example if the file is devfile.yaml and the suffix is 1, the result is devfile1.yaml
func AddSuffixToFileName(fileName string, suffix string) string {
	pos1 := strings.LastIndex(fileName, ".yaml")
	newFileName := fileName[0:pos1] + suffix + ".yaml"
	LogInfoMessage(fmt.Sprintf("Add suffix %s to fileName %s : %s", suffix, fileName, newFileName))
	return newFileName
}

// LogMessage logs the specified message and returns the message logged
func LogMessage(message string) string {
	if testLogger != nil {
		testLogger.Println(message)
	} else {
		fmt.Printf("Logger not available: %s", message)
	}
	return message
}

var errorPrefix = "..... ERROR :"
var warningPrefix = "..... WARNING :"
var infoPrefix = "INFO :"

// LogErrorMessage logs the specified message as an error message and returns the message logged
func LogErrorMessage(message string) string {
	var errMessage []string
	errMessage = append(errMessage, errorPrefix, message)
	return LogMessage(fmt.Sprint(errMessage))
}

// LogWarningMessage logs the specified message as a warning message and returns the message logged
func LogWarningMessage(message string) string {
	var warningMessage []string
	warningMessage = append(warningMessage, warningPrefix, message)
	return LogMessage(fmt.Sprint(warningMessage))
}

// LogInfoMessage logs the specified message as an info message and returns the message logged
func LogInfoMessage(message string) string {
	var infoMessage []string
	infoMessage = append(infoMessage, infoPrefix, message)
	return LogMessage(fmt.Sprint(infoMessage))
}

// TestDevfile is a structure used to track a test devfile and its contents
type TestDevfile struct {
	SchemaDevFile schema.Devfile
	FileName      string
	GroupDefaults map[schema.CommandGroupKind]bool
	UsedPorts     map[int]bool
	Follower      DevfileFollower
	Validator     DevfileValidator
}

var StringCount int = 0

// GetRandomUniqueString returns a unique random string which is n characters long plus an integer to ensure uniqueness
// If lower is set to true a lower case string is returned.
func GetRandomUniqueString(n int, lower bool) string {
	StringCount++
	countAsString := fmt.Sprintf("%05d", StringCount)
	if n < len(countAsString) {
		n += len(countAsString)
	}
	return fmt.Sprintf("%s%s", GetRandomString(n-len(countAsString), lower), countAsString)
}

const schemaBytes = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// GetRandomString returns a random string which is n characters long.
// If lower is set to true a lower case string is returned.
func GetRandomString(n int, lower bool) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = schemaBytes[rand.Intn(len(schemaBytes)-1)]
	}
	randomString := string(b)
	if lower {
		randomString = strings.ToLower(randomString)
	}
	return randomString
}

var GroupKinds = [...]schema.CommandGroupKind{schema.BuildCommandGroupKind, schema.RunCommandGroupKind, schema.TestCommandGroupKind, schema.DebugCommandGroupKind}

var ComponentTypes = []schema.ComponentType{schema.ContainerComponentType, schema.KubernetesComponentType, schema.OpenshiftComponentType, schema.VolumeComponentType}

var ProjectSourceTypes = []schema.ProjectSourceType{schema.GitProjectSourceType, schema.ZipProjectSourceType}

var CommandTypes = []schema.CommandType{schema.ExecCommandType, schema.CompositeCommandType, schema.ApplyCommandType}

var ImportReferenceTypes = []schema.ImportReferenceType{schema.IdImportReferenceType, schema.KubernetesImportReferenceType, schema.UriImportReferenceType}

// GetRandomGroupKind return random group kind. One of "build", "run", "test" or "debug"
func GetRandomGroupKind() schema.CommandGroupKind {
	return GroupKinds[GetRandomNumber(0, len(GroupKinds))-1]
}

// GetRandomValue returns a value selected at random from an array or slice of typed constants such as the ones declared
// in api/test/v200/utils/common/test_utils.go e.g GroupKinds, ComponentType, ProjectTypes, etc
// If the input is not a reflect.Array or reflect.Slice kind, then the value itself is returned and an error will be logged

func GetRandomValue(typedConstants interface{}) reflect.Value {
	constType := reflect.TypeOf(typedConstants)
	values := reflect.ValueOf(typedConstants)

	if constType.Kind() != reflect.Array && constType.Kind() != reflect.Slice {
		LogErrorMessage(fmt.Sprintf("GetRandomValue: Unexpected input \"%v\" kind, want reflect.Array or reflect.Slice kind", constType.Kind()))
		return values
	}

	length := values.Len()
	if length > 0 {
		return values.Index(rand.Intn(length))
	} else {
		LogErrorMessage(fmt.Sprintf("No values found %v", values))
		return values
	}
}

// GetBinaryDecision randomly returns true or false
func GetBinaryDecision() bool {
	return GetRandomDecision(1, 1)
}

// GetRandomDecision randomly returns true or false, but weighted to one or the other.
// For example if success is set to 2 and failure to 1, true is twice as likely to be returned.
func GetRandomDecision(success int, failure int) bool {
	return rand.Intn(success+failure) > failure-1
}

// GetRandomNumber randomly returns an integer between 1 and the number specified.
func GetRandomNumber(min int, max int) int {
	if min == max {
		return 1
	} else if min > max {
		return rand.Intn(max) + 1
	}
	return rand.Intn(max-min) + min + 1
}

// GetDevfile returns a structure used to represent a specific devfile in a test
func GetDevfile(fileName string, follower DevfileFollower, validator DevfileValidator) (TestDevfile, error) {

	var err error
	testDevfile := TestDevfile{}
	testDevfile.SchemaDevFile = schema.Devfile{}
	testDevfile.FileName = fileName
	testDevfile.SchemaDevFile.SchemaVersion = "2.0.0"
	testDevfile.GroupDefaults = make(map[schema.CommandGroupKind]bool)
	for _, kind := range GroupKinds {
		testDevfile.GroupDefaults[kind] = false
	}
	testDevfile.UsedPorts = make(map[int]bool)
	testDevfile.Validator = validator
	testDevfile.Follower = follower

	return testDevfile, err
}

// Runs a test to create and verify a devfile based on the content of the specified TestContent
func (testDevfile *TestDevfile) RunTest(testContent TestContent, t *testing.T) {

	if len(testContent.CommandTypes) > 0 {
		numCommands := GetRandomNumber(1, maxCommands)
		for i := 0; i < numCommands; i++ {
			commandIndex := GetRandomNumber(1, len(testContent.CommandTypes))
			testDevfile.AddCommand(testContent.CommandTypes[commandIndex-1])
		}
	}

	if len(testContent.ComponentTypes) > 0 {
		numComponents := GetRandomNumber(1, maxComponents)
		for i := 0; i < numComponents; i++ {
			componentIndex := GetRandomNumber(1, len(testContent.ComponentTypes))
			testDevfile.AddComponent(testContent.ComponentTypes[componentIndex-1])
		}
	}

	if len(testContent.ProjectTypes) > 0 {
		numProjects := GetRandomNumber(1, maxProjects)
		for i := 0; i < numProjects; i++ {
			projectIndex := GetRandomNumber(1, len(testContent.ProjectTypes))
			testDevfile.AddProject(testContent.ProjectTypes[projectIndex-1])
		}
	}

	if len(testContent.StarterProjectTypes) > 0 {
		numStarterProjects := GetRandomNumber(1, maxStarterProjects)
		for i := 0; i < numStarterProjects; i++ {
			starterProjectIndex := GetRandomNumber(1, len(testContent.StarterProjectTypes))
			testDevfile.AddStarterProject(testContent.StarterProjectTypes[starterProjectIndex-1])
		}
	}

	if testContent.AddParent {
		testDevfile.AddParent()
	}

	if testContent.AddEvents {
		testDevfile.AddEvents()
	}

	if testContent.AddMetaData {
		testDevfile.AddMetaData()
	}

	err := testDevfile.Validator.WriteAndValidate(testDevfile)
	if err != nil {
		t.Fatalf(LogErrorMessage(fmt.Sprintf("ERROR verifying devfile :  %s : %v", testContent.FileName, err)))
	}

}
