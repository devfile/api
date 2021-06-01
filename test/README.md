# schemaTest

The API tests are intended to provide a comprehensive verification of the devfile schemas. This includes:
- Ensuring every possible attribute is valid.
- Ensuring all optional attributes are indeed optional.
- Ensuring any possible specification errors are invalidated by the schema. For example:
    - Missing mandatory attributes.
    - Multiple use of a one-of attribute.
    - Attribute values of the wrong type.

## Test structure

- `test/v200/devfiles` : contains yaml snippets which are used to generate yaml files for the tests. The names of the sub-directories and files should reflect their purpose.
- `test/v200/schemaTest/schema-test.go` : the go unit test program.
- `test/v200/json` :  contains the json files which define the tests which the test program will run:
    - `test-xxxxxxx.json` : these files are the top level json files, they define the schema to verify and the test files to run.
    - `xxxxxx-tests.json` : these are the test files which contain individual tests which provide the yaml snippets to combine and the expected result.

## Running tests locally

From the `test/v200/schemaTest` directory run 
- `go test -v`

The test will read each of the test-xxxxxx.json files and run the tests defined within. The generated .yaml files used for the tests are created in a `tmp/test-xxxxxx/` directory. These files are not deleted when the test finishes so they can be used to assess any errors, however they will be deleted by a subsequent run of the test. Running the test with the -v option ensures you see a full list of passes and failures. 

## Adding Tests

### add a test for a new schema file

1. Create a new `test/v200/json/test-<schema name>.json` file for the schema. In the json file  specify the location of the schema to test (relative to the root directory of the repository), and the list of the existing tests to use. If the generated yaml files require a schemaVersion attribute include its value in the json file. see - *link to sample schema to be added*
1. Run the test

### add a test for a schema changes

1. Modify an existing yaml snippet or create a new one.
1. If appropriate create a new snippet for any possible error cases, for example to omit a required attribute.
1. If a new yaml snippet was created add a test which uses the snippet to the appropriate `json/xxxxxx-tests.json` file. Be careful to ensure the file name used for the test is unique for all tests - this is the name used for the yaml file which is generated for the test. For failure scenarios you may need to run the test first to set the outcome correctly. 
1. If a new  `json/xxxxxx-tests.json` file is created, any existing `test-xxxxxxx.json` files must be updated to use the new file.

### add test for a new schema version

1. Copy and rename the `test/v200` directory for the new version, for example `test\v201`
1. Update the copied `test/v201/json/test-<schema name>.json` files to point to the new schema.
1. Modify the copied tests as needed for the new version as decsribed above.
1. Add `test/v201/schemaTest/tmp` to the .gitignore file.
1. Run the test


# apiTest

A new test approach, shared with the library repository for testing valid devfiles. Basically the test creates lots of valid devfiles whith different content. The attributes which are set and the values to which they are set are randomized. These tests are a work in progress and the intent is to eventually replace schemaTest.  

## Test structure

- `test/v200/apiTest/api-test.go`: The go unit test program
- `test/v200/utils/api/test-utils.go` : utilites, used by the test, which contain functions uniqiue to the api tests.
- `test/v200/utils/common/*-utils.go` : utilites, used by the test, which are also used by the library tests. Mostly contain the code to generate valid devfile content.


## Running tests locally

from the `test/v200/apiTest/` directory run
- `go test -v`

* The test will generate a set of valid devfile.yaml files in `test/v200/apiTest/tmp/api-test/`
* The test will generate a log file:  `test/v200/apiTest/tmp/test.log`
* Each run of the test removes the  `test/v200/apiTest/tmp` directory from the previous run.

# Run test coverage analysis and reporting locally

This is useful to determine if there are any gaps in testing.  Run these steps at the root directory `/api` to get a report of all the overall tests, including the ones for api and schema

- `go test -coverprofile test/v200/api-test-coverage.out -v ./...`
- `go tool cover -html=test/v200/api-test-coverage.out -o test/v200/api-test-coverage.html`

## Viewing test results from a workflow

The tests run automatically with every PR or Push action.  You can see the results in the `devfile/api` repo's `Actions` view:

1.  Select the `CI` workflow and click on the PR you're interested in
1.  To view the console output, select the `build-and-validate` job and expand the `Run Go Tests` step.  This will give you a summary of the tests that were executed and their respective status 
1.  To view the test coverage report, click on the `Summary` page and you should see an `Artifacts` section with the `api-test-coverage-html` file available for download.