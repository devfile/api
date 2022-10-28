# Kube-native API for cloud development workspaces specification

<div id="header" align="center">

[![Apache License](https://img.shields.io/badge/license-Apache-brightgreen.svg)](LICENSE)
[![Contribute](https://img.shields.io/badge/developer-workspace-525C86?logo=eclipse-che&labelColor=FDB940)](https://workspaces.openshift.com/f?url=https://github.com/devfile/api)

</div>

Sources for this API are defined in Go code, starting from the
[devworkspace_types.go source file](pkg/apis/workspaces/v1alpha2/devworkspace_types.go)

From these Go sources, several files are generated:
- A Kubernetes Custom Resource Definition(CRD) with an embedded OpenApi schema,
- json schemas (in the [schemas](schemas) folder) generated from the above CRD, to specify the syntax of:
  - the DevWorkspace CRD itself;
  - the DevWorkspaceTemplate CRD (a devworkspace content, without runtime information);
  - the Devfile 2.0.0 format, which is generated from the `DevWorkspace` API.

Generated files are created by a build script (see section [How to build](#how-to-build)).

## Devfile 2.0.0 file format

A Subset of this `DevWorkspace` API defines a structure (workspace template content), which is also at the core of the **Devfile 2.0** format specification.
For more information about this, please look into the [Devfile support README](https://github.com/devfile/registry-support/blob/main/README.md)

The generated documentation of the Devfile 2.0 format, based on its json schema, is available [here](https://devfile.github.io).

Typescript model is build on each commit of main branch and available as an [NPM package](https://www.npmjs.com/package/@devfile/api).

## Release

Release details and process are found in [Devfile Release](RELEASE.md)

## How to build

In order to build the CRD and the various schemas, you don't need to install any pre-requisite apart from `docker`.
In the root directory, just run the following command:

```console
bash ./docker-run.sh ./build.sh
```

### Typescript model

Typescript model is generated based on JSON Schema with help of <https://github.com/kubernetes-client/gen>.
To generate them locally run:

```console
bash ./build/typescript-model/generate.sh
```

## Specification status

This work is still in an early stage of specification, and the related API and schemas are still a draft proposal.

## Quickly open and test ?

In order to test existing or new Devfile 2.0 or DevWorkspace sample files in a self-service Che workspace (hosted on che.openshift.io), just click on the button below:

[![Contribute](https://img.shields.io/badge/developer-workspace-525C86?logo=eclipse-che&labelColor=FDB940)](https://workspaces.openshift.com/f?url=https://github.com/devfile/api)

As soon as the devworkspace is opened, you should be able to:
- open the `yaml` files in the following folders:
  - `samples/`
  - `devfile-support/samples`
- have `yaml` language support (completion and documentation) based on the current Json schemas.

## Contributing

Please see our [contributing.md](./CONTRIBUTING.md).

## License

Apache License 2.0, see [LICENSE](./LICENSE) for details.
