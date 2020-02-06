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

command -v operator-sdk >/dev/null 2>&1 || { echo -e $RED"operator-sdk is not installed. Aborting."$NC; exit 1; }

operatorVersion=$(operator-sdk version)
[[ $operatorVersion =~ .*v0.12.0.* ]] || { echo -e $RED"operator-sdk v0.12.0 is required"$NC; exit 1; }

set -e

BASE_DIR=$(cd "$(dirname "$0")" && pwd)

mkdir -p ${BASE_DIR}/generated

operator-sdk generate k8s
operator-sdk generate openapi
yq '.spec.validation.openAPIV3Schema' \
  ${BASE_DIR}/deploy/crds/workspaces.ecd.eclipse.org_devworkspaces_crd.yaml \
  > ${BASE_DIR}/schemas/devworkspace.json

jq ".properties.spec.properties.template" ${BASE_DIR}/schemas/devworkspace.json > ${BASE_DIR}/schemas/devworkspace-template-spec.json

cp ${BASE_DIR}/schemas/devworkspace-template-spec.json ${BASE_DIR}/schemas/devfile.json

onError() {
  echo "Cleaning schemas/devfile.json"
  rm -f ${BASE_DIR}/schemas/devfile.json
}
trap 'onError' ERR

for jsonpatch in $(ls ${BASE_DIR}/devfile-support/transformation-rules/*.json 2> /dev/null)
do
  jsonpatch -i --indent 2 ${BASE_DIR}/schemas/devfile.json ${jsonpatch}
done
