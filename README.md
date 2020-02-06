# DevWorkspace API

K8S-native Api for a cloud develoment workspace specification [Draft proposal]

Sources for this API are defined in Go code, starting from the
[devworkspace_types.go source file](pkg/apis/workspaces/v1alpha1/devworkspace_types.go)

From these Go sources, several files are generated:
- A Kubernetes Custom Resource Definition with an embedded OpenApi schema,
- json schemas (in the [schemas](schemas) folder) generated from the above Custom Resource Definition, to specify the syntax of:
  - the DevWorkspace CRD itself
  - the DevWorkspaceTemplate CRD (a workspace content, without runtime information),
  - the Devfile 2.0.0 format, which is generated from the `DevWorkspace` API.

Generated files are created by a build script (see section [How to build](#how-to-build)).

### Devfile 2.0 support

A Subset of this `DevWorkspace` API defines a structure (workspace template content), which is also at the core of the **Devfile 2.0** format specification.
For more information about this, please look into
the [Devfile support Readme](devfile-support/README.md)

## How to build

In order to build the CRD and the various schemas, you don't need to install any pre-requisite apart from `docker`.
In the root directory, just run the following command:

```
./docker-run.sh build.sh
``` 
