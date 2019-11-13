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
	ResourceKind ResourceKind          `json:"resource_kind"`
	Operation    ResourceOperationSpec `json:"operation,omitempty"`
	ProbeArgs    []string              `json:"probe_args"`
}

type ResourceOperationSpec struct {
	Name       string         `json:"name"`
	Price      float64        `json:"price"`
	HTTPAction HTTPActionSpec `json:"http_action,omitempty"`
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

// ResourcePhase defines all phase of dataset lifecycle.
type ResourcePhase string

const (
	// ResourcePhaseNotready means some fields not set, should not create probe
	ResourcePhaseNotready = "Notready"

	// ResourcePhasePending means probe not ready
	ResourcePhasePending = "Pending"

	// ResourcePhaseSynchronous means probe is ready
	ResourcePhaseSynchronous = "Synchronous"

	// ResourcePhaseFailed means some pods of dataset have failed.
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
	CreateTime *metav1.Time  `json:"create_time,omitempty"`
	StartTime  *metav1.Time  `json:"start_time,omitempty"`
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
