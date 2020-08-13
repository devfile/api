package overriding

import (
	//	"errors"
	"errors"
	"strings"

	workspaces "github.com/devfile/api/pkg/apis/workspaces/v1alpha1"
	unions "github.com/devfile/api/pkg/utils/unions"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/sets"
	strategicpatch "k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func ensureOnlyExistingElementsAreOverriden(keyedSpec *keyedDevWorkspaceTemplateSpecContent, keyedOverlay *keyedDevWorkspaceTemplateSpecContent) error {
	return checkKeys(func(elementType string, keysSets []sets.String) []error {
		specKeys := keysSets[0]
		overlayKeys := keysSets[1]
		newElementsInOverlay := overlayKeys.Difference(specKeys)
		if newElementsInOverlay.Len() > 0 {
			return []error{errors.New("Some " +
				elementType +
				" elements do not override any existing element: " +
				strings.Join(newElementsInOverlay.List(), ", ") +
				". They should be defined in the main body, as new elements, not in the overriding section")}
		}
		return []error{}
	},
		keyedSpec, keyedOverlay)
}

func overrideKeyedDevWorkspaceTemplateSpec(keyedOriginal *keyedDevWorkspaceTemplateSpecContent, keyedPatch *keyedDevWorkspaceTemplateSpecContent) (*keyedDevWorkspaceTemplateSpecContent, error) {
	if err := addKeys(keyedOriginal); err != nil {
		return nil, err
	}
	if err := addKeys(keyedPatch); err != nil {
		return nil, err
	}

	if err:= unions.Normalize(&keyedOriginal); err != nil {
		return nil, err
	}
	if err := unions.Normalize(&keyedPatch); err != nil {
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

	keyedPatchBytes, err := json.Marshal(keyedPatch)
	if err != nil {
		return nil, err
	}
	patchMap, err := handleUnmarshal(keyedPatchBytes)
	if err != nil {
		return nil, err
	}

	if err := ensureOnlyExistingElementsAreOverriden(keyedOriginal, keyedPatch); err != nil {
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

	patched := keyedDevWorkspaceTemplateSpecContent{}
	err = json.Unmarshal(patchedBytes, &patched)
	if err != nil {
		return nil, err
	}

	if err = unions.Simplify(&patched); err != nil {
		return nil, err
	}
	return &patched, nil
}

// OverrideDevWorkspaceTemplateSpecBytes implements the overriding logic for parent devfiles or plugins.
// On an Json or Yaml document that contains the core content of the devfile (without the `apiVersion` and `metadata`),
// it allows applying a `patch` which is a document fragment of the same schema.
//
// The Overriding logic is implemented according to strategic merge patch rules, as defined here:
// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-api-machinery/strategic-merge-patch.md#background
//
// The result is a transformed `DevfileWorkspaceTemplateSpec` object that can be serialized back to Yaml or Json.
func OverrideDevWorkspaceTemplateSpecBytes(originalBytes []byte, patchBytes []byte) (*workspaces.DevWorkspaceTemplateSpecContent, error) {
	originalJson, err := yaml.ToJSON(originalBytes)
	if err != nil {
		return nil, err
	}

	keyedOriginal := keyedDevWorkspaceTemplateSpecContent{}
	err = json.Unmarshal(originalJson, &keyedOriginal)
	if err != nil {
		return nil, err
	}

	patchJson, err := yaml.ToJSON(patchBytes)
	if err != nil {
		return nil, err
	}

	keyedPatch := keyedDevWorkspaceTemplateSpecContent{}
	err = json.Unmarshal(patchJson, &keyedPatch)
	if err != nil {
		return nil, err
	}

	patched, err := overrideKeyedDevWorkspaceTemplateSpec(&keyedOriginal, &keyedPatch)
	if err != nil {
		return nil, err
	}

	return removeKeys(patched), nil
}

// OverrideDevWorkspaceTemplateSpec implements the overriding logic for parent devfiles or plugins.
// On an `original` `DevfileWorkspaceTemplateSpec` (which is the core part of a devfile, without the `apiVersion` and `metadata`),
// it allows applying a `patch` which is an `Overrides` or `PluginOverrides` object.
//
// The Overriding logic is implemented according to strategic merge patch rules, as defined here:
// https://github.com/kubernetes/community/blob/master/contributors/devel/sig-api-machinery/strategic-merge-patch.md#background
//
// The result is a transformed `DevfileWorkspaceTemplateSpec` object.
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
