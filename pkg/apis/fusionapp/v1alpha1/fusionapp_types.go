package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// FusionAppSpec defines the desired state of FusionApp
// +k8s:openapi-gen=true
type FusionAppSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	ResourceClaim   []ResourceClaimSpec   `json:"resourceClaim,omitempty"`
	AliasName       string                `json:"aliasName"`
	Icon            string                `json:"icon,omitempty"`
	Description     map[string]string     `json:"description,omitempty"`
	Labels          map[string]string     `json:"labels,omitempty"`
	ProbeImage      string                `json:"probeImage,omitempty"`
	ProbeArgs       []string              `json:"probeArgs"`
}


// FusionAppStatus defines the observed state of FusionApp
// +k8s:openapi-gen=true
type FusionAppStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FusionApp is the Schema for the fusionapps API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type FusionApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FusionAppSpec   `json:"spec,omitempty"`
	Status FusionAppStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FusionAppList contains a list of FusionApp
type FusionAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []FusionApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FusionApp{}, &FusionAppList{})
}
