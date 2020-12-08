# Devfile Registry Terminology

The following terms may be used throughout design proposals for devfile registries

## Devfile Index Image

The container image containing the devfile stacks and index.json used to bootstrap the OCI registry with devfile stacks. It also runs the webserver that functions as a proxy to the OCI registry, and hosts the index.json.

## Devfile Registry Repository 

The GitHub repository that hosts the devfile stacks for consumption within an OCI registry. For example, [devfile/registry](https://github.com/devfile/registry).

## Devfile Registry Support Repository 

The GitHub repository that hosts toolings for OCI devfile registries. This includes the [index generator tool](https://github.com/johnmcollier/registry-support/tree/master/index/generator), and the [index container base image](https://github.com/johnmcollier/registry-support/tree/master/index/server). It also contains the [build tools](https://github.com/johnmcollier/registry-support/tree/master/build-tools) and Dockerfile for creating the Devfile index image.

## Registry Build 

The process that takes the devfile registry repository, generates the index.json and builds it up into the devfile index image.

## Registry Bootstrap

The process that pushes the devfile stacks on the devfile index container to the OCI registry.
