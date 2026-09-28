//
//
// Copyright Red Hat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"encoding/json"

	"github.com/devfile/api/v2/pkg/attributes"
)

// UnmarshalJSON decodes an Endpoint, converting attribute values that are not JSON strings into
// their string representation.
//
// Endpoint attributes are free-form (`map[string]apiext.JSON`) in v1alpha2 but string-based
// (`map[string]string`) in v1alpha1, so an attribute such as `discoverable: true` is valid in
// v1alpha2 yet cannot be decoded as-is here. Conversion from v1alpha2 is implemented as a JSON
// round-trip, so without this the conversion of an endpoint holding a non-string attribute would
// fail with "cannot unmarshal bool into Go struct field Endpoint.container.endpoints.attributes of
// type string".
func (endpoint *Endpoint) UnmarshalJSON(data []byte) error {
	// The alias prevents json.Unmarshal from recursing into this method, and the shadowing
	// `Attributes` field captures the free-form attributes before they reach the string-based one.
	type endpointAlias Endpoint

	decoded := &struct {
		Attributes attributes.Attributes `json:"attributes,omitempty"`
		*endpointAlias
	}{
		endpointAlias: (*endpointAlias)(endpoint),
	}
	if err := json.Unmarshal(data, decoded); err != nil {
		return err
	}
	endpoint.Attributes = stringifyAttributes(decoded.Attributes)
	return nil
}

// stringifyAttributes converts free-form attributes into the string-based map v1alpha1 expects.
// Strings, booleans and numbers are converted to their string representation; any other value
// (an object or an array) keeps its raw JSON representation so that nothing is silently dropped.
func stringifyAttributes(attrs attributes.Attributes) map[string]string {
	if attrs == nil {
		return nil
	}
	stringAttributes := make(map[string]string, len(attrs))
	for key, value := range attrs {
		var err error
		stringValue := attrs.GetString(key, &err)
		if err != nil {
			stringValue = string(value.Raw)
		}
		stringAttributes[key] = stringValue
	}
	return stringAttributes
}
