#!/bin/bash
# This script adds license headers that are missing from go files


if ! command -v addlicense 2> /dev/null
then
  echo "error addlicense must be installed with this command: go install github.com/google/addlicense@latest" && exit 1
else
  echo 'addlicense -v -f license_header.txt **/*.go'
  addlicense -v -f license_header.txt $(find . -not -path '*/\.*' -not -path '*/vendor/*' -name '*.go')
fi


