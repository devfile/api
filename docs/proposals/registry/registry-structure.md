# Devfile Registry Structure
This document outlines the structure of a Devfile Registry Repository that’s used as the basis for an OCI Devfile Registry, hosted on Kubernetes. It also outlines how individual files in each stack will get pushed up to the OCI registry.

This design proposal is a follow up to [Devfile Registry Packaging](https://github.com/devfile/api/blob/master/docs/proposals/registry/devfile-packaging.md) and I recommend reading that first.

## Terminology

Some of the following terms will be used throughout this design proposal:

**Devfile Registry Repository:** The GitHub repository that hosts the devfile stacks for consumption within an OCI registry. For example, [devfile/registry](https://github.com/devfile/registry).

**Devfile Index Image:** The container image containing the devfile stacks and index.json used to bootstrap the OCI registry with devfile stacks

**Registry Build:** The process that takes the devfile registry repository, generates the index.json and builds it up into the devfile index container image.

**Registry Bootstrap:** The process that pushes the devfile stacks on the devfile index container to the OCI registry.

## As-is Today
Currently, the top-level structure of a devfile registry’s repository is unwritten, but it usually consists of a **devfiles** or **stacks** folder, and any associated files specific to that registry

[devfile/registry](https://github.com/devfile/registry):
<img width="914" alt="Screen Shot 2020-11-30 at 2 12 13 PM" src="https://user-images.githubusercontent.com/6880023/100653219-12a48100-3316-11eb-949c-38073a19acbc.png">

Inside each **stacks** or **devfiles** folder, each folder corresponds to a devfile stack, containing usually just the devfile.yaml and a meta.yaml:
<img width="1232" alt="Screen Shot 2020-11-30 at 2 10 50 PM" src="https://user-images.githubusercontent.com/6880023/100653095-e2f57900-3315-11eb-8c3f-86e56896ef15.png">

<img width="1229" alt="Screen Shot 2020-11-30 at 2 11 20 PM" src="https://user-images.githubusercontent.com/6880023/100653133-f0aafe80-3315-11eb-8393-dbabfd8ce0d6.png">

As part of the registry build, the index.json is generated based off of the stacks in the repository, and a devfile index image is generated containing the index.json and stacks.

When deploying the OCI registry, the registry bootstrap process parses the index.json to find the devfile.yaml for each stack, and pushes the devfile up to the registry as a single layer. No other stack artifacts are pushed up as part of the layer, or as separate layer.

## Problem
Our approach of pushing only the devfile.yaml works fine currently because the stacks we’re pushing only have a devfile.yaml in them. However, devfile stacks may contain more than _just_ the devfile, and may also contains resources that are re-used across multiple stacks (such as certain vsx plugins).

We need a defined way of knowing what files in a stack to push up in a, what layers each file in the stack should belong to, and where the stack's file should located. 

## Proposal

To solve the issue listed above, we should:
1) Formally the expected structure of the Devfile Regisry Repository.
2) Define the layers that compose a devfile stack on an OCI registry, and what each layer contains.

### Repository Structure
The structure of the “Devfile Registry Repository” should impose the following requirements:
1) A top-level folder called `stacks`, which contains folders for each devfile stack.
2) Each devfile stack folder must contain a `devfile.yaml`. Other files such as vsx plugins, stack logos, etc. can be included as needed.
3) A Dockerfile to package the stacks and index.json (see https://github.com/devfile/api/blob/master/docs/proposals/registry/devfile-packaging.md)
4) Build scripts or tools for the registry (see https://github.com/devfile/api/blob/master/docs/proposals/registry/devfile-packaging.md)

### Layer Media types
Currently, when we push devfile stacks to an OCI registry, it's pushed as a single layer, using the `application/vnd.devfileio.devfile.layer.v2+yaml` media type. We should instead be pushing the stack as a multi-layer artifact, adding the additional layers:

**VSX Plugins**

*.vsx - `application/vnd.devfileio.vsx.layer.v1.tar`

**Stack logos**

logo.svg - `image/svg+xml` or

logo.png - `image/png`

**Everything else**

archive.tar - `application/x-tar`

As part of the registry build process that packages the stacks into a container image, any files not belonging to the devfile, vsx or logo media types will be lumped together in a tar archive (using the media type `application/x-tar`).

When the registry bootstrap process pushes the stack up to the OCI registry, each file belonging to one of the above media types (`devfile.yaml`, `*.vsx`, `logo.svg`/`logo.png`, `archive.tar`) will be treated as separate layers in the artifact.


### Index.json Modification
The index.json already includes a link to a stack’s devfile, which the registry bootstrap process parses when pushing devfiles up to the registry. 

To make it easier to bootstrap the registry, and to avoid having to programmatically find which files to push individually; as part of the registry build process, include a new `layers` array in the index.json that tells the registry bootstrap process which layers the stack is composed of:
```
{
    "name": "java-maven",
    "displayName": "Maven Java",
    "description": "Upstream Maven and OpenJDK 11",
    "tags": [
      "Java",
      "Maven"
    ],
    "projectType": "maven",
    "language": "java",
    "layers": [
      "devfile.yaml",
      "java-lsp.vsx",
      "xml-lsp.vsx",
      "archive.tar",
    ],
    "links": {
      "self": “catalog/java-maven:latest”,
    }
  },
```



