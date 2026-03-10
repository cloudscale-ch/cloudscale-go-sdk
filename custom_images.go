package cloudscale

import (
	"fmt"
	"time"
)

const customImagesBasePath = "v1/custom-images"

type UserDataHandling string

const UserDataHandlingPassThrough UserDataHandling = "pass-through"
const UserDataHandlingExtendCloudConfig UserDataHandling = "extend-cloud-config"

type CustomImage struct {
	TaggedResource
	// Just use omitempty everywhere. This makes it easy to use restful. Errors
	// will be coming from the API if something is disabled.
	HREF             string            `json:"href,omitempty"`
	UUID             string            `json:"uuid,omitempty"`
	Name             string            `json:"name,omitempty"`
	Slug             string            `json:"slug,omitempty"`
	SizeGB           int               `json:"size_gb,omitempty"`
	Checksums        map[string]string `json:"checksums,omitempty"`
	UserDataHandling UserDataHandling  `json:"user_data_handling,omitempty"`
	FirmwareType     string            `json:"firmware_type,omitempty"`
	Zones            []ZoneStub        `json:"zones"`
	CreatedAt        time.Time         `json:"created_at"`
}

type CustomImageRequest struct {
	TaggedResourceRequest
	Name             string           `json:"name,omitempty"`
	Slug             string           `json:"slug,omitempty"`
	UserDataHandling UserDataHandling `json:"user_data_handling,omitempty"`
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

var ImportIsSuccessful = func(importInfo *CustomImageImport) (bool, error) {
	if importInfo.Status == "success" {
		return true, nil
	}
	return false, fmt.Errorf("waiting for status: %s, current status: %s", "success", importInfo.Status)
}
