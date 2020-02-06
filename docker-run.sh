#!/bin/bash
#
# Copyright (c) 2019 Red Hat, Inc.
# This program and the accompanying materials are made
# available under the terms of the Eclipse Public License 2.0
# which is available at https://www.eclipse.org/legal/epl-2.0/
#
# SPDX-License-Identifier: EPL-2.0
#
# Contributors:
#   Red Hat, Inc. - initial API and implementation

set

# git ROOT directory used to mount filesystem
GIT_ROOT_DIRECTORY=$(git rev-parse --show-toplevel)
GO_MODULE=$(grep -e 'module ' ${GIT_ROOT_DIRECTORY}/go.mod | sed -e 's/module //')
WORKDIR="/home/user/go/src/${GO_MODULE}"
# Container image
IMAGE_NAME="che-incubator/devworkspace-build-prerequisites"

# Operator SDK
OPERATOR_SDK_VERSION=v0.12.0

init() {
  BLUE='\033[1;34m'
  GREEN='\033[0;32m'
  RED='\033[0;31m'
  NC='\033[0m'
  BOLD='\033[1m'
}

check() {
  if [ $# -eq 0 ]; then
    printf "%bError: %bNo script provided. Command is $ docker-run.sh <script-to-run> [optional-arguments-of-script-to-run]\n" "${RED}" "${NC}"
    exit 1
  fi
  echo "check $1"
  if [ ! -f "$1" ] || [ ! -x "$1"]; then
    printf "%bError: %bscript %b provided does not exist. Command is $ docker-run.sh <script-to-run> [optional-arguments-of-script-to-run]\n" "${RED}" "${NC}" "${1}"
    exit 1
  fi
}

# Build image
build() {
  printf "%bBuilding image %b${IMAGE_NAME}${NC}..." "${BOLD}" "${BLUE}"
  if docker build --build-arg OPERATOR_SDK_VERSION=${OPERATOR_SDK_VERSION} -t ${IMAGE_NAME} > docker-build-log 2>&1 -<<EOF
FROM python:3-alpine
ARG OPERATOR_SDK_VERSION
ENV GOROOT /usr/lib/go

RUN apk add --no-cache --update curl bash jq go \
&& pip3 install yq \
&& pip3 install jsonpatch

RUN curl -JL https://github.com/operator-framework/operator-sdk/releases/download/${OPERATOR_SDK_VERSION}/operator-sdk-${OPERATOR_SDK_VERSION}-x86_64-linux-gnu -o /bin/operator-sdk && chmod a+x /bin/operator-sdk
RUN mkdir -p /home/user/go && chmod -R a+w /home/user
ENV HOME /home/user
ENV GOPATH /home/user/go
WORKDIR ${WORKDIR}
EOF
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
  if docker run --user $(id -u):$(id -g) --rm -it -v "${GIT_ROOT_DIRECTORY}":"${WORKDIR}" --entrypoint=/bin/bash ${IMAGE_NAME} -- "$@"
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
