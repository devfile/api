# Contributing

Contributions are welcome!

## Code of Conduct

Before contributing to this repository for the first time, please review our project's [Code of Conduct](https://github.com/devfile/api/blob/main/CODE_OF_CONDUCT.md).

## Certificate of Origin

By contributing to this project you agree to the Developer Certificate of
Origin (DCO). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the
contribution. See the [DCO](DCO) file for details.

In order to show your agreement with the DCO you should include at the end of the commit message,
the following line:

```console
Signed-off-by: Firstname Lastname <email@email.com>
```

Once you set your user.name and user.email in your git config, you can sign your commit automatically with `git commit -s`.

## How to Contribute:

### Issues

- Open or search for [issues](https://github.com/devfile/api/issues) with the label `area/api`.

- If a related issue doesn't exist, you can open a new issue using a relevant [issue form](https://github.com/devfile/api/issues/new/choose). You can tag issues with `/area api`.

### Submitting Pull Request

When you think the code is ready for review, create a pull request and link the issue associated with it.

[Owners](.github/CODEOWNERS) of the repository will watch out for new PRs and provide reviews to them.

If comments have been given in a review, they have to be addressed before merging.

After addressing review comments, don't forget to add a comment in the PR with the reviewer mentioned afterward, so they get notified by Github to provide a re-review.

### Prerequisites

The following are required to build the CRDs and TypeScript models containing your changes:

- Docker or Podman
- Git
- Command-line JSON processor (jq)
- Node 18 or later

Testing requires Go 1.22+ to be installed.

### Building

To build the CRD and the various schemas, you don't need to install any pre-requisite apart from `docker` or `podman`.
In the root directory, if you are using `podman` first run `export USE_PODMAN=true`. Then for either `docker` or `podman` run the following command:

```console
bash ./docker-run.sh ./build.sh
```

### Typescript model

Typescript model is generated based on JSON Schema with help of <https://github.com/kubernetes-client/gen>.
To generate them locally run:

```console
bash ./build/typescript-model/generate.sh
```

### Testing

Find more information about tests in the [testing document](test/README.md).

```console
# schemaTest approach
cd test/v200/schemaTest
go test -v
```

```console
# apiTest approach
cd test/v200/apiTest
go test -v
```

# Contact us

If you have any questions, please visit us the [`#devfile` channel](https://kubernetes.slack.com/archives/C02SX9E5B55) under the [Kubernetes Slack](https://slack.k8s.io) workspace.
