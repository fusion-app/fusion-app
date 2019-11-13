package resource

import (
	"github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	"github.com/fusion-app/fusion-app/pkg/controller/internal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func DefaultLabels(resource *v1alpha1.Resource) labels.Set {
	l := labels.Set{}
	l[internal.LabelNameKey] = resource.Name

	return l
}

func AddKindLabel(resource *v1alpha1.Resource) {
	if resource.ObjectMeta.Labels == nil {
		resource.ObjectMeta.Labels = make(map[string]string)
	}
	resource.ObjectMeta.Labels[internal.LabelKindKey] = string(resource.Spec.ResourceKind)
}

func SelectorForKind(kind v1alpha1.ResourceKind) labels.Selector {
	selector := &metav1.LabelSelector{
		MatchLabels: map[string]string{
			internal.LabelKindKey: string(kind),
		},
	}

	labelSelector, _ := metav1.LabelSelectorAsSelector(selector)

	return labelSelector
}