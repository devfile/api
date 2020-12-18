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

# Based on https://github.com/che-incubator/chectl/blob/master/make-release.sh

set -e
set -u

usage ()
{   echo "Usage: ./make-release.sh <schema-version> <k8s-api-version>"
    exit
}

if [[ $# -lt 2 ]]; then usage; fi

if ! command -v hub > /dev/null; then
  echo "[ERROR] The hub CLI needs to be installed. See https://github.com/github/hub/releases"
  exit
fi
if [[ -z "${GITHUB_TOKEN}" ]]; then
  echo "[ERROR] The GITHUB_TOKEN environment variable must be set."
  exit
fi

VERSION=$1
API_VERSION=$2

init() {
  BRANCH=$(echo $VERSION | sed 's/.$/x/')
  echo $BRANCH
}

apply_sed() {
    SHORT_UNAME=$(uname -s)
  if [ "$(uname)" == "Darwin" ]; then
    sed -i '' "$1" "$2"
  elif [ "${SHORT_UNAME:0:5}" == "Linux" ]; then
    sed -i "$1" "$2"
  fi
}

resetChanges() {
  echo "[INFO] Reset changes in $1 branch"
  git reset --hard
  git checkout $1
  git fetch origin --prune
  git pull origin $1
}

checkoutToReleaseBranch() {
  echo "[INFO] Checking out to $BRANCH branch."
  local branchExist=$(git ls-remote -q --heads | grep $BRANCH | wc -l)
  if [[ $branchExist == 1 ]]; then
    echo "[INFO] $BRANCH exists."
    resetChanges $BRANCH
  else
    echo "[INFO] $BRANCH does not exist. Will be created a new one from master."
    resetChanges master
    git push origin master:$BRANCH
  fi
  git checkout -B $VERSION
}

release() {
  echo "[INFO] Releasing a new $VERSION version"

  # replace nightly versions by release version
  apply_sed "s#jsonschema:version=.*#jsonschema:version=${VERSION}#g" pkg/apis/workspaces/$API_VERSION/doc.go #src/constants.ts

  # Generate the schema
  ./build.sh
}

commitChanges() {
  echo "[INFO] Pushing changes to $VERSION branch"
  git add -A
  git commit -s -m "chore(release): release version ${VERSION}"
  git push origin $VERSION
}

createReleaseBranch() {
  echo "[INFO] Create the release branch based on $VERSION"
  git push origin $VERSION:release -f
}

createPR() {
  echo "[INFO] Creating a PR"
  hub pull-request --base ${BRANCH} --head ${VERSION} -m "Release version ${VERSION}"
}

run() {
  checkoutToReleaseBranch
  release
  commitChanges
  createReleaseBranch
  createPR
}

init
run
