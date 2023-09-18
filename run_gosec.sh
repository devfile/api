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

# This script runs the gosec scanner locally

if ! command -v gosec 2> /dev/null
then
  echo "error gosec must be installed with this command: go install github.com/securego/gosec/v2/cmd/gosec@v2.14.0" && exit 1
fi

gosec -no-fail -fmt=sarif -out=gosec.sarif -exclude-dir test  -exclude-dir generator  ./...
