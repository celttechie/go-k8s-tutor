package controller_runtime

import (
	"context"
	"maps"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const SyncFinalizer = "tutor.celttechie.io/cleanup"

// ConfigSyncReconciler watches source ConfigMaps and maintains an owned
// "<name>-mirror" ConfigMap in the same namespace.
type ConfigSyncReconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
}

// Reconcile implements reconcile.Reconciler.
//
// Desired reconciliation flow:
// 1. Fetch the source corev1.ConfigMap using r.Client.Get(ctx, req.NamespacedName, &cm).
//    - If apierrors.IsNotFound(err), return ctrl.Result{}, nil (object was deleted).
//    - If another error occurs, return ctrl.Result{}, err (controller-runtime will requeue).
// 2. Finalizer handling:
//    - If !cm.DeletionTimestamp.IsZero():
//        If controllerutil.ContainsFinalizer(&cm, SyncFinalizer):
//          // In a real operator, run external cleanup (e.g., cloud/storage) here first.
//          Remove it with controllerutil.RemoveFinalizer(&cm, SyncFinalizer),
//          call r.Client.Update(ctx, &cm), and return ctrl.Result{}, err.
//        Otherwise return ctrl.Result{}, nil.
//    - If !controllerutil.ContainsFinalizer(&cm, SyncFinalizer):
//        Add it with controllerutil.AddFinalizer(&cm, SyncFinalizer),
//        call r.Client.Update(ctx, &cm), and if err != nil return ctrl.Result{}, err.
// 3. Create or update the mirror ConfigMap ("<cm.Name>-mirror" in cm.Namespace):
//    - Construct mirror := &corev1.ConfigMap{ObjectMeta: metav1.ObjectMeta{Name: cm.Name + "-mirror", Namespace: cm.Namespace}}
//    - Use controllerutil.CreateOrUpdate(ctx, r.Client, mirror, func() error { ... })
//      Inside the mutate function:
//        - Set mirror.Data = maps.Clone(cm.Data)
//        - Call controllerutil.SetControllerReference(&cm, mirror, r.Scheme)
//    - Return ctrl.Result{}, err.
func (r *ConfigSyncReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var cm corev1.ConfigMap
	if err := r.Client.Get(ctx, req.NamespacedName, &cm); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	if !cm.DeletionTimestamp.IsZero() {
		if controllerutil.ContainsFinalizer(&cm, SyncFinalizer) {
			controllerutil.RemoveFinalizer(&cm, SyncFinalizer)
			if err := r.Client.Update(ctx, &cm); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if !controllerutil.ContainsFinalizer(&cm, SyncFinalizer) {
		controllerutil.AddFinalizer(&cm, SyncFinalizer)
		if err := r.Client.Update(ctx, &cm); err != nil {
			return ctrl.Result{}, err
		}
	}

	mirror := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cm.Name + "-mirror",
			Namespace: cm.Namespace,
		},
	}
	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, mirror, func() error {
		mirror.Data = maps.Clone(cm.Data)
		return controllerutil.SetControllerReference(&cm, mirror, r.Scheme)
	})
	return ctrl.Result{}, err
}

// SetupWithManager registers ConfigSyncReconciler with the controller-runtime Manager,
// watching primary ConfigMaps (For) and their owned mirror ConfigMaps (Owns).
func (r *ConfigSyncReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.ConfigMap{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
