package v1alpha1

import (
	"errors"
	"reflect"
)

type Identified interface {
	Id() (string, error)
}

func visitUnion(union interface{}, visitor interface{}) (err error) {
	visitorValue := reflect.ValueOf(visitor)
	unionValue := reflect.ValueOf(union)
	oneMemberPresent := false
	typeOfVisitor := visitorValue.Type()
	for i := 0; i < visitorValue.NumField(); i++ {
		unionMemberToRead := typeOfVisitor.Field(i).Name
		unionMember := unionValue.FieldByName(unionMemberToRead)
		if !unionMember.IsNil() {
			if oneMemberPresent {
				err = errors.New("Only one element should be set in union: " + unionValue.Type().Name())
				return
			}
			oneMemberPresent = true
			visitorFunction := visitorValue.Field(i)
			if visitorFunction.IsNil() {
				return
			}
			results := visitorFunction.Call([]reflect.Value{unionMember})
			if !results[0].IsNil() {
				err = results[0].Interface().(error)
			}
			return
		}
	}
	return
}

// +k8s:deepcopy-gen=false
type ComponentVisitor struct {
	Container  func(*ContainerComponent) error
	Plugin     func(*PluginComponent) error
	Volume     func(*VolumeComponent) error
	Kubernetes func(*KubernetesComponent) error
	Openshift  func(*OpenshiftComponent) error
	Custom     func(*CustomComponent) error
}

func (union Component) Visit(visitor ComponentVisitor) error {
	return visitUnion(union, visitor)
}

func (union Component) Id() (string, error) {
	id := ""
	err := union.Visit(ComponentVisitor{
		Container: func(container *ContainerComponent) error {
			id = container.Name
			return nil
		},
		Plugin: func(plugin *PluginComponent) error {
			if plugin.Name != "" {
				id = plugin.Name
				return nil
			}
			return plugin.ImportReference.ImportReferenceUnion.Visit(ImportReferenceUnionVisitor{
				Uri: func(uri string) error {
					id = uri
					return nil
				},
				Id: func(id string) error {
					id = plugin.Id
					if plugin.RegistryUrl != "" {
						id = plugin.RegistryUrl + "/" + id
					}
					return nil
				},
				Kubernetes: func(cr *KubernetesCustomResourceImportReference) error {
					id = cr.Name
					if cr.Namespace != "" {
						id = cr.Namespace + "/" + id
					}
					return nil
				},
			})
		},
		Kubernetes: func(k8s *KubernetesComponent) error {
			id = k8s.Name
			return nil
		},
		Openshift: func(os *OpenshiftComponent) error {
			id = os.Name
			return nil
		},
		Volume: func(vol *VolumeComponent) error {
			id = vol.Name
			return nil
		},
	})
	if err != nil {
		return id, err
	}
	return id, nil
}

// +k8s:deepcopy-gen=false
type ComponentOverrideVisitor struct {
	Container  func(*ContainerComponent) error
	Volume     func(*VolumeComponent) error
	Kubernetes func(*KubernetesComponent) error
	Openshift  func(*OpenshiftComponent) error
}

func (union ComponentOverride) Visit(visitor ComponentOverrideVisitor) error {
	return visitUnion(union, visitor)
}

func (union ComponentOverride) Id() (string, error) {
	id := ""
	err := union.Visit(ComponentOverrideVisitor{
		Container: func(container *ContainerComponent) error {
			id = container.Name
			return nil
		},
		Kubernetes: func(k8s *KubernetesComponent) error {
			id = k8s.Name
			return nil
		},
		Openshift: func(os *OpenshiftComponent) error {
			id = os.Name
			return nil
		},
		Volume: func(vol *VolumeComponent) error {
			id = vol.Name
			return nil
		},
	})
	if err != nil {
		return id, err
	}
	return id, nil
}

// +k8s:deepcopy-gen=false
type CommandVisitor struct {
	Exec         func(*ExecCommand) error
	VscodeTask   func(*VscodeConfigurationCommand) error
	VscodeLaunch func(*VscodeConfigurationCommand) error
	Composite    func(*CompositeCommand) error
	Custom       func(*CustomCommand) error
}

func (union Command) Visit(visitor CommandVisitor) error {
	return visitUnion(union, visitor)
}

// +k8s:deepcopy-gen=false
type ImportReferenceUnionVisitor struct {
	Uri        func(string) error
	Id         func(string) error
	Kubernetes func(*KubernetesCustomResourceImportReference) error
}

func (union ImportReferenceUnion) Visit(visitor ImportReferenceUnionVisitor) error {
	return visitUnion(union, visitor)
}

// +k8s:deepcopy-gen=false
type K8sLikeComponentLocationVisitor struct {
	Uri     func(string) error
	Inlined func(string) error
}

func (union K8sLikeComponentLocation) Visit(visitor K8sLikeComponentLocationVisitor) error {
	return visitUnion(union, visitor)
}

// +k8s:deepcopy-gen=false
type ProjectSourceVisitor struct {
	Git    func(*GitProjectSource) error
	Github func(*GithubProjectSource) error
	Zip    func(*ZipProjectSource) error
	Custom func(*CustomProjectSource) error
}

func (union ProjectSource) Visit(visitor ProjectSourceVisitor) error {
	return visitUnion(union, visitor)
}
