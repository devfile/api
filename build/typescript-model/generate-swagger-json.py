#
#
# Copyright Red Hat
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os
import json


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

def consolidate_schemas() -> dict:
    """
    Consolidate all schemas into one json object
    """
    schemas_dir = os.path.join('schemas/latest')
    consolidated_schemas_json = {
        'definitions': {},
    }

    with open(os.path.join(schemas_dir, 'k8sApiVersion.txt')) as f:
        k8sApiVersion = f.readline()

    with open(os.path.join(schemas_dir, 'jsonSchemaVersion.txt')) as f:
        devfileVersion = f.readline()
        devfileVersion = 'v' + devfileVersion.replace('-alpha', '')

    definitionName = devfileVersion + '.Devfile'
    devfile_json_schema_path = os.path.join(schemas_dir, 'devfile.json')
    with open(devfile_json_schema_path, 'r') as devfileFile:
        jsonData = json.load(devfileFile)
        consolidated_schemas_json['definitions'][definitionName] = jsonData

    definitionName = k8sApiVersion + '.DevWorkspace'
    dw_json_schema_path = os.path.join(schemas_dir, 'dev-workspace.json')
    with open(dw_json_schema_path, 'r') as devfileFile:
        jsonData = json.load(devfileFile)    
        consolidated_schemas_json['definitions'][definitionName] = jsonData

    definitionName = k8sApiVersion + '.DevWorkspaceTemplate'        
    dwt_json_schema_path = os.path.join(schemas_dir, 'dev-workspace-template.json')
    with open(dwt_json_schema_path, 'r') as devfileFile:
        jsonData = json.load(devfileFile)
        consolidated_schemas_json['definitions'][definitionName] = jsonData

    return consolidated_schemas_json

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
            if 'items' in original_definitions[root]['properties'][prop]:
                new_path += '.items'

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

if __name__ == "__main__":
    # Get the schema and flatten them
    swagger_crds_json = consolidate_schemas()
    flatten(swagger_crds_json)
