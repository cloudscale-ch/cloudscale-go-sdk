//go:build integration
// +build integration

package integration

import (
	"context"
	"github.com/cloudscale-ch/cloudscale-go-sdk/v3"
	"testing"
	"time"
)

func TestIntegrationMetrics_GetBucketMetrics(t *testing.T) {
	integrationTest(t)

	objectsUser, err := createObjectsUser(t)
	if err != nil {
		t.Fatalf("ObjectsUsers.Create returned error %s\n", err)
	}

	request := &cloudscale.BucketMetricsRequest{
		Start:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		End:            time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		BucketNames:    []string{},
		ObjectsUserIDs: []string{objectsUser.ID},
	}

	response, err := client.Metrics.GetBucketMetrics(context.Background(), request)
	if err != nil {
		t.Fatalf("Metrics.GetBucketMetrics returned error %s\n", err)
	}

	// We can't get any metrics data without creating a bucket using the S3
	// API.
	assertEqual(t, time.Date(2019, 12, 31, 23, 0, 0, 0, time.UTC), response.Start)
	assertEqual(t, time.Date(2020, 1, 1, 23, 0, 0, 0, time.UTC), response.End)

	err = client.ObjectsUsers.Delete(context.Background(), objectsUser.ID)
	if err != nil {
		t.Fatalf("ObjectsUsers.Get returned error %s\n", err)
	}
}
