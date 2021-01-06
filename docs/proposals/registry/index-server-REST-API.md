# Index server REST APIs
This documentation explains how to use the index server REST APIs

## Gets registry index file
Gets the registry index file content from HTTP response
### HTTP request
```
GET http://{registry host}/index
```
### Request parameters
| Parameter | Description |
| --------- | ----------- |
| registry host | the URL/ingress that exposes registry service |
### Request body
The request body must be empty.
### Request example
```
curl http://devfile-registry.192.168.1.1.nip.io/index
```
### Response example
```json
[
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
    "links": {
      "self": "devfile-catalog/java-maven:latest"
    },
    "resources": [
      "devfile.yaml",
      "archive.tar"
    ]
  },
  {
    "name": "java-openliberty",
    "displayName": "Open Liberty",
    "description": "Open Liberty microservice in Java",
    "projectType": "docker",
    "language": "java",
    "links": {
      "self": "devfile-catalog/java-openliberty:latest"
    },
    "resources": [
      "devfile.yaml",
      "archive.tar"
    ]
  },
  {
    "name": "java-quarkus",
    "displayName": "Quarkus Java",
    "description": "Upstream Quarkus with Java+GraalVM",
    "tags": [
      "Java",
      "Quarkus"
    ],
    "projectType": "quarkus",
    "language": "java",
    "links": {
      "self": "devfile-catalog/java-quarkus:latest"
    },
    "resources": [
      "devfile.yaml",
      "archive.tar"
    ]
  },
  {
    "name": "java-springboot",
    "displayName": "Spring Boot®",
    "description": "Spring Boot® using Java",
    "tags": [
      "Java",
      "Spring"
    ],
    "projectType": "spring",
    "language": "java",
    "links": {
      "self": "devfile-catalog/java-springboot:latest"
    },
    "resources": [
      "devfile.yaml",
      "archive.tar"
    ]
  },
  {
    "name": "java-vertx",
    "displayName": "Vert.x Java",
    "description": "Upstream Vert.x using Java",
    "tags": [
      "Java",
      "Vert.x"
    ],
    "projectType": "vertx",
    "language": "java",
    "links": {
      "self": "devfile-catalog/java-vertx:latest"
    },
    "resources": [
      "devfile.yaml",
      "archive.tar"
    ]
  },
  {
    "name": "nodejs",
    "displayName": "NodeJS Runtime",
    "description": "Stack with NodeJS 12",
    "tags": [
      "NodeJS",
      "Express",
      "ubi8"
    ],
    "projectType": "nodejs",
    "language": "nodejs",
    "links": {
      "self": "devfile-catalog/nodejs:latest"
    },
    "resources": [
      "devfile.yaml",
      "archive.tar"
    ]
  },
  {
    "name": "python",
    "displayName": "Python",
    "description": "Python Stack with Python 3.7",
    "tags": [
      "Python",
      "pip"
    ],
    "projectType": "python",
    "language": "python",
    "links": {
      "self": "devfile-catalog/python:latest"
    },
    "resources": [
      "devfile.yaml",
      "archive.tar"
    ]
  },
  {
    "name": "python-django",
    "displayName": "Django",
    "description": "Python3.7 with Django",
    "tags": [
      "Python",
      "pip",
      "Django"
    ],
    "projectType": "django",
    "language": "python",
    "links": {
      "self": "devfile-catalog/python-django:latest"
    },
    "resources": [
      "devfile.yaml",
      "archive.tar"
    ]
  }
]
```

## Gets registry stack devfile
Gets the specific registry stack devfile content from HTTP response

Note: this REST API only returns the content of `devfile.yaml`, it won't return other resources in the stack
### HTTP request
```
GET http://{registry host}/devfiles/{stack}
```
### Request parameters
| Parameter | Description |
| ----------| ----------- |
| registry host | the URL/ingress that exposes registry service |
| stack | registry stack name |
### Request body
The request body must be empty.
### Request example
```
curl http://devfile-registry.192.168.1.1.nip.io/devfiles/nodejs
```
### Response example
```yaml
schemaVersion: 2.0.0
metadata:
  name: nodejs
  version: 1.0.0
starterProjects:
  - name: nodejs-starter
    git:
      remotes:
        origin: "https://github.com/odo-devfiles/nodejs-ex.git"
components:
  - name: runtime
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
      component: runtime
      commandLine: npm install
      workingDir: /project
      group:
        kind: build
        isDefault: true
  - id: run
    exec:
      component: runtime
      commandLine: npm start
      workingDir: /project
      group:
        kind: run
        isDefault: true
  - id: debug
    exec:
      component: runtime
      commandLine: npm run debug
      workingDir: /project
      group:
        kind: debug
        isDefault: true
  - id: test
    exec:
      component: runtime
      commandLine: npm test
      workingDir: /project
      group:
        kind: test
        isDefault: true
```