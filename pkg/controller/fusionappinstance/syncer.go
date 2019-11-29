package fusionappinstance

import (
	"github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	"github.com/fusion-app/fusion-app/pkg/syncer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewProbeDeploySyncer(appInstance *v1alpha1.FusionAppInstance, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	template := newDeployForProbe(appInstance)
	return syncer.NewDeploySyncer("prob-deploy", appInstance, template, c, scheme)
}

func NewResourceClaimSyncer(appInstance *v1alpha1.FusionAppInstance, resourceClaim *v1alpha1.ResourceClaimSpec, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	template := newResourceClaim(appInstance, resourceClaim)
	metaobj := &v1alpha1.ResourceClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      template.ObjectMeta.Name,
			Namespace: template.ObjectMeta.Namespace,
		},
	}
	return syncer.NewObjectSyncer(appInstance.Name + resourceClaim.Name, appInstance, metaobj, c, scheme, func(existing runtime.Object) error {
		out := existing.(*v1alpha1.ResourceClaim)
		out.Spec.RefAppInstance = v1alpha1.RefAppInstance{
			UID:       string(appInstance.UID),
			Name:      appInstance.Name,
			Namespace: appInstance.Namespace,
		}
		out.Spec.AccessMode = template.Spec.AccessMode
		out.Spec.Selector = template.Spec.Selector
		return nil
	})
}