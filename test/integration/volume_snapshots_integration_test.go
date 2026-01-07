//go:build integration

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cloudscale-ch/cloudscale-go-sdk/v6"
)

func TestIntegrationVolumeSnapshot_CRUD(t *testing.T) {
	integrationTest(t)

	ctx := context.Background()

	// A source volume is needed to create a snapshot.
	volumeCreateRequest := &cloudscale.VolumeCreateRequest{
		Name:   "test-volume-for-snapshot",
		SizeGB: 50,
		Type:   "ssd",
		ZonalResourceRequest: cloudscale.ZonalResourceRequest{
			Zone: testZone,
		},
	}
	volume, err := client.Volumes.Create(ctx, volumeCreateRequest)
	if err != nil {
		t.Fatalf("Volume.Create: %v", err)
	}

	snapshotCreateRequest := &cloudscale.VolumeSnapshotCreateRequest{
		Name:         "test-snapshot",
		SourceVolume: volume.UUID,
	}
	snapshot, err := client.VolumeSnapshots.Create(ctx, snapshotCreateRequest)
	if err != nil {
		t.Fatalf("VolumeSnapshots.Create: %v", err)
	}

	retrieved, err := client.VolumeSnapshots.Get(ctx, snapshot.UUID)
	if err != nil {
		t.Fatalf("VolumeSnapshots.Get: %v", err)
	}
	if retrieved.UUID != snapshot.UUID {
		t.Errorf("Expected UUID %s, got %s", snapshot.UUID, retrieved.UUID)
	}
	if retrieved.Name != "test-snapshot" {
		t.Errorf("Expected retrieved snapshot name 'test-snapshot', got '%s'", retrieved.Name)
	}

	snapshots, err := client.VolumeSnapshots.List(ctx)
	if err != nil {
		t.Fatalf("VolumeSnapshots.List: %v", err)
	}
	if len(snapshots) == 0 {
		t.Error("Expected at least one snapshot")
	}

	if err := client.VolumeSnapshots.Delete(ctx, snapshot.UUID); err != nil {
		t.Fatalf("Warning: failed to delete snapshot %s: %v", snapshot.UUID, err)
	}

	// Wait for snapshot to be fully deleted before deleting volume
	err = waitForSnapshotDeletion(ctx, snapshot.UUID, 10)
	if err != nil {
		t.Fatalf("Snapshot deletion timeout: %v", err)
	}

	if err := client.Volumes.Delete(ctx, volume.UUID); err != nil {
		t.Fatalf("Warning: failed to delete volume %s: %v", volume.UUID, err)
	}
}

func TestIntegrationVolumeSnapshot_Update(t *testing.T) {
	integrationTest(t)

	ctx := context.Background()

	// A source volume is needed to create a snapshot.
	volumeCreateRequest := &cloudscale.VolumeCreateRequest{
		Name:   "test-volume-for-snapshot",
		SizeGB: 50,
		Type:   "ssd",
		ZonalResourceRequest: cloudscale.ZonalResourceRequest{
			Zone: testZone,
		},
	}
	volume, err := client.Volumes.Create(ctx, volumeCreateRequest)
	if err != nil {
		t.Fatalf("Volume.Create: %v", err)
	}

	snapshotCreateRequest := &cloudscale.VolumeSnapshotCreateRequest{
		Name:         "test-snapshot",
		SourceVolume: volume.UUID,
	}
	snapshot, err := client.VolumeSnapshots.Create(ctx, snapshotCreateRequest)
	if err != nil {
		t.Fatalf("VolumeSnapshots.Create: %v", err)
	}

	snapshotUpdateRequest := &cloudscale.VolumeSnapshotUpdateRequest{
		Name: "updated-snapshot",
	}
	err = client.VolumeSnapshots.Update(ctx, snapshot.UUID, snapshotUpdateRequest)
	if err != nil {
		t.Fatalf("VolumeSnapshots.Update: %v", err)
	}

	// Get snapshot again to verify the update
	updatedSnapshot, err := client.VolumeSnapshots.Get(ctx, snapshot.UUID)
	if err != nil {
		t.Fatalf("VolumeSnapshots.Get after update: %v", err)
	}
	if updatedSnapshot.Name != "updated-snapshot" {
		t.Errorf("Expected updated snapshot name 'updated-snapshot', got '%s'", updatedSnapshot.Name)
	}

	if err := client.VolumeSnapshots.Delete(ctx, snapshot.UUID); err != nil {
		t.Fatalf("Warning: failed to delete snapshot %s: %v", snapshot.UUID, err)
	}

	// Wait for snapshot to be fully deleted before deleting volume
	err = waitForSnapshotDeletion(ctx, snapshot.UUID, 10)
	if err != nil {
		t.Fatalf("Snapshot deletion timeout: %v", err)
	}

	if err := client.Volumes.Delete(ctx, volume.UUID); err != nil {
		t.Fatalf("Warning: failed to delete volume %s: %v", volume.UUID, err)
	}
}

// waitForSnapshotDeletion polls the API until the snapshot no longer exists
func waitForSnapshotDeletion(ctx context.Context, snapshotUUID string, maxWaitSeconds int) error {
	for i := 0; i < maxWaitSeconds; i++ {
		snapshot, err := client.VolumeSnapshots.Get(ctx, snapshotUUID)
		if err != nil {

			if apiErr, ok := err.(*cloudscale.ErrorResponse); ok {
				if apiErr.StatusCode == 404 {
					// if we get a 404 error, snapshot is gone, deletion completed
					return nil
				}
			}
			// some other error occurred
			return err
		}

		// if snapshot still exists, it must be in state deleting
		if snapshot.Status != "deleting" {
			return fmt.Errorf(
				"snapshot %s exists but is in unexpected state %q while waiting for deletion",
				snapshotUUID,
				snapshot.Status,
			)
		}

		// snapshot still exists, wait 1 second and try again
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("snapshot %s still exists after %d seconds", snapshotUUID, maxWaitSeconds)
}
