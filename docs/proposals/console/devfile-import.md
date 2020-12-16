# Console Import from Devfile

This document outlines how the console is going to use a devfile 2.0.0 spec for it's import feature targeted for the upcoming release.

## Dev Preview POC

The Devfile import feature was mocked with a POC [openshift/console PR](https://github.com/openshift/console/pull/6321). The POC PR, required the build guidance devfile spec to be implemented. However, the build guidance spec is still an open discussion in the [devfile/api PR](https://github.com/devfile/api/pull/127). 

With an initial target date for the Dec 04, 2020; the devfile import developer preview should look similar to this [demo video](https://drive.google.com/file/d/1uLzDibVZlkMqbjtKkho04e8k2-Ns5A2W/view).

## Required Data

The information required to build a component from devfile are:
- Git Repo Url
- Git Repo Ref
- Devfile Context Dir
- Docker Build Context Dir
- Dockerfile Location
- Container Port
- Image Name

The console devfile import page has the Git Repo Url, Git Repo Ref & the Devfile Context Dir.

<img src="https://user-images.githubusercontent.com/31771087/99319303-4ae89180-2837-11eb-8933-eaaf41160bcd.png">

There are two context directories info required. The `Devfile Context Dir` present in the Console form is used to find the devfile.yaml in a repository. The `Docker Build Context Dir` present in the devfile.yaml is used to establish a context dir for Docker builds relative to the devfile.yaml.

 For phase 1, only the Dockerfile Build Guidance is to be implemented. SourceToImage(S2I) will be implemented in the future sprints.

<img src="https://user-images.githubusercontent.com/31771087/99319306-4c19be80-2837-11eb-9639-a5c130deb4ba.png">

With the target deadline approaching soon, we may need to provide an alternative path for the OpenShift Console to consume the devfile. With the Build Guidance spec [devfile/api PR](https://github.com/devfile/api/pull/127) still an open discussion, this document proposes a solution within the boundaries of the devfile 2.0.0 spec to achieve this.

## Proposed Changes

Here is a proposed 2.0.0 spec devfile that can support devfile build guidance for OpenShift Console:

```yaml
schemaVersion: 2.0.0
metadata:
  name: nodejs
  version: 1.0.0
  attributes:
    alpha.build-context: mydir # key:value pair that establishes the context dir for Docker builds relative to devfile.yaml
    alpha.build-dockerfile: Dockerfile # key:value pair that specifies the location of the Dockerfile relative to alpha.build-context
components:
  - name: runtime
    attributes:
      tool: console-import # key:value pair used to filter container type component that only the Console Devfile Import is interested in
    container:
      image: imageplaceholder # image which will be used by the buildConfig output but not supported for Dev Preview, defaults to Console's application name
      memoryLimit: 1024Mi
      endpoints:
        - name: http-3000 
          targetPort: 3000 # define endpoints in devfile.yaml, that will be used for the devfile service
  - name: runtime2 # other components ignored by the Console Devfile Import
    container:
      image: registry.access.redhat.com/ubi8/nodejs-12:1-45
      memoryLimit: 1024Mi
      mountSources: true
      sourceMapping: /project
      endpoints:
        - name: http-3000
          targetPort: 3000
commands:
  - id: install
    exec:
      component: runtime2
      commandLine: npm install
      workingDir: /project
      group:
        kind: build
        isDefault: true
  - id: run
    exec:
      component: runtime2
      commandLine: npm start
      workingDir: /project
      group:
        kind: run
        isDefault: true
```

Without the build guidance devfile spec, the proposed change here is to mention the docker build context and the dockerfile location in the metadata attribute. The attributes is a free form key-value pair and in the above example, `alpha.build-context` & `alpha.build-dockerfile` are used for the values; which would have otherwise been specified by the Dockerfile Build Guidance spec.

To allow the Console to filter only devfile component containers, we use a `tool: console-import` key value pair attribute in the container component. Any other components in the devfile are ignored by the Console.

Asides the devfile change, the Dev Preview would need to be updated to use the above devfile's image in the future. 

Console's devfile import Dev Preview is using a `BuildConfig` to build the Dockerfile from the dockerfile path location and build context directory. The `BuildConfig` outputs it to an `ImageStreamTag`. However, for the Dev Preview it assumes the `ImageStream` and `ImageStreamTag` are all the same name as the application name. This is tightly coupled and dependant on the information entered in the Console Devfile Import form and thus does not allow us to mention a custom image name in the devfile.

If we want to decouple `ImageStream` and `ImageStreamTag` from the application name, so that a custom image name can be specified in the devfile's component container; then we would need to update the following information: 
1. `pkg/server/devfile-handler.go`
   - image stream name 
   - build config output image stream tag 
   - deployment's container image
2. `frontend/packages/dev-console/src/components/import/import-submit-utils.ts`
   - the annotation mapping for the `ImageStreamTag` in `getTriggerAnnotation()` 
  
This allows us to use the image name from devfile.yaml rather than always using the application name.

These devfile information would be parsed and returned by the library and thus ensuring a consistent UX.

### Note:
The above devfile proposal and POC PR assumes the `BuildConfig` outputs the build to an `ImageStreamTag` which is used by the OpenShift internal registry. To push the image to a non-OpenShift Image Registry, the `BuiidConfig` output can pushed to a private or Dockerhub registry using `DockerImage`. OpenShift [doc](https://docs.openshift.com/container-platform/4.6/builds/managing-build-output.html) outlines how this can be achieved via a `BuildConfig`. Pushing to a private registry requires secret configuration from the Docker config, and this OpenShift [doc](https://docs.openshift.com/container-platform/3.11/dev_guide/builds/build_inputs.html#using-docker-credentials-for-private-registries) explains how to achieve it.

For phase 1, only the `ImageStreamTag` is to be implemented in the `BuildConfig`, `DockerImage` is to be implemented in the future sprints.
