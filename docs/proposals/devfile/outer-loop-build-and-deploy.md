# Project Proposal - Devfile 2.0 outer-loop devfile build and deploy functions

## Background
Existing devfile 2.0 mainly focuses on providing inner-loop support. This proposal adds support on outer-loop scenarios to expand the scope of usage of devfile support.

### Outer-loop scenario
As an application developer, I would like to build a microservice and deploy it to Kubernetes to do a build that is typically done as part of a pipeline

- Do a full build to build the container image
    - Input:\
    The build may use different technologies for building, e.g. dockerfile, buildpacks
    - Output:\
    A container image that contains the microservice and ready to be deployed to Kubernetes
- Deploy the built image to Kubernetes:
    - Deploy the built image to a cluster:\
    Deployment may use different methods of deployment, e.g. Kubernetes deployment manifest, Operators, Helm, Knative

The focus of this proposal is to provide the ability to build and deploy applications or services in a way similar to what will be done for production. In contrast to the inner-loop scenario, the build will not only focus on the building of the application but will also cover building the runtime container that runs the application.

We will cover two main stages:
1. Build the image
2. Deploy the image to a cluster

## Build the image

This includes the creation of the base image of the runtime container, building the application, and packaging the application as part of the container image. Different mechanisms can be used for building the image, for example:
1. Dockerfile
2. Dockerfile stages
3. Buildpacks
4. Source-to-Image (S2I)
5. Tekton

The initial design will focus on the first two items for discussion. Other ways of building the image will be covered later.

The definition of the image build introduces a new component type `image` for the definition of the image build. The `image` component definition can use different build strategies, e.g. `dockerfile`, `s2i`, `buildpacks`, etc.

__Alternative approach:__ we can consider reusing the `exec` type except that these build types are special and we may not want to incorporate the full image build command into an `exec` command. Therefore, I chose to use a different type of command here.

### Example using a file Dockerfile as part of the app:
```yaml
    variables:
      myimage: myimagename
    components:
      - name: mydockerfileimage
        image:
          imageName: {{myimage}}
          dockerfile:
            buildContext: ${PROJECTS_ROOT}/build
            location: Dockerfile
            args: [ "arg1", "arg2", "arg3" ]
            rootRequired: false
```

`imageName`: name of the generated image (or we may use an id for referencing the image). This imageName can also be fully qualified to allow pushing of the image to a specific image registry.

`buildContext`: Path of source directory to establish build context.  Default to ${PROJECT_ROOT} (optional)

`location`: Dockerfile location which can be an URL or a path relative to buildContext

`args`: Argument list for the docker build (optional)

`rootRequired`: Specifies whether a privileged builder pod is required.  Default is false. (optional)

__Note:__
1. A common pattern will be using a global variable to define the image name (`myimage` in the example above) so that it can be easily referred to in the deploy step later.

### Example using a Dockerfile with registry and secret for the image push:
```yaml
    variables:
      myimage: myimagename
    components:
      - name: mydockerfileimage
        image:
          imageName: {{myimage}}
          dockerfile:
            buildContext: ${PROJECTS_ROOT}/build
            location: https://github.com/redhat-developer/devfile-sample/blob/master/src/Dockerfile
            args: [ "arg1", "arg2", "arg3" ]
            rootRequired: false
            envFrom:
            - secretRef:
                name: my-secret
```

`secretRef`: The reference to the secret for pushing the image to the registry

For the secrets, it should support the same mechanisms as specified in https://github.com/devfile/api/issues/299.

### Example using `apply` command on an image component:
```yaml
    commands:
      - id: deploybuild
        apply:
          component: mydockerfileimage
```

### Example using a Dockerfile stored within the devfile registry as a resource with `apply` command:
```yaml
    variables:
      myimage: myimagename
    components:
      - name: mydockerfileimage
        image:
          imageName: {{myimage}}
          dockerfile:
            id: mycompany/my-node-stack-dockerfile/v2.2.2
            args: [ ]
    commands:
      - id: deploybuild
        apply:
          component: mydockerfileimage
```

__Note:__
1. The `apply` command with the `image` component is optional. If the `apply` command that references an image component does not exist, then the image build will be done at the startup automatically (similar to the behaviour of the existing `kubernetes` components).  If there is an `apply` command that references the `image` component, then the user will need to create a composite command to chain the image build apply command with the deploy command (see the deployment command section) to complete the build and deploy.
 
### Example using a Dockerfile stored in a git repo with push registry without apply command:
```yaml
    variables:
      myimage: myimagename
    components:
      - name: mydockerfileimage
        image:
          imageName: {{myimage}}
          dockerfile:
            git:
              remotes:
                origin: "https://github.com/odo-devfiles/nodejs-ex.git"
            location: Dockerfile
            args: [ ]
```

The git definition will be the same as the one in `starterProjects` definition that supports `checkoutFrom`

`location`: the location of the dockerfile within the repo. If `checkoutFrom` is being used, the location will be relative to the root of the resources after cloning the resources using the `checkoutFrom` setting.

Notes:
1. The build tool/mechanism, e.g. buildah/kaniko, used for building the dockerfile is up to the tools so it is not part of the spec
1. If endpoints definition is needed, it will be defined via the corresponding deployment manifest in the deployment step, e.g. inside the kubernetes deployment manifest.
1. Do we need image push registry info to specify where to push the image to? If needed, then the secret for the push can be specified as `secretKeyRef` in https://github.com/devfile/api/issues/299

