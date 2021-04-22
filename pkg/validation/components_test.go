package validation

import (
	"github.com/devfile/api/v2/pkg/attributes"
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/stretchr/testify/assert"
)

// generateDummyContainerComponent returns a dummy container component for testing
func generateDummyContainerComponent(name string, volMounts []v1alpha2.VolumeMount, endpoints []v1alpha2.Endpoint, envs []v1alpha2.EnvVar) v1alpha2.Component {
	image := "docker.io/maven:latest"
	mountSources := true

	return v1alpha2.Component{
		Name: name,
		ComponentUnion: v1alpha2.ComponentUnion{
			Container: &v1alpha2.ContainerComponent{
				Container: v1alpha2.Container{
					Image:        image,
					Env:          envs,
					VolumeMounts: volMounts,
					MountSources: &mountSources,
				},
				Endpoints: endpoints,
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

// generateDummyPluginComponent returns a dummy Plugin component for testing
func generateDummyPluginComponent(name, url string) v1alpha2.Component {

	return v1alpha2.Component{
		Name: name,
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

	endpointUrl18080 := generateDummyEndpoint("url1", 8080)
	endpointUrl18081 := generateDummyEndpoint("url1", 8081)
	endpointUrl28080 := generateDummyEndpoint("url2", 8080)

	invalidVolMountErr := ".*\nvolume mount myinvalidvol belonging to the container component.*\nvolume mount myinvalidvol2 belonging to the container component.*"
	duplicateComponentErr := "duplicate key: component1"
	reservedEnvErr := "env variable .* is reserved and cannot be customized in component.*"
	invalidSizeErr := "size .* for volume component is invalid"
	sameEndpointNameErr := "devfile contains multiple endpoint entries with same name.*"
	sameTargetPortErr := "devfile contains multiple containers with same TargetPort.*"
	invalidURIErr := ".*invalid URI for request"

	pluginOverridesFromMainDevfile := attributes.Attributes{}.PutString(ImportSourceAttribute,
		"uri: http://127.0.0.1:8080").PutString(PluginOverrideAttribute, "main devfile")
	invalidURIErrWithImportAttributes := ".*invalid URI for request, imported from uri: http://127.0.0.1:8080, in plugin overrides from main devfile"

	tests := []struct {
		name       string
		components []v1alpha2.Component
		wantErr    *string
	}{
		{
			name: "Duplicate components present",
			components: []v1alpha2.Component{
				generateDummyVolumeComponent("component1", "1Gi"),
				generateDummyContainerComponent("component1", volMounts, nil, nil),
			},
			wantErr: &duplicateComponentErr,
		},
		{
			name: "Valid container and volume component",
			components: []v1alpha2.Component{
				generateDummyVolumeComponent("myvol", "1Gi"),
				generateDummyContainerComponent("container", volMounts, nil, nil),
				generateDummyContainerComponent("container2", volMounts, nil, nil),
			},
		},
		{
			name: "Invalid container using reserved env PROJECT_SOURCE",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("container1", nil, nil, projectSourceEnv),
			},
			wantErr: &reservedEnvErr,
		},
		{
			name: "Invalid container using reserved env PROJECTS_ROOT",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("container", nil, nil, projectsRootEnv),
			},
			wantErr: &reservedEnvErr,
		},
		{
			name: "Invalid volume component size",
			components: []v1alpha2.Component{
				generateDummyVolumeComponent("myvol", "invalid"),
				generateDummyContainerComponent("container", nil, nil, nil),
			},
			wantErr: &invalidSizeErr,
		},
		{
			name: "Invalid volume mount referencing a wrong volume component",
			components: []v1alpha2.Component{
				generateDummyVolumeComponent("myvol", "1Gi"),
				generateDummyContainerComponent("container1", invalidVolMounts, nil, nil),
			},
			wantErr: &invalidVolMountErr,
		},
		{
			name: "Invalid containers with the same endpoint names",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("name1", nil, []v1alpha2.Endpoint{endpointUrl18080}, nil),
				generateDummyContainerComponent("name2", nil, []v1alpha2.Endpoint{endpointUrl18081}, nil),
			},
			wantErr: &sameEndpointNameErr,
		},
		{
			name: "Invalid containers with the same endpoint target ports",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("name1", nil, []v1alpha2.Endpoint{endpointUrl18080}, nil),
				generateDummyContainerComponent("name2", nil, []v1alpha2.Endpoint{endpointUrl28080}, nil),
			},
			wantErr: &sameTargetPortErr,
		},
		{
			name: "Valid container with same target ports but different endpoint name",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("name1", nil, []v1alpha2.Endpoint{endpointUrl18080, endpointUrl28080}, nil),
			},
		},
		{
			name: "Invalid Openshift Component with bad URI",
			components: []v1alpha2.Component{
				generateDummyOpenshiftComponent("name1", []v1alpha2.Endpoint{endpointUrl18080, endpointUrl28080}, "http//wronguri"),
			},
			wantErr: &invalidURIErr,
		},
		{
			name: "Valid Kubernetes Component",
			components: []v1alpha2.Component{
				generateDummyKubernetesComponent("name1", []v1alpha2.Endpoint{endpointUrl18080, endpointUrl28080}, "http://uri"),
			},
		},
		{
			name: "Invalid OpenShift Component with same endpoint names",
			components: []v1alpha2.Component{
				generateDummyOpenshiftComponent("name1", []v1alpha2.Endpoint{endpointUrl18080, endpointUrl18081}, "http://uri"),
			},
			wantErr: &sameEndpointNameErr,
		},
		{
			name: "Invalid plugin registry url",
			components: []v1alpha2.Component{
				generateDummyPluginComponent("abc", "http//invalidregistryurl"),
			},
			wantErr: &invalidURIErr,
		},
		{
			name: "Invalid component due to bad URI with import source attributes",
			components: []v1alpha2.Component{
				{
					Attributes: pluginOverridesFromMainDevfile,
					Name: "name",
					ComponentUnion: v1alpha2.ComponentUnion{
						Plugin: &v1alpha2.PluginComponent{
							ImportReference: v1alpha2.ImportReference{
								RegistryUrl: "http//invalidregistryurl",
							},
						},
					},
				},
			},
			wantErr: &invalidURIErrWithImportAttributes,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateComponents(tt.components)

			if tt.wantErr != nil && assert.Error(t, err) {
				assert.Regexp(t, *tt.wantErr, err.Error(), "Error message should match")
			} else {
				assert.NoError(t, err, "Expected error to be nil")
			}
		})
	}

}
