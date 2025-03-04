package cloudscale

import (
	"strconv"
	"strings"
	"time"
)

const floatingIPsBasePath = "v1/floating-ips"

type FloatingIP struct {
	Region *Region `json:"region"` // not using RegionalResource here, as FloatingIP can be regional or global
	TaggedResource
	HREF           string            `json:"href"`
	Network        string            `json:"network"`
	IPVersion      int               `json:"ip_version"`
	NextHop        string            `json:"next_hop"`
	Server         *ServerStub       `json:"server"`
	LoadBalancer   *LoadBalancerStub `json:"load_balancer"`
	Type           string            `json:"type"`
	ReversePointer string            `json:"reverse_ptr,omitempty"`
	CreatedAt      time.Time         `json:"created_at"`
}

type FloatingIPCreateRequest struct {
	RegionalResourceRequest
	TaggedResourceRequest
	IPVersion      int    `json:"ip_version"`
	Server         string `json:"server,omitempty"`
	LoadBalancer   string `json:"load_balancer,omitempty"`
	Type           string `json:"type,omitempty"`
	PrefixLength   int    `json:"prefix_length,omitempty"`
	ReversePointer string `json:"reverse_ptr,omitempty"`
}

func (f FloatingIP) IP() string {
	return strings.Split(f.Network, "/")[0]
}

func (f FloatingIP) PrefixLength() int {
	result, _ := strconv.Atoi(strings.Split(f.Network, "/")[1])
	return result
}

type FloatingIPUpdateRequest struct {
	TaggedResourceRequest
	Server         string `json:"server,omitempty"`
	LoadBalancer   string `json:"load_balancer,omitempty"`
	ReversePointer string `json:"reverse_ptr,omitempty"`
}

type FloatingIPsService interface {
	GenericCreateService[FloatingIP, FloatingIPCreateRequest, FloatingIPUpdateRequest]
	GenericGetService[FloatingIP, FloatingIPCreateRequest, FloatingIPUpdateRequest]
	GenericListService[FloatingIP, FloatingIPCreateRequest, FloatingIPUpdateRequest]
	GenericUpdateService[FloatingIP, FloatingIPCreateRequest, FloatingIPUpdateRequest]
	GenericDeleteService[FloatingIP, FloatingIPCreateRequest, FloatingIPUpdateRequest]
}
