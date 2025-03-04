package cloudscale

const serverGroupsBasePath = "v1/server-groups"

type ServerGroup struct {
	ZonalResource
	TaggedResource
	HREF    string       `json:"href"`
	UUID    string       `json:"uuid"`
	Name    string       `json:"name"`
	Type    string       `json:"type"`
	Servers []ServerStub `json:"servers"`
}

type ServerGroupRequest struct {
	ZonalResourceRequest
	TaggedResourceRequest
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

type ServerGroupService interface {
	GenericCreateService[ServerGroup, ServerGroupRequest]
	GenericGetService[ServerGroup]
	GenericListService[ServerGroup]
	GenericUpdateService[ServerGroup, ServerGroupRequest]
	GenericDeleteService[ServerGroup]
}
