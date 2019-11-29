package resourceclaim

import (
	"context"
	"fmt"
	fusionappv1alpha1 "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)


/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ResourceClaim Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileResourceClaim{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("resourceclaim-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource ResourceClaim
	err = c.Watch(&source.Kind{Type: &fusionappv1alpha1.ResourceClaim{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileResourceClaim implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileResourceClaim{}

// ReconcileResourceClaim reconciles a ResourceClaim object
type ReconcileResourceClaim struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ResourceClaim object and makes changes based on the state read
// and what is in the ResourceClaim.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileResourceClaim) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	// Fetch the ResourceClaim resourceClaim
	resourceClaim := &fusionappv1alpha1.ResourceClaim{}
	err := r.client.Get(context.TODO(), request.NamespacedName, resourceClaim)
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
	// if the resource is terminating, ubound resources and stop reconcile
	log.Printf("Reconciling resourceClaim %s", resourceClaim.Name)
	if resourceClaim.ObjectMeta.DeletionTimestamp != nil {
		rs := &fusionappv1alpha1.Resource{}
		err := r.client.Get(context.TODO(), client.ObjectKey{
			Name: resourceClaim.Spec.RefResource.Name,
			Namespace: resourceClaim.Spec.RefResource.Namespace,
		}, rs)
		if err != nil && !errors.IsNotFound(err) {
			return reconcile.Result{}, err
		}
		if err == nil {
			var refResourceClaim []fusionappv1alpha1.RefResourceClaim
			for _, item := range rs.Spec.RefResourceClaim {
				if item.Name != resourceClaim.Name || item.Namespace!= resourceClaim.Namespace {
					refResourceClaim = append(refResourceClaim, item)
				}
			}
			if len(refResourceClaim) == 0 {
				rs.Status.Bound = false
			}
			rs.Spec.RefResourceClaim = refResourceClaim
			err := r.client.Update(context.TODO(), rs)
			if err != nil && !errors.IsNotFound(err) {
				return reconcile.Result{}, err
			}
		}
		appi := &fusionappv1alpha1.FusionAppInstance{}
		err = r.client.Get(context.TODO(), client.ObjectKey{
			Name: resourceClaim.Spec.RefAppInstance.Name,
			Namespace: resourceClaim.Spec.RefAppInstance.Namespace,
		}, appi)
		if err != nil && !errors.IsNotFound(err) {
			return reconcile.Result{}, err
		}
		if err == nil {
			var refResourceClaim []fusionappv1alpha1.RefResourceClaim
			for _, item := range appi.Spec.RefResourceClaim {
				if item.Name != resourceClaim.Name || item.Namespace!= resourceClaim.Namespace {
					refResourceClaim = append(refResourceClaim, item)
				}
			}
			appi.Spec.RefResourceClaim = refResourceClaim
			err := r.client.Update(context.TODO(), appi)
			if err != nil && !errors.IsNotFound(err) {
				return reconcile.Result{}, err
			}
		}
		return reconcile.Result{}, nil
	}
	if resourceClaim.Status.Bound {
		return reconcile.Result{}, nil
	}
	appInstance := new(fusionappv1alpha1.FusionAppInstance)
	err = r.client.Get(context.TODO(), client.ObjectKey{
		Name: resourceClaim.Spec.RefAppInstance.Name,
		Namespace: resourceClaim.Spec.RefAppInstance.Namespace,
	}, appInstance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}
	rsl := &fusionappv1alpha1.ResourceList{}
	err = r.client.List(context.TODO(), &client.ListOptions{}, rsl)
	if err != nil {
		log.Warningf("failed to list resources", err)
		return reconcile.Result{}, err
	}
	resources := make([]fusionappv1alpha1.Resource, 0)
	for _, item := range rsl.Items {
		if !item.Status.Bound {
			resources = append(resources, item)
		}
	}
	if len(resources) == 0 {
		return reconcile.Result{}, fmt.Errorf("no resources available currently")
	}
	mp := make(labels.Set)
	for _, selector := range resourceClaim.Spec.Selector {
		mp[selector.Key] = selector.Value
	}
	labelSelector := labels.SelectorFromSet(mp)
	var resource *fusionappv1alpha1.Resource
	for _, item := range resources {
		if (item.Status.Bound == false || item.Spec.AccessMode == fusionappv1alpha1.ResourceAccessModeShared) && labelSelector.Matches(labels.Set(item.Spec.Labels)) {
			resource = &item
			break
		}
	}
	if resource == nil {
		return reconcile.Result{}, fmt.Errorf("no suitable resources available currently")
	}
	resource.Spec.RefResourceClaim = append(resource.Spec.RefResourceClaim, fusionappv1alpha1.RefResourceClaim{
		UID:       string(resourceClaim.UID),
		Name:      resourceClaim.Name,
		Namespace: resourceClaim.Namespace,
	})
	resource.Status.Bound = true
	err = r.client.Update(context.TODO(), resource)
	if err != nil {
		return reconcile.Result{}, err
	}
	resourceClaim.Spec.RefResource = fusionappv1alpha1.RefResource {
		Name: resource.Name,
		Namespace: resource.Namespace,
		Kind: string(resource.Spec.ResourceKind),
		UID: string(resource.UID),
	}
	resourceClaim.Status.Bound = true
	appInstance.Spec.RefResource = append(appInstance.Spec.RefResource, fusionappv1alpha1.RefResource{
		Name: resource.Name,
		Namespace: resource.Namespace,
		Kind: string(resource.Spec.ResourceKind),
		UID: string(resource.UID),
	})
	_ = r.client.Update(context.TODO(), resourceClaim)
	_ = r.client.Update(context.TODO(), appInstance)
	return reconcile.Result{}, nil
}

