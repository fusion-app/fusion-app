package resource

import (
	"github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	"github.com/fusion-app/fusion-app/pkg/controller/internal"
	"k8s.io/apimachinery/pkg/labels"
)

func DefaultLabels(bus *v1alpha1.Resource) labels.Set {
	l := labels.Set{}
	l[internal.LabelNameKey] = bus.Name

	return l
}

