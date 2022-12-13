package cloudscale

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const loadBalancerBasePath = "v1/load-balancers"

type LoadBalancer struct {
	ZonalResource
	TaggedResource
	// Just use omitempty everywhere. This makes it easy to use restful. Errors
	// will be coming from the API if something is disabled.
	HREF         string       `json:"href"`
	UUID         string       `json:"uuid"`
	Name         string       `json:"name"`
	Status       string       `json:"status"`
	VIPAddresses []VIPAddress `json:"vip_addresses"`
	CreatedAt    time.Time    `json:"created_at"`
}

type VIPAddress struct {
	Version int        `json:"version"`
	Address string     `json:"address"`
	Subnet  SubnetStub `json:"subnet"`
}

type LoadBalancerRequest struct {
	ZonalResourceRequest
	TaggedResourceRequest
	Name   string `json:"name,omitempty"`
	Flavor string `json:"flavor,omitempty"`
}

type LoadBalancerService interface {
	Create(ctx context.Context, createRequest *LoadBalancerRequest) (*LoadBalancer, error)
	Get(ctx context.Context, loadBalancerID string) (*LoadBalancer, error)
	List(ctx context.Context, modifiers ...ListRequestModifier) ([]LoadBalancer, error)
	Update(ctx context.Context, loadBalancerID string, updateRequest *LoadBalancerRequest) error
	Delete(ctx context.Context, loadBalancerID string) error
}

type LoadBalancerServiceOperations struct {
	client *Client
}

func (s LoadBalancerServiceOperations) Update(ctx context.Context, loadBalancerID string, updateRequest *LoadBalancerRequest) error {
	path := fmt.Sprintf("%s/%s", loadBalancerBasePath, loadBalancerID)

	req, err := s.client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return err
	}

	err = s.client.Do(ctx, req, nil)
	if err != nil {
		return err
	}
	return nil
}

func (s LoadBalancerServiceOperations) Create(ctx context.Context, createRequest *LoadBalancerRequest) (*LoadBalancer, error) {
	path := loadBalancerBasePath

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, err
	}

	loadBalancer := new(LoadBalancer)

	err = s.client.Do(ctx, req, loadBalancer)
	if err != nil {
		return nil, err
	}

	return loadBalancer, nil
}

func (s LoadBalancerServiceOperations) Get(ctx context.Context, loadBalancerID string) (*LoadBalancer, error) {
	path := fmt.Sprintf("%s/%s", loadBalancerBasePath, loadBalancerID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	loadBalancer := new(LoadBalancer)
	err = s.client.Do(ctx, req, loadBalancer)
	if err != nil {
		return nil, err
	}

	return loadBalancer, nil
}

func (s LoadBalancerServiceOperations) Delete(ctx context.Context, loadBalancerID string) error {
	path := fmt.Sprintf("%s/%s", loadBalancerBasePath, loadBalancerID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return s.client.Do(ctx, req, nil)
}

func (s LoadBalancerServiceOperations) List(ctx context.Context, modifiers ...ListRequestModifier) ([]LoadBalancer, error) {
	path := loadBalancerBasePath
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	for _, modifier := range modifiers {
		modifier(req)
	}

	loadBalancers := []LoadBalancer{}
	err = s.client.Do(ctx, req, &loadBalancers)
	if err != nil {
		return nil, err
	}

	return loadBalancers, nil
}
