// +build integration

package integration

import (
	"context"
	"github.com/cloudscale-ch/cloudscale-go-sdk"
)

func getAllZones() ([]cloudscale.Zone, error) {
	allRegions, err := getAllRegions()
	if err != nil {
		return nil, err
	}
	allZones := []cloudscale.Zone{}
	for _, region := range allRegions {
		allZones = append(allZones, region.Zones...)
	}
	return allZones, nil
}

func getAllRegions() ([]cloudscale.Region, error) {
	allRegions, err := client.Regions.List(context.Background())
	return allRegions, err
}
