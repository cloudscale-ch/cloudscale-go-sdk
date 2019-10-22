// +build integration

package integration

import (
	"context"
	"testing"
)

func TestListRegions(t *testing.T) {
	integrationTest(t)

	allRegions, err := client.Regions.List(context.Background())
	if err != nil {
		t.Fatalf("Regions.List returned error %s\n", err)
	}

	if len(allRegions) <= 0 {
		t.Fatal("Regions.List returned empty slice\n", err)
	}

	// Check the result for at least one Zone to keep test case as generic as possible
	foundZone := false
	for _, region := range allRegions {
		if len(region.Zones) >= 1 {
			foundZone = true
		}
	}

	if !foundZone {
		t.Fatal("No zones found.")
	}
}
