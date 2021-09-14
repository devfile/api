package v1alpha2

import (
	"reflect"
)

var commandUnion reflect.Type = reflect.TypeOf(CommandUnionVisitor{})

func (union CommandUnion) Visit(visitor CommandUnionVisitor) error {
	return visitUnion(union, visitor)
}
func (union *CommandUnion) discriminator() *string {
	return (*string)(&union.CommandType)
}
func (union *CommandUnion) Normalize() error {
	return normalizeUnion(union, commandUnion)
}
func (union *CommandUnion) Simplify() {
	simplifyUnion(union, commandUnion)
}

// +k8s:deepcopy-gen=false
type CommandUnionVisitor struct {
	Exec      func(*ExecCommand) error
	Apply     func(*ApplyCommand) error
	Composite func(*CompositeCommand) error
	Custom    func(*CustomCommand) error
}

var imageUnion reflect.Type = reflect.TypeOf(ImageUnionVisitor{})

func (union ImageUnion) Visit(visitor ImageUnionVisitor) error {
	return visitUnion(union, visitor)
}
func (union *ImageUnion) discriminator() *string {
	return (*string)(&union.ImageType)
}
func (union *ImageUnion) Normalize() error {
	return normalizeUnion(union, imageUnion)
}
func (union *ImageUnion) Simplify() {
	simplifyUnion(union, imageUnion)
}

// +k8s:deepcopy-gen=false
type ImageUnionVisitor struct {
	Dockerfile func(*DockerfileImage) error
}

var dockerfileLocation reflect.Type = reflect.TypeOf(DockerfileLocationVisitor{})

func (union DockerfileLocation) Visit(visitor DockerfileLocationVisitor) error {
	return visitUnion(union, visitor)
}
func (union *DockerfileLocation) discriminator() *string {
	return (*string)(&union.LocationType)
}
func (union *DockerfileLocation) Normalize() error {
	return normalizeUnion(union, dockerfileLocation)
}
func (union *DockerfileLocation) Simplify() {
	simplifyUnion(union, dockerfileLocation)
}

// +k8s:deepcopy-gen=false
type DockerfileLocationVisitor struct {
	Uri      func(string) error
	Registry func(*DockerfileDevfileRegistrySource) error
	Git      func(*DockerfileGitProjectSource) error
}

var k8sLikeComponentLocation reflect.Type = reflect.TypeOf(K8sLikeComponentLocationVisitor{})

func (union K8sLikeComponentLocation) Visit(visitor K8sLikeComponentLocationVisitor) error {
	return visitUnion(union, visitor)
}
func (union *K8sLikeComponentLocation) discriminator() *string {
	return (*string)(&union.LocationType)
}
func (union *K8sLikeComponentLocation) Normalize() error {
	return normalizeUnion(union, k8sLikeComponentLocation)
}
func (union *K8sLikeComponentLocation) Simplify() {
	simplifyUnion(union, k8sLikeComponentLocation)
}

// +k8s:deepcopy-gen=false
type K8sLikeComponentLocationVisitor struct {
	Uri     func(string) error
	Inlined func(string) error
}

var componentUnion reflect.Type = reflect.TypeOf(ComponentUnionVisitor{})

func (union ComponentUnion) Visit(visitor ComponentUnionVisitor) error {
	return visitUnion(union, visitor)
}
func (union *ComponentUnion) discriminator() *string {
	return (*string)(&union.ComponentType)
}
func (union *ComponentUnion) Normalize() error {
	return normalizeUnion(union, componentUnion)
}
func (union *ComponentUnion) Simplify() {
	simplifyUnion(union, componentUnion)
}

// +k8s:deepcopy-gen=false
type ComponentUnionVisitor struct {
	Container  func(*ContainerComponent) error
	Kubernetes func(*KubernetesComponent) error
	Openshift  func(*OpenshiftComponent) error
	Volume     func(*VolumeComponent) error
	Image      func(*ImageComponent) error
	Plugin     func(*PluginComponent) error
	Custom     func(*CustomComponent) error
}

var importReferenceUnion reflect.Type = reflect.TypeOf(ImportReferenceUnionVisitor{})

func (union ImportReferenceUnion) Visit(visitor ImportReferenceUnionVisitor) error {
	return visitUnion(union, visitor)
}
func (union *ImportReferenceUnion) discriminator() *string {
	return (*string)(&union.ImportReferenceType)
}
func (union *ImportReferenceUnion) Normalize() error {
	return normalizeUnion(union, importReferenceUnion)
}
func (union *ImportReferenceUnion) Simplify() {
	simplifyUnion(union, importReferenceUnion)
}

