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

if ! command -v yq &> /dev/null
then
  echo
  echo "#### ERROR ####"
  echo "####"
  echo "#### Please install the 'yq' tool before being able to use this script"
  echo "#### see https://github.com/kislyuk/yq"
  echo "#### and https://stedolan.github.io/jq/download"
  echo "####"
  echo "###############"
  exit 1
fi

if ! command -v jsonpatch &> /dev/null
then
  echo
  echo "#### ERROR ####"
  echo "####"
  echo "#### Please install the 'jsonpatch' tool before being able to use this script"
  echo "#### For this use:"
  echo "####    pip3 install --user jsonpatch"
  echo "####"
  echo "###############"
  exit 1
fi

command -v operator-sdk >/dev/null 2>&1 || { echo -e "${RED}operator-sdk is not installed. Aborting.${NC}"; exit 1; }

operatorVersion=$(operator-sdk version)
[[ $operatorVersion =~ .*v0.17.0.* ]] || { echo -e "${RED}operator-sdk v0.17.0 is required${NC}"; exit 1; }

set -e

BASE_DIR=$(cd "$(dirname "$0")" && pwd)

mkdir -p "${BASE_DIR}/generated"

operator-sdk generate k8s
operator-sdk generate crds
yq '.spec.validation.openAPIV3Schema' \
  "${BASE_DIR}/deploy/crds/workspace.devfile.io_devworkspaces_crd.yaml" \
  > "${BASE_DIR}/schemas/devworkspace.json"

yq '.spec.validation.openAPIV3Schema' \
  "${BASE_DIR}/deploy/crds/workspace.devfile.io_devworkspacetemplates_crd.yaml" \
  > "${BASE_DIR}/schemas/devworkspace-template.json"

transform()
{
  local transformType="$1"
  local file="$2"
  local rewriteJsonPaths="$3"
  for patch in "${BASE_DIR}"/schema-transformation-rules/"${transformType}"/*.jq
  do
    [ -e "${patch}" ] || continue
    echo "Applying patch $(basename ${patch}) on schema $(basename ${file})"
    while IFS= read -r line
    do
      [ -n "${line}" ] || continue
      jq "${line}" "${file}" > "${file}.temp"
      mv "${file}.temp" "${file}"
    done < "${patch}"
  done
  for patch in "${BASE_DIR}"/schema-transformation-rules/"${transformType}"/*.json
  do
    [ -e "$patch" ] || continue
    patchContent=$(cat "$patch" | sed -e "${rewriteJsonPaths}")
    echo "Applying patch $(basename ${patch}) on schema $(basename ${file})"
    echo "$patchContent" | jsonpatch -i --indent 2 "${file}"
  done
}

onError() {
  echo "Cleaning schemas"
  rm -f "${BASE_DIR}/schemas/*"
}
trap 'onError' ERR

sed -i -e '/"description":/s/ \\t/\\n/g' -e '/"description":/s/ \\n - /\\n- /g' -e '/"description":/s/ \\n \([^-]\)/\\n\\n\1/g' "${BASE_DIR}/schemas/devworkspace-template.json" "${BASE_DIR}/schemas/devworkspace.json"

transform "devworkspace" "${BASE_DIR}/schemas/devworkspace-template.json" ""
transform "devworkspace" "${BASE_DIR}/schemas/devworkspace.json" 's#"path" *: *"/properties/spec/#"path": "/properties/spec/properties/template/#'

jq ".properties.spec" "${BASE_DIR}/schemas/devworkspace-template.json" > "${BASE_DIR}/schemas/devworkspace-template-spec.json"

jq ".properties.parent" "${BASE_DIR}/schemas/devworkspace-template-spec.json" > "${BASE_DIR}/schemas/override-spec.json"
transform "override" "${BASE_DIR}/schemas/override-spec.json" ''

cp "${BASE_DIR}/schemas/devworkspace-template-spec.json" "${BASE_DIR}/schemas/devfile.json"

transform "devfile" "${BASE_DIR}/schemas/devfile.json" ""

echo "Build of CRDs and schemas is finished"
