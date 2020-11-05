package v1alpha2

func (keyed Preference) Key() string {
	return keyed.Name
}

func (keyed Component) Key() string {
	return keyed.Name
}

func (keyed Project) Key() string {
	return keyed.Name
}

func (keyed StarterProject) Key() string {
	return keyed.Name
}

func (keyed Command) Key() string {
	return keyed.Id
}

func (keyed PreferenceParentOverride) Key() string {
	return keyed.Name
}

func (keyed ComponentParentOverride) Key() string {
	return keyed.Name
}

func (keyed ProjectParentOverride) Key() string {
	return keyed.Name
}

func (keyed StarterProjectParentOverride) Key() string {
	return keyed.Name
}

func (keyed CommandParentOverride) Key() string {
	return keyed.Id
}

func (keyed PreferencePluginOverrideParentOverride) Key() string {
	return keyed.Name
}

func (keyed ComponentPluginOverrideParentOverride) Key() string {
	return keyed.Name
}

func (keyed CommandPluginOverrideParentOverride) Key() string {
	return keyed.Id
}

func (keyed PreferencePluginOverride) Key() string {
	return keyed.Name
}

func (keyed ComponentPluginOverride) Key() string {
	return keyed.Name
}

func (keyed CommandPluginOverride) Key() string {
	return keyed.Id
}
