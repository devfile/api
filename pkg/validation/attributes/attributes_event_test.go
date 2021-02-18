package attributes

import (
	"testing"

	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	apiAttributes "github.com/devfile/api/v2/pkg/attributes"
	"github.com/stretchr/testify/assert"
)

func TestValidateEvents(t *testing.T) {

	tests := []struct {
		name       string
		testFile   string
		expected   v1alpha2.Events
		attributes apiAttributes.Attributes
		wantErr    bool
	}{
		{
			name:     "Good Substitution",
			testFile: "test-fixtures/events/event.yaml",
			expected: v1alpha2.Events{
				WorkspaceEvents: v1alpha2.WorkspaceEvents{
					PreStart: []string{
						"FOO",
					},
					PostStart: []string{
						"BAR",
					},
					PreStop: []string{
						"FOOBAR",
					},
					PostStop: []string{
						"BARFOO",
					},
				},
			},
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"bar": "BAR",
				"foo": "FOO",
			}, nil),
			wantErr: false,
		},
		{
			name:     "Invalid Reference",
			testFile: "test-fixtures/events/event.yaml",
			attributes: apiAttributes.Attributes{}.FromMap(map[string]interface{}{
				"foo": "FOO",
			}, nil),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testEvents := v1alpha2.Events{}

			readFileToStruct(t, tt.testFile, &testEvents)

			err := ValidateEvents(tt.attributes, &testEvents)
			if tt.wantErr == (err == nil) {
				t.Errorf("error: %v", err)
				return
			} else if err == nil {
				assert.Equal(t, tt.expected, testEvents, "The two values should be the same.")
			}
		})
	}
}
