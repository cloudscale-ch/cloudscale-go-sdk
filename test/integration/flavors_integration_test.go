//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/cloudscale-ch/cloudscale-go-sdk/v9"
)

func TestListFlavors(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping: short flag passed")
	}
	t.Parallel()

	allFlavors, err := client.Flavors.List(context.Background())
	if err != nil {
		t.Fatalf("Flavors.List returned error %s\n", err)
	}

	if len(allFlavors) <= 0 {
		t.Fatal("Flavors.List returned empty slice\n", err)
	}

	flavorsBySlug := make(map[string]cloudscale.Flavor)
	for _, f := range allFlavors {
		flavorsBySlug[f.Slug] = f
	}

	// Verify a flex (shared vCPU) flavor.
	assertFlavor(t, flavorsBySlug, "flex-4-1", 1, 4, nil)

	// Verify a plus (dedicated CPU) flavor.
	assertFlavor(t, flavorsBySlug, "plus-192-96", 96, 192, nil)

	// Verify a FlavorGPU flavor.
	assertFlavor(t, flavorsBySlug, "gpu2-128-16-1-200", 16, 128, &cloudscale.FlavorGPU{
		Count: 1,
	})
}

func assertFlavor(t *testing.T, flavors map[string]cloudscale.Flavor, slug string, expectedVCPUs int, expectedMemoryGB int, expectedGPU *cloudscale.FlavorGPU) {
	t.Helper()

	f, ok := flavors[slug]
	if !ok {
		t.Fatalf("expected flavor %q not found in list", slug)
	}
	if f.Name == "" {
		t.Errorf("flavor %q: expected non-empty name", slug)
	}
	if f.VCPUCount != expectedVCPUs {
		t.Errorf("flavor %q: expected vcpu_count=%d, got %d", slug, expectedVCPUs, f.VCPUCount)
	}
	if f.MemoryGB != expectedMemoryGB {
		t.Errorf("flavor %q: expected memory_gb=%d, got %d", slug, expectedMemoryGB, f.MemoryGB)
	}
	if len(f.Zones) == 0 {
		t.Errorf("flavor %q: expected at least one zone", slug)
	}

	if expectedGPU == nil {
		if f.GPU != nil {
			t.Errorf("flavor %q: expected gpu=nil, got %+v", slug, f.GPU)
		}
	} else {
		if f.GPU == nil {
			t.Fatalf("flavor %q: expected gpu to be non-nil", slug)
		}
		if f.GPU.Count != expectedGPU.Count {
			t.Errorf("flavor %q: expected gpu.count=%d, got %d", slug, expectedGPU.Count, f.GPU.Count)
		}
		if f.GPU.Name == "" {
			t.Errorf("flavor %q: expected non-empty gpu.name", slug)
		}
		if f.GPU.VRAMPerGPUGB <= 0 {
			t.Errorf("flavor %q: expected gpu.vram_per_gpu_gb > 0, got %d", slug, f.GPU.VRAMPerGPUGB)
		}
	}
}
