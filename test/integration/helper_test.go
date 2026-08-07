//go:build integration

package integration

import (
	"context"
	"errors"
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v5"

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

// waitForDeleted calls existsFunc in a backoff loop until exists is false.
func waitForDeleted(ctx context.Context, existsFunc func() (exists bool, err error)) error {
	options := []backoff.RetryOption{
		backoff.WithBackOff(backoff.NewConstantBackOff(2 * time.Second)),
		backoff.WithMaxElapsedTime(5 * time.Minute),
	}

	_, err := backoff.Retry(ctx, func() (struct{}, error) {
		exists, err := existsFunc()
		if !exists {
			return struct{}{}, nil
		}
		if err == nil {
			return struct{}{}, errors.New("resource not deleted yet")
		}
		return struct{}{}, err
	}, options...)
	return err
}
