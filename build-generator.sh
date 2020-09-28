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

echo "Building generator"

cd "${BASE_DIR}/generator"
go generate ./...
go build -o build/generator
