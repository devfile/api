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
  echo "#### and https://stedolan.github.io/jq/download"
  echo "####"
  echo "###############"
  exit 1
fi

set -e

file1=$(mktemp --suffix=.json)
file2=$(mktemp --suffix=.json)

onError() {
  rm "$file1"
  rm "$file2"
}
trap 'onError' ERR

jq -S '.' "$1" > "$file1"
jq -S '.' "$2" > "$file2"
diff "$file1" "$file2"
