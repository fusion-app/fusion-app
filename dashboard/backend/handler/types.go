package handler

import (
	fusionappv1alpha1  "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
)

const (
	defaultNamespace = "fusion-app-resources"
)

type Resource struct {
	UID          string      `json:"uid,omitempty"`
 	Namespace    string      `json:"namespace"`
	Kind         string      `json:"kind"`
	Phase        string      `json:"phase,omitempty"`
	Bound        bool        `json:"bound,omitempty"`
	Name         string      `json:"name"`
	AccessMode   string      `json:"accessMode"`
	Labels       map[string]string `json:"labels,omitempty"`
	Operation    fusionappv1alpha1.ResourceOperationSpec `json:"operation,omitempty"`
	Icon         string       `json:"icon,omitempty"`
	Description  string       `json:"description,omitempty"`
}

type ResourceList struct {
	Resources []Resource `json:"resources"`
}