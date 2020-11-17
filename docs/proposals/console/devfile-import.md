# Console Import from Devfile

This document outlines how the console is going to use a devfile 2.0.0 spec for it's import feature targeted for Dec 4, 2020.

## As-Is Today

Currently the Devfile import feature is mocked with a POC [openshift/console PR](https://github.com/openshift/console/pull/6321). The POC PR, requires the build guidance devfile spec to be implemented. However, the build guidance spec is still an open discussion in the [devfile/api PR](https://github.com/devfile/api/pull/127). 

With an initial target date for the Dec 04, 2020; the devfile import developer preview should look similar to this [demo video](https://drive.google.com/file/d/1uLzDibVZlkMqbjtKkho04e8k2-Ns5A2W/view).

The information for required to build a component from build guidances are:
- Git Repo Url
- Git Repo Ref
- Build Context Path
- Container Port

The console devfile import page has the Git Repo Url, Git Repo Ref & the Build Context Path

<img src="https://user-images.githubusercontent.com/31771087/99319303-4ae89180-2837-11eb-8933-eaaf41160bcd.png">

The open devfile spec for Build Guidance allows for the Dockerfile Location to be specified

<img src="https://user-images.githubusercontent.com/31771087/99319306-4c19be80-2837-11eb-9639-a5c130deb4ba.png">

The Container Port in the POC was parsed from the Dockerfile's `EXPOSE` command but it should ideally be defined in the devfile.

### Note:
The context path is present in both the Console import form and the devfile spec. It should ideally be present only at the Console form level, since it will help find the `devfile.yaml` and as well as the `Dockerfile`.

## Proposed Changes

With the target deadline closing soon, we may need to provide an alternative path for the OpenShift Console to consume devfile. With the Build Guidance spec [devfile/api PR](https://github.com/devfile/api/pull/127) still an open discussion, this document proposes a solution within the boundaries of the devfile 2.0.0 spec to achieve this.

Here is a proposed 2.0.0 spec devfile that can support build guidance for OpenShift Console:

```yaml
schemaVersion: 2.0.0
metadata:
  name: nodejs
  version: 1.0.0
  attributes:
    build.guidance-dockerfile: https://raw.githubusercontent.com/odo-devfiles/registry/master/devfiles/nodejs/build/Dockerfile # can also be a path relative to the context provided in the Console devfile import form
components:
- name: console
  container:
    image: buildImage:latest # this is the image which will be used by the Console's buildConfig output
    endpoints:
    - name: http-3000 # define endpoints in devfile.yaml, rather than getting it from the Dockerfile
      targetPort: 3000
    env:
      - name: FOO # set container env through the devfile.yaml for the container
        value: "bar"
```

Console's devfile import POC is using a `BuildConfig` to build the Dockerfile from the dockerfile path location. The POC PR outputs it to an `ImageStream`. However, the above devfile proposal would require an exact image name and tag to build a container from the image. The POC PR could change the `BuiidConfig` output to a docker image. OpenShift [doc](https://docs.openshift.com/container-platform/4.6/builds/managing-build-output.html) outlines how this can be achieved via a `BuildConfig` and pushes the build to a private or dockerhub registry.

The devfile container component can also mention endpoints instead of parsing it from the Dockerfile's `EXPOSE` cmd like the Console's POC currently does. Furthermore, more properties can be mentioned in the devfile container component as per the spec, to customize the pod container.

These devfile information would be parsed and returned by the library and thus ensuring a consistent UX across all the consumners of the devfile.