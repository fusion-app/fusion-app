package handler

import (
	fusionappv1alpha1  "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
)

const (
	defaultNamespace = "fusion-app-resources"
)

type ResourceList struct {
	Resources []fusionappv1alpha1.Resource `json:"resources"`
}