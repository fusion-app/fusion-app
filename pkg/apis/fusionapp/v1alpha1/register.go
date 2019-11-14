// NOTE: Boilerplate only.  Ignore this file.

// Package v1alpha1 contains API Schema definitions for the fusion-app v1alpha1 API group
// +k8s:deepcopy-gen=package,register
// +groupName=fusionapp.io
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/runtime/scheme"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}

	AddToScheme = SchemeBuilder.AddToScheme
)

const (
	// GroupName is the group name use in this package.
	GroupName = "fusionapp.io"
	// GroupVersion is the version.
	GroupVersion = "v1alpha1"
)

// Resource takes an unqualified resource and returns a Group-qualified GroupResource.
func GroupResource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

