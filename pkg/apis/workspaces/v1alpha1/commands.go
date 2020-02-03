package v1alpha1

import runtime "k8s.io/apimachinery/pkg/runtime"

// CommandType describes the type of command.
// Only one of the following command type may be specified.
// +kubebuilder:validation:Enum=Exec;VscodeTask;VscodeLaunch;Custom
type CommandType string

const (
	ExecCommandType         CommandType = "Exec"
	VscodeTaskCommandType   CommandType = "VscodeTask"
	VscodeLaunchCommandType CommandType = "VscodeLaunch"
	CompositeCommandType    CommandType = "Composite"
	CustomCommandType       CommandType = "Custom"
)

type BaseCommand struct {
	Alias      string            `json:"alias,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"` // Additional command attributes
}

type LabeledCommand struct {
	BaseCommand `json:",inline"`
	Label       string `json:"label,omitempty"`
}

type Command struct {
	PolymorphicCommand `json:",inline"`
}

// +k8s:openapi-gen=true
// +union
type PolymorphicCommand struct {
	// Type of workspace command
	// +unionDiscriminator
	// +optional
	Type CommandType `json:"type"`

	// Exec command
	// +optional
	Exec *ExecCommand `json:"exec,omitempty"`

	// VscodeTask command
	// +optional
	VscodeTask *VscodeConfigurationCommand `json:"vscodeTask,omitempty"`

	// VscodeLaunch command
	// +optional
	VscodeLaunch *VscodeConfigurationCommand `json:"vscodeLaunch,omitempty"`

	// Composite command
	// +optional
	Composite *CompositeCommand `json:"composite,omitempty"`

	// Custom command
	// +optional
	Custom *CustomCommand `json:"custom,omitempty"`
}

type ExecCommand struct {
	LabeledCommand `json:",inline"`

	// The actual command-line string
	CommandLine string `json:"commandLine"`

	// Describes component to which given action relates
	Component string `json:"component,omitempty"`

	// Working directory where the command should be executed
	Workdir *string `json:"workdir,omitempty"`
}

type CompositeCommand struct {
	LabeledCommand `json:",inline"`

	// The commands that comprise this composite command
	Commands []string `json:"commands,omitempty"`

	// +optional
	Parallel bool `json:"parallel,omitempty"`
}

// +k8s:openapi-gen=true
// +union
type VscodeConfigurationCommandLocation struct {
	// Type of Vscode configuration command location
	// +
	// +unionDiscriminator
	// +optional
	LocationType string `json:"locationType"`

	// Location as an absolute of relative URL
	// +optional
	Url string `json:"url,omitempty"`

	// Embedded content of the vscode configuration file
	// +optional
	Inlined string `json:"inlined,omitempty"`
}

type VscodeConfigurationCommand struct {
	BaseCommand `json:",inline"`
	Location    VscodeConfigurationCommandLocation `json:",inline"`
}

type CustomCommand struct {
	LabeledCommand `json:",inline"`
	CommandClass   string `json:"commandClass"`

	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:EmbeddedResource
	EmbeddedResource runtime.RawExtension `json:"embeddedResource"`
}
