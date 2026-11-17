package client_go

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const KVMResourceName = corev1.ResourceName("devices.kubevirt.io/kvm")

// LoadKubeConfig builds a *rest.Config from kubeconfigPath if non-empty,
// otherwise falling back to rest.InClusterConfig().
// This is the standard dual-mode config loader used by K8s CLIs and operators.
func LoadKubeConfig(kubeconfigPath string) (*rest.Config, error) {
	// TODO: If kubeconfigPath != "", return clientcmd.BuildConfigFromFlags("", kubeconfigPath).
	// Otherwise return rest.InClusterConfig().
	_ = clientcmd.BuildConfigFromFlags
	_ = rest.InClusterConfig
	return nil, nil
}

// ListKVMNodes lists all Nodes in the cluster and returns the names of nodes
// whose Status.Allocatable contains KVMResourceName with a quantity > 0.
func ListKVMNodes(ctx context.Context, client kubernetes.Interface) ([]string, error) {
	// TODO: Call client.CoreV1().Nodes().List(ctx, metav1.ListOptions{}).
	// Loop over nodeList.Items, check node.Status.Allocatable[KVMResourceName],
	// and if ok && qty.Value() > 0, append node.Name.
	_ = metav1.ListOptions{}
	return nil, nil
}

// EnsureNamespace creates the given Namespace if it does not already exist.
// If the Namespace already exists (apierrors.IsAlreadyExists(err)), it returns nil
// so the operation is safe and idempotent.
func EnsureNamespace(ctx context.Context, client kubernetes.Interface, name string) error {
	// TODO: Call client.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{...}, metav1.CreateOptions{}).
	// Return nil if err == nil or apierrors.IsAlreadyExists(err).
	_ = corev1.Namespace{}
	_ = apierrors.IsAlreadyExists
	return nil
}

// DeletePodIfExists deletes the named Pod in namespace.
// If the Pod is already gone (apierrors.IsNotFound(err)), it returns nil.
func DeletePodIfExists(ctx context.Context, client kubernetes.Interface, namespace, name string) error {
	// TODO: Call client.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{}).
	// Return nil if err == nil or apierrors.IsNotFound(err).
	return nil
}