// +k8s:deepcopy-gen=false
type ImportReferenceUnionVisitor struct {
	Uri        func(string) error
	Id         func(string) error
	Kubernetes func(*KubernetesCustomResourceImportReference) error
}

var projectSource reflect.Type = reflect.TypeOf(ProjectSourceVisitor{})

func (union ProjectSource) Visit(visitor ProjectSourceVisitor) error {
	return visitUnion(union, visitor)
}
func (union *ProjectSource) discriminator() *string {
	return (*string)(&union.SourceType)
}
func (union *ProjectSource) Normalize() error {
	return normalizeUnion(union, projectSource)
}
func (union *ProjectSource) Simplify() {
	simplifyUnion(union, projectSource)
}

// +k8s:deepcopy-gen=false
type ProjectSourceVisitor struct {
	Git    func(*GitProjectSource) error
	Zip    func(*ZipProjectSource) error
	Custom func(*CustomProjectSource) error
}

var componentUnionParentOverride reflect.Type = reflect.TypeOf(ComponentUnionParentOverrideVisitor{})

func (union ComponentUnionParentOverride) Visit(visitor ComponentUnionParentOverrideVisitor) error {
	return visitUnion(union, visitor)
}
func (union *ComponentUnionParentOverride) discriminator() *string {
	return (*string)(&union.ComponentType)
}
func (union *ComponentUnionParentOverride) Normalize() error {
	return normalizeUnion(union, componentUnionParentOverride)
}
func (union *ComponentUnionParentOverride) Simplify() {
	simplifyUnion(union, componentUnionParentOverride)
}

// +k8s:deepcopy-gen=false
type ComponentUnionParentOverrideVisitor struct {
	Container  func(*ContainerComponentParentOverride) error
	Kubernetes func(*KubernetesComponentParentOverride) error
	Openshift  func(*OpenshiftComponentParentOverride) error
	Volume     func(*VolumeComponentParentOverride) error
	Image      func(*ImageComponentParentOverride) error
	Plugin     func(*PluginComponentParentOverride) error
}

var projectSourceParentOverride reflect.Type = reflect.TypeOf(ProjectSourceParentOverrideVisitor{})

func (union ProjectSourceParentOverride) Visit(visitor ProjectSourceParentOverrideVisitor) error {
	return visitUnion(union, visitor)
}
func (union *ProjectSourceParentOverride) discriminator() *string {
	return (*string)(&union.SourceType)
}
func (union *ProjectSourceParentOverride) Normalize() error {
	return normalizeUnion(union, projectSourceParentOverride)
}
func (union *ProjectSourceParentOverride) Simplify() {
	simplifyUnion(union, projectSourceParentOverride)
}

// +k8s:deepcopy-gen=false
type ProjectSourceParentOverrideVisitor struct {
	Git func(*GitProjectSourceParentOverride) error
	Zip func(*ZipProjectSourceParentOverride) error
}

var commandUnionParentOverride reflect.Type = reflect.TypeOf(CommandUnionParentOverrideVisitor{})

func (union CommandUnionParentOverride) Visit(visitor CommandUnionParentOverrideVisitor) error {
	return visitUnion(union, visitor)
}
func (union *CommandUnionParentOverride) discriminator() *string {
	return (*string)(&union.CommandType)
}
func (union *CommandUnionParentOverride) Normalize() error {
	return normalizeUnion(union, commandUnionParentOverride)
}
func (union *CommandUnionParentOverride) Simplify() {
	simplifyUnion(union, commandUnionParentOverride)
}

// +k8s:deepcopy-gen=false
type CommandUnionParentOverrideVisitor struct {
	Exec      func(*ExecCommandParentOverride) error
	Apply     func(*ApplyCommandParentOverride) error
	Composite func(*CompositeCommandParentOverride) error
}

var k8sLikeComponentLocationParentOverride reflect.Type = reflect.TypeOf(K8sLikeComponentLocationParentOverrideVisitor{})

func (union K8sLikeComponentLocationParentOverride) Visit(visitor K8sLikeComponentLocationParentOverrideVisitor) error {
	return visitUnion(union, visitor)
}
func (union *K8sLikeComponentLocationParentOverride) discriminator() *string {
	return (*string)(&union.LocationType)
}
func (union *K8sLikeComponentLocationParentOverride) Normalize() error {
	return normalizeUnion(union, k8sLikeComponentLocationParentOverride)
}
func (union *K8sLikeComponentLocationParentOverride) Simplify() {
	simplifyUnion(union, k8sLikeComponentLocationParentOverride)
}

