//go:build integration

package integration

import (
	"context"
	"math/rand"
	"reflect"
	"testing"

	"github.com/cloudscale-ch/cloudscale-go-sdk/v10"
)

func getAllZones() ([]cloudscale.ZoneStub, error) {
	allRegions, err := getAllRegions()
	if err != nil {
		return nil, err
	}
	allZones := []cloudscale.ZoneStub{}
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
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, length)
	for i := range length {
		bytes[i] = letters[rand.Intn(len(letters))] //gosec:disable G404 - random number not cryptographically relevant
	}
	return string(bytes)
}

// TODO: Maybe add an argument with a description for the assertion.
func assertEqual(t *testing.T, expected any, actual any) {
	t.Helper()

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Assertion failed:\nexpected: %#v\n  actual: %#v", expected, actual)
	}
}
