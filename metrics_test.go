package cloudscale

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestMetrics_GetBucketMetrics(t *testing.T) {
	setup()
	defer teardown()

	metricsRequest := &BucketMetricsRequest{
		Start:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		End:            time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		BucketNames:    []string{},
		ObjectsUserIDs: []string{},
	}

	mux.HandleFunc("/metrics/buckets", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string][]string{
			"start": {"2020-01-01"},
			"end":   {"2020-01-01"},
		}

		assertEqual(t, url.Values(expected), r.URL.Query())

		jsonStr := `
			{
				"start": "2019-12-31T23:00:00Z",
				"end": "2020-01-01T23:00:00Z",
				"data": [
					{
						"subject": {
							"name": "Carmichael",
							"objects_user_id": "2b44ef90ff15f010dc013f8d32612e07d7734e5f5581d0c94be56780b0dcd58a"
						},
						"time_series": [
							{
								"start": "2019-12-31T23:00:00Z",
								"end": "2020-01-01T23:00:00Z",
								"usage": {
									"requests": 561,
									"object_count": 1105,
									"storage_bytes": 1729,
									"received_bytes": 2465,
									"sent_bytes": 2821
								}
							}
						]
					}
				]
			}`
		fmt.Fprintf(w, jsonStr)
	})

	metrics, err := client.Metrics.GetBucketMetrics(ctx, metricsRequest)

	if err != nil {
		t.Errorf("Metrics.GetBucketMetrics returned error: %v", err)
		return
	}

	expectedMetrics := BucketMetrics{
		Start: time.Date(2019, 12, 31, 23, 0, 0, 0, time.UTC),
		End:   time.Date(2020, 1, 1, 23, 0, 0, 0, time.UTC),
		Data: []BucketMetricsData{
			{
				Subject: BucketMetricsDataSubject{
					BucketName:    "Carmichael",
					ObjectsUserID: "2b44ef90ff15f010dc013f8d32612e07d7734e5f5581d0c94be56780b0dcd58a",
				},
				TimeSeries: []BucketMetricsInterval{
					{
						Start: time.Date(2019, 12, 31, 23, 0, 0, 0, time.UTC),
						End:   time.Date(2020, 1, 1, 23, 0, 0, 0, time.UTC),
						Usage: BucketMetricsIntervalUsage{
							Requests:      561,
							ObjectCount:   1105,
							StorageBytes:  1729,
							ReceivedBytes: 2465,
							SentBytes:     2821,
						},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(*metrics, expectedMetrics) {
		t.Errorf("expected metrics '%v', received '%v'", expectedMetrics, metrics)
	}
}

func TestMetrics_GetBucketMetricsAdditionalArgs(t *testing.T) {
	setup()
	defer teardown()

	metricsRequest := &BucketMetricsRequest{
		Start:          time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		End:            time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
		BucketNames:    []string{"guete", "morge"},
		ObjectsUserIDs: []string{"hallo", "tschüss"},
	}

	mux.HandleFunc("/metrics/buckets", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string][]string{
			"start":           {"2020-01-01"},
			"end":             {"2020-01-02"},
			"bucket_name":     {"guete", "morge"},
			"objects_user_id": {"hallo", "tschüss"},
		}

		assertEqual(t, url.Values(expected), r.URL.Query())

		// Dummy response.
		fmt.Fprintf(w, "{}")
	})

	_, err := client.Metrics.GetBucketMetrics(ctx, metricsRequest)

	if err != nil {
		t.Errorf("Metrics.GetBucketMetrics returned error: %v", err)
		return
	}
}
