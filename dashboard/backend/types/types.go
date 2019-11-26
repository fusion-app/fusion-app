package types

import (
	fusionappv1alpha1  "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
)

const (
	DefaultNamespace string  = "fusion-app-resources"
	LeftBound        float64 = 116.326949
	RightBound       float64 = 116.328359
	UpBound          float64 = 39.993584
	DownBound        float64 = 39.993086
)

type Resource struct {
	UID          string      `json:"uid,omitempty"`
 	Namespace    string      `json:"namespace,omitempty"`
	Kind         string      `json:"kind"`
	ProbeEnabled bool        `json:"probeEnabled"`
	Phase        string      `json:"phase,omitempty"`
	Bound        bool        `json:"bound,omitempty"`
	Name         string      `json:"name"`
	AliasName    string      `json:"aliasName,omitempty"`
	AccessMode   string      `json:"accessMode"`
	Labels       map[string]string `json:"labels,omitempty"`
	Operation    []fusionappv1alpha1.ResourceOperationSpec `json:"operation,omitempty"`
	Icon         string       `json:"icon,omitempty"`
	Description  map[string]string       `json:"description,omitempty"`
	ProbeArgs    []string     `json:"probeArgs"`
	ProbeImage   string       `json:"probeImage"`
	Position      Position     `json:"position,omitempty"`
}

type Position struct {
	Longitude   float64       `json:"longitude,omitempty"`
	Latitude    float64       `json:"latitude,omitempty"`
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
	Description  map[string]string     `json:"description,omitempty"`
}

type AppAPICreateBody struct {
	AppSpec     App       `json:"appSpec"`
}

type ResourceAPICreateBody struct {
	ResourceSpec Resource `json:"resourceSpec"`
}

type AppAPIPutBody struct {
	Name        string    `json:"name,omitempty"`
	UID         string    `json:"uid,omitempty"`
	AppSpec     App       `json:"appSpec"`
}

type ResourceAPIPutBody struct {
	AppRefResource    fusionappv1alpha1.RefResource   `json:"refResource"`
	ResourceSpec      ResourceSpec     `json:"resourceSpec"`
}

type ResourceAPIQueryBody struct {
	Kind       string   `json:"kind,omitempty"`
	Phase      string   `json:"phase,omitempty"`
	RefResource fusionappv1alpha1.RefResource `json:"refResource,omitempty"`
	LabelSelector []fusionappv1alpha1.SelectorSpec `json:"labelSelector,omitempty"`
}

type ResourceAPIBindBody struct {
	RefResource    fusionappv1alpha1.RefResource   `json:"refResource"`
	RefAppInstance RefAppInstance   `json:"refAppInstance"`
}

type App struct {
	UID	            string                `json:"uid,omitempty"`
	Name	        string		          `json:"name"`
	ResourceClaim   []fusionappv1alpha1.ResourceClaimSpec   `json:"resourceClaim,omitempty"`
	AliasName       string                `json:"aliasName,omitempty"`
	Icon            string                `json:"icon,omitempty"`
	Description     map[string]string     `json:"description,omitempty"`
	Labels          map[string]string     `json:"labels,omitempty"`
}

type RefAppInstance struct {
	UID	        string  `json:"uid,omitempty"`
	Namespace	string  `json:"namespace,omitempty"`
	Kind	    string  `json:"kind,omitempty"`
	Name	    string	`json:"name"`
}


type AppInstance struct {
	UID	        string              `json:"uid,omitempty"`
	Namespace	string              `json:"namespace,omitempty"`
	Name	    string              `json:"name"`
	RefApp      fusionappv1alpha1.RefApp              `json:"refApp"`
	RefResource	[]fusionappv1alpha1.RefResource    `json:"refResource,omitempty"`
	Status	    []fusionappv1alpha1.Action  `json:"status,omitempty"`
	CreateTime  string              `json:"createTime,omitempty"`
	StartTime	string              `json:"startTime,omitempty"`
	UpdateTime	string              `json:"updateTime,omitempty"`
	EndTime	    string              `json:"endTime,omitempty"`
}

type AppInstanceAPICreateBody struct {
	RefApp      fusionappv1alpha1.RefApp             `json:"refApp"`
	//RefResource []RefResource   `json:"refResource"`
	UserLabel  map[string]string  `json:"userLabel,omitempty"`
}

type AppInstanceAPIQueryBody struct{
	RefAppInstance RefAppInstance  `json:"refAppInstance"`
}

type AppInstanceAPIPutBody struct {
	RefAppInstance  RefAppInstance  `json:"refAppInstance"`
	AppInstanceSpec AppinstanceSpec `json:"appInstanceSpec"`
}

type AppinstanceSpec struct {
	ProbeEnabled  bool               `json:"probeEnabled"`
	Labels        map[string]string  `json:"labels,omitempty"`
}

type AppInstanceAPIDeleteBody struct{
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