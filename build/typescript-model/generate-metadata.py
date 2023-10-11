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
import argparse
import yaml


def write_contents(filename: str, mode: str, contents: str) -> None:
    """
    Write the string to the specified filename
    """
    with open(filename, mode) as out:
        out.write(contents)


def get_crd_metadata(output_path: str) -> None:
    """
    Read in the devworkspace and devworkspace template crds and generate metadata into api/apis.ts
    """
    crd_path = "crds"
    typescript_contents = ""
    devworkspace_crd_path = os.path.join(crd_path, 'workspace.devfile.io_devworkspaces.yaml')
    with open(devworkspace_crd_path, 'r') as devfile_file:
        yaml_data = yaml.load(devfile_file, Loader=yaml.FullLoader)
        spec, group, kind, plural, singular, versions, latest_version, latest_api_version = extract_fields(yaml_data)
        typescript_contents += generate_typescript(latest_api_version, group, kind, plural, singular, versions,
                                                   latest_version)

    devworkspacetemplate_crd_path = os.path.join(crd_path, 'workspace.devfile.io_devworkspacetemplates.yaml')
    with open(devworkspacetemplate_crd_path, 'r') as devfile_file:
        yaml_data = yaml.load(devfile_file, Loader=yaml.FullLoader)
        spec, group, kind, plural, singular, versions, latest_version, latest_api_version = extract_fields(yaml_data)
        typescript_contents += generate_typescript(latest_api_version, group, kind, plural, singular, versions,
                                                   latest_version)

    write_contents(os.path.join(output_path, "constants", "constants.ts"), "w", typescript_contents)


def extract_fields(yaml_data: {}) -> (str, str, str, str, str, [], str, str):
    """
    Extract metadata from the crds
    """
    spec = yaml_data['spec']
    group = spec['group']
    kind = spec['names']['kind']
    plural = spec['names']['plural']
    singular = spec['names']['singular']
    versions = [version['name'] for version in spec['versions']]
    latest_version = versions[len(versions) - 1]
    latest_api_version = "{}/{}".format(group, latest_version)
    return spec, group, kind, plural, singular, versions, latest_version, latest_api_version


def generate_typescript(api_version: str, group: str, kind: str, plural: str, singular: str, versions: [],
                        latest_version: str) -> str:
    """
    Export a string representation of the typescript
    """
    return f"""
export const {singular + "ApiVersion"} = '{api_version}';
export const {singular + "Group"} = '{group}';
export const {singular + "Kind"} = '{kind}';
export const {singular + "Plural"} = '{plural}';
export const {singular + "Singular"} = '{singular}';
export const {singular + "Versions"} = {versions};
export const {singular + "LatestVersion"} = '{latest_version}';
    """


def export_typescript_api(output_path: str) -> None:
    """
    Export constants into api.ts
    """
    export_contents = """
export * from './constants/constants';
    """
    write_contents(os.path.join(output_path, "api.ts"), "a", export_contents)


if __name__ == "__main__":
    # Get any additional metadata we can from the crds
    parser = argparse.ArgumentParser(description='Generate metadata from crds')
    parser.add_argument('-p', '--path', action='store', type=str, help='The path to the constants directory')

    args = parser.parse_args()
    if not args.path:
        parser.print_help()
        parser.exit()

    path = args.path

    # Grab the metadata from the crds and put it into constants/constant.ts in typescript-model
    get_crd_metadata(path)

    # Export constants/constant.ts so that you can import constants from the package
    export_typescript_api(path)
