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

# Allow setting of podman environment var in the script runtime
shopt -s expand_aliases
set -eux

# git ROOT directory used to mount filesystem
GIT_ROOT_DIRECTORY=$(git rev-parse --show-toplevel)
GO_MODULE=$(grep -e 'module ' ${GIT_ROOT_DIRECTORY}/go.mod | sed -e 's/module //')
WORKDIR="/projects/src/${GO_MODULE}"
# Container image
IMAGE_NAME="quay.io/devfile/kubernetes-api-build-prerequisites:latest"

# For users who want to use podman this enables the alias to work throughout the scripts runtime
USE_PODMAN=${USE_PODMAN:-false}
if [[ ${USE_PODMAN} == true ]]; then
  alias docker=podman
  echo "using podman as container engine"
fi

init() {
  BLUE='\033[1;34m'
  GREEN='\033[0;32m'
  RED='\033[0;31m'
  NC='\033[0m'
  BOLD='\033[1m'
}

check() {
  if [ $# -eq 0 ]; then
    printf "%bError: %bNo script provided. Command is $ docker-run.sh push|<script-to-run> [optional-arguments-of-script-to-run]\n" "${RED}" "${NC}"
    exit 1
  fi
  echo "check $1"
  if [ ! -f "$1" ] || [ ! -x "$1" ]; then
    printf "%bError: %bscript %b provided does not exist. Command is $ docker-run.sh <script-to-run> [optional-arguments-of-script-to-run]\n" "${RED}" "${NC}" "${1}"
    exit 1
  fi
}

# Build image
build() {
  printf "%bBuilding image %b${IMAGE_NAME}${NC}..." "${BOLD}" "${BLUE}"
  if docker build -t ${IMAGE_NAME} .devfile/ > docker-build-log 2>&1
  then
    printf "%b[OK]%b\n" "${GREEN}" "${NC}"
    rm docker-build-log
  else
    printf "%bFailure%b\n" "${RED}" "${NC}"
    cat docker-build-log
    exit 1
  fi
}

run() {
  printf "%bRunning%b $*\n" "${BOLD}" "${NC}"
  if docker run --user $(id -u):$(id -g) --rm -v "${GIT_ROOT_DIRECTORY}":"${WORKDIR}" ${IMAGE_NAME} -- bash -c "cd \"${WORKDIR}\" && $@"
  then
    printf "Script execution %b[OK]%b\n" "${GREEN}" "${NC}"
  else
    printf "%bFail to run the script%b\n" "${RED}" "${NC}"
    exit 1
  fi
}

init "$@"
check "$@"
build "$@"
run "$@"
