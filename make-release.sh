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
    echo "[INFO] $BRANCH does not exist. Will create a new one from main."
    resetChanges main
    git push origin main:$BRANCH
  fi
  git checkout -B $SCHEMA_VERSION
}

setVersionAndBuild() {
  # Replace pre-release version with release version
  apply_sed "s#jsonschema:version=.*#jsonschema:version=${SCHEMA_VERSION}#g" pkg/apis/workspaces/$K8S_VERSION/doc.go #src/constants.ts

  # Generate the schema
  bash ./build.sh
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

updateVersionOnMain() {
  # Checkout to a PR branch based on main to make our changes in
  git checkout -b $SCHEMA_VERSION
  
  # Set the schema version to the new version (with -alpha appended) and build the schemas
  setVersionAndBuild
  
  commitChanges "chore(post-release): bump schema version to ${SCHEMA_VERSION}"
}

compareMainVersion() {
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

  # If needed, bump the schema version in main and open a PR against main
  # Switch back to the main branch
  BRANCH=main
  resetChanges $BRANCH
  if compareMainVersion; then
    echo "[INFO] Updating schema version on main to ${SCHEMA_VERSION}"
    bumpVersion
    updateVersionOnMain
    createPR "Bump schema version on main to ${SCHEMA_VERSION}"
  else
    echo "[WARN] The passed in schema version ${SCHEMA_VERSION} is less than the current version on main, so not updating the main branch version"
    exit 
  fi
}

init
run
