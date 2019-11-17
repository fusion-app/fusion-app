package handler

import (
	fusionappv1alpha1  "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
)

const (
	defaultNamespace = "fusion-app-resources"
)

type Resource struct {
	UID          string      `json:"uid,omitempty"`
 	Namespace    string      `json:"namespace,omitempty"`
	Kind         string      `json:"kind"`
	ProbeEnabled bool        `json:"probeEnabled,omitempty"`
	Phase        string      `json:"phase,omitempty"`
	Bound        bool        `json:"bound,omitempty"`
	Name         string      `json:"name"`
	AliasName    string      `json:"aliasName,omitempty"`
	AccessMode   string      `json:"accessMode"`
	Labels       map[string]string `json:"labels,omitempty"`
	Operation    []fusionappv1alpha1.ResourceOperationSpec `json:"operation,omitempty"`
	Icon         string       `json:"icon,omitempty"`
	Description  string       `json:"description,omitempty"`
	ProbeArgs    []string     `json:"probeArgs"`
	ProbeImage   string       `json:"probeImage"`
}

type ResourceSpec struct {
	ProbeEnabled bool        `json:"probeEnabled"`
	ProbeImage   string      `json:"probeImage,omitempty"`
	AliasName    string      `json:"aliasName,omitempty"`
	AccessMode   fusionappv1alpha1.ResourceAccessMode    `json:"accessMode,omitempty"`
	Labels       map[string]string `json:"labels,omitempty"`
	Operation    []fusionappv1alpha1.ResourceOperationSpec `json:"operation,omitempty"`
	ProbeArgs    []string     `json:"probeArgs,omitempty"`
	Icon         string       `json:"icon,omitempty"`
	Description  string       `json:"description,omitempty"`
}

type AppRefResource struct {
	UID	        string   `json:"uid,omitempty"`
	Namespace	string   `json:"namespace,omitempty"`
	Kind	    string   `json:"kind,omitempty"`
	Name	    string   `json:"name,omitempty"`
}

type ResourceAPICreateBody struct {
	ResourceSpec Resource `json:"resourceSpec"`
}

type ResourceAPIPutBody struct {
	AppRefResource    AppRefResource   `json:"refResource"`
	ResourceSpec      ResourceSpec     `json:"resourceSpec"`
}

type ResourceAPIQueryBody struct {
	Kind       string   `json:"kind,omitempty"`
	Phase      string   `json:"phase,omitempty"`
	RefResource AppRefResource `json:"refResource,omitempty"`
	LabelSelector []fusionappv1alpha1.SelectorSpec `json:"labelSelector,omitempty"`
}

type ResourceAPIBindBody struct {
	RefResource    AppRefResource   `json:"refResource"`
	RefAppInstance RefAppInstance   `json:"refAppInstance"`
}

type RefAppInstance struct {
	UID	        string  `json:"uid,omitempty"`
	Namespace	string  `json:"namespace,omitempty"`
	Kind	    string  `json:"kind,omitempty"`
	Name	    string	`json:"name"`
}

type RefApp struct {
	UID         string  `json:"uid"`
	Name        string  `json:"name"`
}

type AppInstance struct {
	UID	        string              `json:"uid"`
	Namespace	string              `json:"namespace"`
	Name	    string              `json:"name"`
	RefApp      RefApp              `json:"refApp"`
	RefResource	[]AppRefResource    `json:"refResource"`
	Status	    string              `json:"status"`
	StartTime	string              `json:"startTime"`
	UpdateTime	string              `json:"updateTime"`
	EndTime	    string              `json:"endTime"`
}

type AppInstanceAPICreateBody struct {
	RefApp      RefApp             `json:"refApp"`
	//RefResource []AppRefResource   `json:"refResource"`
	UserLabel  map[string]string  `json:"userLabel,omitempty"`
}

type AppInstanceAPIQueryBody struct{
	RefAppInstance RefAppInstance  `json:"refAppInstance"`
}

type AppInstanceAPIListBody struct {
	Page       int                 `json:"page"`
	Limit      int                 `json:"limit"`
	SortBy     SortOption          `json:"sortBy"`
}

type SortOption struct {
	Field      string              `json:"field"`
	Order      bool                `json:"order"`
}

type RespBody struct {
	RespData   RespData            `json:"data"`
	Status     int                 `json:"status"`
	Timestamp  int64               `json:"timestamp"`
}

type RespData struct {
	SourceDetail  Resource   `json:"source_detail"`
}