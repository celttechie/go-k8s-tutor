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
	if kubeconfigPath != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}
	return rest.InClusterConfig()
}

// ListKVMNodes lists all Nodes in the cluster and returns the names of nodes
// whose Status.Allocatable contains KVMResourceName with a quantity > 0.
func ListKVMNodes(ctx context.Context, client kubernetes.Interface) ([]string, error) {
	nodeList, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var kvmNodes []string
	for _, node := range nodeList.Items {
		if qty, ok := node.Status.Allocatable[KVMResourceName]; ok && qty.Value() > 0 {
			kvmNodes = append(kvmNodes, node.Name)
		}
	}
	return kvmNodes, nil
}

// EnsureNamespace creates the given Namespace if it does not already exist.
// If the Namespace already exists (apierrors.IsAlreadyExists(err)), it returns nil
// so the operation is safe and idempotent.
func EnsureNamespace(ctx context.Context, client kubernetes.Interface, name string) error {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: name},
	}
	_, err := client.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	if err != nil && !apierrors.IsAlreadyExists(err) {
		return err
	}
	return nil
}

// DeletePodIfExists deletes the named Pod in namespace.
// If the Pod is already gone (apierrors.IsNotFound(err)), it returns nil.
func DeletePodIfExists(ctx context.Context, client kubernetes.Interface, namespace, name string) error {
	err := client.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	return nil
}
