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
  "testing"

  "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
  "github.com/devfile/api/v2/pkg/attributes"
  "github.com/stretchr/testify/assert"
)

// Endpoint attributes are typed `map[string]apiext.JSON` in v1alpha2 but `map[string]string` in
// v1alpha1, so any non-string attribute value (e.g. `discoverable: true`) breaks the JSON round-trip
// the conversion relies on. Non-string values must be converted to their string representation
// instead of being dropped or failing the conversion.
func TestComponentConversionFrom_v1alpha2_EndpointAttributes(t *testing.T) {
	for _, tt := range getEndpointAttributeConversionTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			src := &v1alpha2.Component{
				Name: "postgresql",
				ComponentUnion: v1alpha2.ComponentUnion{
					Container: &v1alpha2.ContainerComponent{
						Container: v1alpha2.Container{
							Image: "postgres:latest",
						},
						Endpoints: []v1alpha2.Endpoint{
							{
								Name:       "postgresql",
								TargetPort: 5432,
								Exposure:   v1alpha2.InternalEndpointExposure,
								Attributes: tt.attributes,
							},
						},
					},
				},
			}

			output := &Component{}

			err := convertComponentFrom_v1alpha2(src, output)
			if !assert.NoError(t, err, "Should not return error when converting from v1alpha2") {
				return
			}

			if !assert.NotNil(t, output.Container, "Container component should be converted") {
				return
			}
			if !assert.Len(t, output.Container.Endpoints, 1, "Endpoint should be converted") {
				return
			}
			assert.Equal(t, tt.expected, output.Container.Endpoints[0].Attributes,
				"Endpoint attributes should be converted to their string representation")
		})
	}
}

// Endpoint attributes of a container overridden by a plugin component go through their own JSON
// round-trip, which checks the unmarshalling error and so fails the whole conversion with
// "cannot unmarshal bool into Go struct field Endpoint.container.endpoints.attributes of type string".
func TestPluginComponentConversionFrom_v1alpha2_EndpointAttributes(t *testing.T) {
	for _, tt := range getEndpointAttributeConversionTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			src := &v1alpha2.Component{
				Name: "my-plugin",
				ComponentUnion: v1alpha2.ComponentUnion{
					Plugin: &v1alpha2.PluginComponent{
						ImportReference: v1alpha2.ImportReference{
							ImportReferenceUnion: v1alpha2.ImportReferenceUnion{
								Uri: "https://example.com/plugin.yaml",
							},
						},
						PluginOverrides: v1alpha2.PluginOverrides{
							Components: []v1alpha2.ComponentPluginOverride{
								{
									Name: "postgresql",
									ComponentUnionPluginOverride: v1alpha2.ComponentUnionPluginOverride{
										Container: &v1alpha2.ContainerComponentPluginOverride{
											ContainerPluginOverride: v1alpha2.ContainerPluginOverride{
												Image: "postgres:latest",
											},
											Endpoints: []v1alpha2.EndpointPluginOverride{
												{
													Name:       "postgresql",
													TargetPort: 5432,
													Attributes: tt.attributes,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			}

			output := &Component{}

			err := convertComponentFrom_v1alpha2(src, output)
			if !assert.NoError(t, err, "Should not return error when converting from v1alpha2") {
				return
			}

			if !assert.Len(t, output.Plugin.Components, 1, "Plugin component override should be converted") {
				return
			}
			overriddenContainer := output.Plugin.Components[0].Container
			if !assert.NotNil(t, overriddenContainer, "Container override should be converted") {
				return
			}
			if !assert.Len(t, overriddenContainer.Endpoints, 1, "Endpoint should be converted") {
				return
			}
			assert.Equal(t, tt.expected, overriddenContainer.Endpoints[0].Attributes,
				"Endpoint attributes should be converted to their string representation")
		})
	}
}

// Endpoint attributes of a component overridden by a parent go through their own JSON round-trip,
// and so have to handle non-string attributes as well. Kubernetes and Openshift components declare
// endpoints too.
func TestParentComponentConversionFrom_v1alpha2_EndpointAttributes(t *testing.T) {
	for _, tt := range getEndpointAttributeConversionTestCases() {
		t.Run(tt.name, func(t *testing.T) {
			src := &v1alpha2.ComponentParentOverride{
				Name: "postgresql",
				ComponentUnionParentOverride: v1alpha2.ComponentUnionParentOverride{
					Kubernetes: &v1alpha2.KubernetesComponentParentOverride{
						K8sLikeComponentParentOverride: v1alpha2.K8sLikeComponentParentOverride{
							K8sLikeComponentLocationParentOverride: v1alpha2.K8sLikeComponentLocationParentOverride{
								Inlined: "kubernetes-resource",
							},
							Endpoints: []v1alpha2.EndpointParentOverride{
								{
									Name:       "postgresql",
									TargetPort: 5432,
									Attributes: tt.attributes,
								},
							},
						},
					},
				},
			}
			output := &Component{}

			err := convertParentComponentFrom_v1alpha2(src, output)
			if !assert.NoError(t, err, "Should not return error when converting from v1alpha2") {
				return
			}

			if !assert.NotNil(t, output.Kubernetes, "Kubernetes component should be converted") {
				return
			}
			if !assert.Len(t, output.Kubernetes.Endpoints, 1, "Endpoint should be converted") {
				return
			}
			assert.Equal(t, tt.expected, output.Kubernetes.Endpoints[0].Attributes,
				"Endpoint attributes should be converted to their string representation")
		})
	}
}

func getEndpointAttributeConversionTestCases() []struct {
	name       string
	attributes attributes.Attributes
	expected   map[string]string
} {
	return []struct {
		name       string
		attributes attributes.Attributes
		expected   map[string]string
	}{
		{
			name:       "boolean attribute value",
			attributes: attributes.Attributes{}.PutBoolean("discoverable", true),
			expected:   map[string]string{"discoverable": "true"},
		},
		{
			name:       "number attribute value",
			attributes: attributes.Attributes{}.PutInteger("weight", 10),
			expected:   map[string]string{"weight": "10"},
		},
		{
			name:       "string attribute value",
			attributes: attributes.Attributes{}.PutString("type", "terminal"),
			expected:   map[string]string{"type": "terminal"},
		},
		{
			name:       "large number attribute value",
			attributes: attributes.Attributes{}.PutInteger("size", 1048576),
			expected:   map[string]string{"size": "1048576"},
		},
		{
			name:       "object attribute value",
			attributes: attributes.Attributes{}.Put("meta", map[string]interface{}{"a": 1}, nil),
			expected:   map[string]string{"meta": `{"a":1}`},
		},
		{
			name:       "array attribute value",
			attributes: attributes.Attributes{}.Put("ports", []int{1, 2}, nil),
			expected:   map[string]string{"ports": "[1,2]"},
		},
		{
			name:       "null attribute value",
			attributes: attributes.Attributes{}.Put("discoverable", nil, nil),
			expected:   map[string]string{"discoverable": "null"},
		},
		{
			name:       "empty attributes",
			attributes: attributes.Attributes{},
			expected:   nil,
		},
	}
}
