package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ResourceClaimSpec defines the desired state of ResourceClaim
// +k8s:openapi-gen=true
type ResourceClaimSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Name       string              `json:"name"`
	AccessMode ResourceAccessMode  `json:"accessMode"`
	Selector   []SelectorSpec      `json:"selector"`
}

type SelectorSpec struct {
	Key      string   `json:"key"`
	Value    string   `json:"value"`
	Operator Operator `json:"op"`
}

type Operator string

const (
	Eq   =   "Eq"
	Gt   =   "Gt"
	Lt   =   "Lt"
)

// ResourceClaimStatus defines the observed state of ResourceClaim
// +k8s:openapi-gen=true
type ResourceClaimStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ResourceClaim is the Schema for the resourceclaims API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type ResourceClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ResourceClaimSpec   `json:"spec,omitempty"`
	Status ResourceClaimStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ResourceClaimList contains a list of ResourceClaim
type ResourceClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ResourceClaim `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ResourceClaim{}, &ResourceClaimList{})
}
