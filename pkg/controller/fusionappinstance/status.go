package fusionappinstance

import (
	"context"
	"github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	"github.com/fusion-app/fusion-app/pkg/controller/internal"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
)

func (r *ReconcileFusionAppInstance) updateStatus(appInstance *v1alpha1.FusionAppInstance) error {
	if appInstance.Status.CreateTime == nil {
		now := metav1.Now()
		appInstance.Status.CreateTime = &now
	}

	labels := DefaultLabels(appInstance)

	pods, err := internal.PodsForLabels(appInstance.Namespace, labels, r.client)
	if err != nil {
		return err
	}

	podStatuses := internal.MappingPodsByPhase(pods)
	if podStatuses[v1.PodRunning] == 1 {
		// All pods are running, set start time
		if appInstance.Status.StartTime == nil {
			now := metav1.Now()
			appInstance.Status.StartTime = &now
		}
		appInstance.Status.Phase = v1alpha1.ResourcePhaseRunning
		appInstance.Status.ProbePhase = v1alpha1.ProbePhaseSynchronous
	} else if podStatuses[v1.PodFailed] > 0 {
		appInstance.Status.Phase = v1alpha1.ResourcePhaseFailed
		appInstance.Status.ProbePhase = v1alpha1.ProbePhaseFailed
	} else {
		appInstance.Status.Phase = v1alpha1.ResourcePhasePending
		if appInstance.Spec.ProbeEnabled {
			appInstance.Status.ProbePhase = v1alpha1.ProbePhasePending
		} else {
			appInstance.Status.ProbePhase = v1alpha1.ProbePhaseNotReady
		}
	}
	return r.syncStatus(appInstance)
}

func (r *ReconcileFusionAppInstance) syncStatus(appInstance *v1alpha1.FusionAppInstance) error {
	oldResource := &v1alpha1.Resource{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      appInstance.Name,
		Namespace: appInstance.Namespace,
	}, oldResource)

	if err != nil {
		 return err
	}
	if !reflect.DeepEqual(oldResource.Status, appInstance.Status) {
		return r.client.Update(context.TODO(), appInstance)
	}

	return nil
}