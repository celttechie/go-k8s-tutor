package interfaces_testing

import "fmt"

// Resource represents a simplified Kubernetes object (similar to client.Object).
type Resource interface {
	GetName() string
	GetNamespace() string
}

// ClusterReader abstracts Kubernetes API reads so business logic can be unit tested
// without a live cluster (the same pattern controller-runtime's client.Reader uses).
type ClusterReader interface {
	GetVM(namespace, name string) (Resource, error)
}

// SimpleVM implements Resource via pointer receivers.
type SimpleVM struct {
	Namespace string
	Name      string
}

// TODO: Implement GetName() and GetNamespace() on *SimpleVM so *SimpleVM satisfies Resource.

// CheckVMExists returns true if ClusterReader finds the VM, or an error if lookup fails.
func CheckVMExists(reader ClusterReader, namespace, name string) (bool, error) {
	// TODO: Call reader.GetVM(namespace, name).
	// If err != nil, wrap and return it: fmt.Errorf("checking vm %s/%s: %w", namespace, name, err)
	// If res == nil, return false, nil.
	// Otherwise return true, nil.
	_ = fmt.Errorf
	return false, nil
}
