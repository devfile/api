# Devfile Packaging for OCI Registries

This document will define how devfiles and their artifacts that are hosted on Github can be packaged up and deployed onto an OCI Devfile registry running on a Kubernetes cluster.

I recommend having a read over the [OCI Registry Lifecycle design](https://docs.google.com/document/d/1rQHCp4SWslWWJv5KK3A_iXHvgbqjsuDkSKDD72ifJio/edit) first, before reading this doc. 


## As-is Today
Currently, devfiles are stored on Github in a “devfile registry repository” (such as [devfile/registry](https://github.com/devfile/registry), [odo-devfiles/registry](https://github.com/odo-devfiles/registry), and [eclipse/che-devfile-registry](https://github.com/eclipse/che-devfile-registry)). Generally, each registry repository will have a devfiles folder:

<img width="742" alt="Screen Shot 2020-10-22 at 1 54 00 PM" src="https://user-images.githubusercontent.com/6880023/97219676-b2637200-17a0-11eb-8465-1063aa048768.png">


Each folder under devfiles corresponds to a stack, with its own devfile and associated artifacts. There may be an index.json present for devfile consumers

From this point, the Github repository may act as self-hosting devfile registry (as in the case of [odo-devfiles/registry](https://github.com/odo-devfiles/registry), or it may be built up into an Apache server container (as in the case of v1 registries like [eclipse/che-devfile-registry](https://github.com/eclipse/che-devfile-registry)).

## Are any changes needed to accommodate OCI?
Storing the devfiles (and any associated artifacts) on a dedicated Github repository (the “devfile registry repository”), is still a good idea.

   - Provides a clear, straightforward and non-cluttered place for stack owners to contribute devfiles to
   - Separates the devfile stacks from the OCI registry code and logic

But, with the move to OCI registries, we need a way to easily package the devfile stacks from GitHub and load them into the OCI registry

   - Needs to be offline
     - Means we can’t have the metadata container git clone the devfiles before pushing them up
   - Needs to be easily built and distributable by registry admins.  
     - Once we support “deploying your own registry”, with custom devfiles this will become especially important
   - Shouldn’t require building and deploying custom components of the registry to distribute (i.e. it shouldn’t be baked into the registry or metadata container images)
     - Currently, the metadata container has the devfiles built into it, meaning users need to build their own metadata container image to customize devfiles

## Proposed Changes

**To address the requirements listed above, rather than bundling the devfiles as part of registry's bootstrap container, build the “devfile registry repository” into its own docker image**

<img width="636" alt="Screen Shot 2020-10-26 at 3 39 33 PM" src="https://user-images.githubusercontent.com/6880023/97220180-767cdc80-17a1-11eb-85fd-cf1a5a623aeb.png">


This container image will contain just the devfiles (and their artifacts), as well as any other required files.

When deploying the OCI registry, the image will be specified as an init container, with an entrypoint command that copies the devfiles (and artifacts) to the bootstrap container’s volume. Then, as-is today, when the bootstrap container starts up, it generates the index.json and pushes the devfiles up into the OCI registry using `oras`.

An example of what the devfile registry repository may look like can be found at https://github.com/johnmcollier/registry/tree/registryDockerfile


### Benefits
- Easy and straightforward method of way and distributing devfiles and their artifacts. 
    - If a user wants to deploy a registry with their own devfiles, this will be the only thing they have to change
- Offline by default. Everything needed to preload the OCI registry is contained in this image
- No need to build custom components of the registry (such as the bootstrap container)