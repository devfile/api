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

// CommandGroupType describes the kind of command group.
// +kubebuilder:validation:Enum=build;run;test;debug
type CommandGroupType string

const (
	BuildCommandGroupType CommandGroupType = "build"
	RunCommandGroupType   CommandGroupType = "run"
	TestCommandGroupType  CommandGroupType = "test"
	DebugCommandGroupType CommandGroupType = "debug"
)

type CommandGroup struct {
	// Kind of group the command is part of
	Kind CommandGroupType `json:"kind"`

	// +optional
	// Identifies the default command for a given group kind
	IsDefault bool `json:"isDefault,omitempty"`
}

type BaseCommand struct {
	// Mandatory identifier that allows referencing
	// this command in composite commands, or from
	// a parent, or in events.
	Id string `json:"id"`

	// +optional
	// Defines the group this command is part of
	Group *CommandGroup `json:"group,omitempty"`

	// Optional map of free-form additional command attributes
	Attributes map[string]string `json:"attributes,omitempty"`
}

type LabeledCommand struct {
	BaseCommand `json:",inline"`

	// +optional
	// Optional label that provides a label for this command
	// to be used in Editor UI menus for example
	Label string `json:"label,omitempty"`
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
	CommandType CommandType `json:"commandType"`

	// CLI Command executed in a component container
	// +optional
	Exec *ExecCommand `json:"exec,omitempty"`

	// Command providing the definition of a VsCode Task
	// +optional
	VscodeTask *VscodeConfigurationCommand `json:"vscodeTask,omitempty"`

	// Command providing the definition of a VsCode launch action
	// +optional
	VscodeLaunch *VscodeConfigurationCommand `json:"vscodeLaunch,omitempty"`

	// Composite command that allows executing several sub-commands
	// either sequentially or concurrently
	// +optional
	Composite *CompositeCommand `json:"composite,omitempty"`

	// Custom command whose logic is implementation-dependant
	// and should be provided by the user
	// possibly through some dedicated plugin
	// +optional
	Custom *CustomCommand `json:"custom,omitempty"`
}

type ExecCommand struct {
	LabeledCommand `json:",inline"`

	// The actual command-line string
	CommandLine string `json:"commandLine,omitempty"`

	// Describes component to which given action relates
	Component string `json:"component,omitempty"`

	// Working directory where the command should be executed
	WorkingDir string `json:"workingDir,omitempty"`

	// +optional
	// Optional list of environment variables that have to be set
	// before running the command
	Env []EnvVar `json:"env,omitempty"`
}

type CompositeCommand struct {
	LabeledCommand `json:",inline"`

	// The commands that comprise this composite command
	Commands []string `json:"commands,omitempty"`

	// Indicates if the sub-commands should be executed concurrently
	// +optional
	Parallel bool `json:"parallel,omitempty"`
}

// VscodeConfigurationCommandLocationType describes the type of
// the location the configuration is fetched from.
// Only one of the following component type may be specified.
// +kubebuilder:validation:Enum=Uri;Inlined
type VscodeConfigurationCommandLocationType string

const (
	UriVscodeConfigurationCommandLocationType     VscodeConfigurationCommandLocationType = "Container"
	InlinedVscodeConfigurationCommandLocationType VscodeConfigurationCommandLocationType = "Kubernetes"
)

// +k8s:openapi-gen=true
// +union
type VscodeConfigurationCommandLocation struct {
	// Type of Vscode configuration command location
	// +
	// +unionDiscriminator
	// +optional
	LocationType VscodeConfigurationCommandLocationType `json:"locationType"`

	// Location as an absolute of relative URI
	// the VsCode configuration will be fetched from
	// +optional
	Uri string `json:"uri,omitempty"`

	// Inlined content of the VsCode configuration
	// +optional
	Inlined string `json:"inlined,omitempty"`
}

type VscodeConfigurationCommand struct {
	BaseCommand                        `json:",inline"`
	VscodeConfigurationCommandLocation `json:",inline"`
}

type CustomCommand struct {
	LabeledCommand `json:",inline"`

	// Class of command that the associated implementation component
	// should use to process this command with the appropriate logic
	CommandClass   string `json:"commandClass"`

	// Additional free-form configuration for this custom command
	// that the implementation component will know how to use
  // 	
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:EmbeddedResource
	EmbeddedResource runtime.RawExtension `json:"embeddedResource"`
}
