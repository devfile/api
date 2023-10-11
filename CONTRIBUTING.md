# Contributing

Contributions are welcome!

## Code of Conduct

Before contributing to this repository for the first time, please review our project's [Code of Conduct](https://github.com/devfile/api/blob/main/CODE_OF_CONDUCT.md).

## Getting Started

### Issues

- Open or search for [issues](https://github.com/devfile/api/issues) with the label `area/api`.

- If a related issue doesn't exist, you can open a new issue using a relevant [issue form](https://github.com/devfile/api/issues/new/choose). You can tag issues with `/area api`.

### Building

To build the CRD and the various schemas, you don't need to install any pre-requisite apart from `docker`.
In the root directory, just run the following command:

```console
bash ./docker-run.sh ./build.sh
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

### Pull Requests

All commits must be signed off with the footer:

```git
Signed-off-by: Firstname Lastname <email@email.com>
```

Once you set your user.name and user.email in your git config, you can sign your commit automatically with git commit -s. When you think the code is ready for review, create a pull request and link the issue associated with it.

Owners of the repository will watch out for and review new PRs.

If comments have been given in a review, they have to be addressed before merging.

After addressing review comments, don’t forget to add a comment in the PR afterward, so everyone gets notified by Github and knows to re-review.
