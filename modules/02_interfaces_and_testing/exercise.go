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

func (v *SimpleVM) GetName() string {
	return v.Name
}

func (v *SimpleVM) GetNamespace() string {
	return v.Namespace
}

// CheckVMExists returns true if ClusterReader finds the VM, or an error if lookup fails.
func CheckVMExists(reader ClusterReader, namespace, name string) (bool, error) {
	res, err := reader.GetVM(namespace, name)
	if err != nil {
		return false, fmt.Errorf("checking vm %s/%s: %w", namespace, name, err)
	}
	return res != nil, nil
}
