package cloudscale

import (
	"time"
)

const networkBasePath = "v1/networks"

type Network struct {
	ZonalResource
	TaggedResource
	// Just use omitempty everywhere. This makes it easy to use restful. Errors
	// will be coming from the API if something is disabled.
	HREF      string       `json:"href,omitempty"`
	UUID      string       `json:"uuid,omitempty"`
	Name      string       `json:"name,omitempty"`
	MTU       int          `json:"mtu,omitempty"`
	Subnets   []SubnetStub `json:"subnets"`
	CreatedAt time.Time    `json:"created_at"`
}

type NetworkStub struct {
	HREF string `json:"href,omitempty"`
	Name string `json:"name,omitempty"`
	UUID string `json:"uuid,omitempty"`
}

type NetworkCreateRequest struct {
	ZonalResourceRequest
	TaggedResourceRequest
	Name                 string `json:"name,omitempty"`
	MTU                  int    `json:"mtu,omitempty"`
	AutoCreateIPV4Subnet *bool  `json:"auto_create_ipv4_subnet,omitempty"`
}

type NetworkUpdateRequest struct {
	ZonalResourceRequest
	TaggedResourceRequest
	Name string `json:"name,omitempty"`
	MTU  int    `json:"mtu,omitempty"`
}

type NetworkService interface {
	GenericCreateService[Network, NetworkCreateRequest]
	GenericGetService[Network]
	GenericListService[Network]
	GenericUpdateService[Network, NetworkUpdateRequest]
	GenericDeleteService[Network]
	GenericWaitForService[Network]
}

type NetworkServiceOperations struct {
	client *Client
}