// +k8s:deepcopy-gen=false
type K8sLikeComponentLocationParentOverrideVisitor struct {
	Uri     func(string) error
	Inlined func(string) error
}

var imageUnionParentOverride reflect.Type = reflect.TypeOf(ImageUnionParentOverrideVisitor{})

func (union ImageUnionParentOverride) Visit(visitor ImageUnionParentOverrideVisitor) error {
	return visitUnion(union, visitor)
}
func (union *ImageUnionParentOverride) discriminator() *string {
	return (*string)(&union.ImageType)
}
func (union *ImageUnionParentOverride) Normalize() error {
	return normalizeUnion(union, imageUnionParentOverride)
}
func (union *ImageUnionParentOverride) Simplify() {
	simplifyUnion(union, imageUnionParentOverride)
}

// +k8s:deepcopy-gen=false
type ImageUnionParentOverrideVisitor struct {
	Dockerfile func(*DockerfileImageParentOverride) error
}

var importReferenceUnionParentOverride reflect.Type = reflect.TypeOf(ImportReferenceUnionParentOverrideVisitor{})

func (union ImportReferenceUnionParentOverride) Visit(visitor ImportReferenceUnionParentOverrideVisitor) error {
	return visitUnion(union, visitor)
}
func (union *ImportReferenceUnionParentOverride) discriminator() *string {
	return (*string)(&union.ImportReferenceType)
}
func (union *ImportReferenceUnionParentOverride) Normalize() error {
	return normalizeUnion(union, importReferenceUnionParentOverride)
}
func (union *ImportReferenceUnionParentOverride) Simplify() {
	simplifyUnion(union, importReferenceUnionParentOverride)
}

// +k8s:deepcopy-gen=false
type ImportReferenceUnionParentOverrideVisitor struct {
	Uri        func(string) error
	Id         func(string) error
	Kubernetes func(*KubernetesCustomResourceImportReferenceParentOverride) error
}

var componentUnionPluginOverrideParentOverride reflect.Type = reflect.TypeOf(ComponentUnionPluginOverrideParentOverrideVisitor{})

func (union ComponentUnionPluginOverrideParentOverride) Visit(visitor ComponentUnionPluginOverrideParentOverrideVisitor) error {
	return visitUnion(union, visitor)
}
func (union *ComponentUnionPluginOverrideParentOverride) discriminator() *string {
	return (*string)(&union.ComponentType)
}
func (union *ComponentUnionPluginOverrideParentOverride) Normalize() error {
	return normalizeUnion(union, componentUnionPluginOverrideParentOverride)
}
func (union *ComponentUnionPluginOverrideParentOverride) Simplify() {
	simplifyUnion(union, componentUnionPluginOverrideParentOverride)
}

// +k8s:deepcopy-gen=false
type ComponentUnionPluginOverrideParentOverrideVisitor struct {
	Container  func(*ContainerComponentPluginOverrideParentOverride) error
	Kubernetes func(*KubernetesComponentPluginOverrideParentOverride) error
	Openshift  func(*OpenshiftComponentPluginOverrideParentOverride) error
	Volume     func(*VolumeComponentPluginOverrideParentOverride) error
	Image      func(*ImageComponentPluginOverrideParentOverride) error
}

var commandUnionPluginOverrideParentOverride reflect.Type = reflect.TypeOf(CommandUnionPluginOverrideParentOverrideVisitor{})

func (union CommandUnionPluginOverrideParentOverride) Visit(visitor CommandUnionPluginOverrideParentOverrideVisitor) error {
	return visitUnion(union, visitor)
}
func (union *CommandUnionPluginOverrideParentOverride) discriminator() *string {
	return (*string)(&union.CommandType)
}
func (union *CommandUnionPluginOverrideParentOverride) Normalize() error {
	return normalizeUnion(union, commandUnionPluginOverrideParentOverride)
}
func (union *CommandUnionPluginOverrideParentOverride) Simplify() {
	simplifyUnion(union, commandUnionPluginOverrideParentOverride)
}

// +k8s:deepcopy-gen=false
type CommandUnionPluginOverrideParentOverrideVisitor struct {
	Exec      func(*ExecCommandPluginOverrideParentOverride) error
	Apply     func(*ApplyCommandPluginOverrideParentOverride) error
	Composite func(*CompositeCommandPluginOverrideParentOverride) error
}

