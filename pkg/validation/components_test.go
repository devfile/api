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
			name: "Case 1: Duplicate volume components present",
			components: []v1alpha2.Component{
				generateDummyVolumeComponent("myvol", "1Gi"),
				generateDummyVolumeComponent("myvol", "1Gi"),
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
				generateDummyContainerComponent("container", invalidVolMounts, nil, nil),
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
