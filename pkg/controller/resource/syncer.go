package resource

import (
	"github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	"github.com/fusion-app/fusion-app/pkg/syncer"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewProbeDeploySyncer(resource *v1alpha1.Resource, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	template := newDeployForProbe(resource)
	return syncer.NewDeploySyncer("prob-deploy", resource, template, c, scheme)
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