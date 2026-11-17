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
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := fn(timeoutCtx); err != nil {
		return fmt.Errorf("fetch failed: %w", err)
	}
	return nil
}

// IsIgnorableNotFound returns true if err wraps ErrNotFound OR wraps an *APIError
// whose StatusCode is 404.
// This mirrors Kubernetes' apierrors.IsNotFound(err) pattern used in reconcilers.
func IsIgnorableNotFound(err error) bool {
	if errors.Is(err, ErrNotFound) {
		return true
	}
	var apiErr *APIError
	return errors.As(err, &apiErr) && apiErr.StatusCode == 404
}

// CheckNamespacesConcurrently runs checkFn(ctx, ns) concurrently for every namespace
// using a sync.WaitGroup. It returns a map of namespace -> error for any namespace
// where checkFn returned a non-nil error.
func CheckNamespacesConcurrently(
	ctx context.Context,
	namespaces []string,
	checkFn func(context.Context, string) error,
) map[string]error {
	errs := make(map[string]error)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, ns := range namespaces {
		wg.Add(1)
		go func(namespace string) {
			defer wg.Done()
			if err := checkFn(ctx, namespace); err != nil {
				mu.Lock()
				errs[namespace] = err
				mu.Unlock()
			}
		}(ns)
	}

	wg.Wait()
	return errs
}
