package cloudscale

import (
	"context"
	"net/http"
)

const flavorsBasePath = "v1/flavors"

// Flavor represents a server flavor, i.e., a combination of vCPUs and memory.
// Flavors are read-only and zonal. Categories include shared vCPU (flex-*),
// dedicated CPU (plus-*), and FlavorGPU (gpu*).
type Flavor struct {
	Slug      string     `json:"slug"`
	Name      string     `json:"name"`
	VCPUCount int        `json:"vcpu_count"`
	MemoryGB  int        `json:"memory_gb"`
	GPU       *FlavorGPU `json:"gpu"`
	Zones     []ZoneStub `json:"zones"`
}

type FlavorGPU struct {
	Name         string `json:"name"`
	Count        int    `json:"count"`
	VRAMPerGPUGB int    `json:"vram_per_gpu_gb"`
}

// FlavorService provides listing of available server flavors.
type FlavorService interface {
	List(ctx context.Context) ([]Flavor, error)
}

// FlavorServiceOperations implements FlavorService.
type FlavorServiceOperations struct {
	client *Client
}

// _ ensures that FlavorServiceOperations implements the FlavorService interface at compile time.
var _ FlavorService = &FlavorServiceOperations{}

// List returns all available flavors (GET /v1/flavors).
func (s FlavorServiceOperations) List(ctx context.Context) ([]Flavor, error) {
	req, err := s.client.NewRequest(ctx, http.MethodGet, flavorsBasePath, nil)
	if err != nil {
		return nil, err
	}
	var flavors []Flavor
	err = s.client.Do(ctx, req, &flavors)
	if err != nil {
		return nil, err
	}
	return flavors, nil
}
