package v1alpha1

import (
	"github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"testing"
)

const fuzzIterations = 500
const fuzzNilChance = 0.2

var DevWorkspaceFuzzFunc = func(workspace *DevWorkspace, c fuzz.Continue) {
	c.Fuzz(&workspace.Status)
	c.Fuzz(&workspace.Spec)
}

var DevWorkspaceTemplateFuzzFunc = func(workspace *DevWorkspaceTemplate, c fuzz.Continue) {
	c.Fuzz(&workspace.Spec)
}

var ComponentFuzzFunc = func(component *Component, c fuzz.Continue) {
	switch c.Intn(6) {
	case 0: // Generate Container
		c.Fuzz(&component.Container)
	case 1: // Generate Plugin
		c.Fuzz(&component.Plugin)
	case 2: // Generate Kubernetes
		c.Fuzz(&component.Kubernetes)
	case 3: // Generate OpenShift
		c.Fuzz(&component.Openshift)
	case 4: // Generate Volume
		c.Fuzz(&component.Volume)
	case 5: // Generate Custom
		c.Fuzz(&component.Custom)
	}
}

var CommandFuzzFunc = func(command *Command, c fuzz.Continue) {
	switch c.Intn(6) {
	case 0:
		c.Fuzz(&command.Apply)
	case 1:
		c.Fuzz(&command.Composite)
	case 2:
		c.Fuzz(&command.Custom)
	case 3:
		c.Fuzz(&command.Exec)
	case 4:
		c.Fuzz(&command.VscodeLaunch)
	case 5:
		c.Fuzz(&command.VscodeTask)
	}
}

var PluginComponentsOverrideFuzzFunc = func(component *PluginComponentsOverride, c fuzz.Continue) {
	switch c.Intn(4) {
	case 0:
		c.Fuzz(&component.Container)
	case 1:
		c.Fuzz(&component.Volume)
	case 2:
		c.Fuzz(&component.Openshift)
	case 3:
		c.Fuzz(&component.Kubernetes)
	}
}

var PluginComponentFuzzFunc = func(plugin *PluginComponent, c fuzz.Continue) {
	// TODO: Figure out how to handle custom commands on plugin components
	c.Fuzz(plugin)
	plugin.Name = c.RandString()
	var filteredCommands []Command
	for _, command := range plugin.Commands {
		if command.Custom == nil {
			filteredCommands = append(filteredCommands, command)
		}
	}
	plugin.Commands = filteredCommands
}

var ParentFuzzFunc = func(parent *Parent, c fuzz.Continue) {
	for i := 0; i < c.Intn(4); i++ {
		component := Component{}
		ParentComponentFuzzFunc(&component, c)
		parent.Components = append(parent.Components, component)
	}
	for i := 0; i < c.Intn(4); i++ {
		command := Command{}
		ParentCommandFuzzFunc(&command, c)
		parent.Commands = append(parent.Commands, command)
	}
	for i := 0; i < c.Intn(4); i++ {
		project := Project{}
		ParentProjectFuzzFunc(&project, c)
		parent.Projects = append(parent.Projects, project)
	}
	for i := 0; i < c.Intn(4); i++ {
		starterProject := StarterProject{}
		starterProject.Description = c.RandString()
		ParentProjectFuzzFunc(&starterProject.Project, c)
		parent.StarterProjects = append(parent.StarterProjects, starterProject)
	}
}

var ConditionFuzzFunc = func(condition *WorkspaceCondition, c fuzz.Continue) {
	condition.Reason = c.RandString()
	condition.Type = WorkspaceConditionType(c.RandString())
	condition.Message = c.RandString()
}

var ParentComponentFuzzFunc = func(component *Component, c fuzz.Continue) {
	// Do not generate custom components when working with Parents
	switch c.Intn(5) {
	case 0: // Generate Container
		c.Fuzz(&component.Container)
	case 1: // Generate Plugin
		c.Fuzz(&component.Plugin)
	case 2: // Generate Kubernetes
		c.Fuzz(&component.Kubernetes)
	case 3: // Generate OpenShift
		c.Fuzz(&component.Openshift)
	case 4: // Generate Volume
		c.Fuzz(&component.Volume)
	}
}

var ParentCommandFuzzFunc = func(command *Command, c fuzz.Continue) {
	// Do not generate Custom commands for Parents
	switch c.Intn(5) {
	case 0:
		c.Fuzz(&command.Apply)
	case 1:
		c.Fuzz(&command.Composite)
	case 2:
		c.Fuzz(&command.Exec)
	case 3:
		c.Fuzz(&command.VscodeLaunch)
	case 4:
		c.Fuzz(&command.VscodeTask)
	}
}

var ParentProjectFuzzFunc = func(project *Project, c fuzz.Continue) {
	// Custom projects are not supported in v1alpha2 parent
	project.Name = c.RandString()
	switch c.Intn(3) {
	case 0:
		c.Fuzz(&project.Git)
	case 1:
		c.Fuzz(&project.Github)
	case 2:
		c.Fuzz(&project.Zip)
	}
}

var ProjectFuzzFunc = func(project *Project, c fuzz.Continue) {
	switch c.Intn(4) {
	case 0:
		c.Fuzz(&project.Git)
	case 1:
		c.Fuzz(&project.Github)
	case 2:
		c.Fuzz(&project.Zip)
	case 3:
		c.Fuzz(&project.Custom)
	}
}

// embeddedResource.Object is an interface and hard to fuzz right now.
var RawExtFuzzFunc = func(embeddedResource *runtime.RawExtension, c fuzz.Continue) {}

