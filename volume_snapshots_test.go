package cloudscale

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestVolumeSnapshots_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/volume-snapshots", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		// Verify the request body has correct field name
		var req VolumeSnapshotRequest
		json.NewDecoder(r.Body).Decode(&req)
		if req.SourceVolume != "volume-uuid" {
			t.Errorf("Expected SourceVolume to be 'volume-uuid', got '%s'", req.SourceVolume)
		}

		fmt.Fprint(w, `{"uuid": "snapshot-uuid", "name": "test-snapshot"}`)
	})

	volumeCreateRequest := &VolumeSnapshotRequest{
		Name:         "test-snapshot",
		SourceVolume: "volume-uuid",
	}

	snapshot, err := client.VolumeSnapshots.Create(context.Background(), volumeCreateRequest)

	if err != nil {
		t.Errorf("VolumeSnapshots.Create returned error: %v", err)
	}

	expected := &VolumeSnapshot{UUID: "snapshot-uuid", Name: "test-snapshot"}
	if !reflect.DeepEqual(snapshot, expected) {
		t.Errorf("VolumeSnapshots.Create\n got=%#v\nwant=%#v", snapshot, expected)
	}
}

func TestVolumeSnapshots_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/volume-snapshots/snapshot-uuid", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
          "uuid": "snapshot-uuid",
          "name": "test-snapshot",
          "size_gb": 50,
          "created_at": "2024-01-15T10:30:00Z",
          "volume": {
             "uuid": "volume-uuid"
          },
          "zone": {
             "slug": "lpg1"
          },
          "tags": {}
       }`)
	})

	snapshot, err := client.VolumeSnapshots.Get(context.Background(), "snapshot-uuid")

	if err != nil {
		t.Errorf("VolumeSnapshots.Get returned error: %v", err)
	}

	expected := &VolumeSnapshot{
		UUID:      "snapshot-uuid",
		Name:      "test-snapshot",
		SizeGB:    50,
		CreatedAt: "2024-01-15T10:30:00Z",
		Volume: VolumeStub{
			UUID: "volume-uuid",
		},
		ZonalResource: ZonalResource{
			Zone{
				Slug: "lpg1",
			},
		},
		TaggedResource: TaggedResource{
			Tags: TagMap{},
		},
	}

	if !reflect.DeepEqual(snapshot, expected) {
		t.Errorf("VolumeSnapshots.Get\n got=%#v\nwant=%#v", snapshot, expected)
	}
}

func TestVolumeSnapshots_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/volume-snapshots/snapshot-uuid", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPatch)

		var req VolumeSnapshotUpdateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if req.Name != "updated-snapshot" {
			t.Errorf("Expected Name to be 'updated-snapshot', got '%s'", req.Name)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	updateRequest := &VolumeSnapshotUpdateRequest{
		Name: "updated-snapshot",
	}

	err := client.VolumeSnapshots.Update(context.Background(), "snapshot-uuid", updateRequest)

	if err != nil {
		t.Errorf("VolumeSnapshots.Update returned error: %v", err)
	}
}

func TestVolumeSnapshots_Update_WithTags(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/volume-snapshots/snapshot-uuid", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPatch)

		var req VolumeSnapshotUpdateRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if req.Tags == nil {
			t.Error("Expected Tags to be set")
		}

		if (*req.Tags)["environment"] != "production" {
			t.Errorf("Expected tag 'environment' to be 'production', got '%s'", (*req.Tags)["environment"])
		}

		w.WriteHeader(http.StatusNoContent)
	})

	tags := TagMap{"environment": "production"}

	updateRequest := &VolumeSnapshotUpdateRequest{
		TaggedResourceRequest: TaggedResourceRequest{
			Tags: &tags,
		},
	}

	err := client.VolumeSnapshots.Update(context.Background(), "snapshot-uuid", updateRequest)

	if err != nil {
		t.Errorf("VolumeSnapshots.Update returned error: %v", err)
	}
}

func TestVolumeSnapshots_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/volume-snapshots/snapshot-uuid", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.VolumeSnapshots.Delete(context.Background(), "snapshot-uuid")

	if err != nil {
		t.Errorf("VolumeSnapshots.Delete returned error: %v", err)
	}
}

func TestVolumeSnapshots_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/volume-snapshots", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[
          {
             "uuid": "snapshot-uuid-1",
             "name": "snapshot-1",
             "size_gb": 50,
             "created_at": "2024-01-15T10:30:00Z",
             "volume": {
                "uuid": "volume-uuid-1"
             },
             "zone": {
                "slug": "lpg1"
             },
             "tags": {}
          },
          {
             "uuid": "snapshot-uuid-2",
             "name": "snapshot-2",
             "size_gb": 100,
             "created_at": "2024-01-16T11:00:00Z",
             "volume": {
                "uuid": "volume-uuid-2"
             },
             "zone": {
                "slug": "rma1"
             },
             "tags": {"environment": "test"}
          }
       ]`)
	})

	snapshots, err := client.VolumeSnapshots.List(context.Background())

	if err != nil {
		t.Errorf("VolumeSnapshots.List returned error: %v", err)
	}

	if len(snapshots) != 2 {
		t.Errorf("Expected 2 snapshots, got %d", len(snapshots))
	}

	expected := []VolumeSnapshot{
		{
			UUID:      "snapshot-uuid-1",
			Name:      "snapshot-1",
			SizeGB:    50,
			CreatedAt: "2024-01-15T10:30:00Z",
			Volume: VolumeStub{
				UUID: "volume-uuid-1",
			},
			ZonalResource: ZonalResource{
				Zone{
					Slug: "lpg1",
				},
			},
			TaggedResource: TaggedResource{
				Tags: TagMap{},
			},
		},
		{
			UUID:      "snapshot-uuid-2",
			Name:      "snapshot-2",
			SizeGB:    100,
			CreatedAt: "2024-01-16T11:00:00Z",
			Volume: VolumeStub{
				UUID: "volume-uuid-2",
			},
			ZonalResource: ZonalResource{
				Zone{
					Slug: "rma1",
				},
			},
			TaggedResource: TaggedResource{
				Tags: TagMap{"environment": "test"},
			},
		},
	}

	if !reflect.DeepEqual(snapshots, expected) {
		t.Errorf("VolumeSnapshots.List\n got=%#v\nwant=%#v", snapshots, expected)
	}
}
func TestVolumeSnapshots_List_WithFilters(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/volume-snapshots", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)

		// Check if filter query params are present
		query := r.URL.Query()
		if query.Get("tag:environment") != "production" {
			t.Errorf("Expected tag:environment filter to be 'production', got '%s'", query.Get("tag:environment"))
		}

		fmt.Fprint(w, `[
          {
             "uuid": "snapshot-uuid-1",
             "name": "snapshot-1",
             "size_gb": 50,
             "created_at": "2024-01-15T10:30:00Z",
             "volume": {
                "uuid": "volume-uuid-1",
                "name": "volume-1"
             },
             "zone": {
                "slug": "lpg1"
             },
             "tags": {"environment": "production"}
          }
       ]`)
	})

	tagFilter := TagMap{"environment": "production"}
	snapshots, err := client.VolumeSnapshots.List(
		context.Background(),
		WithTagFilter(tagFilter),
	)

	if err != nil {
		t.Errorf("VolumeSnapshots.List returned error: %v", err)
	}

	if len(snapshots) != 1 {
		t.Errorf("Expected 1 snapshot, got %d", len(snapshots))
	}

	if snapshots[0].UUID != "snapshot-uuid-1" {
		t.Errorf("Expected snapshot UUID 'snapshot-uuid-1', got '%s'", snapshots[0].UUID)
	}
}
