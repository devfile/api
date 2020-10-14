#!/bin/bash
#
# Copyright (c) 2020 Red Hat, Inc.
# This program and the accompanying materials are made
# available under the terms of the Eclipse Public License 2.0
# which is available at https://www.eclipse.org/legal/epl-2.0/
#
# SPDX-License-Identifier: EPL-2.0
#
# Contributors:
#   Red Hat, Inc. - initial API and implementation

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

if ! command -v jsonschema-cli &> /dev/null
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
    if ! jsonschema-cli validate "${BASE_DIR}/${schemaPath}" "${BASE_DIR}/${devfile}" >> validate-output.txt
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
