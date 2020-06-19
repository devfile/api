package v1alpha1

// Component that allows the developer to add a configured container into his workspace
type ContainerComponent struct {
	BaseComponent `json:",inline"`
	Container     `json:",inline"`
	MemoryLimit   string     `json:"memoryLimit,omitempty"`
	Endpoints     []Endpoint `json:"endpoints,omitempty"`
}

// EndpointExposure describes the way an endpoint is exposed on the network.
// Only one of the following exposures may be specified: public, internal, none.
// +kubebuilder:validation:Enum=public;internal;none
type EndpointExposure string

const (
	// Endpoint will be exposed on the public network, typically through
	// a K8S ingress or an OpenShift route 
	PublicEndpointExposure  EndpointExposure = "public"
	// Endpoint will be exposed internally outside of the main workspace POD,
	// typically by K8S services, to be consumed by other elements running
	// on the same cloud internal network.
	InternalEndpointExposure EndpointExposure = "internal"
	// Endpoint will not be exposed and will only be accessible
	// inside the main workspace POD, on a local address.
	NoneEndpointExposure EndpointExposure = "none"
)

type Endpoint struct {
	Name string `json:"name"`

	// +optional
	TargetPort int `json:"targetPort,omitempty"`

	// Describes how the endpoint should be exposed on the network.
	//
	//
	// - `public` means that the endpoint will be exposed on the public network, typically through
	// a K8S ingress or an OpenShift route.
	//
	//
	// - `internal` means that the endpoint will be exposed internally outside of the main workspace POD,
	// typically by K8S services, to be consumed by other elements running
	// on the same cloud internal network.
	//
	//
	// - `none` means that the endpoint will not be exposed and will only be accessible
	// inside the main workspace POD, on a local address.
	// +optional
	Exposure EndpointExposure `json:"exposure,omitempty"`

	// Low-level protocol of traffic coming through this endpoint.
	// Default value is "tcp"
	// +optional
	Protocol string `json:"protocol,omitmepty"`

	// Describes whether the endpoint should be secured and protected by some
	// authentication process
	// +optional
	Secure bool `json:"secure"`

	// Map of implementation-dependant string-based free-form attributes.
	// Examples of Che-specific attributes:
	// cookiesAuthEnabled ("true"/"false"), scheme ("http", "ws"), type ("terminal", "ide")
	// +optional
	Attributes map[string]string `json:"attributes,omitempty"`
}

type Container struct {
	Name string `json:"name"`

	Image string `json:"image,omitempty"`

	// +optional
	// Environment variables used in this container
	Env []EnvVar `json:"env,omitempty"`

	// +optional
	// List of volumes mounts that should be mounted is this container.
	VolumeMounts []VolumeMount `json:"volumeMounts,omitempty"`

	//+optional
	MemoryLimit string `json:"memoryLimit,omitempty"`

	// The command to run in the dockerimage component instead of the default one provided in the image.
	// Defaults to an empty array, meaning use whatever is defined in the image.
	//+optional
	Command []string `json:"command,omitempty"`

	// The arguments to supply to the command running the dockerimage component. The arguments are supplied either to the default command provided in the image or to the overridden command.
	// Defaults to an empty array, meaning use whatever is defined in the image.
	//+optional
	Args []string `json:"args,omitempty"`

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

	// The path in the component container where the volume should be mounted.
	// If not path is mentioned, default path is the is `/<name>`.
	// +optional
	Path string `json:"path,omitempty"`
}
