//go:build integration
// +build integration

package integration

import (
	"context"
	"github.com/cloudscale-ch/cloudscale-go-sdk"
	"testing"
	"time"
)

func TestIntegrationMetrics_GetBuckets(t *testing.T) {
	integrationTest(t)

	request := &cloudscale.BucketMetricsRequest{
		Start:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		End:            time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		BucketNames:    []string{},
		ObjectsUserIDs: []string{},
	}

	response, err := client.Metrics.GetBuckets(context.Background(), request)
	if err != nil {
		t.Fatalf("Metrics.GetBuckets returned error %s\n", err)
	}

	// We can't get any metrics data without creating a bucket using the S3
	// API.
	assertEqual(t, time.Date(2019, 12, 31, 23, 0, 0, 0, time.UTC), response.Start)
	assertEqual(t, time.Date(2020, 1, 1, 23, 0, 0, 0, time.UTC), response.End)
}
