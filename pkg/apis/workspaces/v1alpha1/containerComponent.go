package v1alpha1

// Component that allows the developer to add a configured container into his workspace
type ContainerComponent struct {
	BaseComponent  `json:",inline"`
	Container  `json:",inline"`
	MemoryLimit string     `json:"memoryLimit,omitempty"`
	Endpoints   []Endpoint `json:"endpoints,omitempty"`
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
	// +optional
	// +kubebuilder:default="tcp"
	Protocol string `json:"protocol,omitmepty"`
	// +optional
	// +kubebuilder:default="http"
	Scheme string `json:"scheme,omitmepty"`
	// +optional
	Secure bool `json:"secure"`
	// +optional
	CookiesAuthEnabled bool `json:"cookiesAuthEnabled"`
	// +optional
	Path string `json:"path,omitempty"`

	// +kubebuilder:validation:Enum=ide;terminal;ide-dev
	Type string `json:"type,omitmepty"`
}

type Container struct {
	Name  string   `json:"name"`
	Image string   `json:"image"`
	// +optional
	Env   []EnvVar `json:"env,omitempty"`
	// +optional
	Volumes []Volume `json:"volumes,omitempty"`
	//+optional
	MemoryLimit  string `json:"memoryLimit,omitempty"`
	//+optional
	MountSources bool   `json:"mountSources"`
}

type EnvVar struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

// Volume that should be mounted to a component container
type Volume struct {
	// The volume name.
	// If several components mount the same volume then they will reuse the volume
	// and will be able to access to the same files
	Name string `json:"name"`

	MountPath string `json:"mountPath"`
}
