package cloudscale

import (
	"context"
)

type VolumeSnapshot struct {
	ZonalResource
	TaggedResource
	HREF      string     `json:"href,omitempty"`
	UUID      string     `json:"uuid,omitempty"`
	Name      string     `json:"name,omitempty"`
	SizeGB    int        `json:"size_gb,omitempty"`
	CreatedAt string     `json:"created_at,omitempty"`
	Volume    VolumeStub `json:"volume,omitempty"`
}

type VolumeSnapshotRequest struct {
	TaggedResourceRequest
	Name         string `json:"name,omitempty"`
	SourceVolume string `json:"source_volume,omitempty"`
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
