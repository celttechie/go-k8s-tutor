package context_errors

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestFetchWithTimeout(t *testing.T) {
	// Case 1: fast function succeeds
	err := FetchWithTimeout(context.Background(), 100*time.Millisecond, func(ctx context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// Case 2: slow function hits deadline
	err = FetchWithTimeout(context.Background(), 15*time.Millisecond, func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected wrapped context.DeadlineExceeded, got %v", err)
	}
}

func TestIsIgnorableNotFound(t *testing.T) {
	wrappedSentinel := fmt.Errorf("getting vm: %w", ErrNotFound)
	if !IsIgnorableNotFound(wrappedSentinel) {
		t.Fatalf("expected true for wrapped ErrNotFound")
	}

	wrapped404 := fmt.Errorf("getting pod: %w", &APIError{StatusCode: 404, Message: "pods \"web\" not found"})
	if !IsIgnorableNotFound(wrapped404) {
		t.Fatalf("expected true for wrapped 404 *APIError")
	}

	wrapped500 := fmt.Errorf("getting pod: %w", &APIError{StatusCode: 500, Message: "internal error"})
	if IsIgnorableNotFound(wrapped500) {
		t.Fatalf("expected false for 500 *APIError")
	}

	if IsIgnorableNotFound(nil) {
		t.Fatalf("expected false for nil error")
	}
}

func TestCheckNamespacesConcurrently(t *testing.T) {
	namespaces := []string{"default", "kube-system", "virt-prod", "broken-ns"}
	var started atomic.Int32
	boom := errors.New("rbacs denied")

	errs := CheckNamespacesConcurrently(context.Background(), namespaces, func(ctx context.Context, ns string) error {
		started.Add(1)
		time.Sleep(10 * time.Millisecond)
		if ns == "broken-ns" {
			return boom
		}
		return nil
	})

	if int(started.Load()) != len(namespaces) {
		t.Fatalf("expected %d checks to run, got %d", len(namespaces), started.Load())
	}
	if len(errs) != 1 || !errors.Is(errs["broken-ns"], boom) {
		t.Fatalf("expected only broken-ns error, got %v", errs)
	}
}
