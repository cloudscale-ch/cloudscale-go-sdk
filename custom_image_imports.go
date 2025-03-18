package cloudscale

const customImageImportsBasePath = "v1/custom-images/import"

type CustomImageStub struct {
	HREF string `json:"href,omitempty"`
	UUID string `json:"uuid,omitempty"`
	Name string `json:"name,omitempty"`
}

type CustomImageImport struct {
	TaggedResource
	// Just use omitempty everywhere. This makes it easy to use restful. Errors
	// will be coming from the API if something is disabled.
	HREF         string          `json:"href,omitempty"`
	UUID         string          `json:"uuid,omitempty"`
	CustomImage  CustomImageStub `json:"custom_image,omitempty"`
	URL          string          `json:"url,omitempty"`
	Status       string          `json:"status,omitempty"`
	ErrorMessage string          `json:"error_message,omitempty"`
}

type CustomImageImportRequest struct {
	TaggedResourceRequest
	URL              string   `json:"url,omitempty"`
	Name             string   `json:"name,omitempty"`
	Slug             string   `json:"slug,omitempty"`
	UserDataHandling string   `json:"user_data_handling,omitempty"`
	FirmwareType     string   `json:"firmware_type,omitempty"`
	SourceFormat     string   `json:"source_format,omitempty"`
	Zones            []string `json:"zones,omitempty"`
}

type CustomImageImportsService interface {
	GenericCreateService[CustomImageImport, CustomImageImportRequest]
	GenericGetService[CustomImageImport]
	GenericListService[CustomImageImport]
	GenericWaitForService[CustomImageImport]
}
