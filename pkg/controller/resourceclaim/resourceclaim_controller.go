package resourceclaim

import (
	fusionappv1alpha1 "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
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
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileResourceClaim) Reconcile(request reconcile.Request) (reconcile.Result, error) {

	// Fetch the ResourceClaim instance
	//instance := &fusionappv1alpha1.ResourceClaim{}
	//err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	//if err != nil {
	//	if errors.IsNotFound(err) {
	//		// Request object not found, could have been deleted after reconcile request.
	//		// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
	//		// Return and don't requeue
	//		return reconcile.Result{}, nil
	//	}
	//	// Error reading the object - requeue the request.
	//	return reconcile.Result{}, err
	//}
	//rsl := &fusionappv1alpha1.ResourceList{}
	//err = r.client.List(context.TODO(), &client.ListOptions{}, rsl)
	//if err != nil {
	//	log.Warningf("failed to list resources", err)
	//	return reconcile.Result{}, err
	//}
	//resources := make([]fusionappv1alpha1.Resource, 0)
	//for _, item := range rsl.Items {
	//	if !item.Status.Bound {
	//		resources = append(resources, item)
	//	}
	//}
	//if len(resources) == 0 {
	//	return reconcile.Result{}, fmt.Errorf("no resources available currently")
	//}
	//mp := make(labels.Set)
	//for _, selector := range instance.Spec.Selector {
	//	mp[selector.Key] = selector.Value
	//}
	//labelSelector := labels.SelectorFromSet(mp)
	//var resource *fusionappv1alpha1.Resource
	//for _, item := range resources {
	//	if (item.Status.Bound == false || item.Spec.AccessMode == fusionappv1alpha1.ResourceAccessModeShared) && labelSelector.Matches(labels.Set(item.Spec.Labels)) {
	//		resource = &item
	//		break
	//	}
	//}
	//if resource == nil {
	//	return reconcile.Result{}, fmt.Errorf("no suitable resources available currently")
	//}
	//if resource.Status.Bound == false {
	//	resource.Status.Bound = true
	//	err = r.client.Update(context.TODO(), resource)
	//	if err != nil {
	//		return reconcile.Result{}, err
	//	}
	//}
	return reconcile.Result{}, nil
}

