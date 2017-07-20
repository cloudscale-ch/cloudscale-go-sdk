package cloudscale

import "testing"

func acceptanceTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test")
	}
}

func TestAcceptanceServer_Create(t *testing.T) {
	acceptanceTest(t)
}
