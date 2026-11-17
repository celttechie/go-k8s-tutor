package controller_runtime

import (
	"context"
	"maps"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func newTestScheme(t *testing.T) *runtime.Scheme {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := corev1.AddToScheme(scheme); err != nil {
		t.Fatalf("failed to add corev1 to scheme: %v", err)
	}
	return scheme
}

func TestReconcileMissingConfigMap(t *testing.T) {
	scheme := newTestScheme(t)
	c := fake.NewClientBuilder().WithScheme(scheme).Build()
	r := &ConfigSyncReconciler{Client: c, Scheme: scheme}

	res, err := r.Reconcile(context.Background(), ctrl.Request{
		NamespacedName: types.NamespacedName{Namespace: "default", Name: "does-not-exist"},
	})
	if err != nil {
		t.Fatalf("expected nil error for NotFound, got %v", err)
	}
	if res != (ctrl.Result{}) {
		t.Fatalf("expected empty ctrl.Result, got %+v", res)
	}
}

func TestReconcileAddsFinalizerAndCreatesMirror(t *testing.T) {
	ctx := context.Background()
	scheme := newTestScheme(t)

	src := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "vm-defaults",
			Namespace: "default",
			UID:       types.UID("source-uid-123"),
		},
		Data: map[string]string{
			"os":   "rocky9",
			"cpus": "2",
		},
	}

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(src).Build()
	r := &ConfigSyncReconciler{Client: c, Scheme: scheme}

	if _, err := r.Reconcile(ctx, ctrl.Request{
		NamespacedName: types.NamespacedName{Namespace: "default", Name: "vm-defaults"},
	}); err != nil {
		t.Fatalf("Reconcile failed: %v", err)
	}

	var updated corev1.ConfigMap
	if err := c.Get(ctx, types.NamespacedName{Namespace: "default", Name: "vm-defaults"}, &updated); err != nil {
		t.Fatalf("failed to get source ConfigMap: %v", err)
	}
	if !controllerutil.ContainsFinalizer(&updated, SyncFinalizer) {
		t.Fatalf("expected finalizer %q on source ConfigMap, got %v", SyncFinalizer, updated.Finalizers)
	}

	var mirror corev1.ConfigMap
	if err := c.Get(ctx, types.NamespacedName{Namespace: "default", Name: "vm-defaults-mirror"}, &mirror); err != nil {
		t.Fatalf("expected mirror ConfigMap vm-defaults-mirror to be created: %v", err)
	}
	if !maps.Equal(mirror.Data, src.Data) {
		t.Fatalf("expected mirror Data %v, got %v", src.Data, mirror.Data)
	}
	if len(mirror.OwnerReferences) != 1 || mirror.OwnerReferences[0].UID != src.UID {
		t.Fatalf("expected OwnerReference pointing to source UID %s, got %+v", src.UID, mirror.OwnerReferences)
	}
}

func TestReconcileUpdatesMirrorOnSourceChange(t *testing.T) {
	ctx := context.Background()
	scheme := newTestScheme(t)

	src := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:       "vm-defaults",
			Namespace:  "default",
			UID:        types.UID("source-uid-123"),
			Finalizers: []string{SyncFinalizer},
		},
		Data: map[string]string{"os": "rocky9", "cpus": "2"},
	}

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(src).Build()
	r := &ConfigSyncReconciler{Client: c, Scheme: scheme}

	req := ctrl.Request{NamespacedName: types.NamespacedName{Namespace: "default", Name: "vm-defaults"}}
	if _, err := r.Reconcile(ctx, req); err != nil {
		t.Fatalf("initial Reconcile failed: %v", err)
	}

	// Simulate user updating the source ConfigMap
	var latest corev1.ConfigMap
	if err := c.Get(ctx, req.NamespacedName, &latest); err != nil {
		t.Fatalf("failed to get source ConfigMap: %v", err)
	}
	latest.Data["cpus"] = "4"
	latest.Data["memory"] = "8Gi"
	if err := c.Update(ctx, &latest); err != nil {
		t.Fatalf("failed to update source ConfigMap: %v", err)
	}

	if _, err := r.Reconcile(ctx, req); err != nil {
		t.Fatalf("update Reconcile failed: %v", err)
	}

	var mirror corev1.ConfigMap
	if err := c.Get(ctx, types.NamespacedName{Namespace: "default", Name: "vm-defaults-mirror"}, &mirror); err != nil {
		t.Fatalf("failed to get mirror ConfigMap: %v", err)
	}
	if !maps.Equal(mirror.Data, latest.Data) {
		t.Fatalf("expected mirror Data %v after source update, got %v", latest.Data, mirror.Data)
	}
}

func TestReconcileRemovesFinalizerOnDelete(t *testing.T) {
	ctx := context.Background()
	scheme := newTestScheme(t)
	now := metav1.Now()

	src := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:              "vm-defaults",
			Namespace:         "default",
			Finalizers:        []string{SyncFinalizer},
			DeletionTimestamp: &now,
		},
	}

	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(src).Build()
	r := &ConfigSyncReconciler{Client: c, Scheme: scheme}

	if _, err := r.Reconcile(ctx, ctrl.Request{
		NamespacedName: types.NamespacedName{Namespace: "default", Name: "vm-defaults"},
	}); err != nil {
		t.Fatalf("Reconcile failed: %v", err)
	}

	var updated corev1.ConfigMap
	err := c.Get(ctx, types.NamespacedName{Namespace: "default", Name: "vm-defaults"}, &updated)
	// Note: fake client deletes the object once its last finalizer is removed,
	// so either NotFound OR finalizer removed is valid.
	if err == nil && controllerutil.ContainsFinalizer(&updated, SyncFinalizer) {
		t.Fatalf("expected finalizer %q to be removed on deletion", SyncFinalizer)
	}
}
