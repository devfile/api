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

package overriding

import (
	"reflect"
	"testing"
)

func TestUnionStrings(t *testing.T) {
	tests := []struct {
		name     string
		a        []string
		b        []string
		expected []string
	}{
		{
			name:     "both empty",
			a:        []string{},
			b:        []string{},
			expected: []string{},
		},
		{
			name:     "first is nil, second is empty",
			a:        nil,
			b:        []string{},
			expected: []string{},
		},
		{
			name:     "first is empty, second is nil",
			a:        []string{},
			b:        nil,
			expected: []string{},
		},
		{
			name:     "both are nil",
			a:        nil,
			b:        nil,
			expected: []string{},
		},
		{
			name:     "no overlap",
			a:        []string{"x", "y"},
			b:        []string{"a", "b"},
			expected: []string{"x", "y", "a", "b"},
		},
		{
			name:     "partial overlap",
			a:        []string{"a", "b"},
			b:        []string{"b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "duplicate in same slice",
			a:        []string{"a", "a", "b"},
			b:        []string{"b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "only duplicates",
			a:        []string{"a", "b"},
			b:        []string{"a", "b"},
			expected: []string{"a", "b"},
		},
		{
			name:     "preserves order",
			a:        []string{"c", "a"},
			b:        []string{"b", "a", "d"},
			expected: []string{"c", "a", "b", "d"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := UnionStrings(tt.a, tt.b)

			if actual == nil {
				t.Errorf("expected non-nil slice, got nil")
			}

			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("UnionStrings(%v, %v) = %v; want %v", tt.a, tt.b, actual, tt.expected)
			}
		})
	}
}
