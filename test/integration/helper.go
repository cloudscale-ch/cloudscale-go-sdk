// +build integration

package integration

import (
	"context"
	"github.com/cloudscale-ch/cloudscale-go-sdk"
	"math/rand"
	"time"
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

func randomNotVerySecurePassword(length int) string {
	// based on: https://stackoverflow.com/a/12321192
	rand.Seed(time.Now().UTC().UnixNano())

	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