var dockerfileLocationParentOverride reflect.Type = reflect.TypeOf(DockerfileLocationParentOverrideVisitor{})

func (union DockerfileLocationParentOverride) Visit(visitor DockerfileLocationParentOverrideVisitor) error {
	return visitUnion(union, visitor)
}
func (union *DockerfileLocationParentOverride) discriminator() *string {
	return (*string)(&union.LocationType)
}
func (union *DockerfileLocationParentOverride) Normalize() error {
	return normalizeUnion(union, dockerfileLocationParentOverride)
}
func (union *DockerfileLocationParentOverride) Simplify() {
	simplifyUnion(union, dockerfileLocationParentOverride)
}

// +k8s:deepcopy-gen=false
type DockerfileLocationParentOverrideVisitor struct {
	Uri      func(string) error
	Registry func(*DockerfileDevfileRegistrySourceParentOverride) error
	Git      func(*DockerfileGitProjectSourceParentOverride) error
}

var k8sLikeComponentLocationPluginOverrideParentOverride reflect.Type = reflect.TypeOf(K8sLikeComponentLocationPluginOverrideParentOverrideVisitor{})

func (union K8sLikeComponentLocationPluginOverrideParentOverride) Visit(visitor K8sLikeComponentLocationPluginOverrideParentOverrideVisitor) error {
	return visitUnion(union, visitor)
}
func (union *K8sLikeComponentLocationPluginOverrideParentOverride) discriminator() *string {
	return (*string)(&union.LocationType)
}
func (union *K8sLikeComponentLocationPluginOverrideParentOverride) Normalize() error {
	return normalizeUnion(union, k8sLikeComponentLocationPluginOverrideParentOverride)
}
func (union *K8sLikeComponentLocationPluginOverrideParentOverride) Simplify() {
	simplifyUnion(union, k8sLikeComponentLocationPluginOverrideParentOverride)
}

// +k8s:deepcopy-gen=false
type K8sLikeComponentLocationPluginOverrideParentOverrideVisitor struct {
	Uri     func(string) error
	Inlined func(string) error
}

var imageUnionPluginOverrideParentOverride reflect.Type = reflect.TypeOf(ImageUnionPluginOverrideParentOverrideVisitor{})

func (union ImageUnionPluginOverrideParentOverride) Visit(visitor ImageUnionPluginOverrideParentOverrideVisitor) error {
	return visitUnion(union, visitor)
}
func (union *ImageUnionPluginOverrideParentOverride) discriminator() *string {
	return (*string)(&union.ImageType)
}
func (union *ImageUnionPluginOverrideParentOverride) Normalize() error {
	return normalizeUnion(union, imageUnionPluginOverrideParentOverride)
}
func (union *ImageUnionPluginOverrideParentOverride) Simplify() {
	simplifyUnion(union, imageUnionPluginOverrideParentOverride)
}

// +k8s:deepcopy-gen=false
type ImageUnionPluginOverrideParentOverrideVisitor struct {
	Dockerfile func(*DockerfileImagePluginOverrideParentOverride) error
}

var dockerfileLocationPluginOverrideParentOverride reflect.Type = reflect.TypeOf(DockerfileLocationPluginOverrideParentOverrideVisitor{})

func (union DockerfileLocationPluginOverrideParentOverride) Visit(visitor DockerfileLocationPluginOverrideParentOverrideVisitor) error {
	return visitUnion(union, visitor)
}
func (union *DockerfileLocationPluginOverrideParentOverride) discriminator() *string {
	return (*string)(&union.LocationType)
}
func (union *DockerfileLocationPluginOverrideParentOverride) Normalize() error {
	return normalizeUnion(union, dockerfileLocationPluginOverrideParentOverride)
}
func (union *DockerfileLocationPluginOverrideParentOverride) Simplify() {
	simplifyUnion(union, dockerfileLocationPluginOverrideParentOverride)
}

// +k8s:deepcopy-gen=false
type DockerfileLocationPluginOverrideParentOverrideVisitor struct {
	Uri      func(string) error
	Registry func(*DockerfileDevfileRegistrySourcePluginOverrideParentOverride) error
	Git      func(*DockerfileGitProjectSourcePluginOverrideParentOverride) error
}

var componentUnionPluginOverride reflect.Type = reflect.TypeOf(ComponentUnionPluginOverrideVisitor{})

