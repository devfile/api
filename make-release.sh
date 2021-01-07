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

SCHEMA_VERSION=$1
K8S_VERSION=$2

if ! command -v hub > /dev/null; then
  echo "[ERROR] The hub CLI needs to be installed. See https://github.com/github/hub/releases"
  exit
fi
if [[ -z "${GITHUB_TOKEN}" ]]; then
  echo "[ERROR] The GITHUB_TOKEN environment variable must be set."
  exit
fi

if ! [[ "$SCHEMA_VERSION" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
	echo >&2 "$SCHEMA_VERSION isn't a valid semver tag for the schema. Aborting..."
	exit 1
fi


init() {
  BRANCH="${SCHEMA_VERSION%.*}.x"
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
  git checkout -B $SCHEMA_VERSION
}

setVersionAndBuild() {
  # Replace pre-release version with release version
  apply_sed "s#jsonschema:version=.*#jsonschema:version=${SCHEMA_VERSION}#g" pkg/apis/workspaces/$K8S_VERSION/doc.go #src/constants.ts

  # Generate the schema
  ./build.sh
}

commitChanges() {
  echo "[INFO] Pushing changes to $SCHEMA_VERSION branch"
  git add -A
  git commit -s -m "$1"
  git push origin $SCHEMA_VERSION
}

createReleaseBranch() {
  echo "[INFO] Create the release branch based on $SCHEMA_VERSION"
  git push origin $SCHEMA_VERSION
}

createPR() {
  echo "[INFO] Creating a PR"
  hub pull-request --base ${BRANCH} --head ${SCHEMA_VERSION} -m "$1"
}

bumpVersion() {
  IFS='.' read -a semver <<< "$SCHEMA_VERSION"
  MAJOR=${semver[0]}
  MINOR=${semver[1]}
  SCHEMA_VERSION=$MAJOR.$((MINOR+1)).0-alpha
}

updateVersionOnMaster() {
  # Checkout to a PR branch based on master to make our changes in
  git checkout -b $SCHEMA_VERSION
  
  # Set the schema version to the new version (with -alpha appended) and build the schemas
  setVersionAndBuild
  
  commitChanges "chore(post-release): bump schema version to ${SCHEMA_VERSION}"
}

compareMasterVersion() {
  # Parse the version passed in.
  IFS='.' read -a semver <<< "$SCHEMA_VERSION"
  MAJOR=${semver[0]}
  MINOR=${semver[1]}
  BUGFIX=${semver[2]}
  
  # Parse the version currently set in the schema
  latestVersion=`cat schemas/latest/jsonSchemaVersion.txt`
  IFS='.' read -a latestSemVer <<< "$latestVersion"
  local latestMajor=${latestSemVer[0]}
  local latestMinor=${latestSemVer[1]}
  local latestBugfix=$(echo ${latestSemVer[2]} | awk -F '-' '{print $1}')
  
  # Compare the new vers
  if ((latestMajor <= MAJOR)) && ((latestMinor <= MINOR)) && ((latestBugfix <= BUGFIX)); then
    return 0
  else
    return 1
  fi
}

run() {
  # Create the release branch and open a PR against the release branch, updating the release version
  echo "[INFO] Releasing a new ${SCHEMA_VERSION} version"
  checkoutToReleaseBranch
  setVersionAndBuild
  commitChanges "chore(release): release version ${SCHEMA_VERSION}"
  createReleaseBranch
  createPR "Release version ${SCHEMA_VERSION}"

  # If needed, bump the schema version in master and open a PR against master
  # Switch back to the master branch
  BRANCH=master
  resetChanges $BRANCH
  if compareMasterVersion; then
    echo "[INFO] Updating schema version on master to ${SCHEMA_VERSION}"
    bumpVersion
    updateVersionOnMaster
    createPR "Bump schema version on master to ${SCHEMA_VERSION}"
  else
    echo "[WARN] The passed in schema version ${SCHEMA_VERSION} is less than the current version on master, so not updating the master branch version"
    exit 
  fi
}

init
run
