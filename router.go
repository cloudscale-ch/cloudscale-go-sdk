package cloudscale

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const routerBasePath = "v1/routers"

type Router struct {
	ZonalResource
	TaggedResource
	HREF                     string            `json:"href"`
	UUID                     string            `json:"uuid"`
	Name                     string            `json:"name"`
	CreatedAt                time.Time         `json:"created_at"`
	Status                   string            `json:"status"`
	InternetGateway          bool              `json:"internet_gateway"`
	InternetGatewayAddresses []IPAddress       `json:"internet_gateway_addresses,omitempty"`
	Interfaces               []RouterInterface `json:"interfaces,omitempty"`
}

type IPAddress struct {
	Address    string     `json:"address"`
	Subnet     SubnetStub `json:"subnet"`
	Version    int        `json:"version"`
	ReversePTR *string    `json:"reverse_ptr"`
}
type RouterInterface struct {
	UUID       string      `json:"uuid"`
	Network    NetworkStub `json:"network"`
	Addresses  []IPAddress `json:"addresses"`
	Type       string      `json:"type"`
	MACAddress string      `json:"mac_address"`
}

type RouterCreateRequest struct {
	ZonalResourceRequest
	TaggedResourceRequest
	Name            string `json:"name"`
	InternetGateway bool   `json:"internet_gateway"`
}

// RouterUpdateRequest is not implemented yet because the API is not implemented yet
type RouterUpdateRequest struct{}

type RouterService interface {
	GenericCreateService[Router, RouterCreateRequest]
	GenericGetService[Router]
	GenericListService[Router]
	// GenericUpdateService[Router, RouterUpdateRequest]
	GenericDeleteService[Router]
	GenericWaitForService[Router]
	// CreateInterface creates a new interface attached to this router
	CreateInterface(ctx context.Context, routerUUID string, createReq CreateInterfaceRequest) (*RouterInterface, error)
	// DeleteInterface removes an interface attached to this router
	DeleteInterface(ctx context.Context, routerUUID, interfaceUUID string) error
}

type CreateInterfaceRequest struct {
	Network   string                 `json:"network"`
	Addresses []CreateAddressRequest `json:"addresses"`
}
type CreateAddressRequest struct {
	Subnet  string `json:"subnet"`
	Address string `json:"address"`
}

type RouterServiceOperations struct {
	GenericServiceOperations[Router, RouterCreateRequest, RouterUpdateRequest]
	client *Client
}

func (r RouterServiceOperations) CreateInterface(ctx context.Context, routerUUID string, createReq CreateInterfaceRequest) (*RouterInterface, error) {
	path := fmt.Sprintf("%s/%s/interfaces", routerBasePath, routerUUID)
	ctx = WithOperationPath(ctx, routerBasePath+"/:id/interfaces")
	req, err := r.client.NewRequest(ctx, http.MethodPost, path, createReq)
	if err != nil {
		return nil, err
	}
	res := &RouterInterface{}
	if err := r.client.Do(ctx, req, res); err != nil {
		return nil, err
	}
	return res, nil
}

func (r RouterServiceOperations) DeleteInterface(ctx context.Context, routerUUID, interfaceUUID string) error {
	path := fmt.Sprintf("%s/%s/interfaces/%s", routerBasePath, routerUUID, interfaceUUID)
	ctx = WithOperationPath(ctx, routerBasePath+"/:id/interfaces/:interface_id")

	req, err := r.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return r.client.Do(ctx, req, nil)
}

const (
	RouterActive = "active"
)

var RouterIsActive = func(router *Router) (bool, error) {
	if router.Status == RouterActive {
		return true, nil
	}
	return false, fmt.Errorf("waiting for status: %s, current status: %s", RouterActive, router.Status)
}
