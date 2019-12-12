#!/bin/bash

operator-sdk-v0.12.0 generate k8s
operator-sdk-v0.12.0 generate openapi
yq '.spec.validation.openAPIV3Schema.properties.spec.properties.template' deploy/crds/workspaces.ecd.eclipse.org_devworkspaces_crd.yaml > generated/schema.json
