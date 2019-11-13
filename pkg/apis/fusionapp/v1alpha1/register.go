// NOTE: Boilerplate only.  Ignore this file.

// Package v1alpha1 contains API Schema definitions for the fusion-app v1alpha1 API group
// +k8s:deepcopy-gen=package,register
// +groupName=fusionapp.io
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/runtime/scheme"
)

var (
	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}

	AddToScheme = SchemeBuilder.AddToScheme
)

const (
	// GroupName is the group name use in this package.
	GroupName = "fusionapp.io"
	// Kind is the kind name.
	Kind = "Resource"
	// GroupVersion is the version.
	GroupVersion = "v1alpha1"
	// Plural is the Plural for Workspace.
	Plural = "resources"
	// Singular is the singular for Workspace.
	Singular = "resource"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

	// SchemeGroupVersionKind is the GroupVersionKind of the resource.
	SchemeGroupVersionKind = SchemeGroupVersion.WithKind(Kind)
)

// Resource takes an unqualified resource and returns a Group-qualified GroupResource.
func GroupResource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

// addKnownTypes adds the set of types defined in this package to the supplied scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Resource{},
		&ResourceList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}
