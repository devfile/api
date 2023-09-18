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

# This script adds license headers that are missing from go files


if ! command -v addlicense 2> /dev/null
then
  echo "error addlicense must be installed with this command: go install github.com/google/addlicense@latest" && exit 1
else
  echo 'addlicense -v -f license_header.txt **/*.go'
  addlicense -v -f license_header.txt $(find . -not -path '*/\.*' -not -path '*/vendor/*' -not -name 'zz_generated.*.go' -name '*.go')
fi


