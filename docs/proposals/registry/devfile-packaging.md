# Devfile Packaging for OCI Registries

This document will define how devfiles and their artifacts that are hosted on Github can be packaged up and deployed onto an OCI Devfile registry running on a Kubernetes cluster.

I recommend having a read over the [OCI Registry Lifecycle design](https://docs.google.com/document/d/1rQHCp4SWslWWJv5KK3A_iXHvgbqjsuDkSKDD72ifJio/edit) first, before reading this doc. 

## Terminology

Some of the following terms will be used throughout this design proposal:

**Devfile Registry Repository:** The GitHub repository that hosts the devfile stacks for consumption within an OCI registry. For example, [devfile/registry](https://github.com/devfile/registry).

**Registry Support Repository:** The GitHub repository that hosts tooling for OCI devfile registries. This includes the [index generator tool](https://github.com/johnmcollier/registry-support/tree/master/index/generator), and the [bootstrap container base image](https://github.com/johnmcollier/registry-support/tree/master/oci-devfile-registry-metadata)

**Bootstrap Container:** The sidecar container deployed alongside the OCI registry that loads the devfile stacks into the OCI registry and hosts the index.json for consumption by tools such as Che and Odo.

**Registry Dockerfile / Image:** The Dockerfile and resulting container image that the `Devfile Registry Repository` is built up into. It's based upon a base image provided by us.

## As-is Today
Currently, devfiles are stored on Github in a "devfile registry repository” (such as [devfile/registry](https://github.com/devfile/registry), [odo-devfiles/registry](https://github.com/odo-devfiles/registry), and [eclipse/che-devfile-registry](https://github.com/eclipse/che-devfile-registry)). Generally, each registry repository will have a devfiles folder:

<img width="742" alt="Screen Shot 2020-10-22 at 1 54 00 PM" src="https://user-images.githubusercontent.com/6880023/97219676-b2637200-17a0-11eb-8465-1063aa048768.png">


Each folder under devfiles corresponds to a stack, with its own devfile and associated artifacts. There may be an index.json present for devfile consumers

From this point, the Github repository may act as self-hosting devfile registry (as in the case of [odo-devfiles/registry](https://github.com/odo-devfiles/registry), or it may be built up into an Apache server container (as in the case of v1 registries like [eclipse/che-devfile-registry](https://github.com/eclipse/che-devfile-registry)).

## Are any changes needed to accommodate OCI?
Storing the devfiles (and any associated artifacts) on a dedicated Github repository (the “devfile registry repository”), is still a good idea.

   - Provides a clear, straightforward and non-cluttered place for stack owners to contribute devfiles to
   - Separates the devfile stacks from the OCI registry code and logic

But, with the move to OCI registries, we need a way to easily package the devfile stacks from GitHub and load them into the OCI registry

   - Needs to retrieve the devfiles without an internet connection
     - Means we can’t have the bootstrap container git clone the devfiles before pushing them up
   - Needs to be easily built and distributable by registry admins.  
     - Once we support “deploying your own registry”, with custom devfiles this will become especially important
   - Shouldn’t require building multiple components of the registry to distribute

## Proposed Changes

**To address the requirements listed above, have the “devfile registry repository” built into a custom bootstrap image, based on a base-image provided by us.**

<img width="633" alt="Screen Shot 2020-10-29 at 12 25 53 PM" src="https://user-images.githubusercontent.com/6880023/97602641-eaa5c300-19e1-11eb-8c3b-1c3ee0cb6f11.png">


The base-image will be the "bootstrap container" image that's currently in the [devfile/registry-support](https://github.com/devfile/registry-support) repo and will be hosted on the `devfile` org on quay.io. It contains the logic for loading the devfiles into the OCI registry as well as hosting the index.json.

The Dockerfile will just contain lines to copy the devfiles and index.json into the container, nothing else will be required (unless the registry administrator chooses to add more).

As part of the registry build in the registry repository (such as [devfile/registry](https://github.com/devfile/registry)), the index.json will be generated and any tests/validation specific to the repository will be run. After this, the Registry Dockerfile will be built and an image will be produced and pushed up to a container registry.

When deploying the OCI registry, the registry's image will be specified as the bootstrap container. Then, when the bootstrap container starts, it pushes the devfiles up into the OCI registry using `oras` and serves the pre-generated index.json.

An example of what the devfile registry repository may look like can be found at https://github.com/johnmcollier/registry/tree/registryDockerfile


### Benefits
- Easy and straightforward method of way and distributing devfiles and their artifacts. 
    - If a user wants to deploy a registry with their own devfiles, this will be the only thing they have to change
- Everything needed to preload the OCI registry is contained in this image