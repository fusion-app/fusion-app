package fusionappinstance

import (
	"context"
	"github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	"github.com/fusion-app/fusion-app/pkg/controller/internal"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *ReconcileFusionAppInstance) updateStatus(appInstance *v1alpha1.FusionAppInstance) error {
	if appInstance.Status.CreateTime == nil {
		now := metav1.Now()
		appInstance.Status.CreateTime = &now
	}
	if appInstance.Status.UpdateTime == nil {
		now := metav1.Now()
		appInstance.Status.UpdateTime = &now
	}
	labels := DefaultLabels(appInstance)
	pods, err := internal.PodsForLabels(appInstance.Namespace, labels, r.client)
	if err != nil {
		return err
	}
	app := new(v1alpha1.FusionApp)
	err = r.client.Get(context.TODO(), client.ObjectKey{Name: appInstance.Spec.RefApp.Name}, app)
	if errors.IsNotFound(err) {
		log.Warningf("failed to get fusionApp %v: %v", appInstance.Spec.RefApp.Name, err)
	} else if err != nil {
		return err
	} else {
		if len(appInstance.Spec.RefResource) == len(app.Spec.ResourceClaim) {
			appInstance.Status.Phase = v1alpha1.FusionAppInstancePhaseReady
		} else {
			appInstance.Status.Phase = v1alpha1.FusionAppInstancePhaseNotReady
		}
	}
	podStatuses := internal.MappingPodsByPhase(pods)
	if podStatuses[v1.PodRunning] == 1 {
		// All pods are running, set start time
		if appInstance.Status.StartTime == nil {
			now := metav1.Now()
			appInstance.Status.StartTime = &now
		}
		appInstance.Status.ProbePhase = v1alpha1.ProbePhaseSynchronous
	} else if podStatuses[v1.PodFailed] > 0 {
		appInstance.Status.ProbePhase = v1alpha1.ProbePhaseFailed
	} else {
		if appInstance.Spec.ProbeEnabled {
			appInstance.Status.ProbePhase = v1alpha1.ProbePhasePending
		} else {
			appInstance.Status.ProbePhase = v1alpha1.ProbePhaseNotReady
		}
	}
	return r.syncStatus(appInstance)
}

func (r *ReconcileFusionAppInstance) syncStatus(appInstance *v1alpha1.FusionAppInstance) error {
	oldAppInstance := &v1alpha1.FusionAppInstance{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      appInstance.Name,
		Namespace: appInstance.Namespace,
	}, oldAppInstance)

	if err != nil {
		 return err
	}
	if !reflect.DeepEqual(oldAppInstance.Status, appInstance.Status) {
		return r.client.Update(context.TODO(), appInstance)
	}

	return nil
}