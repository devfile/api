//
//
// Copyright Red Hat
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/google/go-cmp/cmp"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
)

const fuzzIterations = 500
const fuzzNilChance = 0.2

var devWorkspaceFuzzFunc = func(workspace *DevWorkspace, c fuzz.Continue) {
	c.Fuzz(&workspace.Status)
	c.Fuzz(&workspace.Spec)
}

var devWorkspaceTemplateFuzzFunc = func(workspace *DevWorkspaceTemplate, c fuzz.Continue) {
	c.Fuzz(&workspace.Spec)
}

var componentFuzzFunc = func(component *Component, c fuzz.Continue) {
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

var commandFuzzFunc = func(command *Command, c fuzz.Continue) {
	switch c.Intn(4) {
	case 0:
		c.Fuzz(&command.Apply)
	case 1:
		c.Fuzz(&command.Composite)
	case 2:
		c.Fuzz(&command.Custom)
	case 3:
		c.Fuzz(&command.Exec)
	}
}

var pluginComponentsOverrideFuzzFunc = func(component *PluginComponentsOverride, c fuzz.Continue) {
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

var pluginComponentFuzzFunc = func(plugin *PluginComponent, c fuzz.Continue) {
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

var parentFuzzFunc = func(parent *Parent, c fuzz.Continue) {
	for i := 0; i < c.Intn(4); i++ {
		component := Component{}
		parentComponentFuzzFunc(&component, c)
		parent.Components = append(parent.Components, component)
	}
	for i := 0; i < c.Intn(4); i++ {
		command := Command{}
		parentCommandFuzzFunc(&command, c)
		parent.Commands = append(parent.Commands, command)
	}
	for i := 0; i < c.Intn(4); i++ {
		project := Project{}
		parentProjectFuzzFunc(&project, c)
		parent.Projects = append(parent.Projects, project)
	}
	for i := 0; i < c.Intn(4); i++ {
		starterProject := StarterProject{}
		starterProject.Description = c.RandString()
		parentProjectFuzzFunc(&starterProject.Project, c)
		parent.StarterProjects = append(parent.StarterProjects, starterProject)
	}
}

var conditionFuzzFunc = func(condition *WorkspaceCondition, c fuzz.Continue) {
	condition.Reason = c.RandString()
	condition.Type = WorkspaceConditionType(c.RandString())
	condition.Message = c.RandString()
}

var parentComponentFuzzFunc = func(component *Component, c fuzz.Continue) {
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

var parentCommandFuzzFunc = func(command *Command, c fuzz.Continue) {
	// Do not generate Custom commands for Parents
	// Also: commands in Parents cannot have attributes.
	switch c.Intn(3) {
	case 0:
		c.Fuzz(&command.Apply)
		if command.Apply != nil {
			command.Apply.Attributes = nil
		}
	case 1:
		c.Fuzz(&command.Composite)
		if command.Composite != nil {
			command.Composite.Attributes = nil
		}
	case 2:
		c.Fuzz(&command.Exec)
		if command.Exec != nil {
			command.Exec.Attributes = nil
		}
	}
}

var parentProjectFuzzFunc = func(project *Project, c fuzz.Continue) {
	// Custom projects are not supported in v1alpha2 parent
	project.Name = c.RandString()
	switch c.Intn(3) {
	case 0:
		c.Fuzz(&project.Git)
		if project.Git != nil {
			project.Git.SparseCheckoutDir = ""
		}
	case 1:
		c.Fuzz(&project.Github)
		if project.Github != nil {
			project.Github.SparseCheckoutDir = ""
		}
	case 2:
		c.Fuzz(&project.Zip)
		if project.Zip != nil {
			project.Zip.SparseCheckoutDir = ""
		}
	}
}

var projectFuzzFunc = func(project *Project, c fuzz.Continue) {
	switch c.Intn(4) {
	case 0:
		c.Fuzz(&project.Git)
		if project.Git != nil {
			project.Git.SparseCheckoutDir = ""
		}
	case 1:
		c.Fuzz(&project.Github)
		if project.Github != nil {
			project.Github.SparseCheckoutDir = ""
		}
	case 2:
		c.Fuzz(&project.Zip)
		if project.Zip != nil {
			project.Zip.SparseCheckoutDir = ""
		}
	case 3:
		c.Fuzz(&project.Custom)
	}
}

// embeddedResource.Object is an interface and hard to fuzz right now.
var rawExtFuzzFunc = func(embeddedResource *runtime.RawExtension, c fuzz.Continue) {}

func TestDevWorkspaceConversion_v1alpha1(t *testing.T) {
	f := fuzz.New().NilChance(fuzzNilChance).MaxDepth(100).Funcs(
		devWorkspaceFuzzFunc,
		conditionFuzzFunc,
		parentFuzzFunc,
		componentFuzzFunc,
		commandFuzzFunc,
		projectFuzzFunc,
		pluginComponentsOverrideFuzzFunc,
		pluginComponentFuzzFunc,
		rawExtFuzzFunc,
	)
	for i := 0; i < fuzzIterations; i++ {
		original := &DevWorkspace{}
		intermediate := &v1alpha2.DevWorkspace{}
		output := &DevWorkspace{}
		f.Fuzz(original)
		input := original.DeepCopy()
		err := convertDevWorkspaceTo_v1alpha2(input, intermediate)
		if !assert.NoError(t, err, "Should not return error when converting to v1alpha2") {
			return
		}
		err = convertDevWorkspaceFrom_v1alpha2(intermediate, output)
		if !assert.NoError(t, err, "Should not return error when converting from v1alpha2") {
			return
		}
		if !assert.True(t, cmp.Equal(original, output), "DevWorkspace should not be changed when converting between v1alpha1 and v1alpha2") {
			t.Logf("Diff: \n%s\n", cmp.Diff(original, output))
		}
	}
}

func TestDevWorkspaceTemplateConversion_v1alpha1(t *testing.T) {
	f := fuzz.New().NilChance(fuzzNilChance).MaxDepth(100).Funcs(
		devWorkspaceTemplateFuzzFunc,
		conditionFuzzFunc,
		parentFuzzFunc,
		componentFuzzFunc,
		commandFuzzFunc,
		projectFuzzFunc,
		pluginComponentsOverrideFuzzFunc,
		pluginComponentFuzzFunc,
		rawExtFuzzFunc,
	)
	for i := 0; i < fuzzIterations; i++ {
		original := &DevWorkspaceTemplate{}
		intermediate := &v1alpha2.DevWorkspaceTemplate{}
		output := &DevWorkspaceTemplate{}
		f.Fuzz(original)
		input := original.DeepCopy()
		err := convertDevWorkspaceTemplateTo_v1alpha2(input, intermediate)
		if !assert.NoError(t, err, "Should not return error when converting to v1alpha2") {
			return
		}
		err = convertDevWorkspaceTemplateFrom_v1alpha2(intermediate, output)
		if !assert.NoError(t, err, "Should not return error when converting from v1alpha2") {
			return
		}
		if !assert.True(t, cmp.Equal(original, output), "DevWorkspaceTemplate should not be changed when converting between v1alpha1 and v1alpha2") {
			t.Logf("Diff: \n%s\n", cmp.Diff(original, output))
		}
	}
}

func BenchmarkDevWorkspaceConversion(b *testing.B) {
	f := fuzz.New().NilChance(fuzzNilChance).MaxDepth(100).Funcs(
		devWorkspaceFuzzFunc,
		conditionFuzzFunc,
		parentFuzzFunc,
		componentFuzzFunc,
		commandFuzzFunc,
		projectFuzzFunc,
		pluginComponentsOverrideFuzzFunc,
		pluginComponentFuzzFunc,
		rawExtFuzzFunc,
	)
	b.ResetTimer()
	b.Run("Convert to v1alpha2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			v1alpha1DW := &DevWorkspace{}
			v1alpha2DW := &v1alpha2.DevWorkspace{}
			f.Fuzz(v1alpha1DW)
			b.StartTimer()
			err := convertDevWorkspaceTo_v1alpha2(v1alpha1DW, v1alpha2DW)
			if err != nil {
				b.FailNow()
			}
		}
	})
	b.Run("Convert from v1alpha2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			v1alpha1DW := &DevWorkspace{}
			v1alpha2DW := &v1alpha2.DevWorkspace{}
			f.Fuzz(v1alpha1DW)
			err := convertDevWorkspaceTo_v1alpha2(v1alpha1DW, v1alpha2DW)
			if err != nil {
				b.FailNow()
			}
			b.StartTimer()
			err = convertDevWorkspaceFrom_v1alpha2(v1alpha2DW, v1alpha1DW)
			if err != nil {
				b.FailNow()
			}
		}
	})
}
