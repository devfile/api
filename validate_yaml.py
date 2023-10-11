#!/usr/bin/env python
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
import json
import yaml
import sys

from typing import Any
from jsonschema import validate, ValidationError


class ParseSchemaError(Exception):
    pass


class YamlValidationError(Exception):
    pass


class OpenFileError(Exception):
    pass


class NotEnoughArgsError(Exception):
    pass


class YamlValidator:
    """
    Yaml validator validates a given yaml file against
    a chosen template.
    """
    def __init__(self, schema_path: str) -> None:
        self.schema = self._parse_json_file(schema_path)

    def _open_file(self, path: str) -> Any:
        try:
            return open(path)
        except OSError as exc:
            raise OpenFileError(f"::error:: failed to open file {path}: {exc}")

    def _parse_json_file(self, json_path: str) -> dict[str, Any]:
        return json.load(self._open_file(json_path))

    def _get_yaml_file(self, yaml_path: str):
        return yaml.load(self._open_file(yaml_path), Loader=yaml.SafeLoader)

    def validate(self, path: str) -> bool:
        try:
            _ = validate(instance=self._get_yaml_file(path), schema=self.schema)
        except ValidationError as exc:
            raise YamlValidationError(f"error:: validation failed: {str(exc.message)}")
        return True


def parse_arg(index: int) -> str:
    try:
        return sys.argv[index]
    except IndexError:
        raise NotEnoughArgsError(
            "Missing Args: Example usage -> validate-yaml.py <schema_path> <yaml_path>"
        )


if __name__ == "__main__":
    schema_path = parse_arg(1)
    yaml_path = parse_arg(2)
    validator = YamlValidator(schema_path=schema_path)
    validator.validate(yaml_path)