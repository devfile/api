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
