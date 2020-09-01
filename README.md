# Kube-native API for cloud development workspaces specification

<div id="header" align="center">

[![Eclipse License](https://img.shields.io/badge/license-Eclipse-brightgreen.svg)](LICENSE)
[![Contribute](https://www.eclipse.org/che/contribute.svg)](https://che.openshift.io/f/?url=https://github.com/devfile/api)

</div>

Sources for this API are defined in Go code, starting from the
[devworkspace_types.go source file](pkg/apis/workspaces/v1alpha2/devworkspace_types.go)

From these Go sources, several files are generated:
- A Kubernetes Custom Resource Definition(CRD) with an embedded OpenApi schema,
- json schemas (in the [schemas](schemas) folder) generated from the above CRD, to specify the syntax of:
  - the DevWorkspace CRD itself;
  - the DevWorkspaceTemplate CRD (a workspace content, without runtime information);
  - the Devfile 2.0.0 format, which is generated from the `DevWorkspace` API.

Generated files are created by a build script (see section [How to build](#how-to-build)).

### Devfile 2.0.0 file format

A Subset of this `DevWorkspace` API defines a structure (workspace template content), which is also at the core of the **Devfile 2.0** format specification.
For more information about this, please look into the [Devfile support Readme](devfile-support/README.md)

The generated documentation of the Devfile 2.0 format, based on its json schema, is available here: https://devfile.github.io

## How to build

In order to build the CRD and the various schemas, you don't need to install any pre-requisite apart from `docker`.
In the root directory, just run the following command:

```
./docker-run.sh ./build.sh
``` 

## Specification status

This work is still in an early stage of specification, and the related API and schemas are still a draft proposal.

## Quickly open and test ?

In order to test existing or new Devfile 2.0 or DevWorkspace sample files in a self-service Che workspace (hosted on che.openshift.io), just click on the button below:

[![Contribute](https://che.openshift.io/factory/resources/factory-contribute.svg)](https://che.openshift.io/f/?url=https://github.com/devfile/api)

As soon as the workspace is opened, you should be able to:
- open the `yaml` files in the following folders:
  - `samples/`
  - `devfile-support/samples`
- have `yaml` language support (completion and documentation) based on the current Json schemas.
