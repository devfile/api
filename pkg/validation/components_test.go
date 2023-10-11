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

package validation

import (
	"testing"

	"github.com/devfile/api/v2/pkg/attributes"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/assert"
)

// generateDummyContainerComponent returns a dummy container component for testing
func generateDummyContainerComponent(name string, volMounts []v1alpha2.VolumeMount, endpoints []v1alpha2.Endpoint, envs []v1alpha2.EnvVar, annotation v1alpha2.Annotation, dedicatedPod bool) v1alpha2.Component {
	image := "docker.io/maven:latest"
	mountSources := true

	return v1alpha2.Component{
		Name: name,
		ComponentUnion: v1alpha2.ComponentUnion{
			Container: &v1alpha2.ContainerComponent{
				Container: v1alpha2.Container{
					Image:        image,
					Annotation:   &annotation,
					Env:          envs,
					VolumeMounts: volMounts,
					MountSources: &mountSources,
					DedicatedPod: &dedicatedPod,
				},
				Endpoints: endpoints,
			}}}
}

// generateDummyContainerComponentWithResourceRequirement returns a dummy container component with resource requirement for testing
func generateDummyContainerComponentWithResourceRequirement(name, memoryLimit, memoryRequest, cpuLimit, cpuRequest string) v1alpha2.Component {
	image := "docker.io/maven:latest"

	return v1alpha2.Component{
		Name: name,
		ComponentUnion: v1alpha2.ComponentUnion{
			Container: &v1alpha2.ContainerComponent{
				Container: v1alpha2.Container{
					Image:         image,
					MemoryLimit:   memoryLimit,
					MemoryRequest: memoryRequest,
					CpuLimit:      cpuLimit,
					CpuRequest:    cpuRequest,
				},
			}}}
}

// generateDummyVolumeComponent returns a dummy volume component for testing
func generateDummyVolumeComponent(name, size string) v1alpha2.Component {

	return v1alpha2.Component{

		Name: name,
		ComponentUnion: v1alpha2.ComponentUnion{
			Volume: &v1alpha2.VolumeComponent{
				Volume: v1alpha2.Volume{
					Size: size,
				},
			},
		},
	}
}

// generateDummyOpenshiftComponent returns a dummy Openshift component for testing
func generateDummyOpenshiftComponent(name string, endpoints []v1alpha2.Endpoint, uri string) v1alpha2.Component {

	return v1alpha2.Component{
		Name: name,
		ComponentUnion: v1alpha2.ComponentUnion{
			Openshift: &v1alpha2.OpenshiftComponent{
				K8sLikeComponent: v1alpha2.K8sLikeComponent{
					K8sLikeComponentLocation: v1alpha2.K8sLikeComponentLocation{
						Uri: uri,
					},
					Endpoints: endpoints,
				},
			},
		},
	}
}

// generateDummyKubernetesComponent returns a dummy Kubernetes component for testing
func generateDummyKubernetesComponent(name string, endpoints []v1alpha2.Endpoint, uri string) v1alpha2.Component {

	return v1alpha2.Component{
		Name: name,
		ComponentUnion: v1alpha2.ComponentUnion{
			Kubernetes: &v1alpha2.KubernetesComponent{
				K8sLikeComponent: v1alpha2.K8sLikeComponent{
					K8sLikeComponentLocation: v1alpha2.K8sLikeComponentLocation{
						Uri: uri,
					},
					Endpoints: endpoints,
				},
			},
		},
	}
}

// generateDummyImageComponent returns a dummy Image Dockerfile Component for testing
func generateDummyImageComponent(name string, src v1alpha2.DockerfileSrc) v1alpha2.Component {

	return v1alpha2.Component{
		Name: name,
		ComponentUnion: v1alpha2.ComponentUnion{
			Image: &v1alpha2.ImageComponent{
				Image: v1alpha2.Image{
					ImageName: "image:latest",
					ImageUnion: v1alpha2.ImageUnion{
						Dockerfile: &v1alpha2.DockerfileImage{
							DockerfileSrc: src,
							Dockerfile: v1alpha2.Dockerfile{
								BuildContext: "/path",
							},
						},
					},
				},
			},
		},
	}
}

