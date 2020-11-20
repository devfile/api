# Console Import from Devfile

This document outlines how the console is going to use a devfile 2.0.0 spec for it's import feature targeted for the upcoming release.

## As-Is Today

Currently the Devfile import feature is mocked with a POC [openshift/console PR](https://github.com/openshift/console/pull/6321). The POC PR, requires the build guidance devfile spec to be implemented. However, the build guidance spec is still an open discussion in the [devfile/api PR](https://github.com/devfile/api/pull/127). 

With an initial target date for the Dec 04, 2020; the devfile import developer preview should look similar to this [demo video](https://drive.google.com/file/d/1uLzDibVZlkMqbjtKkho04e8k2-Ns5A2W/view).

The information required to build a component from build guidance are:
- Git Repo Url
- Git Repo Ref
- Build Context Path
- Container Port
- Image Name

The console devfile import page has the Git Repo Url, Git Repo Ref & the Build Context Path

<img src="https://user-images.githubusercontent.com/31771087/99319303-4ae89180-2837-11eb-8933-eaaf41160bcd.png">

The Dockerfile is assumed to be present in the context directory. The Container Port is assumed to be 8080 and the Image Name is assumed to be the same as the application name.

### Note:
The context path is present in both the Console import form and the open devfile spec(Refer to the OpenShift Console Devfile Import screenshot above and the open Build Guidance PR screenshot below). Do we require two context dir information? Is the context dir in the Console form for devfile.yaml and the context dir in the build guidance spec for Dockerfile relative to devfile.yaml? Or do we always assume both the devfile.yaml and Dockerfile are always at the same dir level.

 For phase 1, only the Dockerfile Build Guidance is to be implemented. SourceToImage(S2I) will be implemented in the future sprints.

<img src="https://user-images.githubusercontent.com/31771087/99319306-4c19be80-2837-11eb-9639-a5c130deb4ba.png">

## Proposed Changes

With the target deadline approaching soon, we may need to provide an alternative path for the OpenShift Console to consume the devfile. With the Build Guidance spec [devfile/api PR](https://github.com/devfile/api/pull/127) still an open discussion, this document proposes a solution within the boundaries of the devfile 2.0.0 spec to achieve this.

Here is a proposed 2.0.0 spec devfile that can support build guidance for OpenShift Console:

```yaml
schemaVersion: 2.0.0
metadata:
  name: nodejs
  version: 1.0.0
  attributes:
    alpha.build-dockerfile: https://raw.githubusercontent.com/odo-devfiles/registry/master/devfiles/nodejs/build/Dockerfile # can also be a path relative to the context
components:
- name: myapplication
  attributes:
    tool: console-import # key:value pair used to filter container type component that only the Console Devfile Import is interested in
  container:
    image: buildContextImageOutput:latest # this is the image which will be used by the Console's buildConfig output
    endpoints:
    - name: http-3000 # define endpoints in devfile.yaml, rather than hardcoding a default
      targetPort: 3000
    env:
      - name: FOO # set container env through the devfile.yaml for the container
        value: "bar"
```

Console's devfile import POC is using a `BuildConfig` to build the Dockerfile from the dockerfile path location, currently found from the context dir. The POC PR `BuildConfig` outputs it to an `ImageStreamTag`. However, the POC PR assumes the `ImageStream` and `ImageStreamTag` are all the same name as the application name. This is tightly coupled and dependant on the information entered in the Console Devfile Import form and thus does not allow us to mention a custom image name in the devfile.

If we want to decouple `ImageStream` and `ImageStreamTag` from the application name, so that a custom image name can be specified in the devfile's component container; then we would need to update the following files: 
1. `pkg/server/devfile-handler.go`
   - image stream name 
   - build config output image stream tag 
   - deployment's container image
2. `frontend/packages/dev-console/src/components/import/import-submit-utils.ts`
   - the annotation mapping for the `ImageStreamTag` in `getTriggerAnnotation()` 
  
This allows us to use the image name from devfile.yaml rather than always using the application name.

The above proposed devfile container component can also mentions an endpoint instead of hardcoding a default 8080. This can be updated in `createOrUpdateResources()` in `frontend/packages/dev-console/src/components/import/import-submit-utils.ts`

These devfile information would be parsed and returned by the library and thus ensuring a consistent UX.

### Note:
The above devfile proposal and POC PR assumes the `BuildConfig` outputs the build to an `ImageStreamTag` which is used by the OpenShift internal registry. To push the image to a non-OpenShift Image Registry, the `BuiidConfig` output can pushed to a private or Dockerhub registry using `DockerImage`. OpenShift [doc](https://docs.openshift.com/container-platform/4.6/builds/managing-build-output.html) outlines how this can be achieved via a `BuildConfig`. Pushing to a private registry requires secret configuration from the Docker config, and this OpenShift [doc](https://docs.openshift.com/container-platform/3.11/dev_guide/builds/build_inputs.html#using-docker-credentials-for-private-registries) explains how to achieve it.

For phase 1, only the `ImageStreamTag` is to be implemented in the `BuildConfig`, `DockerImage` is to be implemented in the future sprints.
