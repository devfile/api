#
# Copyright (c) 2021 Red Hat, Inc.
# This program and the accompanying materials are made
# available under the terms of the Eclipse Public License 2.0
# which is available at https://www.eclipse.org/legal/epl-2.0/
#
# SPDX-License-Identifier: EPL-2.0
#
# Contributors:
#   Red Hat, Inc. - initial API and implementation
#

import os
import json
import yaml

def write_json(filename: str, object: dict) -> None:
    """
    Write the json object to the specified filename
    """
    with open(filename, 'w') as out:
        json.dump(object, out, sort_keys=False, indent=2, separators=(',', ': '), ensure_ascii=True)

def create_ref(path):
    """
    Create a json definition reference to a specific path 
    """
    return '#/definitions/' + path

def retrieve_metadata() -> object:
    """
    Retrieve the metadata json field from the latest devworkspace json schema 
    """

    devworkspace_json_schema_path = os.path.join('schemas', 'latest', 'dev-workspace.json')
    with open(devworkspace_json_schema_path, 'r') as f:
        devworkspace_json_schema = json.load(f)
        metadata = devworkspace_json_schema['properties']['metadata']
        return metadata

def consolidate_crds() -> object:
    """
    Consolidate all crds in /crds into one json object
    """
    crds_dir = os.path.join('crds')
    crds = os.listdir(crds_dir)
    consolidated_crds_json = {
        'definitions': {},
    }
    additional_metadata = retrieve_metadata()
    for file in crds:
        crd_file_path = os.path.join(crds_dir, file) 
        with open(crd_file_path) as file:
            yamlData = yaml.load(file, Loader=yaml.FullLoader)
            crd_name = yamlData['spec']['names']['kind']
            
            # Add all the available schema versions 
            for version in yamlData['spec']['versions']:
                new_json_name = version['name'] + '.' + crd_name
                version['schema']['openAPIV3Schema']['properties']['metadata'] = additional_metadata
                new_schema = version['schema']['openAPIV3Schema']
                consolidated_crds_json['definitions'][new_json_name] = new_schema

    return consolidated_crds_json

def add_property_definition(root_definitions_object: dict, current_path: str, curr_object: dict, queue: list) -> None:
    """
    Given an object, convert the child properties into references with new definitions at the root of root_definitions_object.
    Also removes oneOf references since they aren't supported by openapi-generator.

    Converts:
        {
            "properties": {
                "foo": {
                    "type": "object",
                    "properties": {
                        "bar": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    
    into:
        {
            "definitions": {
                "foo": {
                    "type": "object",
                    "properties": {
                        "bar": {
                            "$ref": "#/definitions/bar"
                        }
                    }
                },
                "bar": {
                    "type": string
                }
            }
        }
    """
    for prop in curr_object['properties']:
        new_path = current_path + '.' + prop
        new_object = curr_object['properties'][prop]

        # openapi-generator doesn't accept oneOf so we have to remove them
        if 'oneOf' in new_object:
            del new_object['oneOf']

        root_definitions_object[new_path] = new_object

        # openapi-generator doesn't accept oneOf so we have to remove them
        if 'items' in new_object:
            if 'oneOf' in new_object['items']:
                del new_object['items']['oneOf']
            new_path += ".items"

        queue.append({
            new_path: new_object
        })
        curr_object['properties'][prop] = {
            '$ref': create_ref(new_path)
        }

def add_item_definition(root_definitions_object: dict, current_path: str, curr_object: dict, queue: list) -> None:
    """
    Given an object, convert the child properties into references with new definitions at the root of root_definitions_object.
    Also removes oneOf references since they aren't supported by openapi-generator.

    Converts:
        {
            "v1devworkspace": {
                "properties": {
                    "spec": {
                        "items": {
                            "type": "object",
                            "properties": {
                                "foo": {
                                    "type": "string",
                                    "description": "Type of funding or platform through which funding is possible."
                                },
                            }
                        },
                        "type": "array"    
                    }
                }
                
            }
        }
    
    into:
        {
            "definitions": {
                "v1devworkspace": {
                    "properties": {
                        "spec": {
                            "$ref": "#/definitions/v1devworkspace.spec.items"
                        }
                    }
                },
                "v1devworkspace.spec.items": {
                    "items": {
                        "$ref": "#/definitions/v1devworkspace.spec"
                    },
                    "type": "array"
                },
                "v1devworkspace.spec": {
                    "type": "object",
                    "properties": {
                        "foo": {
                            "$ref": "#/definitions/v1devworkspace.spec.items.foo"
                        },
                    }
                }
                "v1devworkspace.spec.items.foo": {
                    "type": "string",
                    "description": "Type of funding or platform through which funding is possible."
                },
            }
        }
    """
    if 'properties' in curr_object['items']:
        root_definitions_object[current_path] = curr_object

        path = current_path
        pathList = current_path.split('.')
        if pathList[-1] == 'items':
            pathList = pathList[:-1]
            path = '.'.join(pathList)

        for prop in curr_object['items']['properties']:
            new_path = current_path + '.' + prop
            new_object = curr_object['items']['properties'][prop]

            # openapi-generator doesn't accept oneOf so we have to remove them
            if 'oneOf' in new_object:
                del new_object['oneOf']
            root_definitions_object[new_path] = new_object
            queue.append({
                new_path: new_object
            })
            curr_object['items']['properties'][prop] = {
                '$ref': create_ref(new_path)
            }
        root_definitions_object[path] = curr_object['items']
        curr_object['items'] = {
            '$ref': create_ref(path)
        }
    else:
        root_definitions_object[current_path] = curr_object

