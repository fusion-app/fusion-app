package resource

import (
	"context"
	"github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	"github.com/fusion-app/fusion-app/pkg/controller/internal"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
)

func (r *ReconcileResource) updateStatus(resource *v1alpha1.Resource) error {
	if resource.Status.CreateTime == nil {
		now := metav1.Now()
		resource.Status.CreateTime = &now
	}

	labels := DefaultLabels(resource)

	pods, err := internal.PodsForLabels(resource.Namespace, labels, r.client)
	if err != nil {
		return err
	}

	podStatuses := internal.MappingPodsByPhase(pods)
	if podStatuses[v1.PodRunning] == 2 {
		// All pods are running, set start time
		if resource.Status.StartTime == nil {
			now := metav1.Now()
			resource.Status.StartTime = &now
		}
		resource.Status.Phase = v1alpha1.ResourcePhaseRunning
	} else if podStatuses[v1.PodFailed] > 0 {
		resource.Status.Phase = v1alpha1.ResourcePhaseFailed
	} else {
		resource.Status.Phase = v1alpha1.ResourcePhasePending
	}

	return r.syncStatus(resource)
}

func (r *ReconcileResource) syncStatus(resource *v1alpha1.Resource) error {
	oldResource := &v1alpha1.Resource{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      resource.Name,
		Namespace: resource.Namespace,
	}, oldResource)

	if err != nil {
		 return err
	}
	if !reflect.DeepEqual(oldResource.Status, resource.Status) {
		return r.client.Update(context.TODO(), resource)
	}

	return nil
}