#!/bin/bash
set -e

SCRIPT_DIR=`dirname $( readlink -m $( type -p ${0} ))`
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
    $WORK_DIR/gen/openapi/typescript.sh $WORK_DIR/typescript-models $WORK_DIR/config.sh

    sed -i 's/\"name\": \".*\"/"name": "@devfile\/api"/g' $WORK_DIR/typescript-models/package.json
    sed -i 's/\"description\": \".*\"/"description": "Typescript types for devfile api"/g' $WORK_DIR/typescript-models/package.json
    sed -i 's/\"repository\": \".*\"/"repository": "devfile\/api"/g' $WORK_DIR/typescript-models/package.json
    sed -i 's/\"license\": \".*\"/"license": "Apache-2.0"/g' $WORK_DIR/typescript-models/package.json
    sed -i 's/\"@types\/bluebird\": \".*\"/"@types\/bluebird": "3.5.21"/g' $WORK_DIR/typescript-models/package.json
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

generate_swagger_json
k8s_client_gen
generate_typescript_metadata
build_typescript_model

echo "[INFO] Typescript model is successfully generated and verified."
echo "[INFO] It can be accessed at $WORK_DIR/typescript-models"
