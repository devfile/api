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

// componentsWithEndpoints lists the keys of the component types that declare endpoints.
var componentsWithEndpoints = []string{"container", "kubernetes", "openshift"}

// stringifyComponentEndpointAttributes rewrites the endpoint attributes of a marshalled v1alpha2
// component so that it can be decoded into its v1alpha1 counterpart.
//
// Endpoint attributes are free-form (`map[string]apiext.JSON`) in v1alpha2 but string-based
// (`map[string]string`) in v1alpha1, so an attribute such as `discoverable: true` is valid in
// v1alpha2 yet cannot be decoded as-is here. Conversion from v1alpha2 is implemented as a JSON
// round-trip, so without this the conversion of a component holding an endpoint with a non-string
// attribute would fail with "cannot unmarshal bool into Go struct field
// Endpoint.container.endpoints.attributes of type string".
//
// All the v1alpha2 component flavours (`Component`, `ComponentPluginOverride` and
// `ComponentParentOverride`) share the same JSON representation, so this operates on the
// marshalled bytes rather than on the Go types.
func stringifyComponentEndpointAttributes(data []byte) ([]byte, error) {
	component := map[string]json.RawMessage{}
	if err := json.Unmarshal(data, &component); err != nil {
		return nil, err
	}

	rewritten := false
	for _, componentType := range componentsWithEndpoints {
		body, found := component[componentType]
		if !found {
			continue
		}
		updatedBody, updated, err := stringifyEndpointAttributes(body)
		if err != nil {
			return nil, err
		}
		if updated {
			component[componentType] = updatedBody
			rewritten = true
		}
	}

	// Leave the document untouched when there is nothing to convert.
	if !rewritten {
		return data, nil
	}
	return json.Marshal(component)
}

// stringifyEndpointAttributes rewrites the attributes of the endpoints declared by a marshalled
// component body, and reports whether anything was rewritten.
func stringifyEndpointAttributes(body json.RawMessage) (json.RawMessage, bool, error) {
	componentBody := map[string]json.RawMessage{}
	if err := json.Unmarshal(body, &componentBody); err != nil {
		return nil, false, err
	}
	rawEndpoints, found := componentBody["endpoints"]
	if !found {
		return nil, false, nil
	}

	var endpoints []map[string]json.RawMessage
	if err := json.Unmarshal(rawEndpoints, &endpoints); err != nil {
		return nil, false, err
	}

	rewritten := false
	for _, endpoint := range endpoints {
		rawAttributes, found := endpoint["attributes"]
		if !found {
			continue
		}
		freeFormAttributes := attributes.Attributes{}
		if err := json.Unmarshal(rawAttributes, &freeFormAttributes); err != nil {
			return nil, false, err
		}
		if len(freeFormAttributes) == 0 {
			continue
		}
		stringAttributes, err := json.Marshal(stringifyAttributes(freeFormAttributes))
		if err != nil {
			return nil, false, err
		}
		endpoint["attributes"] = stringAttributes
		rewritten = true
	}

	if !rewritten {
		return nil, false, nil
	}

	updatedEndpoints, err := json.Marshal(endpoints)
	if err != nil {
		return nil, false, err
	}
	componentBody["endpoints"] = updatedEndpoints

	updatedBody, err := json.Marshal(componentBody)
	if err != nil {
		return nil, false, err
	}
	return updatedBody, true, nil
}

// stringifyAttributes converts free-form attributes into the string-based map v1alpha1 expects.
// A JSON string is unquoted; every other value (boolean, number, object, array) keeps its verbatim
// JSON text, so that no precision is lost: `1048576` becomes "1048576", and not the "1.048576e+06"
// that `attributes.Attributes.GetString` would produce.
//
// The conversion is one way. Nothing parses these strings back when converting to v1alpha2 again,
// so `discoverable: true` comes out of a v1alpha2 -> v1alpha1 -> v1alpha2 round-trip as the string
// "true". Scalars stay usable, because `GetBoolean` and `GetNumber` fall back to strconv when the
// attribute holds a string, but an object or an array comes back as a string and no longer decodes
// with `GetInto`. Parsing the strings back is not an option: a string attribute the user actually
// authored as "true" cannot be told apart from a stringified boolean.
func stringifyAttributes(attrs attributes.Attributes) map[string]string {
	stringAttributes := make(map[string]string, len(attrs))
	for key, value := range attrs {
		// A JSON `null` is decoded into an empty Raw by apiext.JSON.
		if len(value.Raw) == 0 {
			stringAttributes[key] = "null"
			continue
		}
		var stringValue string
		if err := json.Unmarshal(value.Raw, &stringValue); err != nil {
			stringValue = string(value.Raw)
		}
		stringAttributes[key] = stringValue
	}
	return stringAttributes
}
