package v1alpha1

// Component that allows the developer to add a configured container into his workspace
type ContainerComponent struct {
	BaseComponent `json:",inline"`
	Container     `json:",inline"`
	MemoryLimit   string     `json:"memoryLimit,omitempty"`
	Endpoints     []Endpoint `json:"endpoints,omitempty"`
}

type Endpoint struct {
	Name          string                 `json:"name"`
	TargetPort    int                    `json:"targetPort"`
	Configuration *EndpointConfiguration `json:"configuration,omitEmpty"`
	Attributes    map[string]string      `json:"attributes,omitempty"`
}

type EndpointConfiguration struct {
	// +optional
	Public bool `json:"public"`
	// +optional
	Discoverable bool `json:"discoverable"`
	// The is the low-level protocol of traffic coming through this endpoint.
	// Default value is "tcp"
	// +optional
	Protocol string `json:"protocol,omitmepty"`
	// The is the URL scheme to use when accessing the endpoint.
	// Default value is "http"
	// +optional
	Scheme string `json:"scheme,omitmepty"`
	// +optional
	Secure bool `json:"secure"`
	// +optional
	CookiesAuthEnabled bool `json:"cookiesAuthEnabled"`
	// +optional
	Path string `json:"path,omitempty"`

	// +kubebuilder:validation:Enum=ide;terminal;ide-dev
	// +optional
	Type string `json:"type,omitmepty"`
}

type Container struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	// +optional
	// Environment variables used in this container
	Env []EnvVar `json:"env,omitempty"`

	// +optional
	// List of volumes mounts that should be mounted is this container.
	VolumeMounts []VolumeMount `json:"volumeMounts,omitempty"`

	//+optional
	MemoryLimit string `json:"memoryLimit,omitempty"`

	//+optional
	MountSources bool `json:"mountSources"`

	//+optional
	//
	// Optional specification of the path in the container where
	// project sources should be transferred/mounted when `mountSources` is `true`.
	// When omitted, the value of the `PROJECTS_ROOT` environment variable is used.
	SourceMapping string `json:"sourceMapping"`
}

type EnvVar struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

// Volume that should be mounted to a component container
type VolumeMount struct {
	// The volume mount name is the name of an existing `Volume` component.
	// If no corresponding `Volume` component exist it is implicitly added.
	// If several containers mount the same volume name
	// then they will reuse the same volume and will be able to access to the same files.
	Name string `json:"name"`

	// The path in the component container where the volume should be mounted
	Path string `json:"path"`
}
