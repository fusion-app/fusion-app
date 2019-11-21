package fusionappinstance

import (
	"github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	"github.com/fusion-app/fusion-app/pkg/controller/internal"
	"k8s.io/apimachinery/pkg/labels"
)

func DefaultLabels(appInstance *v1alpha1.FusionAppInstance) labels.Set {
	l := labels.Set{}
	l[internal.LabelNameKey] = appInstance.Name

	return l
}
