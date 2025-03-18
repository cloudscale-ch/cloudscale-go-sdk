package cloudscale

import (
	"time"
)

const customImagesBasePath = "v1/custom-images"

const UserDataHandlingPassThrough = "pass-through"
const UserDataHandlingExtendCloudConfig = "extend-cloud-config"

type CustomImage struct {
	ZonalResource
	TaggedResource
	// Just use omitempty everywhere. This makes it easy to use restful. Errors
	// will be coming from the API if something is disabled.
	HREF             string            `json:"href,omitempty"`
	UUID             string            `json:"uuid,omitempty"`
	Name             string            `json:"name,omitempty"`
	Slug             string            `json:"slug,omitempty"`
	SizeGB           int               `json:"size_gb,omitempty"`
	Checksums        map[string]string `json:"checksums,omitempty"`
	UserDataHandling string            `json:"user_data_handling,omitempty"`
	FirmwareType     string            `json:"firmware_type,omitempty"`
	Zones            []Zone            `json:"zones"`
	CreatedAt        time.Time         `json:"created_at"`
}

type CustomImageRequest struct {
	TaggedResourceRequest
	Name             string `json:"name,omitempty"`
	Slug             string `json:"slug,omitempty"`
	UserDataHandling string `json:"user_data_handling,omitempty"`
}

type CustomImageService interface {
	GenericGetService[CustomImage]
	GenericListService[CustomImage]
	GenericUpdateService[CustomImage, CustomImageRequest]
	GenericDeleteService[CustomImage]
	GenericWaitForService[CustomImage]
}

type CustomImageServiceOperations struct {
	client *Client
}
