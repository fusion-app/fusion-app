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
	RefApp       RefApp            `json:"refApp,omitempty"`
	RefResource  []RefResource     `json:"refResource,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
	ProbeImage   string            `json:"probeImage,omitempty"`
	ProbeArgs    []string          `json:"probeArgs,omitempty"`
	ProbeEnabled bool              `json:"probeEnabled"`
	AliasName    string            `json:"aliasName,omitempty"`
	RefResourceClaim []RefResourceClaim `json:"refResourceClaim,omitempty"`
}

type RefApp struct {
	UID    string    `json:"uid,omitempty"`
	Name   string    `json:"name"`
}

type RefResource struct {
	UID	        string    `json:"uid,omitempty"`
	Namespace	string    `json:"namespace,omitempty"`
	Kind	    string    `json:"kind,omitempty"`
	Name	    string    `json:"name"`
	AliasName   string   `json:"aliasName,omitempty"`
	Icon        string   `json:"icon,omitempty"`
	Description map[string]string `json:"description,omitempty"`
}

// FusionAppInstanceStatus defines the observed state of FusionAppInstance
// +k8s:openapi-gen=true
type FusionAppInstanceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	CreateTime *metav1.Time  `json:"createTime"`
	StartTime  *metav1.Time  `json:"startTime"`
	EndTime    *metav1.Time  `json:"endTime,omitempty"`
	UpdateTime *metav1.Time  `json:"updateTime"`
	Phase      FusionAppInstancePhase `json:"phase"`
	ProbePhase ProbePhase    `json:"probePhase"`
	ActionStatus []Action    `json:"actionStatus,omitempty"`
}


type Action struct {
	ActionID    string      `json:"actionID"`
	ActionName  string      `json:"actionName"`
	RefResource RefResource `json:"refResource,omitempty"`
	State       ActionState `json:"state"`
}

type ActionState string

type FusionAppInstancePhase string

const (
	FusionAppInstancePhaseNotReady FusionAppInstancePhase = "NotReady"
	FusionAppInstancePhaseReady FusionAppInstancePhase = "Ready"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FusionAppInstance is the Schema for the fusionappinstances API
// +genclient
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
