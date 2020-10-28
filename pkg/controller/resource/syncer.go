package resource

import (
	"github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	"github.com/fusion-app/fusion-app/pkg/syncer"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewProbeAndMSDeploySyncer(resource *v1alpha1.Resource, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	template := newDeployForProbeAndMS(resource)
	return syncer.NewDeploySyncer("probe-deploy", resource, template, c, scheme)
}

func NewMSServiceSyncer(resource *v1alpha1.Resource, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	template := newServiceForMS(resource)
	return syncer.NewServiceSyncer("ms-svc", resource, template, c, scheme)
}

func NewPatcherConfigmapSyncer(resource *v1alpha1.Resource, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	template := newConfigmapForPatcher(resource)
	return syncer.NewConfigmapSyncer("patcher-cm", resource, template, c, scheme)
}
//func NewConsumerPodSyncer(resource *v1alpha1.Resource, c client.Client, scheme *runtime.Scheme) syncer.Interface {
//	pod := newPodForConsumer(resource)
//	return syncer.NewObjectSyncer("consumer-pod", resource, pod, c, scheme, func(existing runtime.Object) error {
//		out := existing.(*v1.Pod)
//		if !reflect.DeepEqual(out.Spec, pod.Spec) {
//			out.Spec = pod.Spec
//		}
//		return nil
//	})
//}