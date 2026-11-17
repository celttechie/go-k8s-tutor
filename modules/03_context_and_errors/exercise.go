package context_errors

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrNotFound is a sentinel error representing a missing Kubernetes resource.
var ErrNotFound = errors.New("resource not found")

// APIError represents a structured HTTP/API error returned by the Kubernetes API server.
type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("api error (%d): %s", e.StatusCode, e.Message)
}

// FetchWithTimeout derives a child context with the given timeout, invokes fn,
// and wraps any returned error with context ("fetch failed: %w").
// Always defer cancel() so timer resources aren't leaked.
func FetchWithTimeout(ctx context.Context, timeout time.Duration, fn func(context.Context) error) error {
	// TODO: Create a timeout context with context.WithTimeout(ctx, timeout) and defer cancel().
	// Call fn(timeoutCtx) and wrap any non-nil error with fmt.Errorf("fetch failed: %w", err).
	_ = context.WithTimeout
	_ = fmt.Errorf
	return nil
}

// IsIgnorableNotFound returns true if err wraps ErrNotFound OR wraps an *APIError
// whose StatusCode is 404.
// This mirrors Kubernetes' apierrors.IsNotFound(err) pattern used in reconcilers.
func IsIgnorableNotFound(err error) bool {
	// TODO: Use errors.Is for ErrNotFound and errors.As for *APIError.
	_ = errors.Is
	return false
}

// CheckNamespacesConcurrently runs checkFn(ctx, ns) concurrently for every namespace
// using a sync.WaitGroup. It returns a map of namespace -> error for any namespace
// where checkFn returned a non-nil error.
func CheckNamespacesConcurrently(
	ctx context.Context,
	namespaces []string,
	checkFn func(context.Context, string) error,
) map[string]error {
	// TODO: Launch one goroutine per namespace, track completion with sync.WaitGroup,
	// protect map writes with sync.Mutex, and record non-nil errors.
	_ = sync.WaitGroup{}
	return map[string]error{}
}
