#!/bin/bash
# This script runs the gosec scanner locally

if ! command -v gosec 2> /dev/null
then
  echo "error gosec must be installed with this command: go install github.com/securego/gosec/v2/cmd/gosec@v2.14.0" && exit 1
fi

gosec -no-fail -fmt=sarif -out=gosec.sarif -exclude-dir test  -exclude-dir generator  ./...
