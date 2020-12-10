# Devfile Versioning and Release Process
This design document outlines the proposed versioning and release process for the Devfile spec. 

This document summarizes parts from he Devfile API technical meeting slides presented by @davidfestal back in October, and I recommend having a read through of it: https://docs.google.com/presentation/d/1ohM1HzPB59a3ajvB7rVWJkKONLBAuT0mb5P20_A6vLs/edit#slide=id.g8fc722bef9_1_8

## Versioning

The following sections outline how we version the Devfile Kubernetes API, as well as the JSON schema generated from the API.

### Kubernetes API
The Devfile Kubernetes API (defined in https://github.com/devfile/api/) is the single source of truth for the Devfile spec. It’s defined under `pkg/apis/` and contains the Go structs that are used by the devworkspace operator and devfile library. It's also where the Devfile JSON schema is generated from. 

**Versioning Scheme**: [Kubernetes](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definition-versioning) (e.g. v1alpha1, v1beta2, v1, v2, etc)

**Update Process**: 

   1) Add a new folder in the source repository (`pkg/apis/workspaces` in `devfile/api`)
   2) Add a schema and version in the CRD manifests
   3) Generate the JSON schema from the API. New JSON schema will be located under `schemas/latest`.
   4) Update the devworkspace operator and devfile library to consume the Go structs in the new API version, as needed.

**When to Update?**  As incrementing the Kubernetes API version for Devfile is a relatively heavy process, and affects the library, only update the K8s API version when absolutelyneeded  - for **big** changes or backwards incompatible changes.

### Devfile JSON Schema

As mentioned above, the Devfile JSON schema is generated from the Go structs defined in the Devfile Kubernetes API. The latest JSON schema for a given API version is located under `schemas/<api-version>` in the [devfile/api repo](https://github.com/devfile/api/). 

**Versioning Scheme**: Semver (major.minor.bugfix)

**Update Process**: 
  
   1) Update the schema version string in the `// +devfile:jsonschema:version=<schema-version>` annotation in `pkg/apis/workspace/<k8s-api-version>/doc.go`
   2) Re-generate the json schema
   3) Publish new schema to devfile.io website

**When to Update?** On each release of the schema, incremented based on the changes going in to the release (e.g. bugfix = comments / bugfixes, minor = larger and / or backwards compatible changes, major == breaking / backwards incompatible changes). K8S API version updates should also result in an appropriate increment of the schema version.


### Relationship Between K8s API version and JSON Schema Version

The Devfile JSON schema is generated from the Kubernetes API, and the version for the JSON schema is set in the doc.go file in the K8S API (`pkg/apis/workspace/<api-version>/doc.go).

As we’re only updating the K8S API version when needed, but incrementing the schema version more frequently, this means that any given API version may point to multiple, backwards-compatible, schema versions over its lifespan. The schema version under `schemas/<api-version>` in [devfile/api repo](https://github.com/devfile/api/) points to the latest JSON schema generated from the K8S API.

## Release Process
The following steps outline the steps done to release a new version of the Devfile schema and publish its schemas to the devfile.io website

   1) The release engineer tasked with creating the release clones the repository (and checks out the release branch if one already exists)

   2) `make-release.sh <schema-version> <k8s-api-version>` is run:

      i) A release branch (the name corresponding to the schema version) is created, if one does not already exist.

      ii) The schema-version is updated in `pkg/apis/workspace/<api-version>/doc.go`.

      iii) New JSON schemas are generated to include the new schema version
      
      iv) A new commit including the changes

      v) Finally, a PR is opened to merge these changes into the release branch

   3) Once the release PR is approved and merged, the release engineer creates a new release on GitHub based off the release branch that was just created and includes the generated `devfile.json` as a release artifact. 
       - The tag `v{major}.{minor}.{bugfix}` is used, to enable the new version of the API to be pulled in as a Go module.

   4) Once the release is published, GitHub actions run to publish the new schema to devfile.io. The “stable” schema is also updated to point to the new schema.

An example pull request, `make-release.sh` script and GitHub action can be found here:
- [Release Pull Request](https://github.com/johnmcollier/api/pull/7)
- [make-release.sh](https://github.com/johnmcollier/api/blob/master/make-release.sh)
- [release-schema.yaml](https://github.com/johnmcollier/api/blob/master/.github/workflows/release-schema.yaml)
