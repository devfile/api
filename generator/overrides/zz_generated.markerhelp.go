//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Generated for the devfile generator

// Code generated by helpgen. DO NOT EDIT.

package overrides

import (
	"sigs.k8s.io/controller-tools/pkg/markers"
)

func (FieldOverridesInclude) Help() *markers.DefinitionHelp {
	return &markers.DefinitionHelp{
		Category: "Overrides",
		DetailedHelp: markers.DetailedHelp{
			Summary: "drives whether a field should be overriden in devfile parent or plugins",
			Details: "",
		},
		FieldHelp: map[string]markers.DetailedHelp{
			"Omit": {
				Summary: "indicates that this field cannot be overridden at all.",
				Details: "",
			},
			"OmitInPlugin": {
				Summary: "OmmitInPlugin indicates that this field cannot be overridden in a devfile plugin.",
				Details: "",
			},
			"Description": {
				Summary: "indicates the description that should be added as Go documentation on the generated structs.",
				Details: "",
			},
		},
	}
}

func (Generator) Help() *markers.DefinitionHelp {
	return &markers.DefinitionHelp{
		Category: "",
		DetailedHelp: markers.DetailedHelp{
			Summary: "generates additional GO code for the overriding of elements in devfile parent or plugins.",
			Details: "",
		},
		FieldHelp: map[string]markers.DetailedHelp{
			"IsForPluginOverrides": {
				Summary: "indicates that the generated code should be done for plugin overrides. When false, the parent overrides are generated",
				Details: "",
			},
			"suffix": {
				Summary: "",
				Details: "",
			},
			"rootTypeToProcess": {
				Summary: "",
				Details: "",
			},
		},
	}
}
