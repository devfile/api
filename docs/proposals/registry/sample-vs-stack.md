# Proposal for sample vs stack in a devfile registry

## Overview

The existing devfile registry focuses on the stack support to provide generic language/framework/runtime, e.g. Node, Java maven, Quarkus, etc.,  support to build and run user applications. These devfiles are called ***stack devfiles***.  There is another kind of devfile in a devfile registry that is tailored for a building and running a specific application. These devfiles are referred to as ***sample devfiles*** (example: https://github.com/redhat-developer/devfile-sample).

This proposal covers the design to support sample devfiles as a first-class citizen and differentiates between stack vs samples so that the tools can consume them properly. Also, it allows the source of a given stack and samples to be stored under a different repository other than the main source registry repo.

## Adding samples or stacks to the devfile registry

A file called ***extraDevfileEntries.yaml*** is added under the root of the devfile registry source repository to add samples and stacks from other repositories to the registry. This file contains the location information on where the extra samples and stacks can be found during a registry build.

### Sample extraDevfileEntires.yaml:

__Note:___ Proposal #2 will be used

#### Proposal 1 (deferred: will only implement later when needed):
If we assume the devfile in a given sample contains the information on the sample, we can future simplify it by just referring the sample to the location. All the metadata information is extracted directly from the devfile contained in the sample/stack repository.
```yaml
    schemaVersion: 1.0.0
    samples:
    - name: nodejs-basic
      git:
        remotes:
          origin: https://github.com/redhat-developer/devfile-sample/
    - name: vertx-secured-http
      git:
        remotes:
          origin: https://github.com/openshift-vertx-examples/vertx-secured-http-example-redhat
    - name: my-maven-sample
      zip:
          location: https://my.company.com/samples/my-maven-sample.zip
    stacks:
    - name: my-maven
      git:
        remotes:
          origin: https://github.com/eystacks/my-maven
```

__Note:__ the location of the sample supports the same configuration as the `starterProjects` definition, i.e. `git` and `zip`. Refer to the definitions of the existing `git` and `zip` elements for supported settings.

#### Proposal 2:

Alternatively, if we assume the devfile contained in the sample folder may refer to the original stack instead of the specific example, we may need to include a way for the user to specify the metadata associated with the sample as part of the sample definition.
```yaml
    schemaVersion: 1.0.0
    samples:
      - name: nodejs-basic
        displayName: Basic NodeJS
        description: A simple Hello World Node.js application
        icon: https://raw.githubusercontent.com/maysunfaisal/node-bulletin-board-2/main/nodejs-icon.png
        tags: ["NodeJS", "Express"]
        projectType: nodejs
        language: nodejs
        git:
          remotes:
            origin: https://github.com/redhat-developer/devfile-sample.git
      - name: code-with-quarkus
        displayName: Basic Quarkus
        description: A simple Hello World Java application using Quarkus
        icon: https://raw.githubusercontent.com/elsony/devfile-sample-code-with-quarkus/main/.devfile/icon/quarkus.png
        tags: ["Java", "Quarkus"]
        projectType: quarkus
        language: java
        git:
          remotes:
            origin: https://github.com/elsony/devfile-sample-code-with-quarkus.git
      - name: java-springboot-basic
        displayName: Basic Spring Boot
        description: A simple Hello World Java Spring Boot application using Maven
        icon: https://raw.githubusercontent.com/elsony/devfile-sample-java-springboot-basic/main/.devfile/icon/spring-logo.png
        tags: ["Java", "Spring"]
        projectType: springboot
        language: java
        git:
          remotes:
            origin: https://github.com/elsony/devfile-sample-java-springboot-basic.git
      - name: python-basic
        displayName: Basic Python
        description: A simple Hello World application using Python
        icon: https://raw.githubusercontent.com/elsony/devfile-sample-python-basic/main/.devfile/icon/python.png
        tags: ["Python"]
        projectType: python
        language: python
        git:
          remotes:
            origin: https://github.com/elsony/devfile-sample-python-basic.git
```

#### Proposal 3: Include the source of the samples as part of the registry source repository (deferred: will only implement later when needed)
The existing devfile registry source for the stacks have the following structure:

    ── stacks
      │── java-maven
      │   └── devfile.yaml
      │
      │── java-openliberty
      │   └── devfile.yaml
      │
      │── java-quarkus
      │   └── devfile.yaml
      │
      │── nodejs
      │   │── build
      │   │   └── Dockerfile
      │   │── deploy
      │   │   └── deployment-manifest.yaml
      │   └── devfile.yaml
                …


With the introduction of the samples in the registry, the samples will be stored under a similar directory structure if the sample source is stored under the registry source repo:

    ── stacks
      │── java-maven
      │   └── devfile.yaml
      │
      │── java-openliberty
      │   └── devfile.yaml
      │
      │── java-quarkus
      │   └── devfile.yaml
      │
      │── nodejs
      │   │── build
      │   │   └── Dockerfile
      │   │── deploy
      │   │   └── deployment-manifest.yaml
      │   └── devfile.yaml
                …
    ── samples
      │── nodejs-sample
      │   │── .gitignore
      │   │       LICENSE
      │   │       README.md
      │   │       devfile.yaml
      │   └── src
      │              └── Dockerfile
      │                       package-lock.json
      │                       package.json
      │                       server.js
      │
      │── springboot-sample
      │   │── devfile.yaml
      │   │       .gitignore
      │   │       LICENSE
      │   │       README.md
      │   │       mvnw
      │   │       mvnw.cmd
      │   │       pom.xml
      │   └── src
      │              │── main
      │              │      java
      │              │           │── com
          │              │             ...
                …


## Registry query results
The information on the registry index (index.json) is different between a stack and a sample. The differences are based on how the two will be used.
```json
    [
        {
            "name": "java-maven",
            "version": "1.1.0",
            "displayName": "Maven Java",
            "description": "Upstream Maven and OpenJDK 11",
            "type": "stack",
            "tags": [
                "Java",
                "Maven"
            ],
            "projectType": "maven",
            "language": "java",
            "links": {
                "self": "devfile-catalog/java-maven:latest"
            },
            "resources": [
                "devfile.yaml"
            ],
            "starterProjects": [
                "springbootproject"
            ]
        },
        {
            "name": "nodejs-basic",
            "displayName": "Basic NodeJS",
            "description": "A simple Hello world NodeJS application",
            "icon": "nodejsIcon.svg",
            "type": "sample",
            "tags": [
                "NodeJS",
                "Express"
            ],
            "projectType": "nodejs",
            "language": "nodejs",
            "git": {
                "remotes": {
                    "origin": "https://github.com/redhat-developer/devfile-sample/"
                }
            }
        }
    ]
```

The differences between a stack entry vs a sample entry in the devfile index registry are:
1. The `type` field
2. A stack has the `resources` and `starterProjects` while the sample has the resource location information, e.g. `git` or `zip` location information

## Registry library change
The registry library and the registry REST API will provide a filtering mechanism to specify which type (sample or stack) is needed in the query.