### Example using SourceToImage (S2I):
```yaml
    variables:
      myimage: myimagename
    components:
      - name: mys2iimage
        image:
          imageName: {{myimage}}
          s2i:
            builderImageNamespace: mynamespace
            builderImageStreamTag: mytag
            scriptLocation:
              remotes:
                origin: "https://github.com/odo-devfiles/nodejs-ex.git"
    commands:
      - id: deploybuild
      apply:
        component: mys2iimage
```

`builderImageNamespace`: Namespace where builder image is present

`builderImageStreamTag`: Builder image name with tag

`scriptionLocation`: Script URL to override default scripts provided by builder image

`incrementalBuild`: Flag that indicates whether to perform incremental builds. Default is true (optional)

## Deployment of the image
Deploy the image to a cluster. Different technologies can be used for deploying the image, for example:
1. Kubernetes deployment manifest
1. Operators
1. Helm

The initial design will focus on the first two items for discussion. Other ways of building the image will be covered later.

To specify the deploy step, we reuse the existing `apply` command type under the `commands` to apply `kubernetes` components. The deploy step uses a new `group` kind called `deploy`.

### Deploying technologies
The current design is to try to reuse the existing `kubernetes` component as much as possible for all three deployment technologies, namely Kubernetes deployment manifest, Operators, and Helm given that all of them can be represented as `kubernetes` objects. We’ll only introduce `operator` or `helm` component types only if we find that there are extra requirements needed for Operators and Helm.

#### Examples of using `deploy` group command:
##### Kubernetes deployment manifest:
```yaml
    components:
      - name: myk8sdeploy
        kubernetes:
          uri: deploy/deployment-manifest.yaml

    commands:
      - id: deployk8s
        apply:
          component: myk8sdeploy
          group:
            kind: deploy
            isDefault: true
        attributes: 
          - name: CONTAINER_IMAGE
            value: {{myimage}}
```

`uri`: Kubernetes manifest location (can use `kubectl` to deploy) which can be an URL or a path relative to the devfile. [Example of Kubernetes deployment](#markdown-header-example-of-kubernetes-deployment-manifest) manifest and [example of Operator deployment manifest](#markdown-header-example-of-operator-deployment-manifest).

Variables that need to be replaced during the deployment can be specified using a new `attributes` definition under the `apply` command. One usage example is to pass the image name along to the deployment manifest. A common practice is to use a global variable, e.g. `myimage` in the example above, to refer to the image built in the image built stage.
 
##### Kubernetes deployment manifest (inlined):
```yaml
    components:
      - name: myk8deploy
        kubernetes:
          inlined: |
            apiVersion: batch/v1
            kind: Job
            metadata:
              name: pi
            spec:
              template:
                spec:
                  containers:
                  - name: job
                    image: {{myimage}}
                    command: ["some",  "command", "with", "parameters"]
                  restartPolicy: Never
              backoffLimit: 4
    commands:
      - id: deployk8s
        apply:
          component: myk8sdeploy
          group:
            kind: deploy
            isDefault: true
```

#### Example of Operator deployment manifest:
```yaml
    apiVersion: app.stacks/v1beta1
    kind: RuntimeComponent
    metadata:
      name: {{.COMPONENT_NAME}}
    spec:
      applicationImage: {{.CONTAINER_IMAGE}}
      service:
        type: ClusterIP
        port: {{.PORT}}
      expose: true
      storage:
        size: 2Gi
        mountPath: "/logs"
```

##### Helm:
```yaml
    components:
      - name: myhelmdeploy
        helm:
          chart: http://helm-chart-url
          values: 
            image: quay.io/sample/hello-world # A chart may expose a well-known values.yaml parameter called "image".
            replicas: 3
    commands:
      - id: deployHelm
        apply:
          component: myhelmdeploy
          group:
            kind: deploy
            isDefault: true
        variables: 
          - name: CONTAINER_IMAGE
            value: {{myimage}}
```

### Deployment manifest examples
#### Example of Kubernetes deployment manifest:
```yaml
    ---
    kind: Deployment
    apiVersion: apps/v1
    metadata:
      name: {{.COMPONENT_NAME}}
    spec:
      replicas: 1
      selector:
        matchLabels:
          app: {{.COMPONENT_NAME}}
    template:
      metadata:
        creationTimestamp: null
        labels:
          app: {{.COMPONENT_NAME}}
      spec:
        containers:
          - name: {{.COMPONENT_NAME}}
            image: {{.CONTAINER_IMAGE}}
            ports:
              - name: http
                containerPort: {{.PORT}}
                protocol: TCP
    ---
    kind: Service
    apiVersion: v1
    metadata:
      name: {{.COMPONENT_NAME}}
    spec:
      ports:
        - protocol: TCP
          port: {{.PORT}}
          targetPort: {{.PORT}}
      selector:
        app: {{.COMPONENT_NAME}}
      type: ClusterIP
      sessionAffinity: None
    ---
    kind: Route
    apiVersion: route.openshift.io/v1
    metadata:
      name: {{.COMPONENT_NAME}}
      annotations:
        openshift.io/host.generated: 'true'
    spec:
      to:
        kind: Service
        name: {{.COMPONENT_NAME}}
        weight: 100
      port:
        targetPort: {{.PORT}}
      wildcardPolicy: None
```

This example is using variables for tools to replace some of the info during deployment.

