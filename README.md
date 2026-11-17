# go-k8s-tutor

A hands-on, self-contained project for learning the **20% of Go that powers 90% of Kubernetes tooling, CLIs, and operators** (including `client-go`, `controller-runtime`, and KubeVirt CRDs).

Instead of generic Go web-server tutorials, every module focuses on real Kubernetes idioms:
1. **Structs, Pointers, & Struct Tags** (`*int32`, `omitempty`, `json`/`yaml` marshaling of K8s manifests)
2. **Interfaces & Unit Testing** (`client.Object`, `Reconciler`, mocking without a cluster)
3. **`context.Context` & Error Wrapping** (`%w`, `errors.Is`, timeouts, `NotFound` handling)
4. **`k8s.io/client-go`** (connecting via `$KUBECONFIG`, listing/creating typed & dynamic resources)
5. **`sigs.k8s.io/controller-runtime`** (Manager, `Reconcile(ctx, req)`, Informer caches, OwnerReferences, Finalizers, and `fake.NewClientBuilder()`)

## How to Use This Repo

### Option A: With an AI Coding Tutor (Antigravity / Cursor / Claude / etc.)
Open this folder in your AI IDE. The included [`AGENTS.md`](./AGENTS.md) configures the assistant to act as a **Socratic Go-for-Kubernetes tutor**—guiding you one module at a time, giving you exercises to implement yourself, and running `go test` to verify your work without spoiling the answers.

### Option B: Solo Self-Study
1. Read [`CURRICULUM.md`](./CURRICULUM.md) for the module-by-module roadmap and open-source code reading pointers.
2. Work through `modules/01_*`, `modules/02_*`, etc.
3. Run tests as you go:
   ```bash
   go test ./modules/01_structs_pointers_tags/...
   ```

## Prerequisites
* Go 1.26+ (`go version`)
* Optional for Module 4+: a local Kubernetes cluster (`kind`, `k3s`, `minikube`) with `$KUBECONFIG` set.
