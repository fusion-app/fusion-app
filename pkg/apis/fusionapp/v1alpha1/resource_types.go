package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ResourceSpec defines the desired state of Resource
// +k8s:openapi-gen=true
type ResourceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ResourceKind ResourceKind `json:"resourceKind"`
	Icon         string       `json:"icon,omitempty"`
	Description  string       `json:"description,omitempty"`
	AccessMode   ResourceAccessMode `json:"accessMode"`
	ProbeArgs    []string     `json:"probeArgs"`
}

type ResourceKind string

const (
	ResourceKindHuman   = "Human"
	ResourceKindService = "Service"
	ResourceKindEdge    = "Edge"
)

type ResourceAccessMode string

const (
	ResourceAccessModeExclusive = "Exclusive"
	ResourceAccessModeShared = "Shared"
)

// ResourcePhase defines all phase of dataset lifecycle.
type ResourcePhase string

const (
	// ResourcePhaseNotReady means some fields not set, should not create probe
	ResourcePhaseNotReady = "NotReady"

	// ResourcePhasePending means probe not ready
	ResourcePhasePending = "Pending"

	// ResourcePhaseSynchronous means probe is ready
	ResourcePhaseSynchronous = "Synchronous"

	// ResourcePhaseFailed means some pods of Resource have failed.
	ResourcePhaseFailed = "Failed"
)

// ResourceStatus defines the observed state of Resource
// +k8s:openapi-gen=true
type ResourceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Phase      ResourcePhase `json:"phase"`
	Bound      bool          `json:"bound"`
	CreateTime *metav1.Time  `json:"createTime,omitempty"`
	StartTime  *metav1.Time  `json:"startTime,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Resource is the Schema for the resources API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Resource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceSpec   `json:"spec,omitempty"`
	Status ResourceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ResourceList contains a list of Resource
type ResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Resource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Resource{}, &ResourceList{})
}
