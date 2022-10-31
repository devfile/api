# Devfile Release

## Release Process

The current process can be found here:

[Devfile Versioning and Release Process](docs/proposals/versioning-and-release.md).

## Levels of Support

Devfile releases are considered **_stable_**  even though the JSON schema is derived from an **_alpha_** version
of the K8s API which is under development and managed by the `devworkspace` team.

This means the scope of support will only apply to the subset of APIs that are part of the Devfile spec.

The following table summarizes the versioning relationship of a Devfile release and it's corresponding components:

|    Components/Versions of a Devfile Release    | Description| Release Stage  |
| :---        | :---            |:--- |
| [K8s API](docs/proposals/versioning-and-release.md#kubernetes-api)     |   The `devworkspace` API that Devfiles is based upon.  It is independently versioned from the Devfile release.  <br> `e.g. K8s API, v1alpha2 is part of the Devfile 2.1.0 release`   | alpha
| [JSON Schema](docs/proposals/versioning-and-release.md#devfile-json-schema)   |     The  Devfile JSON structure that is generated from the K8s API.   The version is in sync with the Devfiles release.  <br> `e.g. JSON schema v2.1.0 is part of the Devfile 2.1.0 release`  |stable
| API release version | This is the software release version of the K8s API.   Similar to the JSON schema, this version is also in sync with the Devfile release.  <br> `e.g. K8s API v1alpha2 version 2.1.0 is part of the Devfile 2.1.0 release`| stable |

## Decoupling the API Release Version (Future Consideration)

Currently, our JSON Schema and API releases are kept in sync with the Devfile version, but we could run in a situation
where we would need to branch off a separate API release and maintain out of sync versions that can be delivered within
a Devfile release or outside of one.

This scenario can happen when we introduce breaking API changes that would impact our consumers while continuing to
maintain backward compatibility with the schema.

We discussed the following approach in a team meeting, and it's documented here in case we ever need to take action:

1. The **main** branch will continue to be the branch for active development.  If there is a breaking API change, then
we would create a new API release branch and bump up the major version.
2. There's the potential for dual maintenance if we are working towards a major release and find a breaking API change
in the current release.  We would then create a new API branch to work on the fix and then deliver the changes back to **main**.
3. Since dev is done on the **main** branch, consumers will run into a known versioning bug:
[Incorrect api version is generated from "go get github.com/devfile/api/v2@main"](https://github.com/devfile/api/issues/599).
which will cause interim, unreleased versions to incorrectly appear as `v2xxx` rather than the latest `vNext` API version.
We will need to document this limitation and make it clear that even though the versioning is incorrect,
they are getting the latest code.
