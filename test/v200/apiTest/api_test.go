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

package apiTest

import (
	"testing"

	schema "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiUtils "github.com/devfile/api/v2/test/v200/utils/api"
	commonUtils "github.com/devfile/api/v2/test/v200/utils/common"
)

func Test_ExecCommand(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.CommandTypes = []schema.CommandType{schema.ExecCommandType}
	testContent.EditContent = false
	testContent.FileName = commonUtils.GetDevFileName()
	apiUtils.RunTest(testContent, t)
}

func Test_ApplyCommand(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.CommandTypes = []schema.CommandType{schema.ApplyCommandType}
	testContent.EditContent = false
	testContent.FileName = commonUtils.GetDevFileName()
	apiUtils.RunTest(testContent, t)
}

func Test_CompositeCommand(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.CommandTypes = []schema.CommandType{schema.CompositeCommandType}
	testContent.EditContent = false
	testContent.FileName = commonUtils.GetDevFileName()
	apiUtils.RunTest(testContent, t)
}

func Test_MultiCommand(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.CommandTypes = []schema.CommandType{schema.ExecCommandType,
		schema.CompositeCommandType,
		schema.ApplyCommandType}
	testContent.EditContent = true
	testContent.FileName = commonUtils.GetDevFileName()
	apiUtils.RunTest(testContent, t)
}

func Test_ContainerComponent(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.ComponentTypes = []schema.ComponentType{schema.ContainerComponentType}
	testContent.EditContent = false
	testContent.FileName = commonUtils.GetDevFileName()
	apiUtils.RunTest(testContent, t)
}

func Test_KubernetesComponent(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.ComponentTypes = []schema.ComponentType{schema.KubernetesComponentType}
	testContent.EditContent = false
	testContent.FileName = commonUtils.GetDevFileName()
	apiUtils.RunTest(testContent, t)
}

func Test_OpenshiftComponent(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.ComponentTypes = []schema.ComponentType{schema.OpenshiftComponentType}
	testContent.EditContent = false
	testContent.FileName = commonUtils.GetDevFileName()
	apiUtils.RunTest(testContent, t)
}

func Test_VolumeComponent(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.ComponentTypes = []schema.ComponentType{schema.VolumeComponentType}
	testContent.FileName = commonUtils.GetDevFileName()
	apiUtils.RunTest(testContent, t)
}

func Test_MultiComponent(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.ComponentTypes = commonUtils.ComponentTypes
	testContent.FileName = commonUtils.GetDevFileName()
	apiUtils.RunTest(testContent, t)
}

func Test_Projects(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.ProjectTypes = commonUtils.ProjectSourceTypes
	testContent.FileName = commonUtils.GetDevFileName()
	apiUtils.RunTest(testContent, t)
}

func Test_StarterProjects(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.StarterProjectTypes = commonUtils.ProjectSourceTypes
	testContent.FileName = commonUtils.GetDevFileName()
	apiUtils.RunTest(testContent, t)
}

func Test_Events(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.AddEvents = true
	testContent.FileName = commonUtils.GetDevFileName()
	apiUtils.RunTest(testContent, t)
}

func Test_MetaData(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.AddMetaData = true
	testContent.FileName = commonUtils.GetDevFileName()
	apiUtils.RunTest(testContent, t)
}

func Test_Parent(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.AddParent = true
	testContent.FileName = commonUtils.GetDevFileName()
	apiUtils.RunTest(testContent, t)
}

func Test_Everything(t *testing.T) {
	testContent := commonUtils.TestContent{}
	testContent.CommandTypes = commonUtils.CommandTypes
	testContent.ComponentTypes = commonUtils.ComponentTypes
	testContent.ProjectTypes = commonUtils.ProjectSourceTypes
	testContent.StarterProjectTypes = commonUtils.ProjectSourceTypes
	testContent.AddMetaData = true
	testContent.AddEvents = true
	testContent.AddParent = true
	testContent.FileName = commonUtils.GetDevFileName()
	apiUtils.RunTest(testContent, t)
}
