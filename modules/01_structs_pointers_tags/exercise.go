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
//
// Note: This exercise returns a VirtualMachine value for simplicity; in real
// client-go and controller-runtime code, API objects are almost always passed
// as pointers (*VirtualMachine) so they satisfy client.Object and avoid copying.
func NewVirtualMachine(name, namespace string, cpus int32, memory string) VirtualMachine {
	return VirtualMachine{
		APIVersion: "kubevirt.io/v1",
		Kind:       "VirtualMachine",
		Metadata: ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: VMSpec{
			CPUs:   &cpus,
			Memory: memory,
		},
	}
}

// EffectiveCPUs returns the configured CPU count, or defaultCPUs if Spec.CPUs is nil.
func EffectiveCPUs(vm VirtualMachine, defaultCPUs int32) int32 {
	if vm.Spec.CPUs != nil {
		return *vm.Spec.CPUs
	}
	return defaultCPUs
}

// MarshalVM serializes a VirtualMachine to JSON bytes.
func MarshalVM(vm VirtualMachine) ([]byte, error) {
	return json.Marshal(vm)
}
