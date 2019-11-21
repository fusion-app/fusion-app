package fusionappinstance

import (
	"github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	"github.com/fusion-app/fusion-app/pkg/syncer"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewProbeDeploySyncer(appInstance *v1alpha1.FusionAppInstance, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	template := newDeployForProbe(appInstance)
	return syncer.NewDeploySyncer("prob-deploy", appInstance, template, c, scheme)
}
