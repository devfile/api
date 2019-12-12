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
