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

VERSION=$1
API_VERSION=$2

if ! command -v hub > /dev/null; then
  echo "[ERROR] The hub CLI needs to be installed. See https://github.com/github/hub/releases"
  exit
fi
if [[ -z "${GITHUB_TOKEN}" ]]; then
  echo "[ERROR] The GITHUB_TOKEN environment variable must be set."
  exit
fi

if ! [[ "$VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
	echo >&2 "$VERSION isn't a valid semver tag for the schema. Aborting..."
	exit 1
fi


init() {
  BRANCH="${VERSION%.*}.x"
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
  if git ls-remote -q --heads | grep -q $BRANCH ; then
    echo "[INFO] $BRANCH exists."
    resetChanges $BRANCH
  else
    echo "[INFO] $BRANCH does not exist. Will create a new one from master."
    resetChanges master
    git push origin master:$BRANCH
  fi
  git checkout -B $VERSION
}

setVersionAndBuild() {
  # Replace pre-release version with release version
  apply_sed "s#jsonschema:version=.*#jsonschema:version=${VERSION}#g" pkg/apis/workspaces/$API_VERSION/doc.go #src/constants.ts

  # Generate the schema
  ./build.sh
}

commitChanges() {
  echo "[INFO] Pushing changes to $VERSION branch"
  git add -A
  git commit -s -m "$1"
  git push origin $VERSION
}

createReleaseBranch() {
  echo "[INFO] Create the release branch based on $VERSION"
  git push origin $VERSION
}

createPR() {
  echo "[INFO] Creating a PR"
  hub pull-request --base ${BRANCH} --head ${VERSION} -m "Release version ${VERSION}"
}

bumpVersion() {
  IFS='.' read -a semver <<< "$VERSION"
  MAJOR=${semver[0]}
  MINOR=${semver[1]}
  VERSION=$MAJOR.$((MINOR+1)).0-alpha
}

updateVersionOnMaster() {
  # Switch back to the master branch
  BRANCH=master
  resetChanges $BRANCH
  git checkout -b $VERSION

  # Set the schema version on master to the new version (with -alpha appended) and build the schemas
  setVersionAndBuild
  
  commitChanges "chore(post-release): bump schema version to ${VERSION}"
}

run() {
  # Create the release branch and open a PR against the release branch, updating the release version
  echo "[INFO] Releasing a new $VERSION version"
  checkoutToReleaseBranch
  setVersionAndBuild
  commitChanges "chore(release): release version ${VERSION}"
  createReleaseBranch
  createPR "Release version ${VERSION}"

  # Bump the schema version in master and open a PR against master
  echo "[INFO] Updating schema version on master to $VERSION"
  bumpVersion
  updateVersionOnMaster
  createPR "Bump schema version to ${VERSION}"
}

init
run
