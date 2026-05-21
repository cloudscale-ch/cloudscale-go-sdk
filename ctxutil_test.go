package cloudscale

import (
	"context"
	"testing"
)

func TestOperationPath(t *testing.T) {
	ctx := context.Background()

	// No path attached.
	if got := OperationPath(ctx); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}

	ctx = WithOperationPath(ctx, "v1/servers/:id")
	if got := OperationPath(ctx); got != "v1/servers/:id" {
		t.Fatalf("expected %q, got %q", "v1/servers/:id", got)
	}

	// Nested context should inherit.
	child := context.WithValue(ctx, struct{}{}, "other")
	if got := OperationPath(child); got != "v1/servers/:id" {
		t.Fatalf("expected %q, got %q", "v1/servers/:id", got)
	}
}
