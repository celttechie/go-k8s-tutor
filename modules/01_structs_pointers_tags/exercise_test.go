package structs_pointers_tags

import (
	"encoding/json"
	"testing"
)

func TestNewVirtualMachine(t *testing.T) {
	vm := NewVirtualMachine("rocky9-test", "default", 2, "2Gi")

	if vm.APIVersion != "kubevirt.io/v1" {
		t.Fatalf("expected APIVersion kubevirt.io/v1, got %q", vm.APIVersion)
	}
	if vm.Kind != "VirtualMachine" {
		t.Fatalf("expected Kind VirtualMachine, got %q", vm.Kind)
	}
	if vm.Metadata.Name != "rocky9-test" || vm.Metadata.Namespace != "default" {
		t.Fatalf("unexpected metadata: %+v", vm.Metadata)
	}
	if vm.Spec.CPUs == nil || *vm.Spec.CPUs != 2 {
		t.Fatalf("expected Spec.CPUs == 2, got %v", vm.Spec.CPUs)
	}
	if vm.Spec.Running != nil {
		t.Fatalf("expected Spec.Running to be nil (unset), got %v", *vm.Spec.Running)
	}
	if vm.Spec.Memory != "2Gi" {
		t.Fatalf("expected Spec.Memory 2Gi, got %q", vm.Spec.Memory)
	}
}

func TestEffectiveCPUs(t *testing.T) {
	unset := VirtualMachine{}
	if got := EffectiveCPUs(unset, 4); got != 4 {
		t.Fatalf("expected default 4 when nil, got %d", got)
	}

	vm := NewVirtualMachine("vm1", "default", 8, "4Gi")
	if got := EffectiveCPUs(vm, 4); got != 8 {
		t.Fatalf("expected 8 when set, got %d", got)
	}
}

func TestMarshalVMOmitsEmptyRunning(t *testing.T) {
	vm := NewVirtualMachine("vm1", "default", 2, "2Gi")
	data, err := MarshalVM(vm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	spec, ok := decoded["spec"].(map[string]any)
	if !ok {
		t.Fatalf("expected spec object in JSON")
	}
	if _, exists := spec["running"]; exists {
		t.Fatalf("expected 'running' to be omitted when nil")
	}
}
