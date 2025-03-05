package cloudscale

import (
	"encoding/json"
	"reflect"
)

const subnetBasePath = "v1/subnets"

var UseCloudscaleDefaults = []string{"CLOUDSCALE_DEFAULTS"}

type Subnet struct {
	TaggedResource
	// Just use omitempty everywhere. This makes it easy to use restful. Errors
	// will be coming from the API if something is disabled.
	HREF           string      `json:"href,omitempty"`
	UUID           string      `json:"uuid,omitempty"`
	CIDR           string      `json:"cidr,omitempty"`
	Network        NetworkStub `json:"network,omitempty"`
	GatewayAddress string      `json:"gateway_address,omitempty"`
	DNSServers     []string    `json:"dns_servers,omitempty"`
}

type SubnetStub struct {
	HREF string `json:"href,omitempty"`
	CIDR string `json:"cidr,omitempty"`
	UUID string `json:"uuid,omitempty"`
}

type SubnetCreateRequest struct {
	TaggedResourceRequest
	CIDR           string    `json:"cidr,omitempty"`
	Network        string    `json:"network,omitempty"`
	GatewayAddress string    `json:"gateway_address,omitempty"`
	DNSServers     *[]string `json:"dns_servers,omitempty"`
}

type SubnetUpdateRequest struct {
	TaggedResourceRequest
	GatewayAddress string    `json:"gateway_address,omitempty"`
	DNSServers     *[]string `json:"dns_servers"`
}

func (request SubnetUpdateRequest) MarshalJSON() ([]byte, error) {
	type Alias SubnetUpdateRequest // Create an alias to avoid recursion

	if request.DNSServers == nil {
		return json.Marshal(&struct {
			Alias
			DNSServers []string `json:"dns_servers,omitempty"`
		}{
			Alias: (Alias)(request),
		})
	}

	if reflect.DeepEqual(*request.DNSServers, UseCloudscaleDefaults) {
		return json.Marshal(&struct {
			Alias
			DNSServers []string `json:"dns_servers"` // important: no omitempty
		}{
			Alias:      (Alias)(request),
			DNSServers: nil,
		})
	}

	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(request),
	})
}

func (request SubnetCreateRequest) MarshalJSON() ([]byte, error) {
	type Alias SubnetCreateRequest // Create an alias to avoid recursion

	if request.DNSServers == nil {
		return json.Marshal(&struct {
			Alias
			DNSServers []string `json:"dns_servers,omitempty"`
		}{
			Alias: (Alias)(request),
		})
	}

	if reflect.DeepEqual(*request.DNSServers, UseCloudscaleDefaults) {
		return json.Marshal(&struct {
			Alias
			DNSServers []string `json:"dns_servers"` // important: no omitempty
		}{
			Alias:      (Alias)(request),
			DNSServers: nil,
		})
	}

	return json.Marshal(&struct {
		Alias
	}{
		Alias: (Alias)(request),
	})
}

type SubnetService interface {
	GenericCreateService[Subnet, SubnetCreateRequest]
	GenericGetService[Subnet]
	GenericListService[Subnet]
	GenericUpdateService[Subnet, SubnetUpdateRequest]
	GenericDeleteService[Subnet]
}

type SubnetServiceOperations struct {
	client *Client
}
