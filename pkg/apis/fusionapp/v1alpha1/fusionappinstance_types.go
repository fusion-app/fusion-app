package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FusionAppInstanceSpec defines the desired state of FusionAppInstance
// +k8s:openapi-gen=true
type FusionAppInstanceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	RefApp      RefApp            `json:"refApp,omitempty"`
	RefResource []AppRefResource  `json:"refResource,omitempty"`
}

type RefApp struct {
	UID    string    `json:"uid,omitempty"`
	Name   string    `json:"name"`
}

type AppRefResource struct {
	UID	        string    `json:"uid,omitempty"`
	Namespace	string    `json:"namespace,omitempty"`
	Kind	    string    `json:"kind"`
	Name	    string    `json:"name"`
}

// FusionAppInstanceStatus defines the observed state of FusionAppInstance
// +k8s:openapi-gen=true
type FusionAppInstanceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	CreateTime *metav1.Time  `json:"createTime,omitempty"`
	StartTime  *metav1.Time  `json:"startTime,omitempty"`
	EndTime    *metav1.Time  `json:"endTime,omitempty"`
	UpdateTime *metav1.Time  `json:"updateTime,omitempty"`
	Phase      FusionAppInstancePhase `json:"phase"`
}

type FusionAppInstancePhase string
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FusionAppInstance is the Schema for the fusionappinstances API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type FusionAppInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FusionAppInstanceSpec   `json:"spec,omitempty"`
	Status FusionAppInstanceStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FusionAppInstanceList contains a list of FusionAppInstance
type FusionAppInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FusionAppInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FusionAppInstance{}, &FusionAppInstanceList{})
}
