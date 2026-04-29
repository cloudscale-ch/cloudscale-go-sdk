package cloudscale

const volumeSnapshotsBasePath = "v1/volume-snapshots"

type SourceVolumeStub struct {
	HREF string `json:"href"`
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type VolumeSnapshot struct {
	ZonalResource
	TaggedResource
	HREF         string           `json:"href,omitempty"`
	UUID         string           `json:"uuid,omitempty"`
	Name         string           `json:"name,omitempty"`
	SizeGB       int              `json:"size_gb,omitempty"`
	CreatedAt    string           `json:"created_at,omitempty"`
	SourceVolume SourceVolumeStub `json:"source_volume,omitempty"`
	Status       string           `json:"status,omitempty"`
}

type VolumeSnapshotCreateRequest struct {
	TaggedResourceRequest
	Name         string `json:"name,omitempty"`
	SourceVolume string `json:"source_volume,omitempty"`
}

type VolumeSnapshotUpdateRequest struct {
	TaggedResourceRequest
	Name string `json:"name,omitempty"`
}

type VolumeSnapshotService interface {
	GenericCreateService[VolumeSnapshot, VolumeSnapshotCreateRequest]
	GenericGetService[VolumeSnapshot]
	GenericListService[VolumeSnapshot]
	GenericUpdateService[VolumeSnapshot, VolumeSnapshotUpdateRequest]
	GenericDeleteService[VolumeSnapshot]
	GenericWaitForService[VolumeSnapshot]
}
