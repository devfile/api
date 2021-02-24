package attributes

import (
	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
)

// ValidateAndReplaceForCommands validates the commands data for global attribute references and replaces them with the attribute value
func ValidateAndReplaceForCommands(attributes apiAttributes.Attributes, commands []v1alpha2.Command) error {

	for i := range commands {
		var err error

		// Validate various command types
		switch {
		case commands[i].Exec != nil:
			if err = validateAndReplaceForExecCommand(attributes, commands[i].Exec); err != nil {
				return err
			}
		case commands[i].Composite != nil:
			if err = validateAndReplaceForCompositeCommand(attributes, commands[i].Composite); err != nil {
				return err
			}
		case commands[i].Apply != nil:
			if err = validateAndReplaceForApplyCommand(attributes, commands[i].Apply); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateAndReplaceForExecCommand validates the exec command data for global attribute references and replaces them with the attribute value
func validateAndReplaceForExecCommand(attributes apiAttributes.Attributes, exec *v1alpha2.ExecCommand) error {
	var err error

	if exec != nil {
		// Validate exec command line
		if exec.CommandLine, err = validateAndReplaceDataWithAttribute(exec.CommandLine, attributes); err != nil {
			return err
		}

		// Validate exec component
		if exec.Component, err = validateAndReplaceDataWithAttribute(exec.Component, attributes); err != nil {
			return err
		}

		// Validate exec working dir
		if exec.WorkingDir, err = validateAndReplaceDataWithAttribute(exec.WorkingDir, attributes); err != nil {
			return err
		}

		// Validate exec label
		if exec.Label, err = validateAndReplaceDataWithAttribute(exec.Label, attributes); err != nil {
			return err
		}

		// Validate exec env
		if len(exec.Env) > 0 {
			if err = validateAndReplaceForEnv(attributes, exec.Env); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateAndReplaceForCompositeCommand validates the composite command data for global attribute references and replaces them with the attribute value
func validateAndReplaceForCompositeCommand(attributes apiAttributes.Attributes, composite *v1alpha2.CompositeCommand) error {
	var err error

	if composite != nil {
		// Validate composite label
		if composite.Label, err = validateAndReplaceDataWithAttribute(composite.Label, attributes); err != nil {
			return err
		}

		// Validate composite commands
		for i := range composite.Commands {
			if composite.Commands[i], err = validateAndReplaceDataWithAttribute(composite.Commands[i], attributes); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateAndReplaceForApplyCommand validates the apply command data for global attribute references and replaces them with the attribute value
func validateAndReplaceForApplyCommand(attributes apiAttributes.Attributes, apply *v1alpha2.ApplyCommand) error {
	var err error

	if apply != nil {
		// Validate composite label
		if apply.Label, err = validateAndReplaceDataWithAttribute(apply.Label, attributes); err != nil {
			return err
		}

		// Validate apply component
		if apply.Component, err = validateAndReplaceDataWithAttribute(apply.Component, attributes); err != nil {
			return err
		}
	}

	return nil
}
