package cloudscale

import (
	"context"
)

type VolumeSnapshot struct {
	ZonalResource
	TaggedResource
	HREF      string     `json:"href"`
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	SizeGB    int        `json:"size_gb"`
	CreatedAt string     `json:"created_at"`
	Volume    VolumeStub `json:"volume"`
}

type VolumeSnapshotRequest struct {
	TaggedResourceRequest
	Name         string `json:"name"`
	SourceVolume string `json:"source_volume"`
}

type VolumeSnapshotUpdateRequest struct {
	TaggedResourceRequest
	Name string `json:"name,omitempty"`
}

const volumeSnapshotsBasePath = "v1/volume-snapshots"

type VolumeSnapshotService interface {
	Create(ctx context.Context, createRequest *VolumeSnapshotRequest) (*VolumeSnapshot, error)
	Get(ctx context.Context, snapshotID string) (*VolumeSnapshot, error)
	Update(ctx context.Context, snapshotID string, updateRequest *VolumeSnapshotUpdateRequest) error
	Delete(ctx context.Context, snapshotID string) error
	List(ctx context.Context, opts ...ListRequestModifier) ([]VolumeSnapshot, error)
}
