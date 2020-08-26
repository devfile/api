package v1alpha2

func (keyed Component) Key() string {
	return keyed.Name
}

func (keyed PluginComponentsOverride) Key() string {
	return keyed.Name
}

func (keyed Command) Key() string {
	return keyed.Id
}

func (keyed Project) Key() string {
	return keyed.Name
}

func (container DevWorkspaceTemplateSpecContent) GetToplevelLists() TopLevelLists {
	return TopLevelLists {
		"Projects": extractKeys(container.Projects),		
		"Components": extractKeys(container.Components),		
		"Commands": extractKeys(container.Commands),		
	}
}

func (container PluginOverrides) GetToplevelLists() TopLevelLists {
	return TopLevelLists {
		"Components": extractKeys(container.Components),		
		"Commands": extractKeys(container.Commands),		
	}
}

func (container ParentOverrides) GetToplevelLists() TopLevelLists {
	return TopLevelLists {
		"Projects": extractKeys(container.Projects),		
		"Components": extractKeys(container.Components),		
		"Commands": extractKeys(container.Commands),		
	}
}
