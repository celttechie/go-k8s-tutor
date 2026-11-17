package client_go

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func TestLoadKubeConfig(t *testing.T) {
	kubeconfig := `apiVersion: v1
kind: Config
clusters:
- cluster:
    server: https://127.0.0.1:6443
  name: test-cluster
contexts:
- context:
    cluster: test-cluster
    user: test-user
  name: test-ctx
current-context: test-ctx
users:
- name: test-user
  user:
    token: dummy-token
`
	path := filepath.Join(t.TempDir(), "kubeconfig")
	if err := os.WriteFile(path, []byte(kubeconfig), 0o600); err != nil {
		t.Fatalf("failed to write temp kubeconfig: %v", err)
	}

	cfg, err := LoadKubeConfig(path)
	if err != nil {
		t.Fatalf("unexpected error loading kubeconfig: %v", err)
	}
	if cfg == nil || cfg.Host != "https://127.0.0.1:6443" {
		t.Fatalf("expected Host https://127.0.0.1:6443, got %+v", cfg)
	}
}

func TestListKVMNodes(t *testing.T) {
	client := fake.NewClientset(
		&corev1.Node{
			ObjectMeta: metav1.ObjectMeta{Name: "worker-kvm-1"},
			Status: corev1.NodeStatus{
				Allocatable: corev1.ResourceList{
					KVMResourceName: resource.MustParse("110"),
				},
			},
		},
		&corev1.Node{
			ObjectMeta: metav1.ObjectMeta{Name: "worker-no-kvm"},
			Status: corev1.NodeStatus{
				Allocatable: corev1.ResourceList{
					corev1.ResourceCPU: resource.MustParse("8"),
				},
			},
		},
		&corev1.Node{
			ObjectMeta: metav1.ObjectMeta{Name: "worker-zero-kvm"},
			Status: corev1.NodeStatus{
				Allocatable: corev1.ResourceList{
					KVMResourceName: resource.MustParse("0"),
				},
			},
		},
	)

	nodes, err := ListKVMNodes(context.Background(), client)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !slices.Equal(nodes, []string{"worker-kvm-1"}) {
		t.Fatalf("expected [worker-kvm-1], got %v", nodes)
	}
}

func TestEnsureNamespace(t *testing.T) {
	ctx := context.Background()
	client := fake.NewClientset()

	if err := EnsureNamespace(ctx, client, "virt-vms"); err != nil {
		t.Fatalf("first EnsureNamespace failed: %v", err)
	}
	// Second call should be idempotent (AlreadyExists -> nil)
	if err := EnsureNamespace(ctx, client, "virt-vms"); err != nil {
		t.Fatalf("second EnsureNamespace should ignore AlreadyExists, got: %v", err)
	}

	ns, err := client.CoreV1().Namespaces().Get(ctx, "virt-vms", metav1.GetOptions{})
	if err != nil || ns.Name != "virt-vms" {
		t.Fatalf("expected namespace virt-vms to exist, got err=%v", err)
	}
}

func TestDeletePodIfExists(t *testing.T) {
	ctx := context.Background()
	client := fake.NewClientset(&corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "launcher-1", Namespace: "default"},
	})

	if err := DeletePodIfExists(ctx, client, "default", "launcher-1"); err != nil {
		t.Fatalf("expected nil deleting existing pod, got: %v", err)
	}
	// Second delete should ignore NotFound
	if err := DeletePodIfExists(ctx, client, "default", "launcher-1"); err != nil {
		t.Fatalf("expected nil deleting missing pod, got: %v", err)
	}

	pods, err := client.CoreV1().Pods("default").List(ctx, metav1.ListOptions{})
	if err != nil || len(pods.Items) != 0 {
		t.Fatalf("expected 0 pods remaining, got %d (err=%v)", len(pods.Items), err)
	}
}
