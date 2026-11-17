# Go for Kubernetes Curriculum

This curriculum focuses on the exact Go patterns used in production Kubernetes tools, operators, and virtualization platforms (such as `kubevirt`, `cdi`, `forklift`, and `zarf`).

---

## Module 1: Structs, Pointers, & Struct Tags
**Location:** `modules/01_structs_pointers_tags/`

* **Concepts:**
  * Nested Go structs representing Kubernetes `TypeMeta`, `ObjectMeta`, `Spec`, and `Status`.
  * Why Kubernetes APIs use pointer fields (`*int32`, `*bool`) to distinguish "unset" from `0` / `false`.
  * Struct tags (`json:"replicas,omitempty"`) and `encoding/json` marshaling.
* **Open-Source Reading:**
  * `kubevirt/containerized-data-importer` (`staging/src/kubevirt.io/containerized-data-importer-api/pkg/apis/core/v1beta1/types.go`)

---

## Module 2: Interfaces & Unit Testing Without a Cluster
**Location:** `modules/02_interfaces_and_testing/`

* **Concepts:**
  * Implicit interface satisfaction in Go (no `implements` keyword).
  * Value vs. pointer receivers and how they affect interface satisfaction (`client.Object`).
  * Dependency injection via interfaces so Kubernetes tools can be unit-tested with fakes/mocks.
* **Open-Source Reading:**
  * `zarf-dev/zarf` (`src/pkg/cluster/`)

---

## Module 3: `context.Context`, Error Wrapping, & Concurrency
**Location:** `modules/03_context_and_errors/`

* **Concepts:**
  * Passing `ctx context.Context` for timeouts and cancellation (`context.WithTimeout`).
  * Wrapping errors with `fmt.Errorf("...: %w", err)` and inspecting them with `errors.Is` / `errors.As`.
  * Goroutines, channels, and `sync.WaitGroup` for concurrent cluster/VM checks.

---

## Module 4: `k8s.io/client-go` (Talking to a Real Cluster)
**Location:** `modules/04_client_go/`

* **Concepts:**
  * Loading `$KUBECONFIG` via `clientcmd.BuildConfigFromFlags` (and in-cluster config).
  * Listing Nodes (checking `/dev/kvm` allocatable resources), Pods, and `StorageClass`es.
  * Handling `apierrors.IsNotFound(err)` and `apierrors.IsAlreadyExists(err)`.

---

## Module 5: `sigs.k8s.io/controller-runtime` (The Operator Loop)
**Location:** `modules/05_controller_runtime/`

* **Concepts:**
  * `Manager`, `Reconcile(ctx, ctrl.Request) (ctrl.Result, error)`, and informer caches.
  * `Get`, `Create`, `Patch`, `Status().Update()`, `OwnerReference`s, and `Finalizer`s.
  * Writing fast, deterministic controller unit tests using `sigs.k8s.io/controller-runtime/pkg/client/fake`.
* **Open-Source Reading:**
  * `kubev2v/forklift` (`pkg/controller/`)
  * `kubevirt/containerized-data-importer` (`pkg/controller/`)
