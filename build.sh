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

bash "${BASE_DIR}"/build-generator.sh

cd "${BASE_DIR}"

# We have to generate plugin overrides before generating parent overrides, as the parent overrides
# require the overrides generated for plugins

echo "Generating Plugin Overrides"

generator/build/generator "overrides:isForPluginOverrides=true" "paths=./pkg/apis/workspaces/v1alpha2"

echo "Generating Parent Overrides"

generator/build/generator "overrides:isForPluginOverrides=false" "paths=./pkg/apis/workspaces/v1alpha2"

echo "Validating K8S API Source code"

generator/build/generator "validate" "paths=./pkg/apis/workspaces/v1alpha2"

echo "Generating Interface Implementations"

generator/build/generator "interfaces" "paths=./pkg/apis/workspaces/v1alpha2"

echo "Generating K8S CRDs"

generator/build/generator "crds" "output:crds:artifacts:config=crds" "paths=./pkg/apis/workspaces/v1alpha2;./pkg/apis/workspaces/v1alpha1"

echo "Generating DeepCopy implementations"

generator/build/generator "deepcopy" "paths=./pkg/apis/workspaces/v1alpha2;./pkg/apis/workspaces/v1alpha1"

echo "Generating JsonSchemas"

generator/build/generator "schemas" "output:schemas:artifacts:config=schemas" "paths=./pkg/apis/workspaces/v1alpha2"

echo "Generating Getter Implementations"

generator/build/generator "getters" "paths=./pkg/apis/workspaces/v1alpha2"

echo "Finished generation of required GO sources, K8S CRDs, and Json Schemas"
