package validation

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
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

	endpoint1 := generateDummyEndpoint("url1", 8080)
	endpoint2 := generateDummyEndpoint("url1", 8081)
	endpoint3 := generateDummyEndpoint("url2", 8080)

	tests := []struct {
		name        string
		components  []v1alpha2.Component
		wantErr     bool
		wantErrType error
	}{
		{
			name: "Case 1: Duplicate components present",
			components: []v1alpha2.Component{
				generateDummyVolumeComponent("component1", "1Gi"),
				generateDummyContainerComponent("component1", volMounts, nil, nil),
			},
			wantErr: true,
		},
		{
			name: "Case 2: Valid container and volume component",
			components: []v1alpha2.Component{
				generateDummyVolumeComponent("myvol", "1Gi"),
				generateDummyContainerComponent("container", volMounts, nil, nil),
				generateDummyContainerComponent("container2", volMounts, nil, nil),
			},
			wantErr: false,
		},
		{
			name: "Case 3: Invalid container using reserved env PROJECT_SOURCE",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("container1", nil, nil, projectSourceEnv),
			},
			wantErr: true,
		},
		{
			name: "Case 4: Invalid container using reserved env PROJECTS_ROOT",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("container", nil, nil, projectsRootEnv),
			},
			wantErr: true,
		},
		{
			name: "Case 5: Invalid volume component size",
			components: []v1alpha2.Component{
				generateDummyVolumeComponent("myvol", "invalid"),
				generateDummyContainerComponent("container", nil, nil, nil),
			},
			wantErr: true,
		},
		{
			name: "Case 6: Invalid volume mount",
			components: []v1alpha2.Component{
				generateDummyVolumeComponent("myvol", "1Gi"),
				generateDummyContainerComponent("container1", invalidVolMounts, nil, nil),
				generateDummyContainerComponent("container2", invalidVolMounts, nil, nil),
			},
			wantErr: true,
		},
		{
			name: "Case 7: Invalid container with same endpoint names",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("name1", nil, []v1alpha2.Endpoint{endpoint1}, nil),
				generateDummyContainerComponent("name2", nil, []v1alpha2.Endpoint{endpoint2}, nil),
			},
			wantErr: true,
		},
		{
			name: "Case 8: Invalid container with same endpoint target ports",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("name1", nil, []v1alpha2.Endpoint{endpoint1}, nil),
				generateDummyContainerComponent("name2", nil, []v1alpha2.Endpoint{endpoint3}, nil),
			},
			wantErr: true,
		},
		{
			name: "Case 9: Valid container with multiple same target ports but different endpoint name",
			components: []v1alpha2.Component{
				generateDummyContainerComponent("name1", nil, []v1alpha2.Endpoint{endpoint1, endpoint3}, nil),
			},
			wantErr: false,
		},
		{
			name: "Case 10: Invalid Openshift Component with bad URI",
			components: []v1alpha2.Component{
				generateDummyOpenshiftComponent("name1", []v1alpha2.Endpoint{endpoint1, endpoint3}, "http//wronguri"),
			},
			wantErr: true,
		},
		{
			name: "Case 11: Valid Kubernetes Component",
			components: []v1alpha2.Component{
				generateDummyKubernetesComponent("name1", []v1alpha2.Endpoint{endpoint1, endpoint3}, "http://uri"),
			},
			wantErr: false,
		},
		{
			name: "Case 12: Invalid OpenShift Component endpoints",
			components: []v1alpha2.Component{
				generateDummyOpenshiftComponent("name1", []v1alpha2.Endpoint{endpoint1, endpoint2}, "http://uri"),
			},
			wantErr: true,
		},
		{
			name: "Case 13: Invalid component name with all numeric values",
			components: []v1alpha2.Component{
				generateDummyVolumeComponent("123", "1Gi"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateComponents(tt.components)

			if tt.wantErr && got == nil {
				t.Errorf("TestValidateComponents error - expected an err but got nil")
			} else if !tt.wantErr && got != nil {
				t.Errorf("TestValidateComponents error - unexpected err %v", got)
			}
		})
	}

}
