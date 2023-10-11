#!/bin/bash
#
#
# Copyright Red Hat
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

BLUE='\033[1;34m'
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'
BOLD='\033[1m'

if ! command -v jq &> /dev/null
then
  echo
  echo "#### ERROR ####"
  echo "####"
  echo "#### Please install the 'jq' tool before being able to use this script"
  echo "#### see https://stedolan.github.io/jq/download"
  echo "####"
  echo "###############"
  exit 1
fi

if ! command -v jsonschema &> /dev/null
then
  echo
  echo "#### ERROR ####"
  echo "####"
  echo "#### Please install the 'jsonschema-cli' tool before being able to use this script"
  echo "#### see https://pypi.org/project/jsonschema-cli/"
  echo "####"
  echo "###############"
  exit 1
fi

BASE_DIR=$(cd "$(dirname "$0")" && pwd)

rm -f validate-output.txt &> /dev/null

for schema in "devfile" "dev-workspace" "dev-workspace-template"
do
  schemaPath="./schemas/latest/ide-targeted/${schema}.json"
  devfiles=$(jq -r '."yaml.schemas"."'${schemaPath}'"[]' .theia/settings.json)
  if [ $? -ne 0 ]; then
    exit 1
  fi 
  echo "Validating $schema files against ${schemaPath}"
  for devfile in $devfiles
  do
    python3 validate_yaml.py "${BASE_DIR}/${schemaPath}" "${BASE_DIR}/${devfile}" >> validate-output.txt
    if [ "$(cat validate-output.txt)" != "" ]
    then
      echo "  - $devfile => INVALID"
    else 
      echo "  - $devfile => OK"
    fi
  done
done
if [ "$(cat validate-output.txt)" != "" ]
then
  echo
  echo "Some files are invalid according to the related json schema."  
  echo "Detailed information can be found in the 'validate-output.txt' file."
  exit 1
fi

rm -r validate-output.txt &> /dev/null
echo "Validation of devfiles is successful"
