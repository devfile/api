package v1alpha1

type Command struct {
	Actions    []CommandAction   `json:"actions,omitempty"`    // List of the actions of given command. Now the only one command must be specified in list; but there are plans to implement supporting multiple actions commands.
	Attributes map[string]string `json:"attributes,omitempty"` // Additional command attributes
	Name       string            `json:"name"`                 // Describes the name of the command. Should be unique per commands set.
}

type CommandAction struct {
	Command          *string `json:"command,omitempty"`          // The actual action command-line string
	Component        *string `json:"component,omitempty"`        // Describes component to which given action relates
	Type             string  `json:"type"`                       // Describes action type
	Workdir          *string `json:"workdir,omitempty"`          // Working directory where the command should be executed
	Reference        *string `json:"reference,omitempty"`        // Working directory where the command should be executed
	ReferenceContent *string `json:"referenceContent,omitempty"` // Working directory where the command should be executed
}

type Project struct {
	Name   string        `json:"name"`
	Source ProjectSource `json:"source"` // Describes the project's source - type and location
}

// Describes the project's source - type and location
type ProjectSource struct {
	Location string `json:"location"` // Project's source location address. Should be URL for git and github located projects, or; file:// for zip.
	Type     string `json:"type"`     // Project's source type.
}

// Workspace component: Anything that will bring additional features / tooling / behaviour / context
// to the workspace, in order to make working in it easier.
type Component struct {
	Name string `json:"name"`
}

// Component that allows the developer to add a configured container into his workspace
type DeveloperRuntime struct {
	Component
	MemoryLimit string `json:"memoryLimit,omitempty"`
	Endpoints []Endpoint `json:"endpoints,omitempty"`
	Container
}

type Endpoint struct {
	Name          string                 `json:"name"`
	TargetPort    int                    `json:"targetPort"`
	Configuration *EndpointConfiguration `json:"configuration,omitEmpty"`
	attributes    map[string]string      `json:"attributes,omitempty"`
}

type EndpointConfiguration struct {
	Public             bool   `json:"public"`
	Discoverable       bool   `json:"discoverable"`
	Protocol           string `json:"protocol,omitmepty"`
	Schema             string `json:"schema,omitmepty"`
	Secure             bool   `json:"secure"`
	CookiesAuthEnabled bool   `json:"public"`
	Path               string `json:"path",omitempty`

	// +kubebuilder:validation:Enum=ide,terminal
	Type string `json:"type,omitmepty"`
}


type Container struct {
	Name           string          `json:"name" yaml:"name"`
	Image          string          `json:"image" yaml:"image"`
//	Env            []EnvVar        `json:"env" yaml:"env"`
//	EditorCommands []EditorCommand `json:"editorCommands" yaml:"editorCommands"`
//	Volumes        []Volume        `json:"volumes" yaml:"volumes"`
//	Ports          []ExposedPort   `json:"ports" yaml:"ports"`
	MemoryLimit    string          `json:"memoryLimit" yaml:"memoryLimit"`
	MountSources   bool            `json:"mountSources" yaml:"mountSources"`
}

