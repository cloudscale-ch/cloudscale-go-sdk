package cloudscale

import (
	"context"
	"fmt"
	"net/http"
)

type GenericService[TResource any, TRequest any] interface {
	Create(ctx context.Context, createRequest *TRequest) (*TResource, error)
	Get(ctx context.Context, resourceID string) (*TResource, error)
	List(ctx context.Context, modifiers ...ListRequestModifier) ([]TResource, error)
	Update(ctx context.Context, resourceID string, updateRequest *TRequest) error
	Delete(ctx context.Context, resourceID string) error
}

type GenericServiceOperations[TResource any, TRequest any] struct {
	client *Client
	path   string
}

func (g GenericServiceOperations[TResource, TRequest]) Create(ctx context.Context, createRequest *TRequest) (*TResource, error) {
	req, err := g.client.NewRequest(ctx, http.MethodPost, g.path, createRequest)
	if err != nil {
		return nil, err
	}

	resource := new(TResource)

	err = g.client.Do(ctx, req, resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (g GenericServiceOperations[TResource, TRequest]) Get(ctx context.Context, resourceID string) (*TResource, error) {
	path := fmt.Sprintf("%s/%s", g.path, resourceID)

	req, err := g.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	resource := new(TResource)
	err = g.client.Do(ctx, req, resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (g GenericServiceOperations[TResource, TRequest]) List(ctx context.Context, modifiers ...ListRequestModifier) ([]TResource, error) {
	req, err := g.client.NewRequest(ctx, http.MethodGet, g.path, nil)
	if err != nil {
		return nil, err
	}

	for _, modifier := range modifiers {
		modifier(req)
	}

	resources := []TResource{}
	err = g.client.Do(ctx, req, &resources)
	if err != nil {
		return nil, err
	}

	return resources, nil
}

func (g GenericServiceOperations[TResource, TRequest]) Update(ctx context.Context, resourceID string, updateRequest *TRequest) error {
	path := fmt.Sprintf("%s/%s", g.path, resourceID)

	req, err := g.client.NewRequest(ctx, http.MethodPatch, path, updateRequest)
	if err != nil {
		return err
	}

	err = g.client.Do(ctx, req, nil)
	if err != nil {
		return err
	}
	return nil
}

func (g GenericServiceOperations[TResource, TRequest]) Delete(ctx context.Context, resourceID string) error {
	path := fmt.Sprintf("%s/%s", g.path, resourceID)

	req, err := g.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return g.client.Do(ctx, req, nil)
}
