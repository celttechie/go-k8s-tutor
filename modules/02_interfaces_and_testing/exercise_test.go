package interfaces_testing

import (
	"errors"
	"testing"
)

// fakeReader is an in-memory mock implementing ClusterReader for unit tests.
type fakeReader struct {
	vms map[string]Resource
	err error
}

func (f *fakeReader) GetVM(namespace, name string) (Resource, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.vms[namespace+"/"+name], nil
}

func TestSimpleVMImplementsResource(t *testing.T) {
	var r Resource = &SimpleVM{Namespace: "prod", Name: "app-01"}
	if r.GetName() != "app-01" || r.GetNamespace() != "prod" {
		t.Fatalf("unexpected Resource values: %s/%s", r.GetNamespace(), r.GetName())
	}
}

func TestCheckVMExists(t *testing.T) {
	reader := &fakeReader{
		vms: map[string]Resource{
			"default/web-01": &SimpleVM{Namespace: "default", Name: "web-01"},
		},
	}

	ok, err := CheckVMExists(reader, "default", "web-01")
	if err != nil || !ok {
		t.Fatalf("expected (true, nil), got (%v, %v)", ok, err)
	}

	ok, err = CheckVMExists(reader, "default", "missing")
	if err != nil || ok {
		t.Fatalf("expected (false, nil), got (%v, %v)", ok, err)
	}

	boom := errors.New("api server unreachable")
	errReader := &fakeReader{err: boom}
	_, err = CheckVMExists(errReader, "default", "web-01")
	if !errors.Is(err, boom) {
		t.Fatalf("expected wrapped boom error, got %v", err)
	}
}
