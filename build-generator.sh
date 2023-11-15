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

set -e

CURRENT_DIR=$(pwd)
BASE_DIR=$(cd "$(dirname "$0")" && pwd)

onError() {
  cd "${CURRENT_DIR}"
}
trap 'onError' ERR

echo "Building generator"

cd "${BASE_DIR}/generator"
go generate ./...
GOFLAGS="-buildvcs=false" go build -o build/generator
