package cloudscale

import "context"

type VolumeSnapshot struct {
	HREF      string     `json:"href"`
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	SizeGB    int        `json:"size_gb"`
	CreatedAt string     `json:"created_at"`
	Volume    VolumeStub `json:"volume"`
	Zone      Zone       `json:"zone"`
	Tags      *TagMap    `json:"tags"`
}

type VolumeSnapshotRequest struct {
	Name         string  `json:"name"`
	SourceVolume string  `json:"source_volume"`
	Tags         *TagMap `json:"tags,omitempty"`
}

type VolumeSnapshotUpdateRequest struct {
	Name string  `json:"name,omitempty"`
	Tags *TagMap `json:"tags,omitempty"`
}

const volumeSnapshotsBasePath = "v1/volume-snapshots"

type VolumeSnapshotService interface {
	Create(ctx context.Context, createRequest *VolumeSnapshotRequest) (*VolumeSnapshot, error)
	Get(ctx context.Context, snapshotID string) (*VolumeSnapshot, error)
	Update(ctx context.Context, snapshotID string, updateRequest *VolumeSnapshotUpdateRequest) error
	Delete(ctx context.Context, snapshotID string) error
	List(ctx context.Context, opts ...ListRequestModifier) ([]VolumeSnapshot, error)
}
