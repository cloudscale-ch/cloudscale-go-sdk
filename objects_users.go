package cloudscale

const objectsUsersBasePath = "v1/objects-users"

// ObjectsUser contains information
type ObjectsUser struct {
	TaggedResource
	HREF        string              `json:"href,omitempty"`
	ID          string              `json:"id,omitempty"`
	DisplayName string              `json:"display_name,omitempty"`
	Keys        []map[string]string `json:"keys,omitempty"`
}

// ObjectsUserRequest is used to create and update Objects Users
type ObjectsUserRequest struct {
	TaggedResourceRequest
	DisplayName string `json:"display_name,omitempty"`
}

// ObjectsUsersService manages users of the S3-compatible objects storage
type ObjectsUsersService interface {
	GenericCreateService[ObjectsUser, ObjectsUserRequest]
	GenericGetService[ObjectsUser]
	GenericListService[ObjectsUser]
	GenericUpdateService[ObjectsUser, ObjectsUserRequest]
	GenericDeleteService[ObjectsUser]
	GenericWaitForService[ObjectsUser]
}
