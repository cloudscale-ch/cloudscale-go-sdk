// +build integration

package integration

import (
	"context"
	"github.com/cloudscale-ch/cloudscale-go-sdk"
	"testing"
)

func getAllZones(t *testing.T) []cloudscale.Zone {
	allRegions := getAllRegions(t)
	allZones := []cloudscale.Zone{}
	for _, region := range allRegions {
		allZones = append(allZones, region.Zones...)
	}
	return allZones
}

func getAllRegions(t *testing.T) []cloudscale.Region {
	allRegions, err := client.Regions.List(context.Background())
	if err != nil {
		t.Fatalf("Regions.List returned error %s\n", err)
	}
	return allRegions
}
