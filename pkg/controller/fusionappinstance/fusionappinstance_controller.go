package fusionappinstance

import (
	"context"
	fusionappv1alpha1 "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	"github.com/fusion-app/fusion-app/pkg/syncer"
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const controllerName = "fusionappinstance-controller"

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new FusionAppInstance Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileFusionAppInstance{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New(controllerName, mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource FusionAppInstance
	err = c.Watch(&source.Kind{Type: &fusionappv1alpha1.FusionAppInstance{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	subResources := []runtime.Object{
		&corev1.Pod{},
		&appsv1.Deployment{},
	}

	// Watch for changes to secondary resource Pods and requeue the owner Resource
	for _, subResource := range subResources {
		err = c.Watch(&source.Kind{Type: subResource}, &handler.EnqueueRequestForOwner{
			IsController: true,
			OwnerType:    &fusionappv1alpha1.FusionAppInstance{},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// blank assignment to verify that ReconcileFusionAppInstance implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileFusionAppInstance{}

// ReconcileFusionAppInstance reconciles a FusionAppInstance object
type ReconcileFusionAppInstance struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
	recorder record.EventRecorder
}

// Reconcile reads that state of the cluster for a FusionAppInstance object and makes changes based on the state read
// and what is in the FusionAppInstance.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileFusionAppInstance) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the FusionAppInstance instance
	instance := &fusionappv1alpha1.FusionAppInstance{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	log.Printf("Reconciling FusionAppInstance %s", instance.Name)
	// if the resource is terminating, stop reconcile
	if instance.ObjectMeta.DeletionTimestamp != nil {
		for _, refResource := range instance.Spec.RefResource {
			rs := &fusionappv1alpha1.Resource{}
			err := r.client.Get(context.TODO(), client.ObjectKey{Name: refResource.Name, Namespace: refResource.Namespace}, rs)
			if err != nil && !errors.IsNotFound(err) {
				return reconcile.Result{}, err
			}
			if rs.Status.Bound  {
				rs.Status.Bound = false
				err := r.client.Update(context.TODO(), rs)
				if err != nil && !errors.IsNotFound(err) {
					return reconcile.Result{}, err
				}
			}
		}
		return reconcile.Result{}, nil
	}
	probeEnabled := os.Getenv("APPINSTANCE_PROBE_ENABLED")
	if instance.Spec.ProbeEnabled || probeEnabled == "true" {
		if len(instance.Spec.ProbeImage) == 0 || len(instance.Spec.ProbeArgs) == 0 {
			app := &fusionappv1alpha1.FusionApp{}
			err := r.client.Get(context.TODO(),
				client.ObjectKey{Name: instance.Spec.RefApp.Name}, app)
			if err != nil {
				if errors.IsNotFound(err) {
					log.Printf("can't find referred App")
					return reconcile.Result{}, nil
				}
				// Error reading the object - requeue the request.
				return reconcile.Result{}, err
			}
			instance.Spec.ProbeImage = app.Spec.ProbeImage
			in, out := &app.Spec.ProbeArgs, &instance.Spec.ProbeArgs
			*out = make([]string, len(*in))
			for i := range *in {
				(*in)[i] = (*out)[i]
			}
		}
	}
	var syncers []syncer.Interface
	if instance.Spec.ProbeEnabled || probeEnabled == "true" {
		syncers = append(syncers, NewProbeDeploySyncer(instance, r.client, r.scheme))
	} else {
		deploy := &appsv1.Deployment{}
		err := r.client.Get(context.TODO(), client.ObjectKey{Name: instance.Name + "-probe-deploy", Namespace: instance.Namespace}, deploy)
		if err == nil {
			if err := r.client.Delete(context.TODO(), deploy); err != nil {
				return reconcile.Result{}, err
			}
		}
	}
	if len(syncers) > 0 {
		if err := r.sync(syncers); err != nil {
			return reconcile.Result{}, err
		}
	}
	return reconcile.Result{}, r.updateStatus(instance)
}

func (r *ReconcileFusionAppInstance) sync(syncers []syncer.Interface) error {
	for _, s := range syncers {
		if err := syncer.Sync(context.TODO(), s, r.recorder); err != nil {
			return err
		}
	}
	return nil
}