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

set -e

SHORT_NAME="$(uname -s)"
if [ "$(uname)" == "Darwin" ]; then
    SCRIPT_DIR=`dirname $( realpath $( type -p ${0} ))`
else
    SCRIPT_DIR=`dirname $( readlink -m $( type -p ${0} ))`
fi
WORK_DIR=${SCRIPT_DIR}/workdir
echo "[INFO] Using the following folder to store all build files ${SCRIPT_DIR}/workdir"
mkdir -p $WORK_DIR

GEN_REVISION=b32dcd6dc9c1c0c4fcf227c9539ae9ff0530b936

k8s_client_gen() {
    [ ! -d $WORK_DIR/gen ] && git clone https://github.com/kubernetes-client/gen.git $WORK_DIR/gen || echo "kubernetes-client/gen is already cloned into $WORK_DIR/gen"

    echo "[INFO] Checking out gen to ${GEN_REVISION}"
    pushd "$WORK_DIR/gen"
    git checkout ${GEN_REVISION}
    popd

    echo "[INFO] preparing config files for gen"
    # Remove the contents of custom objects spec so that we aren't bundling any extra objects
    echo "{}" > $WORK_DIR/gen/openapi/custom_objects_spec.json

    cat <<EOF > ${WORK_DIR}/config.sh
export KUBERNETES_BRANCH=''
export CLIENT_VERSION=''
export PACKAGE_NAME=''
export USERNAME=''
export REPOSITORY=''
EOF
    echo "[INFO] Lauching gen to generate typescript files based on swagger json"
    export OPENAPI_SKIP_FETCH_SPEC=true
    export OPENAPI_GENERATOR_COMMIT="v6.3.0"
    bash $WORK_DIR/gen/openapi/typescript.sh $WORK_DIR/typescript-models $WORK_DIR/config.sh

    apply_sed 's/\"name\": \".*\"/"name": "@devfile\/api"/g' $WORK_DIR/typescript-models/package.json
    apply_sed 's/\"description\": \".*\"/"description": "Typescript types for devfile api"/g' $WORK_DIR/typescript-models/package.json
    apply_sed 's/\"repository\": \".*\"/"repository": "devfile\/api"/g' $WORK_DIR/typescript-models/package.json
    apply_sed 's/\"license\": \".*\"/"license": "Apache-2.0"/g' $WORK_DIR/typescript-models/package.json
    apply_sed 's/\"@types\/bluebird\": \".*\"/"@types\/bluebird": "3.5.21"/g' $WORK_DIR/typescript-models/package.json
    
    echo "" > $WORK_DIR/typescript-models/.npmignore
    echo "[INFO] Generated typescript model which now is available in $WORK_DIR/typescript-models"
}

generate_swagger_json() {
    echo "[INFO] Generating Swagger JSON..."
    python3 $SCRIPT_DIR/generate-swagger-json.py
    rm -rf $WORK_DIR/typescript-models
    mkdir -p $WORK_DIR/typescript-models
    mv swagger.json $WORK_DIR/typescript-models/swagger.json.unprocessed
    echo "[INFO] Generating Swagger JSON. It's in $WORK_DIR/typescript-models/swagger.json.unprocessed"
}

generate_typescript_metadata() {
    echo "[INFO] Generating typescript constants from crds ..."
    mkdir -p $WORK_DIR/typescript-models/constants
    python3 $SCRIPT_DIR/generate-metadata.py -p $WORK_DIR/typescript-models
    echo "[INFO] Finished generating typescript constant from crds. They are available in $WORK_DIR/typescript-models/constants"
}

build_typescript_model() {
    echo "[INFO] Verify that generated model is buildable..."
    cd $WORK_DIR/typescript-models
    yarn && yarn build || "[ERROR] Generated typescript model failed to build. Check it at $WORK_DIR/typescript-models."
    echo "[INFO] Done."
}

apply_sed(){
    if [ "$(uname)" == "Darwin" ]; then
        sed -i '' "$@"
    else
        sed -i "$@"
    fi
}

generate_swagger_json
k8s_client_gen
generate_typescript_metadata
build_typescript_model

echo "[INFO] Typescript model is successfully generated and verified."
echo "[INFO] It can be accessed at $WORK_DIR/typescript-models"
