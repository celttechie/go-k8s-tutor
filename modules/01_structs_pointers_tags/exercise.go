package structs_pointers_tags

import "encoding/json"

// ObjectMeta mirrors Kubernetes metav1.ObjectMeta (simplified).
type ObjectMeta struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// VMSpec mirrors a simplified KubeVirt VirtualMachineSpec.
// Notice Running is a *bool so callers can distinguish "unset" (nil)
// from explicitly false.
type VMSpec struct {
	Running *bool  `json:"running,omitempty"`
	CPUs    *int32 `json:"cpus,omitempty"`
	Memory  string `json:"memory"`
}

// VirtualMachine mirrors a Kubernetes CRD object.
type VirtualMachine struct {
	APIVersion string     `json:"apiVersion"`
	Kind       string     `json:"kind"`
	Metadata   ObjectMeta `json:"metadata"`
	Spec       VMSpec     `json:"spec"`
}

// NewVirtualMachine constructs a VirtualMachine with:
// - APIVersion: "kubevirt.io/v1"
// - Kind: "VirtualMachine"
// - Metadata.Name and Metadata.Namespace set
// - Spec.Memory set, Spec.CPUs set to a pointer to cpus, and Spec.Running left nil (unset).
func NewVirtualMachine(name, namespace string, cpus int32, memory string) VirtualMachine {
	// TODO: Implement this function.
	// Hint: In Go, you can take the address of a local variable (&cpus) or use a helper.
	return VirtualMachine{}
}

// EffectiveCPUs returns the configured CPU count, or defaultCPUs if Spec.CPUs is nil.
func EffectiveCPUs(vm VirtualMachine, defaultCPUs int32) int32 {
	// TODO: Safely dereference vm.Spec.CPUs if non-nil, otherwise return defaultCPUs.
	return 0
}

// MarshalVM serializes a VirtualMachine to JSON bytes.
func MarshalVM(vm VirtualMachine) ([]byte, error) {
	// TODO: Use json.Marshal(vm) and return the result.
	_ = json.Marshal
	return nil, nil
}
