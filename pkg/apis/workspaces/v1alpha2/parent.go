package v1alpha2

type Parent struct {
	ImportReference `json:",inline" yaml:",inline"`
	ParentOverrides `json:",inline" yaml:",inline"`
}