func (union ComponentUnionPluginOverride) Visit(visitor ComponentUnionPluginOverrideVisitor) error {
	return visitUnion(union, visitor)
}
func (union *ComponentUnionPluginOverride) discriminator() *string {
	return (*string)(&union.ComponentType)
}
func (union *ComponentUnionPluginOverride) Normalize() error {
	return normalizeUnion(union, componentUnionPluginOverride)
}
func (union *ComponentUnionPluginOverride) Simplify() {
	simplifyUnion(union, componentUnionPluginOverride)
}

// +k8s:deepcopy-gen=false
type ComponentUnionPluginOverrideVisitor struct {
	Container  func(*ContainerComponentPluginOverride) error
	Kubernetes func(*KubernetesComponentPluginOverride) error
	Openshift  func(*OpenshiftComponentPluginOverride) error
	Volume     func(*VolumeComponentPluginOverride) error
	Image      func(*ImageComponentPluginOverride) error
}

var commandUnionPluginOverride reflect.Type = reflect.TypeOf(CommandUnionPluginOverrideVisitor{})

func (union CommandUnionPluginOverride) Visit(visitor CommandUnionPluginOverrideVisitor) error {
	return visitUnion(union, visitor)
}
func (union *CommandUnionPluginOverride) discriminator() *string {
	return (*string)(&union.CommandType)
}
func (union *CommandUnionPluginOverride) Normalize() error {
	return normalizeUnion(union, commandUnionPluginOverride)
}
func (union *CommandUnionPluginOverride) Simplify() {
	simplifyUnion(union, commandUnionPluginOverride)
}

// +k8s:deepcopy-gen=false
type CommandUnionPluginOverrideVisitor struct {
	Exec      func(*ExecCommandPluginOverride) error
	Apply     func(*ApplyCommandPluginOverride) error
	Composite func(*CompositeCommandPluginOverride) error
}

var k8sLikeComponentLocationPluginOverride reflect.Type = reflect.TypeOf(K8sLikeComponentLocationPluginOverrideVisitor{})

func (union K8sLikeComponentLocationPluginOverride) Visit(visitor K8sLikeComponentLocationPluginOverrideVisitor) error {
	return visitUnion(union, visitor)
}
func (union *K8sLikeComponentLocationPluginOverride) discriminator() *string {
	return (*string)(&union.LocationType)
}
func (union *K8sLikeComponentLocationPluginOverride) Normalize() error {
	return normalizeUnion(union, k8sLikeComponentLocationPluginOverride)
}
func (union *K8sLikeComponentLocationPluginOverride) Simplify() {
	simplifyUnion(union, k8sLikeComponentLocationPluginOverride)
}

// +k8s:deepcopy-gen=false
type K8sLikeComponentLocationPluginOverrideVisitor struct {
	Uri     func(string) error
	Inlined func(string) error
}

var imageUnionPluginOverride reflect.Type = reflect.TypeOf(ImageUnionPluginOverrideVisitor{})

func (union ImageUnionPluginOverride) Visit(visitor ImageUnionPluginOverrideVisitor) error {
	return visitUnion(union, visitor)
}
func (union *ImageUnionPluginOverride) discriminator() *string {
	return (*string)(&union.ImageType)
}
func (union *ImageUnionPluginOverride) Normalize() error {
	return normalizeUnion(union, imageUnionPluginOverride)
}
func (union *ImageUnionPluginOverride) Simplify() {
	simplifyUnion(union, imageUnionPluginOverride)
}

// +k8s:deepcopy-gen=false
type ImageUnionPluginOverrideVisitor struct {
	Dockerfile func(*DockerfileImagePluginOverride) error
}

var dockerfileLocationPluginOverride reflect.Type = reflect.TypeOf(DockerfileLocationPluginOverrideVisitor{})

func (union DockerfileLocationPluginOverride) Visit(visitor DockerfileLocationPluginOverrideVisitor) error {
	return visitUnion(union, visitor)
}
func (union *DockerfileLocationPluginOverride) discriminator() *string {
	return (*string)(&union.LocationType)
}
func (union *DockerfileLocationPluginOverride) Normalize() error {
	return normalizeUnion(union, dockerfileLocationPluginOverride)
}
func (union *DockerfileLocationPluginOverride) Simplify() {
	simplifyUnion(union, dockerfileLocationPluginOverride)
}

// +k8s:deepcopy-gen=false
type DockerfileLocationPluginOverrideVisitor struct {
	Uri      func(string) error
	Registry func(*DockerfileDevfileRegistrySourcePluginOverride) error
	Git      func(*DockerfileGitProjectSourcePluginOverride) error
}
