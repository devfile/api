//
// Copyright (c) 2019-2020 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//

package unions

import (
	"reflect"

	workspaces "github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	"github.com/mitchellh/reflectwalk"
)

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
				u.Simplify()
			}
		}
	}
	return nil
}
func (n *simplifier) StructField(reflect.StructField, reflect.Value) error {
	return nil
}

// Normalize allows normalizing all the unions
// encountered while walking through the whole struct tree.
// Union normalizing works according to the following rules:
// - When only one field of the union is set and no discriminator is set, set the discriminator according to the union value.
// - When several fields are set and a discriminator is set, remove (== reset to zero value) all the values that do not match the discriminator.
// - When only one union value is set and it matches discriminator, just do nothing.
// - In other case, something is inconsistent or ambiguous: an error is thrown.
func Normalize(tree interface{}) error {
	return reflectwalk.Walk(tree, &normalizer{})
}

// Simplify allows removing the discriminator of all unions
// encountered while walking through the whole struct tree,
// but after normalizing them if necessary.
func Simplify(tree interface{}) error {
	return reflectwalk.Walk(tree, &simplifier{})
}
