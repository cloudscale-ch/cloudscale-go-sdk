package cloudscale

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"
)

const serverBasePath = "v1/servers"

const ServerRunning = "running"
const ServerStopped = "stopped"
const ServerRebooted = "rebooted"

type Server struct {
	ZonalResource
	TaggedResource
	HREF            string            `json:"href"`
	UUID            string            `json:"uuid"`
	Name            string            `json:"name"`
	Status          string            `json:"status"`
	Flavor          Flavor            `json:"flavor"`
	Image           Image             `json:"image"`
	Volumes         []VolumeStub      `json:"volumes"`
	Interfaces      []Interface       `json:"interfaces"`
	SSHFingerprints []string          `json:"ssh_fingerprints"`
	SSHHostKeys     []string          `json:"ssh_host_keys"`
	AntiAfinityWith []ServerStub      `json:"anti_affinity_with"`
	ServerGroups    []ServerGroupStub `json:"server_groups"`
	CreatedAt       time.Time         `json:"created_at"`
}

type ServerStub struct {
	HREF string `json:"href"`
	UUID string `json:"uuid"`
}

type ServerGroupStub struct {
	HREF string `json:"href"`
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

type Flavor struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	VCPUCount int    `json:"vcpu_count"`
	MemoryGB  int    `json:"memory_gb"`
}

type Image struct {
	Slug            string `json:"slug"`
	Name            string `json:"name"`
	OperatingSystem string `json:"operating_system"`
	DefaultUsername string `json:"default_username"`
}

type VolumeStub struct {
	Type       string `json:"type"`
	DevicePath string `json:"device_path"`
	SizeGB     int    `json:"size_gb"`
	UUID       string `json:"uuid"`
}

type Interface struct {
	Type      string      `json:"type,omitempty"`
	Network   NetworkStub `json:"network,omitempty"`
	Addresses []Address   `json:"addresses,omitempty"`
}

type Address struct {
	Version      int        `json:"version"`
	Address      string     `json:"address"`
	PrefixLength int        `json:"prefix_length"`
	Gateway      string     `json:"gateway"`
	ReversePtr   string     `json:"reverse_ptr"`
	Subnet       SubnetStub `json:"subnet"`
}

type ServerRequest struct {
	ZonalResourceRequest
	TaggedResourceRequest
	Name              string                 `json:"name"`
	Flavor            string                 `json:"flavor"`
	Image             string                 `json:"image"`
	Zone              string                 `json:"zone,omitempty"`
	VolumeSizeGB      int                    `json:"volume_size_gb,omitempty"`
	Volumes           *[]ServerVolumeRequest `json:"volumes,omitempty"`
	Interfaces        *[]InterfaceRequest    `json:"interfaces,omitempty"`
	BulkVolumeSizeGB  int                    `json:"bulk_volume_size_gb,omitempty"`
	SSHKeys           []string               `json:"ssh_keys"`
	Password          string                 `json:"password,omitempty"`
	UsePublicNetwork  *bool                  `json:"use_public_network,omitempty"`
	UsePrivateNetwork *bool                  `json:"use_private_network,omitempty"`
	UseIPV6           *bool                  `json:"use_ipv6,omitempty"`
	AntiAffinityWith  string                 `json:"anti_affinity_with,omitempty"`
	ServerGroups      []string               `json:"server_groups,omitempty"`
	UserData          string                 `json:"user_data,omitempty"`
}

type ServerUpdateRequest struct {
	TaggedResourceRequest
	Name       string              `json:"name,omitempty"`
	Status     string              `json:"status,omitempty"`
	Flavor     string              `json:"flavor,omitempty"`
	Interfaces *[]InterfaceRequest `json:"interfaces,omitempty"`
}

type ServerVolumeRequest struct {
	SizeGB int    `json:"size_gb,omitempty"`
	Type   string `json:"type,omitempty"`
}

type InterfaceRequest struct {
	Network   string            `json:"network,omitempty"`
	Addresses *[]AddressRequest `json:"addresses,omitempty"`
}

type AddressRequest struct {
	Subnet  string `json:"subnet,omitempty"`
	Address string `json:"address,omitempty"`
}

type ServerService interface {
	GenericCreateService[Server, ServerRequest]
	GenericGetService[Server]
	GenericListService[Server]
	GenericUpdateService[Server, ServerUpdateRequest]
	GenericDeleteService[Server]
	Reboot(ctx context.Context, serverID string) error
	Start(ctx context.Context, serverID string) error
	Stop(ctx context.Context, serverID string) error
}

type ServerServiceOperations struct {
	GenericServiceOperations[Server, ServerRequest, ServerUpdateRequest]
	client *Client
}

func (s ServerServiceOperations) Reboot(ctx context.Context, serverID string) error {
	path := fmt.Sprintf("%s/%s/reboot", serverBasePath, serverID)
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return err
	}
	return s.client.Do(ctx, req, nil)
}

func (s ServerServiceOperations) Start(ctx context.Context, serverID string) error {
	path := fmt.Sprintf("%s/%s/start", serverBasePath, serverID)
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return err
	}
	return s.client.Do(ctx, req, nil)
}

func (s ServerServiceOperations) Stop(ctx context.Context, serverID string) error {
	path := fmt.Sprintf("%s/%s/stop", serverBasePath, serverID)
	req, err := s.client.NewRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return err
	}
	return s.client.Do(ctx, req, nil)
}

func (s ServerServiceOperations) Update(ctx context.Context, id string, req *ServerUpdateRequest) error {
	if req.Status != "" {
		err := error(nil)
		switch req.Status {
		case ServerRunning:
			err = s.Start(ctx, id)
		case ServerStopped:
			err = s.Stop(ctx, id)
		case ServerRebooted:
			err = s.Reboot(ctx, id)
		default:
			return fmt.Errorf("Status Not Supported %s", req.Status)
		}
		if err != nil {
			return err
		}
		// Get rid of status
		req.Status = ""
	}

	emptyRequest := ServerUpdateRequest{}
	if reflect.DeepEqual(emptyRequest, *req) {
		return nil
	}

	// Call the generic Update method directly using the embedded field name
	return s.GenericServiceOperations.Update(ctx, id, req)
}