// generateDummyPluginComponent returns a dummy Plugin component for testing
func generateDummyPluginComponent(name, url string, compAttribute attributes.Attributes) v1alpha2.Component {

	return v1alpha2.Component{
		Attributes: compAttribute,
		Name:       name,
		ComponentUnion: v1alpha2.ComponentUnion{
			Plugin: &v1alpha2.PluginComponent{
				ImportReference: v1alpha2.ImportReference{
					RegistryUrl: url,
				},
			},
		},
	}
}

func TestValidateComponents(t *testing.T) {

	volMounts := []v1alpha2.VolumeMount{
		{
			Name: "myvol",
			Path: "/some/path/",
		},
	}

	invalidVolMounts := []v1alpha2.VolumeMount{
		{
			Name: "myinvalidvol",
		},
		{
			Name: "myinvalidvol2",
		},
	}

	projectSourceEnv := []v1alpha2.EnvVar{
		{
			Name:  EnvProjectsSrc,
			Value: "/some/path/",
		},
	}

	projectsRootEnv := []v1alpha2.EnvVar{
		{
			Name:  EnvProjectsRoot,
			Value: "/some/path/",
		},
	}

	twoRemotesGitSrc := v1alpha2.DockerfileSrc{
		Git: &v1alpha2.DockerfileGitProjectSource{
			GitProjectSource: v1alpha2.GitProjectSource{
				GitLikeProjectSource: v1alpha2.GitLikeProjectSource{
					Remotes: map[string]string{
						"a": "abc",
						"x": "xyz",
					},
					CheckoutFrom: &v1alpha2.CheckoutFrom{
						Remote: "a",
					},
				},
			},
		},
	}

	zeroRemoteGitSrc := v1alpha2.DockerfileSrc{
		Git: &v1alpha2.DockerfileGitProjectSource{
			GitProjectSource: v1alpha2.GitProjectSource{
				GitLikeProjectSource: v1alpha2.GitLikeProjectSource{
					CheckoutFrom: &v1alpha2.CheckoutFrom{
						Remote: "a",
					},
				},
			},
		},
	}

	invalidRemoteGitSrc := v1alpha2.DockerfileSrc{
		Git: &v1alpha2.DockerfileGitProjectSource{
			GitProjectSource: v1alpha2.GitProjectSource{
				GitLikeProjectSource: v1alpha2.GitLikeProjectSource{
					Remotes: map[string]string{
						"a": "abc",
					},
					CheckoutFrom: &v1alpha2.CheckoutFrom{
						Remote: "b",
					},
				},
			},
		},
	}

	validRemoteGitSrc := v1alpha2.DockerfileSrc{
		Git: &v1alpha2.DockerfileGitProjectSource{
			GitProjectSource: v1alpha2.GitProjectSource{
				GitLikeProjectSource: v1alpha2.GitLikeProjectSource{
					Remotes: map[string]string{
						"a": "abc",
					},
					CheckoutFrom: &v1alpha2.CheckoutFrom{
						Remote: "a",
					},
				},
			},
		},
	}

	validUriSrc := v1alpha2.DockerfileSrc{
		Uri: "uri",
	}

	endpointUrl18080 := generateDummyEndpoint("url1", 8080)
	endpointUrl18081 := generateDummyEndpoint("url1", 8081)
	endpointUrl28080 := generateDummyEndpoint("url2", 8080)
	endpointUrl28081 := generateDummyEndpoint("url2", 8081)

	invalidVolMountErr := ".*\nvolume mount myinvalidvol belonging to the container component.*\nvolume mount myinvalidvol2 belonging to the container component.*"
	duplicateComponentErr := "duplicate key: component1"
	reservedEnvErr := "env variable .* is reserved and cannot be customized in component.*"
	invalidSizeErr := "size .* for volume component is invalid"
	sameEndpointNameErr := "devfile contains multiple endpoint entries with same name.*"
	sameTargetPortErr := "devfile contains multiple containers with same endpoint targetPort.*"
	invalidURIErr := ".*invalid URI for request"
	imageCompTwoRemoteErr := "component .* should have one remote only"
	imageCompNoRemoteErr := "component .* should have at least one remote"
	imageCompInvalidRemoteErr := "unable to find the checkout remote .* in the remotes for component .*"
	DeploymentAnnotationConflictErr := "deployment annotation: deploy-key1 has been declared multiple times and with different values"
	ServiceAnnotationConflictErr := "service annotation: svc-key1 has been declared multiple times and with different values"

	pluginOverridesFromMainDevfile := attributes.Attributes{}.PutString(ImportSourceAttribute,
		"uri: http://127.0.0.1:8080").PutString(PluginOverrideAttribute, "main devfile")
	invalidURIErrWithImportAttributes := ".*invalid URI for request, imported from uri: http://127.0.0.1:8080, in plugin overrides from main devfile"
	invalidCpuRequest := ".*cpuRequest is greater than cpuLimit."
	invalidMemoryRequest := ".*memoryRequest is greater than memoryLimit."
	quantityParsingErr := "error parsing .* requirement for component.*"

	tests := []struct {
		name       string
		components []v1alpha2.Component
		wantErr    []string
	}{
		{
			name: "Duplicate components present",
			components: []v1alpha2.Component{
				generateDummyVolumeComponent("component1", "1Gi"),
				generateDummyContainerComponent("component1", nil, nil, nil, v1alpha2.Annotation{}, false),
			},
			wantErr: []string{duplicateComponentErr},
		},
		{
			name: "Valid container and volume component",
			components: []v1alpha2.Component{
				generateDummyVolumeComponent("myvol", "1Gi"),
				generateDummyContainerComponent("container", volMounts, nil, nil, v1alpha2.Annotation{}, false),
				generateDummyContainerComponent("container2", volMounts, nil, nil, v1alpha2.Annotation{}, false),
			},
		},
		{
			name: "Invalid container using reserved env PROJECT_SOURCE",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("container1", nil, nil, projectSourceEnv, v1alpha2.Annotation{}, false),
			},
			wantErr: []string{reservedEnvErr},
		},
		{
			name: "Invalid container using reserved env PROJECTS_ROOT",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("container", nil, nil, projectsRootEnv, v1alpha2.Annotation{}, false),
			},
			wantErr: []string{reservedEnvErr},
		},
		{
			name: "Invalid volume component size",
			components: []v1alpha2.Component{
				generateDummyVolumeComponent("myvol", "invalid"),
				generateDummyContainerComponent("container", nil, nil, nil, v1alpha2.Annotation{}, false),
			},
			wantErr: []string{invalidSizeErr},
		},
		{
			name: "Invalid volume mount referencing a wrong volume component",
			components: []v1alpha2.Component{
				generateDummyVolumeComponent("myvol", "1Gi"),
				generateDummyContainerComponent("container1", invalidVolMounts, nil, nil, v1alpha2.Annotation{}, false),
			},
			wantErr: []string{invalidVolMountErr},
		},
		{
			name: "Invalid containers with the same endpoint names",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("name1", nil, []v1alpha2.Endpoint{endpointUrl18080}, nil, v1alpha2.Annotation{}, false),
				generateDummyContainerComponent("name2", nil, []v1alpha2.Endpoint{endpointUrl18081}, nil, v1alpha2.Annotation{}, false),
			},
			wantErr: []string{sameEndpointNameErr},
		},
		{
			name: "Invalid containers with the same endpoint target ports",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("name1", nil, []v1alpha2.Endpoint{endpointUrl18080}, nil, v1alpha2.Annotation{}, false),
				generateDummyContainerComponent("name2", nil, []v1alpha2.Endpoint{endpointUrl28080}, nil, v1alpha2.Annotation{}, false),
			},
			wantErr: []string{sameTargetPortErr},
		},
		{
			name: "Valid container with same target ports in a single component",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("name1", nil, []v1alpha2.Endpoint{endpointUrl18080, endpointUrl28080}, nil, v1alpha2.Annotation{}, false),
			},
		},
		{
			name: "Invalid Kube components with the same endpoint names",
			components: []v1alpha2.Component{
				generateDummyKubernetesComponent("name1", []v1alpha2.Endpoint{endpointUrl18080}, ""),
				generateDummyKubernetesComponent("name2", []v1alpha2.Endpoint{endpointUrl18081}, ""),
			},
			wantErr: []string{sameEndpointNameErr},
		},
		{
			name: "Valid Kube component with the same endpoint target ports as the container component's",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("name1", nil, []v1alpha2.Endpoint{endpointUrl18080}, nil, v1alpha2.Annotation{}, false),
				generateDummyKubernetesComponent("name2", []v1alpha2.Endpoint{endpointUrl28080}, ""),
			},
		},
		{
			name: "Invalid Kube components with the same endpoint names",
			components: []v1alpha2.Component{
				generateDummyKubernetesComponent("name1", []v1alpha2.Endpoint{endpointUrl18080}, ""),
				generateDummyKubernetesComponent("name2", []v1alpha2.Endpoint{endpointUrl28080}, ""),
			},
		},
		{
			name: "Valid containers with valid resource requirement",
			components: []v1alpha2.Component{
				generateDummyContainerComponentWithResourceRequirement("name1", "1024Mi", "512Mi", "1024Mi", "512Mi"),
				generateDummyContainerComponentWithResourceRequirement("name2", "", "512Mi", "", "512Mi"),
				generateDummyContainerComponentWithResourceRequirement("name3", "1024Mi", "", "1024Mi", ""),
				generateDummyContainerComponentWithResourceRequirement("name4", "", "", "", ""),
			},
		},
		{
			name: "Invalid containers with resource limit smaller than resource requested",
			components: []v1alpha2.Component{
				generateDummyContainerComponentWithResourceRequirement("name1", "512Mi", "1024Mi", "", ""),
				generateDummyContainerComponentWithResourceRequirement("name2", "", "", "512Mi", "1024Mi"),
			},
			wantErr: []string{invalidMemoryRequest, invalidCpuRequest},
		},
		{
			name: "Invalid container with resource quantity parsing error",
			components: []v1alpha2.Component{
				generateDummyContainerComponentWithResourceRequirement("name1", "512invalid", "", "", ""),
			},
			wantErr: []string{quantityParsingErr},
		},
		{
			name: "Valid container with deployment, service and ingress annotations",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("name1", nil, nil, nil, v1alpha2.Annotation{
					Deployment: map[string]string{
						"deploy-key1": "deploy-value1",
						"deploy-key2": "deploy-value2",
					},
					Service: map[string]string{
						"svc-key1": "svc-value1",
						"svc-key2": "svc-value2",
					},
				}, false),
			},
		},
		{
			name: "Valid containers with different key and value pairs for deployment, service and ingress annotations",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("name1", nil, nil, nil, v1alpha2.Annotation{
					Deployment: map[string]string{
						"deploy-key1": "deploy-value1",
					},
					Service: map[string]string{
						"svc-key1": "svc-value1",
					},
				}, false),
				generateDummyContainerComponent("name2", nil, nil, nil, v1alpha2.Annotation{
					Deployment: map[string]string{
						"deploy-key2": "deploy-value2",
					},
					Service: map[string]string{
						"svc-key2": "svc-value2",
					},
				}, false),
			},
		},
		{
			name: "Valid containers with same key and value pairs for deployment, service and ingress annotations",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("name1", nil, nil, nil, v1alpha2.Annotation{
					Deployment: map[string]string{
						"deploy-key1": "deploy-value1",
					},
					Service: map[string]string{
						"svc-key1": "svc-value1",
					},
				}, false),
				generateDummyContainerComponent("name2", nil, nil, nil, v1alpha2.Annotation{
					Deployment: map[string]string{
						"deploy-key1": "deploy-value1",
					},
					Service: map[string]string{
						"svc-key1": "svc-value1",
					},
				}, false),
			},
		},
		{
			name: "Valid containers with conflict key and value pairs for deployment and service annotations when dedicatedPod is set to true",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("name1", nil, nil, nil, v1alpha2.Annotation{
					Deployment: map[string]string{
						"deploy-key1": "deploy-value1",
					},
					Service: map[string]string{
						"svc-key1": "svc-value1",
					},
				}, false),
				generateDummyContainerComponent("name2", nil, nil, nil, v1alpha2.Annotation{
					Deployment: map[string]string{
						"deploy-key1": "deploy-value2",
					},
					Service: map[string]string{
						"svc-key1": "svc-value2",
					},
				}, true),
			},
		},
		{
			name: "Invalid containers with conflict key and value pairs for deployment and service annotations when dedicatedPod is set to false",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("name1", nil, nil, nil, v1alpha2.Annotation{
					Deployment: map[string]string{
						"deploy-key1": "deploy-value1",
					},
					Service: map[string]string{
						"svc-key1": "svc-value1",
					},
				}, false),
				generateDummyContainerComponent("name2", nil, nil, nil, v1alpha2.Annotation{
					Deployment: map[string]string{
						"deploy-key1": "deploy-value2",
					},
					Service: map[string]string{
						"svc-key1": "svc-value2",
					},
				}, false),
				generateDummyContainerComponent("name3", nil, nil, nil, v1alpha2.Annotation{
					Deployment: map[string]string{
						"deploy-key1": "deploy-value2",
					},
					Service: map[string]string{
						"svc-key1": "svc-value2",
					},
				}, false),
			},
			wantErr: []string{DeploymentAnnotationConflictErr, ServiceAnnotationConflictErr},
		},
		{
			name: "Invalid Openshift Component with bad URI",
			components: []v1alpha2.Component{
				generateDummyOpenshiftComponent("name1", []v1alpha2.Endpoint{endpointUrl18080, endpointUrl28081}, "http//wronguri"),
			},
			wantErr: []string{invalidURIErr},
		},
		{
			name: "Valid Kubernetes Component",
			components: []v1alpha2.Component{
				generateDummyKubernetesComponent("name1", []v1alpha2.Endpoint{endpointUrl18080, endpointUrl28081}, "http://uri"),
			},
		},
		{
			name: "Invalid OpenShift Component with same endpoint names",
			components: []v1alpha2.Component{
				generateDummyOpenshiftComponent("name1", []v1alpha2.Endpoint{endpointUrl18080, endpointUrl18081}, "http://uri"),
			},
			wantErr: []string{sameEndpointNameErr},
		},
		{
			name: "Multiple errors: Duplicate component name, invalid plugin registry url, bad URI with import source attributes",
			components: []v1alpha2.Component{
				generateDummyVolumeComponent("component1", "1Gi"),
				generateDummyPluginComponent("component1", "http//invalidregistryurl", attributes.Attributes{}),
				generateDummyPluginComponent("abc", "http//invalidregistryurl", pluginOverridesFromMainDevfile),
			},
			wantErr: []string{duplicateComponentErr, invalidURIErr, invalidURIErrWithImportAttributes},
		},
		{
			name: "Invalid image dockerfile component with more than one remote",
			components: []v1alpha2.Component{
				generateDummyImageComponent("name1", twoRemotesGitSrc),
			},
			wantErr: []string{imageCompTwoRemoteErr},
		},
		{
			name: "Invalid image dockerfile component with zero remote",
			components: []v1alpha2.Component{
				generateDummyImageComponent("name1", zeroRemoteGitSrc),
			},
			wantErr: []string{imageCompNoRemoteErr},
		},
		{
			name: "Invalid image dockerfile component with wrong checkout",
			components: []v1alpha2.Component{
				generateDummyImageComponent("name1", invalidRemoteGitSrc),
			},
			wantErr: []string{imageCompInvalidRemoteErr},
		},
		{
			name: "Valid image dockerfile component with correct remote",
			components: []v1alpha2.Component{
				generateDummyImageComponent("name1", validRemoteGitSrc),
			},
		},
		{
			name: "Valid image dockerfile component with non git src",
			components: []v1alpha2.Component{
				generateDummyImageComponent("name1", validUriSrc),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateComponents(tt.components)

			merr, ok := err.(*multierror.Error)
			if ok {
				if tt.wantErr != nil {
					if assert.Equal(t, len(tt.wantErr), len(merr.Errors), "Error list length should match") {
						for i := 0; i < len(merr.Errors); i++ {
							assert.Regexp(t, tt.wantErr[i], merr.Errors[i].Error(), "Error message should match")
						}
					}
				} else {
					t.Errorf("Error should be nil, got %v", err)
				}
			} else if tt.wantErr != nil {
				t.Errorf("Error should not be nil, want %v, got %v", tt.wantErr, err)
			}
		})
	}

}
