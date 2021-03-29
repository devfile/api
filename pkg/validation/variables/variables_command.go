package variables

import (
	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
)

// ValidateAndReplaceForCommands validates the commands data for global variable references and replaces them with the variable value
func ValidateAndReplaceForCommands(variables map[string]string, commands []v1alpha2.Command) error {

	for i := range commands {
		var err error

		// Validate various command types
		switch {
		case commands[i].Exec != nil:
			if err = validateAndReplaceForExecCommand(variables, commands[i].Exec); err != nil {
				return err
			}
		case commands[i].Composite != nil:
			if err = validateAndReplaceForCompositeCommand(variables, commands[i].Composite); err != nil {
				return err
			}
		case commands[i].Apply != nil:
			if err = validateAndReplaceForApplyCommand(variables, commands[i].Apply); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateAndReplaceForExecCommand validates the exec command data for global variable references and replaces them with the variable value
func validateAndReplaceForExecCommand(variables map[string]string, exec *v1alpha2.ExecCommand) error {
	var err error

	if exec != nil {
		// Validate exec command line
		if exec.CommandLine, err = validateAndReplaceDataWithVariable(exec.CommandLine, variables); err != nil {
			return err
		}

		// Validate exec working dir
		if exec.WorkingDir, err = validateAndReplaceDataWithVariable(exec.WorkingDir, variables); err != nil {
			return err
		}

		// Validate exec label
		if exec.Label, err = validateAndReplaceDataWithVariable(exec.Label, variables); err != nil {
			return err
		}

		// Validate exec env
		if len(exec.Env) > 0 {
			if err = validateAndReplaceForEnv(variables, exec.Env); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateAndReplaceForCompositeCommand validates the composite command data for global variable references and replaces them with the variable value
func validateAndReplaceForCompositeCommand(variables map[string]string, composite *v1alpha2.CompositeCommand) error {
	var err error

	if composite != nil {
		// Validate composite label
		if composite.Label, err = validateAndReplaceDataWithVariable(composite.Label, variables); err != nil {
			return err
		}
	}

	return nil
}

// validateAndReplaceForApplyCommand validates the apply command data for global variable references and replaces them with the variable value
func validateAndReplaceForApplyCommand(variables map[string]string, apply *v1alpha2.ApplyCommand) error {
	var err error

	if apply != nil {
		// Validate apply label
		if apply.Label, err = validateAndReplaceDataWithVariable(apply.Label, variables); err != nil {
			return err
		}
	}

	return nil
}
