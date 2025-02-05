# Kube-native API for cloud development workspaces specification

<div id="header" align="center">

[![Apache License](https://img.shields.io/badge/license-Apache-brightgreen.svg)](LICENSE)
[![Contribute](https://img.shields.io/badge/developer-workspace-525C86?logo=eclipse-che&labelColor=FDB940)](https://workspaces.openshift.com/f?url=https://github.com/devfile/api)
[![OpenSSF Best Practices](https://www.bestpractices.dev/projects/8179/badge)](https://www.bestpractices.dev/projects/8179)
[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/devfile/api/badge)](https://securityscorecards.dev/viewer/?uri=github.com/devfile/api)
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

You can read the available generated [documentation of the Devfile 2.0 format](https://devfile.io/docs/2.3.0/devfile-schema), based on its json schema.

Typescript model is build on each commit of main branch and available as an [NPM package](https://www.npmjs.com/package/@devfile/api).

## Release

Release details and process are found in [Devfile Release](RELEASE.md)

## How to build

For information about building this project visit [CONTRIBUTING.md](./CONTRIBUTING.md#building).

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

### Adding License Headers

[`license_header`](./license_header.txt) contains the license header to be contained under all source files. For Go sources, this can be included by running `bash add_licenses.sh`.

Ensure `github.com/google/addlicense` is installed by running `go install github.com/google/addlicense@latest`.
