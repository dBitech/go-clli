package main

import (
	"testing"

	"github.com/dbitech/go-clli/pkg/clli"
)

// TestRootBuild ensures the root package builds under `go test` without conflicting mains.
func TestRootBuild(t *testing.T) {
	if _, err := clli.Parse("LSANCA12"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
