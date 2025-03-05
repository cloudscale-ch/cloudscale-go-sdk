package cloudscale

import (
	"fmt"
	"net/http"
	"time"
)

const volumeBasePath = "v1/volumes"

type Volume struct {
	ZonalResource
	TaggedResource
	// Just use omitempty everywhere. This makes it easy to use restful. Errors
	// will be coming from the API if something is disabled.
	HREF        string    `json:"href,omitempty"`
	UUID        string    `json:"uuid,omitempty"`
	Name        string    `json:"name,omitempty"`
	SizeGB      int       `json:"size_gb,omitempty"`
	Type        string    `json:"type,omitempty"`
	ServerUUIDs *[]string `json:"server_uuids,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type VolumeRequest struct {
	ZonalResourceRequest
	TaggedResourceRequest
	Name        string    `json:"name,omitempty"`
	SizeGB      int       `json:"size_gb,omitempty"`
	Type        string    `json:"type,omitempty"`
	ServerUUIDs *[]string `json:"server_uuids,omitempty"`
}

type VolumeService interface {
	GenericCreateService[Volume, VolumeRequest]
	GenericGetService[Volume]
	GenericListService[Volume]
	GenericUpdateService[Volume, VolumeRequest]
	GenericDeleteService[Volume]
}

// WithNameFilter uses an undocumented feature of the cloudscale.ch API
func WithNameFilter(name string) ListRequestModifier {
	return func(request *http.Request) {
		query := request.URL.Query()
		query.Add(fmt.Sprintf("name"), name)
		request.URL.RawQuery = query.Encode()
	}
}
