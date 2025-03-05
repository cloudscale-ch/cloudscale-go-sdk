package cloudscale

import (
	"context"
	"fmt"
	"github.com/cenkalti/backoff/v5"
	"net/http"
	"time"
)

type GenericCreateService[TResource any, TCreateRequest any] interface {
	Create(ctx context.Context, createRequest *TCreateRequest) (*TResource, error)
}

type GenericGetService[TResource any] interface {
	Get(ctx context.Context, resourceID string) (*TResource, error)
}

type GenericListService[TResource any] interface {
	List(ctx context.Context, modifiers ...ListRequestModifier) ([]TResource, error)
}

type GenericUpdateService[TResource any, TUpdateRequest any] interface {
	Update(ctx context.Context, resourceID string, updateRequest *TUpdateRequest) error
}

type GenericDeleteService[TResource any] interface {
	Delete(ctx context.Context, resourceID string) error
}

type GenericWaitForService[TResource any] interface {
	WaitFor(ctx context.Context, resourceID string, condition func(resource *TResource) (bool, error), opts ...backoff.RetryOption) (*TResource, error)
}

type GenericServiceOperations[TResource any, TCreateRequest any, TUpdateRequest any] struct {
	client *Client
	path   string
}

func (g GenericServiceOperations[TResource, TCreateRequest, TUpdateRequest]) Create(ctx context.Context, createRequest *TCreateRequest) (*TResource, error) {
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

func (g GenericServiceOperations[TResource, TCreateRequest, TUpdateRequest]) Get(ctx context.Context, resourceID string) (*TResource, error) {
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

func (g GenericServiceOperations[TResource, TCreateRequest, TUpdateRequest]) List(ctx context.Context, modifiers ...ListRequestModifier) ([]TResource, error) {
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

func (g GenericServiceOperations[TResource, TCreateRequest, TUpdateRequest]) Update(ctx context.Context, resourceID string, updateRequest *TUpdateRequest) error {
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

func (g GenericServiceOperations[TResource, TCreateRequest, TUpdateRequest]) Delete(ctx context.Context, resourceID string) error {
	path := fmt.Sprintf("%s/%s", g.path, resourceID)

	req, err := g.client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	return g.client.Do(ctx, req, nil)
}

func (g GenericServiceOperations[TResource, TCreateRequest, TUpdateRequest]) WaitFor(
	ctx context.Context,
	resourceID string,
	condition func(resource *TResource) (bool, error),
	opts ...backoff.RetryOption,
) (*TResource, error) {
	// Prepend the default backoff option.
	// If a user passes their own WithBackOff option, it will override this default.
	options := append([]backoff.RetryOption{
		backoff.WithBackOff(backoff.NewConstantBackOff(2 * time.Second)),
		backoff.WithMaxElapsedTime(2 * time.Minute),
	}, opts...)

	return backoff.Retry(ctx, func() (*TResource, error) {
		resource, err := g.Get(ctx, resourceID)
		if err != nil {
			return nil, err
		}

		ok, condErr := condition(resource)
		if ok {
			return resource, nil // Exit when the condition is met.
		}

		// If the condition provided an error, return it as our retry error message.
		if condErr != nil {
			return nil, condErr // Continue retrying
		}
		return nil, fmt.Errorf("condition not met yet") // Continue retrying
	}, options...)
}
