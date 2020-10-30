package v1alpha2

type PluginComponent struct {
	BaseComponent   `json:",inline" yaml:",inline"`
	ImportReference `json:",inline" yaml:",inline"`
	PluginOverrides `json:",inline" yaml:",inline"`
}
