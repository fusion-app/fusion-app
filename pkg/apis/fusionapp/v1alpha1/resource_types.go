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
	Labels       map[string]string     `json:"labels"`
	ResourceKind ResourceKind          `json:"resourceKind"`
	Icon         string                `json:"icon,omitempty"`
	Description  map[string]string     `json:"description,omitempty"`
	AccessMode   ResourceAccessMode    `json:"accessMode"`
	Operation    []ResourceOperationSpec `json:"operation,omitempty"`
	ProbeImage   string                `json:"probeImage,omitempty"`
	ProbeArgs    []string              `json:"probeArgs"`
	ProbeEnabled bool                  `json:"probeEnabled"`
	AliasName    string                `json:"aliasName,omitempty"`
	RefResourceClaim []RefResourceClaim `json:"refResourceClaim,omitempty"`
}

type RefResourceClaim struct {
	UID          string                `json:"uid"`
	Name         string                `json:"name"`
	Namespace    string                `json:"namespace"`
} 

type ResourceOperationSpec struct {
	Name        string         `json:"name"`
	Price       float64        `json:"price"`
	Description string         `json:"description"`
	HTTPAction  HTTPActionSpec `json:"httpAction,omitempty"`
}

type HTTPActionSpec struct {
	Action  string            `json:"action"`
	URL     string            `json:"url"`
	Query   map[string]string `json:"query"`
	Headers map[string]string `json:"headers"`
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
	ResourceAccessModeShared    = "Shared"
)

// ResourcePhase defines all phase of dataset lifecycle.
type ResourcePhase string

const (
	// ResourcePhasePending means Resource not ready
	ResourcePhasePending   = "Pending"

	// ResourcePhaseSynchronous means Resource is ready
	ResourcePhaseRunning   = "Running"

	// ResourcePhaseFailed means some pods of Resource have failed.
	ResourcePhaseFailed    = "Failed"
)

// ResourcePhase defines all phase of dataset lifecycle.
type ProbePhase string

const (
	// ProbePhaseNotReady means probe not ready
	ProbePhaseNotReady       = "NotReady"

	// ProbePhasePending means probe pod not started
	ProbePhasePending        = "Pending"

	// ProbePhaseSynchronous means probe is ready
	ProbePhaseSynchronous    = "Synchronous"

	// ProbePhaseFailed means the probe pod has failed.
	ProbePhaseFailed         = "Failed"
)

// ResourceStatus defines the observed state of Resource
// +k8s:openapi-gen=true
type ResourceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Phase      ResourcePhase `json:"phase"`
	ProbePhase ProbePhase    `json:"probePhase,omitempty"`
	Bound      bool          `json:"bound"`
	CreateTime *metav1.Time  `json:"createTime,omitempty"`
	StartTime  *metav1.Time  `json:"startTime,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Resource is the Schema for the resources API
// +genclient
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
