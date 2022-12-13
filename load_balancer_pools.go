package cloudscale

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const loadBalancerPoolBasePath = "v1/load-balancers/pools"

type LoadBalancerStub struct {
	HREF string `json:"href,omitempty"`
	UUID string `json:"uuid,omitempty"`
	Name string `json:"name,omitempty"`
}

type LoadBalancerPool struct {
	TaggedResource
	// Just use omitempty everywhere. This makes it easy to use restful. Errors
	// will be coming from the API if something is disabled.
	HREF         string           `json:"href,omitempty"`
	UUID         string           `json:"uuid,omitempty"`
	Name         string           `json:"name,omitempty"`
	CreatedAt    time.Time        `json:"created_at,omitempty"`
	LoadBalancer LoadBalancerStub `json:"load_balancer,omitempty"`
	Algorithm    string           `json:"algorithm,omitempty"`
	Protocol     string           `json:"protocol,omitempty"`
}

type LoadBalancerPoolRequest struct {
	TaggedResourceRequest
	Name         string `json:"name,omitempty"`
	LoadBalancer string `json:"load_balancer,omitempty"`
	Algorithm    string `json:"algorithm,omitempty"`
	Protocol     string `json:"protocol,omitempty"`
}

type LoadBalancerPoolService interface {
	Create(ctx context.Context, createRequest *LoadBalancerPoolRequest) (*LoadBalancerPool, error)
	Get(ctx context.Context, loadBalancerPoolID string) (*LoadBalancerPool, error)
	List(ctx context.Context, modifiers ...ListRequestModifier) ([]LoadBalancerPool, error)
	Update(ctx context.Context, loadBalancerPoolID string, updateRequest *LoadBalancerPoolRequest) error
	Delete(ctx context.Context, loadBalancerPoolID string) error
}

type LoadBalancerPoolServiceOperations struct {
	client *Client
}

func (s LoadBalancerPoolServiceOperations) Update(ctx context.Context, loadBalancerPoolID string, updateRequest *LoadBalancerPoolRequest) error {
	path := fmt.Sprintf("%s/%s", loadBalancerPoolBasePath, loadBalancerPoolID)

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

func (s LoadBalancerPoolServiceOperations) Create(ctx context.Context, createRequest *LoadBalancerPoolRequest) (*LoadBalancerPool, error) {
	path := loadBalancerPoolBasePath

	req, err := s.client.NewRequest(ctx, http.MethodPost, path, createRequest)
	if err != nil {
		return nil, err
	}

	loadBalancerPool := new(LoadBalancerPool)

	err = s.client.Do(ctx, req, loadBalancerPool)
	if err != nil {
		return nil, err
	}

	return loadBalancerPool, nil
}

func (s LoadBalancerPoolServiceOperations) Get(ctx context.Context, loadBalancerPoolID string) (*LoadBalancerPool, error) {
	path := fmt.Sprintf("%s/%s", loadBalancerPoolBasePath, loadBalancerPoolID)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	loadBalancerPool := new(LoadBalancerPool)
	err = s.client.Do(ctx, req, loadBalancerPool)
	if err != nil {
		return nil, err
	}

	return loadBalancerPool, nil
}

func (s LoadBalancerPoolServiceOperations) Delete(ctx context.Context, loadBalancerPoolID string) error {
	path := fmt.Sprintf("%s/%s", loadBalancerPoolBasePath, loadBalancerPoolID)

	req, err := s.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return s.client.Do(ctx, req, nil)
}

func (s LoadBalancerPoolServiceOperations) List(ctx context.Context, modifiers ...ListRequestModifier) ([]LoadBalancerPool, error) {
	path := loadBalancerPoolBasePath
	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	for _, modifier := range modifiers {
		modifier(req)
	}

	loadBalancerPools := []LoadBalancerPool{}
	err = s.client.Do(ctx, req, &loadBalancerPools)
	if err != nil {
		return nil, err
	}

	return loadBalancerPools, nil
}
