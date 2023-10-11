# Devfile Versioning and Release Process
This document outlines the versioning and release process for the Devfile spec.

## Reference  
This process summarizes parts from the Devfile API technical meeting slides presented by @davidfestal back in October, and I recommend having a read through of it: https://docs.google.com/presentation/d/1ohM1HzPB59a3ajvB7rVWJkKONLBAuT0mb5P20_A6vLs/edit#slide=id.g8fc722bef9_1_8

## Versioning

The following sections outline how we version the Devfile Kubernetes API, as well as the JSON schema generated from the API.

### Kubernetes API
The Devfile Kubernetes API (defined in https://github.com/devfile/api/) is the single source of truth for the Devfile spec. It’s defined under `pkg/apis/` and contains the Go structs that are used by the devworkspace operator and devfile library. It's also where the Devfile JSON schema is generated from. 

**Versioning Scheme**: [Kubernetes](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definition-versioning) (e.g. v1alpha1, v1beta2, v1, v2, etc)

**How to Update?**

   1) Add a new folder for the version in the [devfile/api](https://github.com/devfile/api/) repository under [pkg/apis/workspaces](https://github.com/devfile/api/tree/main/pkg/apis/workspaces). For example `pkg/apis/workspaces/v1` if bumping the K8S API version to `v1`.
   2) Add a schema and version in the CRD manifests
   3) Go through the JSON schema update process outlined below to update the JSON schema version.
   4) Update the devworkspace operator and devfile library to consume the Go structs in the new K8S API version, as needed.

**When to Update?**

As incrementing the Kubernetes API version for Devfile is a relatively heavy process, and affects the library, only update the K8s API version when absolutely needed (for **big** changes or backwards incompatible changes).
   - Backwards incompatible changes are defined as any change in the schema that would cause K8S API validation errors on the old resource version (e.g. removed fields or new mandatory fields without a default)
   - New, optional fields to the API do not necessarily require a version bump, but if it's a large change, it may be wise to bump the version.

### Devfile JSON Schema

As mentioned above, the Devfile JSON schema is generated from the Go structs defined in the Devfile Kubernetes API. The latest JSON schema for a given K8S API version is located under `schemas/<api-version>` in the [devfile/api repo](https://github.com/devfile/api/). 

**Versioning Scheme**: Semantic Versioning (major.minor.bugfix)

**How to Update?**
  
   1) Update the schema version string in the `// +devfile:jsonschema:version=<schema-version>` annotation in `pkg/apis/workspace/<k8s-api-version>/doc.go`
   2) Re-generate the json schema

**When to Update?** 

On each release of the schema, incremented based on the changes going in to the release. E.g.:
   
   - major == breaking / backwards incompatible changes. 
   - minor == larger and / or backwards compatible changes
   - bugfix == comments / bugfixes

K8S API version updates should also result in an appropriate increment of the schema version.


### Relationship Between K8s API version and JSON Schema Version

The Devfile JSON schema is generated from the Kubernetes API, and the version for the JSON schema is set in the doc.go file in the K8S API (`pkg/apis/workspace/<api-version>/doc.go`).

As we’re only updating the K8S API version when needed, but incrementing the schema version more frequently, this means that any given K8S API version may point to multiple, backwards-compatible, schema versions over its lifespan. The schema version under `schemas/<api-version>` in [devfile/api repo](https://github.com/devfile/api/) points to the matching JSON schema generated from the K8S API.

## Release Process
The following steps outline the steps done to release a new version of the Devfile schema and publish its schemas to the devfile.io website

   1) The release engineer tasked with creating the release clones the repository (and checks out the release branch if one already exists)

   2) The release engineer installs the `hub` CLI from https://github.com/github/hub if it is not already installed on their machine.

   3) `export GITHUB_TOKEN=<token>` is run, where `<token>` is a GitHub personal access token with `repo` permissions created from https://github.com/settings/tokens.

   4) `./make-release.sh <schema-version> <k8s-api-version>` is run:

      i) A release branch (the name corresponding to the schema version) is created, if one does not already exist.

      ii) The schema-version is updated in `pkg/apis/workspace/<api-version>/doc.go`.

      iii) New JSON schemas are generated to include the new schema version
      
      iv) A new commit including the changes

      v) A PR is opened to merge these changes into the release branch

      vi) The schema version on the main branch is bumped and a PR is opened, provided the schema version passed in is **not** older than the main branch schema version. 

   5) Once the release PR is approved and merged, the release engineer creates a new release on GitHub based off the release branch that was just created and includes the generated `devfile.json` as a release artifact. 
       - The tag `v{major}.{minor}.{bugfix}`, where the `{major}.{minor}.{bugfix}` corresponds to the devfile schema version, is used to enable the new version of the API to be pulled in as a Go module.

   6) Once the release is published, GitHub actions run to publish the new schema to devfile.io. The “stable” schema is also updated to point to the new schema.

   7) Make a release announcement on the devfile mailing list and slack channel

An example pull request, `make-release.sh` script and GitHub action can be found here:
- [Release Pull Request](https://github.com/devfile/api/pull/958)
- [make-release.sh](./make-release.sh)
- [release-schema.yaml](./.github/workflows/release-devfile-schema.yaml)

**Schema Store**

After releasing a new version, for example 2.2.0, we will also need to update the schemastore's [catalog.json](https://github.com/SchemaStore/schemastore/blob/master/src/api/json/catalog.json#L1119-L1132).

Open a PR to update the devfile entry to include the new devfile version and also update the default url to the latest version.

## Post Release

Create a pre-release tag pointing to the first commit hash following a published release.  In addition to preparing for the next release, this will also allow devfile clients to update their dependencies without the unintended side effects of version downgrading.  See [api issue 559](https://github.com/devfile/api/issues/599) for the background discussion.

1)  Create a new release tag following the [Github instructions](https://docs.github.com/en/repositories/releasing-projects-on-github/managing-releases-in-a-repository)
2)  Sync the pre-release tag name to the version found in the [latest schemas](https://github.com/devfile/api/blob/main/schemas/latest/devfile.json#L4) e.g. `{major}.{minor}.{bugfix}-alpha`
3)  Select the option `Set as a pre-release` before publishing

