# Kube-native API for cloud development workspaces specification

<div id="header" align="center">

[![Eclipse License](https://img.shields.io/badge/license-Eclipse-brightgreen.svg)](LICENSE)
[![Contribute](https://che.openshift.io/factory/resources/factory-contribute.svg)](https://che.openshift.io/f/?url=https://github.com/devfile/kubernetes-api)

</div>

Sources for this API are defined in Go code, starting from the
[devworkspace_types.go source file](pkg/apis/workspaces/v1alpha1/devworkspace_types.go)

From these Go sources, several files are generated:
- A Kubernetes Custom Resource Definition(CRD) with an embedded OpenApi schema,
- json schemas (in the [schemas](schemas) folder) generated from the above CRD, to specify the syntax of:
  - the DevWorkspace CRD itself;
  - the DevWorkspaceTemplate CRD (a workspace content, without runtime information);
  - the Devfile 2.0.0 format, which is generated from the `DevWorkspace` API.

Generated files are created by a build script (see section [How to build](#how-to-build)).

### Compliance Table

| v2 feature  | odo | devworkspace controller | eclipse che |
| ------------- | ------------- | ------------- | ------------- |
| [Starter Projects](https://github.com/che-incubator/devworkspace-api/issues/42)  |  |  |  |
| [Component is a polymophic type](https://github.com/che-incubator/devworkspace-api/issues/4)  |  |  |  |
| [Shared Volumes Across Components](https://github.com/che-incubator/devworkspace-api/issues/19)  |  |  |  |
| [Out of Main Pod Compoenents](https://github.com/devfile/kubernetes-api/issues/48)  |  |  |  |
| [Replace Alias with Name](https://github.com/che-incubator/devworkspace-api/issues/9)  |  |  |  |
| [Renaming dockerimage component type](https://github.com/che-incubator/devworkspace-api/issues/8)  |  |  |  |
| [Specify sources path for containers](https://github.com/che-incubator/devworkspace-api/issues/17)  |  |  |  |
| [Specify size of volume for component](https://github.com/che-incubator/devworkspace-api/issues/14)  |  |  |  |
| [Containers endpoints (routes/ingresses)](https://github.com/che-incubator/devworkspace-api/issues/33)  |  |  |  |
| [New plugins spec](https://github.com/che-incubator/devworkspace-api/issues/31)  |  |  |  |
| [Command Groups: build,run,test,debug](https://github.com/che-incubator/devworkspace-api/issues/27)  |  |  |  |
| [Apply Command](https://github.com/devfile/kubernetes-api/issues/56)  |  |  |  |
| [Environment Varibables for a Specific Command](https://github.com/che-incubator/devworkspace-api/issues/21)  |  |  |  |
| [Renaming workdir into workingDir](https://github.com/che-incubator/devworkspace-api/issues/22)  |  |  |  |
| [Id and label for Composite Commands](https://github.com/che-incubator/devworkspace-api/issues/18)  |  |  |  |
| [Run exec Commands as Specific User](https://github.com/che-incubator/devworkspace-api/issues/34)  |  |  |  |
| [Devfile metadata: add a link to an external website](https://github.com/che-incubator/devworkspace-api/issues/38)  |  |  |  |
| [Stack/Devfile Matching Rules](https://github.com/che-incubator/devworkspace-api/issues/40)  |  |  |  |
| [Devfile parents](https://github.com/che-incubator/devworkspace-api/issues/25) |  |  |  |
| [Events](https://github.com/che-incubator/devworkspace-api/issues/32)  |  |  |  |

### Devfile 2.0.0 file format

A Subset of this `DevWorkspace` API defines a structure (workspace template content), which is also at the core of the **Devfile 2.0** format specification.
For more information about this, please look into the [Devfile support Readme](devfile-support/README.md)

The generated documentation of the Devfile 2.0 format, based on its json schema, is available here: https://devfile.github.io/website/

## How to build

In order to build the CRD and the various schemas, you don't need to install any pre-requisite apart from `docker`.
In the root directory, just run the following command:

```
./docker-run.sh build.sh
``` 

## Specification status

This work is still in an early stage of specification, and the related API and schemas are still a draft proposal.

## Quickly open and test ?

In order to test existing or new Devfile 2.0 or DevWorkspace sample files in a self-service Che workspace (hosted on che.openshift.io), just click on the button below:

[![Contribute](https://che.openshift.io/factory/resources/factory-contribute.svg)](https://che.openshift.io/f/?url=https://github.com/devfile/kubernetes-api)

As soon as the workspace is opened, you should be able to:
- open the `yaml` files in the following folders:
  - `samples/`
  - `devfile-support/samples`
- have `yaml` language support (completion and documentation) based on the current Json schemas.
