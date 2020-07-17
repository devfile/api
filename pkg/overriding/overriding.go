package overriding

import (
	//	"errors"
	"reflect"
	workspaces "github.com/devfile/kubernetes-api/pkg/apis/workspaces/v1alpha1"
	"k8s.io/apimachinery/pkg/util/json"
	strategicpatch "k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/yaml"
	"github.com/mitchellh/reflectwalk"
)

func handleUnmarshal(j []byte) (map[string]interface{}, error) {
	if j == nil {
		j = []byte("{}")
	}

	m := map[string]interface{}{}
	err := json.Unmarshal(j, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

type keyedCommand struct {
	Id string `json:"id,omitempty"`

	workspaces.Command `json:",inline"`
}

type keyedComponent struct {
	Name string `json:"name,omitempty"`

	workspaces.Component `json:",inline"`
}

type keyedDevWorkspaceTemplateSpecContent struct {
	// Projects worked on in the workspace, containing names and sources locations
	// +optional
	Projects []workspaces.Project `json:"projects,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// List of the workspace components, such as editor and plugins,
	// user-provided containers, or other types of components.
	// +optional
	Components []keyedComponent `json:"components,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// Predefined, ready-to-use, workspace-related commands.
	// +optional
	Commands []keyedCommand `json:"commands,omitempty" patchStrategy:"merge" patchMergeKey:"id"`

	// Bindings of commands to events.
	// Each command is referred-to by its name.
	// +optional
	Events *workspaces.Events `json:"events,omitempty"`
}

func addKeys(keyedSpec *keyedDevWorkspaceTemplateSpecContent) error {
	for idx := range keyedSpec.Commands {
		key, err := keyedSpec.Commands[idx].Key()
		if err != nil {
			return err
		}
		keyedSpec.Commands[idx].Id = key
	}
	for idx := range keyedSpec.Components {
		key, err := keyedSpec.Components[idx].Key()
		if err != nil {
			return err
		}
		keyedSpec.Components[idx].Name = key
	}
	return nil
}

func removeKeys(keyedSpec *keyedDevWorkspaceTemplateSpecContent) *workspaces.DevWorkspaceTemplateSpecContent {
	result := workspaces.DevWorkspaceTemplateSpecContent {
		Projects: keyedSpec.Projects,
		Events: keyedSpec.Events,
	}

	if keyedSpec.Commands != nil && len(keyedSpec.Commands) > 0 {
		result.Commands = []workspaces.Command{}
	}

	if keyedSpec.Components != nil && len(keyedSpec.Components) > 0 {
		result.Components = []workspaces.Component{}
	}

	for _, keyedCommand := range keyedSpec.Commands {
		result.Commands = append(result.Commands, keyedCommand.Command)
	}
	for _, keyedComponent := range keyedSpec.Components {
		result.Components = append(result.Components, keyedComponent.Component)
	}
	return &result
}

type normalizer struct {
}
func (n *normalizer) Struct(s reflect.Value) error {
	if s.CanAddr() {
		addr := s.Addr()
		if addr.CanInterface() {
			i := addr.Interface()
			if u, ok := i.(workspaces.Union); ok {
				u.Normalize()
			}
		}
	}
	return nil
}
func (n *normalizer) StructField(reflect.StructField, reflect.Value) error {
	return nil
}

type simplifier struct {
}
func (n *simplifier) Struct(s reflect.Value) error {
	if s.CanAddr() {
		addr := s.Addr()
		if addr.CanInterface() {
			i := addr.Interface()
			if u, ok := i.(workspaces.Union); ok {
				u.Normalize()
				u.Simplify()
			}
		}
	}
	return nil
}
func (n *simplifier) StructField(reflect.StructField, reflect.Value) error {
	return nil
}

func Normalize(keyedSpec *keyedDevWorkspaceTemplateSpecContent) error {
	return reflectwalk.Walk(keyedSpec, &normalizer{})
}

func Simplify(keyedSpec *keyedDevWorkspaceTemplateSpecContent) error {
	return reflectwalk.Walk(keyedSpec, &simplifier{})
}

func OverrideDevWorkspaceTemplateSpecBytes(originalBytes []byte, patchBytes []byte) (*workspaces.DevWorkspaceTemplateSpecContent, error) {
	originalJson, err := yaml.ToJSON(originalBytes)
	if err != nil {
		return nil, err
	}

	keyedOriginal := keyedDevWorkspaceTemplateSpecContent {}
	err = json.Unmarshal(originalJson, &keyedOriginal)
	if err != nil {
		return nil, err
	}

	if err = Normalize(&keyedOriginal); err != nil {
		return nil, err
	}

	if err = addKeys(&keyedOriginal); err != nil {
		return nil, err
	}

	keyedOriginalBytes, err := json.Marshal(keyedOriginal)
	if err != nil {
		return nil, err
	}
	
	originalMap, err := handleUnmarshal(keyedOriginalBytes)
	if err != nil {
		return nil, err
	}

	patchJson, err := yaml.ToJSON(patchBytes)
	if err != nil {
		return nil, err
	}
	
	keyedPatch := keyedDevWorkspaceTemplateSpecContent {}
	err = json.Unmarshal(patchJson, &keyedPatch)
	if err != nil {
		return nil, err
	}

	if err = Normalize(&keyedPatch); err != nil {
		return nil, err
	}

	if err = addKeys(&keyedPatch); err != nil {
		return nil, err
	}

	keyedPatchBytes, err := json.Marshal(keyedPatch)
	if err != nil {
		return nil, err
	}
	patchMap, err := handleUnmarshal(keyedPatchBytes)
	if err != nil {
		return nil, err
	}


	schema, err := strategicpatch.NewPatchMetaFromStruct(keyedOriginal)
	if err != nil {
		return nil, err
	}

	patchedMap, err := strategicpatch.StrategicMergeMapPatchUsingLookupPatchMeta(originalMap, patchMap, schema)
	if err != nil {
		return nil, err
	}

	patchedBytes, err := json.Marshal(patchedMap)
	if err != nil {
		return nil, err
	}

	patched := keyedDevWorkspaceTemplateSpecContent {}
	err = json.Unmarshal(patchedBytes, &patched)
	if err != nil {
		return nil, err
	}

	if err = Simplify(&patched); err != nil {
		return nil, err
	}

	return removeKeys(&patched), nil
}

func OverrideDevWorkspaceTemplateSpec(original *workspaces.DevWorkspaceTemplateSpecContent, patch interface{}) (*workspaces.DevWorkspaceTemplateSpecContent, error) {
	originalBytes, err := json.Marshal(original)
	if err != nil {
		return nil, err
	}
	
	patchBytes, err := json.Marshal(patch)
	if err != nil {
		return nil, err
	}
	
	return OverrideDevWorkspaceTemplateSpecBytes(originalBytes, patchBytes)
}
