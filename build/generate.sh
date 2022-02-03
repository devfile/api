#!/bin/bash
set -e

SCRIPT_DIR=`dirname $( readlink -m $( type -p ${0} ))`
WORK_DIR=${SCRIPT_DIR}/workdir
echo "[INFO] Using the following folder to store all build files ${SCRIPT_DIR}/workdir"
mkdir -p $WORK_DIR

GEN_REVISION=a3aef4de7a1d5dab72021aa282fffd8bc8a022ca

LANG=typescript
PACKAGE=''
while getopts l:p: flag; do
	case "$flag" in
		l) LANG=${OPTARG};;
		p) PACKAGE=${OPTARG};; 
	esac
done

MODELS_DIR="${WORK_DIR}/${LANG}-models"
mkdir -p $MODELS_DIR

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
export PACKAGE_NAME="$PACKAGE"
export USERNAME=''
export REPOSITORY=''
EOF
    echo "[INFO] Lauching gen to generate $LANG files based on swagger json"
    export OPENAPI_SKIP_FETCH_SPEC=true
    $WORK_DIR/gen/openapi/$LANG.sh $MODELS_DIR $WORK_DIR/config.sh

    if [ "$LANG" = "typescript" ]; then
        sed -i 's/\"name\": \".*\"/"name": "@devfile\/api"/g' $MODELS_DIR/package.json
        sed -i 's/\"description\": \".*\"/"description": "Typescript types for devfile api"/g' $MODELS_DIR/package.json
        sed -i 's/\"repository\": \".*\"/"repository": "devfile\/api"/g' $MODELS_DIR/package.json
        sed -i 's/\"license\": \".*\"/"license": "EPL-2.0"/g' $MODELS_DIR/package.json
        sed -i 's/\"@types\/bluebird\": \".*\"/"@types\/bluebird": "3.5.21"/g' $MODELS_DIR/package.json
        echo "" > $MODELS_DIR/.npmignore
    fi
    echo "[INFO] Generated $LANG model which now is available in $MODELS_DIR"
}

generate_swagger_json() {
    echo "[INFO] Generating Swagger JSON..."
    python3 $SCRIPT_DIR/generate-swagger-json.py
    rm -rf $MODELS_DIR
    mkdir -p $MODELS_DIR
    mv swagger.json $MODELS_DIR/swagger.json.unprocessed
    echo "[INFO] Generating Swagger JSON. It's in $MODELS_DIR/swagger.json.unprocessed"
}

generate_typescript_metadata() {
    echo "[INFO] Generating typescript constants from crds ..."
    mkdir -p $MODELS_DIR/constants
    python3 $SCRIPT_DIR/typescript/generate-metadata.py -p $MODELS_DIR
    echo "[INFO] Finished generating typescript constant from crds. They are available in $MODELS_DIR/constants"
}

build_typescript_model() {
    echo "[INFO] Verify that generated model is buildable..."
    cd $MODELS_DIR
    yarn && yarn build || "[ERROR] Generated typescript model failed to build. Check it at $MODELS_DIR."
    echo "[INFO] Done."
}

generate_swagger_json
k8s_client_gen
if [ "$LANG" = "typescript" ]; then
    generate_typescript_metadata
    build_typescript_model
fi

echo "[INFO] $LANG model is successfully generated and verified."
echo "[INFO] It can be accessed at $MODELS_DIR"
