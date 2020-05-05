package v1alpha1

type Parent struct {
	ImportReference                 `json:",inline"`
	DevWorkspaceTemplateSpecContent `json:",inline"`
}
