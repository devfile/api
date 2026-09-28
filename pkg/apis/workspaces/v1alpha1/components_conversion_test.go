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
	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
)

func TestComponentConversion_v1alpha1(t *testing.T) {
	f := fuzz.New().NilChance(fuzzNilChance).MaxDepth(100).Funcs(
		componentFuzzFunc,
		commandFuzzFunc,
		pluginComponentsOverrideFuzzFunc,
		pluginComponentFuzzFunc,
		rawExtFuzzFunc,
	)
	for i := 0; i < fuzzIterations; i++ {
		original := &Component{}
		intermediate := &v1alpha2.Component{}
		output := &Component{}
		f.Fuzz(original)
		input := original.DeepCopy()
		err := convertComponentTo_v1alpha2(input, intermediate)
		if !assert.NoError(t, err, "Should not return error when converting to v1alpha2") {
			return
		}
		err = convertComponentFrom_v1alpha2(intermediate, output)
		if !assert.NoError(t, err, "Should not return error when converting from v1alpha2") {
			return
		}
		if !assert.True(t, cmp.Equal(original, output), "Component should not be changed when converting between v1alpha1 and v1alpha2") {
			t.Logf("Diff: \n%s\n", cmp.Diff(original, output))
		}
	}
}

func TestComponentConversionFrom_v1alpha2(t *testing.T) {

	src := &v1alpha2.Component{
		Name: "test1",
		ComponentUnion: v1alpha2.ComponentUnion{
			Image: &v1alpha2.ImageComponent{
				Image: v1alpha2.Image{
					ImageName: "image:latest",
				},
			},
		},
	}
	output := &Component{}

	err := convertComponentFrom_v1alpha2(src, output)
	if !assert.NoError(t, err, "Should not return error when converting from v1alpha2") {
		return
	}

	assert.Equal(t, &Component{}, output, "Conversion from v1alpha2 should be skipped for Image Component")
}

// Endpoint attributes are typed `map[string]apiext.JSON` in v1alpha2 but `map[string]string` in
// v1alpha1, so any non-string attribute value (e.g. `discoverable: true`) breaks the JSON round-trip
// the conversion relies on. Non-string values must be converted to their string representation
// instead of being dropped or failing the conversion.
func TestComponentConversionFrom_v1alpha2_NonStringEndpointAttributes(t *testing.T) {
	tests := []struct {
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
	}

	for _, tt := range tests {
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
func TestPluginComponentConversionFrom_v1alpha2_NonStringEndpointAttributes(t *testing.T) {
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
											Attributes: attributes.Attributes{}.PutBoolean("discoverable", true),
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
	assert.Equal(t, map[string]string{"discoverable": "true"}, overriddenContainer.Endpoints[0].Attributes,
		"Endpoint attributes should be converted to their string representation")
}
