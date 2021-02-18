package attributes

import (
	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
)

// ValidateCommands validates the commands data for a global attribute
func ValidateCommands(attributes apiAttributes.Attributes, commands *[]v1alpha2.Command) error {

	if commands != nil {
		for i := range *commands {
			var err error

			// Validate various command types
			switch {
			case (*commands)[i].Exec != nil:
				if err = validateExecCommand(attributes, (*commands)[i].Exec); err != nil {
					return err
				}
			case (*commands)[i].Composite != nil:
				if err = validateCompositeCommand(attributes, (*commands)[i].Composite); err != nil {
					return err
				}
			case (*commands)[i].Apply != nil:
				if err = validateApplyCommand(attributes, (*commands)[i].Apply); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// validateExecCommand validates the exec command data for a global attribute
func validateExecCommand(attributes apiAttributes.Attributes, exec *v1alpha2.ExecCommand) error {
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
			if err = validateEnv(attributes, &exec.Env); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateExecCommand validates the composite command data for a global attribute
func validateCompositeCommand(attributes apiAttributes.Attributes, composite *v1alpha2.CompositeCommand) error {
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

// validateApplyCommand validates the apply command data for a global attribute
func validateApplyCommand(attributes apiAttributes.Attributes, apply *v1alpha2.ApplyCommand) error {
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
