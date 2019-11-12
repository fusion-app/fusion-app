package resource

import (
	"github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	"github.com/fusion-app/fusion-app/pkg/syncer"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewProbePodSyncer(resource *v1alpha1.Resource, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	pod := newPodForProbe(resource)
	return syncer.NewObjectSyncer("probe-pod", resource, pod, c, scheme, func(existing runtime.Object) error {
		out := existing.(*v1.Pod)
		if !reflect.DeepEqual(out.Spec, pod.Spec) {
			out.Spec = pod.Spec
		}
		return nil
	})
}

func NewConsumerPodSyncer(resource *v1alpha1.Resource, c client.Client, scheme *runtime.Scheme) syncer.Interface {
	pod := newPodForConsumer(resource)
	return syncer.NewObjectSyncer("consumer-pod", resource, pod, c, scheme, func(existing runtime.Object) error {
		out := existing.(*v1.Pod)
		if !reflect.DeepEqual(out.Spec, pod.Spec) {
			out.Spec = pod.Spec
		}
		return nil
	})
}