def add_definition(root_definitions_object: dict, current_path: str, curr_object: dict, queue: list) -> None:
    """
    Create a property or item definition depending on if property or items is in the current_object
    """
    if 'properties' in curr_object:
        add_property_definition(root_definitions_object, current_path, curr_object, queue)
    elif 'items' in curr_object:
        add_item_definition(root_definitions_object, current_path, curr_object, queue)

def flatten(consolidated_crds_object: dict) -> None:
    """
    Flatten and then produce a new swagger.json file that can be processed by open-api-generator
    """
    original_definitions = consolidated_crds_object['definitions']
    flattened_swagger_object = {
        'definitions': {},
        'paths': {},
        'info': {
            'title': 'Kubernetes',
            'version': 'unversioned'
        },
        'swagger': '2.0'
    }
    for root in original_definitions:
        flattened_swagger_object['definitions'][root] = original_definitions[root]

        queue = []

        # Add in all the initial properties to the queue
        for prop in original_definitions[root]['properties']:
            new_path = root + '.' + prop
            queue.append({
                new_path: original_definitions[root]['properties'][prop]
            })

            # Create a new definition so that the properties are pulled out correctly
            flattened_swagger_object['definitions'][new_path] = original_definitions[root]['properties'][prop]

            # Create a ref from the property such as spec to the new path such as v1alpha1.devworkspaces.workspace.devfile.io_spec
            original_definitions[root]['properties'][prop] = {
                '$ref': create_ref(new_path)
            }
        
        # Continue until all properties have been flattened
        while len(queue) != 0:
            next_item = queue.pop().popitem()
            path = next_item[0]
            new_object = next_item[1]
            add_definition(flattened_swagger_object['definitions'], path, new_object, queue)

    write_json('swagger.json', flattened_swagger_object)

def devfile_schema_to_crd() -> object:
    """
    Convert the devfile schema to a crd so that we can generate the types 
    """

    devfile_json_schema_path = os.path.join('schemas', 'latest', 'devfile.json')
    with open(devfile_json_schema_path, 'r') as devfileFile:
        devfile_json_schema = json.load(devfileFile)
        devfile_crd = {
            'apiVersion': 'apiextensions.k8s.io/v1',
            'kind': 'CustomResourceDefinition',
            'metadata': {
                'creationTimestamp': None,
                'name': 'devfile.workspace.devfile.io'
            },
            'spec': {
                'group': 'devfile.devfile.io',
                'names': {
                    'kind': 'Devfile',
                    'listKind': 'DevfileList',
                    'plural': 'Devfiles',
                    'shortNames': [],
                    'singular': 'devfile'
                },
                'scope': 'Namespaced',
                'versions': [
                    {
                        'name': 'test',
                        'schema': {
                            'openAPIV3Schema': devfile_json_schema
                        },
                        'served': True,
                        'storage': True,
                        'subresources': {
                            'status': {

                            }
                        }
                    }
                ]
            },
            'status': {
                'acceptedNames': {
                    'kind': '',
                    'plural': ''
                },
                'conditions': [],
                'storedVersions': []
            }
        }
        with open(os.path.join('crds', 'workspace.devfile.io_devfile.yaml'), 'w') as crdFile:
            yaml.dump(devfile_crd, crdFile)

if __name__ == "__main__":
    # Create a devfile crd that will be used to generate types
    devfile_schema_to_crd()

    # Get the crds and flatten them
    swagger_crds_json = consolidate_crds()
    flatten(swagger_crds_json)
