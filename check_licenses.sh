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

# This script checks if license headers that are missing/invalid from go files

if ! command -v addlicense 2> /dev/null
then
    echo "error addlicense must be installed with this command: go install github.com/google/addlicense@latest" && exit 1
else
    files=($(addlicense -check -v -f license_header.txt $(find . -not -path '*/\.*' -not -path '*/vendor/*' -not -name 'zz_generated.*.go' -name '*.go')))
    if [[ $? != 0 ]] && [[ ${#files[@]} -eq 0 ]]
    then
        echo "addheader check failed to run "
        exit 1
    elif [[ ${#files[@]} -gt 0 ]]
    then
        echo "The following files do not have valid license headers:"
        for file in ${files[@]}
        do
            echo ${file}
        done
        exit 1
    else
        echo "license headers are valid"
    fi
fi
