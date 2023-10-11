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
	"reflect"
)

/*

This file implements the `Union` methods on all the struct types
that are defined as Kubernetes unions.

The implementations here mainly delegate to generic implementation functions.
so, in the future, we should probably produce this file
by some code generation mechanism based on API source code parsing, especially
based on the `+union` comments in the API GO source code.

*/

// +k8s:deepcopy-gen=false
type ComponentVisitor struct {
	Container  func(*ContainerComponent) error
	Plugin     func(*PluginComponent) error
	Volume     func(*VolumeComponent) error
	Kubernetes func(*KubernetesComponent) error
	Openshift  func(*OpenshiftComponent) error
	Custom     func(*CustomComponent) error
}

var componentVisitorType reflect.Type = reflect.TypeOf(ComponentVisitor{})

func (union Component) Visit(visitor ComponentVisitor) error {
	return visitUnion(union, visitor)
}
func (union *Component) discriminator() *string {
	return (*string)(&union.ComponentType)
}
func (union *Component) Normalize() error {
	return normalizeUnion(union, componentVisitorType)
}
func (union *Component) Simplify() {
	simplifyUnion(union, componentVisitorType)
}

// +k8s:deepcopy-gen=false
type PluginComponentsOverrideVisitor struct {
	Container  func(*ContainerComponent) error
	Volume     func(*VolumeComponent) error
	Kubernetes func(*KubernetesComponent) error
	Openshift  func(*OpenshiftComponent) error
}

var pluginComponentsOverrideVisitorType reflect.Type = reflect.TypeOf(PluginComponentsOverrideVisitor{})

func (union PluginComponentsOverride) Visit(visitor PluginComponentsOverrideVisitor) error {
	return visitUnion(union, visitor)
}
func (union *PluginComponentsOverride) discriminator() *string {
	return (*string)(&union.ComponentType)
}
func (union *PluginComponentsOverride) Normalize() error {
	return normalizeUnion(union, pluginComponentsOverrideVisitorType)
}
func (union *PluginComponentsOverride) Simplify() {
	simplifyUnion(union, pluginComponentsOverrideVisitorType)
}

// +k8s:deepcopy-gen=false
type CommandVisitor struct {
	Apply        func(*ApplyCommand) error
	Exec         func(*ExecCommand) error
	VscodeTask   func(*VscodeConfigurationCommand) error
	VscodeLaunch func(*VscodeConfigurationCommand) error
	Composite    func(*CompositeCommand) error
	Custom       func(*CustomCommand) error
}

var commandVisitorType reflect.Type = reflect.TypeOf(CommandVisitor{})

func (union Command) Visit(visitor CommandVisitor) error {
	return visitUnion(union, visitor)
}
func (union *Command) discriminator() *string {
	return (*string)(&union.CommandType)
}
func (union *Command) Normalize() error {
	return normalizeUnion(union, commandVisitorType)
}
func (union *Command) Simplify() {
	simplifyUnion(union, commandVisitorType)
}

// +k8s:deepcopy-gen=false
type ImportReferenceUnionVisitor struct {
	Uri        func(string) error
	Id         func(string) error
	Kubernetes func(*KubernetesCustomResourceImportReference) error
}

var importReferenceUnionVisitorType reflect.Type = reflect.TypeOf(ImportReferenceUnionVisitor{})

func (union ImportReferenceUnion) Visit(visitor ImportReferenceUnionVisitor) error {
	return visitUnion(union, visitor)
}
func (union *ImportReferenceUnion) discriminator() *string {
	return (*string)(&union.ImportReferenceType)
}
func (union *ImportReferenceUnion) Normalize() error {
	return normalizeUnion(union, importReferenceUnionVisitorType)
}
func (union *ImportReferenceUnion) Simplify() {
	simplifyUnion(union, importReferenceUnionVisitorType)
}

// +k8s:deepcopy-gen=false
type K8sLikeComponentLocationVisitor struct {
	Uri     func(string) error
	Inlined func(string) error
}

var k8sLikeComponentLocationVisitorType reflect.Type = reflect.TypeOf(K8sLikeComponentLocationVisitor{})

func (union K8sLikeComponentLocation) Visit(visitor K8sLikeComponentLocationVisitor) error {
	return visitUnion(union, visitor)
}
func (union *K8sLikeComponentLocation) discriminator() *string {
	return (*string)(&union.LocationType)
}
func (union *K8sLikeComponentLocation) Normalize() error {
	return normalizeUnion(union, k8sLikeComponentLocationVisitorType)
}
func (union *K8sLikeComponentLocation) Simplify() {
	simplifyUnion(union, k8sLikeComponentLocationVisitorType)
}

// +k8s:deepcopy-gen=false
type VscodeConfigurationCommandLocationVisitor struct {
	Uri     func(string) error
	Inlined func(string) error
}

var vscodeConfigurationCommandLocationVisitorType reflect.Type = reflect.TypeOf(VscodeConfigurationCommandLocationVisitor{})

func (union VscodeConfigurationCommandLocation) Visit(visitor VscodeConfigurationCommandLocation) error {
	return visitUnion(union, visitor)
}
func (union *VscodeConfigurationCommandLocation) discriminator() *string {
	return (*string)(&union.LocationType)
}
func (union *VscodeConfigurationCommandLocation) Normalize() error {
	return normalizeUnion(union, vscodeConfigurationCommandLocationVisitorType)
}
func (union *VscodeConfigurationCommandLocation) Simplify() {
	simplifyUnion(union, vscodeConfigurationCommandLocationVisitorType)
}

// +k8s:deepcopy-gen=false
type ProjectSourceVisitor struct {
	Git    func(*GitProjectSource) error
	Github func(*GithubProjectSource) error
	Zip    func(*ZipProjectSource) error
	Custom func(*CustomProjectSource) error
}

var projectSourceVisitorType reflect.Type = reflect.TypeOf(ProjectSourceVisitor{})

func (union ProjectSource) Visit(visitor ProjectSourceVisitor) error {
	return visitUnion(union, visitor)
}
func (union *ProjectSource) discriminator() *string {
	return (*string)(&union.SourceType)
}
func (union *ProjectSource) Normalize() error {
	return normalizeUnion(union, projectSourceVisitorType)
}
func (union *ProjectSource) Simplify() {
	simplifyUnion(union, projectSourceVisitorType)
}
