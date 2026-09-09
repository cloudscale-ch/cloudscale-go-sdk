//go:build integration

package integration

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/cloudscale-ch/cloudscale-go-sdk/v10"
)

func checkIfRouterIsDeleted(ctx context.Context, router *cloudscale.Router) error {
	r, err := client.Routers.Get(ctx, router.UUID)
	if err != nil {
		if cerr, ok := errors.AsType[*cloudscale.ErrorResponse](err); ok && cerr.StatusCode == http.StatusNotFound {
			// The router cannot be found, which means it was successfully deleted
			return nil
		}
		// A different API error code was returned
		return err
	}

	return fmt.Errorf("router %q still exists with status %q", r.UUID, r.Status)
}

func checkIfRouterInterfaceIsDeleted(ctx context.Context, router *cloudscale.Router, routerInterface *cloudscale.RouterInterface) error {
	r, err := client.Routers.Get(ctx, router.UUID)
	if err != nil {
		return err
	}
	for _, i := range r.Interfaces {
		if i.UUID == routerInterface.UUID {
			return fmt.Errorf("interface %q still attached to router %q", routerInterface.UUID, router.UUID)
		}
	}
	// The router interface is not in the list, which means it was successfully detached
	return nil
}

// newTestRouter creates a router for use in a single test and registers a
// cleanup that deletes it (and waits for the deletion to complete) once the
// test finishes.
func newTestRouter(t *testing.T) *cloudscale.Router {
	t.Helper()

	createRouterRequest := &cloudscale.RouterCreateRequest{
		Name:                 "offline-router",
		InternetGateway:      false,
		ZonalResourceRequest: cloudscale.ZonalResourceRequest{Zone: testZone},
	}

	router, err := client.Routers.Create(t.Context(), createRouterRequest)
	if err != nil {
		t.Fatalf("Routers.Create returned error %s", err)
	}

	t.Cleanup(func() {
		router, err := client.Routers.Get(context.Background(), router.UUID)
		if err != nil {
			t.Errorf("Routers.Get returned error: %v", err)
			return
		}

		for _, routerInterface := range router.Interfaces {
			err = client.Routers.DeleteInterface(context.Background(), router.UUID, routerInterface.UUID)
			if err != nil {
				t.Errorf("Routers.DeleteInterface returned error: %v", err)
				return
			}

			err = checkIfRouterInterfaceIsDeleted(context.Background(), router, &routerInterface)
			if err != nil {
				t.Errorf("check if router interface is deleted failed: %v", err)
				return
			}
		}

		err = client.Routers.Delete(context.Background(), router.UUID)
		if err != nil {
			t.Errorf("Routers.Delete returned error: %s", err)
			return
		}

		err = checkIfRouterIsDeleted(context.Background(), router)
		if err != nil {
			t.Errorf("check if router is deleted failed: %v", err)
			return
		}
	})

	return router
}

// newTestNetworkWithSubnet creates a network and a subnet within it for use
// in a single test and registers a cleanup that deletes the network once the
// test finishes.
func newTestNetworkWithSubnet(t *testing.T) (*cloudscale.Network, *cloudscale.Subnet) {
	t.Helper()

	createNetworkRequest := &cloudscale.NetworkCreateRequest{
		Name:                 testRunPrefix,
		AutoCreateIPV4Subnet: new(false),
		ZonalResourceRequest: cloudscale.ZonalResourceRequest{Zone: testZone},
	}
	network, err := client.Networks.Create(t.Context(), createNetworkRequest)
	if err != nil {
		t.Fatalf("Networks.Create returned error: %s", err)
	}

	t.Cleanup(func() {
		if err := client.Networks.Delete(context.Background(), network.UUID); err != nil {
			t.Errorf("Networks.Delete returned error: %s", err)
		}
	})

	createSubnetRequest := &cloudscale.SubnetCreateRequest{
		Network: network.UUID,
		CIDR:    "192.168.99.0/24",
	}
	subnet, err := client.Subnets.Create(t.Context(), createSubnetRequest)
	if err != nil {
		t.Fatalf("Subnets.Create returned error: %s", err)
	}

	return network, subnet
}

func TestIntegrationRouter_Create(t *testing.T) {
	t.Parallel()

	createRouterRequest := &cloudscale.RouterCreateRequest{
		Name:                 "offline-router",
		InternetGateway:      false,
		ZonalResourceRequest: cloudscale.ZonalResourceRequest{Zone: testZone},
	}

	router, err := client.Routers.Create(t.Context(), createRouterRequest)
	if err != nil {
		t.Fatalf("Routers.Create returned error: %s", err)
	}
	t.Cleanup(func() {
		err := client.Routers.Delete(context.Background(), router.UUID)
		if err != nil {
			t.Errorf("Routers.Delete returned error: %s", err)
			return
		}

		err = checkIfRouterIsDeleted(context.Background(), router)
		if err != nil {
			t.Errorf("check if router is deleted failed: %v", err)
		}
	})

	if router.Name != createRouterRequest.Name {
		t.Errorf("Router.Name got=%s\nwant=%s", router.Name, createRouterRequest.Name)
	}

	if router.InternetGateway != createRouterRequest.InternetGateway {
		t.Errorf("Router.InternetGateway got=%t\nwant=%t", router.InternetGateway, createRouterRequest.InternetGateway)
	}

	// The number of InternetGatewayAddresses should be 0, since InternetGateway is initially set to false
	if numInternetGatewayAddresses := len(router.InternetGatewayAddresses); numInternetGatewayAddresses != 0 {
		t.Errorf("Number of InternetGatewayAddresses got=%d\nwant=%d", numInternetGatewayAddresses, 0)
	}
}

func TestIntegrationRouter_Update_Name(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t)

	updateNameRouterRequest := &cloudscale.RouterUpdateRequest{
		Name: "router.example.com",
	}

	err := client.Routers.Update(t.Context(), router.UUID, updateNameRouterRequest)
	if err != nil {
		t.Fatalf("Routers.Update returned error: %s", err)
	}

	updated, err := client.Routers.Get(t.Context(), router.UUID)
	if err != nil {
		t.Fatalf("Routers.Get returned error: %s", err)
	}

	if updateNameRouterRequest.Name != updated.Name {
		t.Errorf("router.Name got=%s\nwant=%s", updated.Name, updateNameRouterRequest.Name)
	}
}

func TestIntegrationRouter_Update_InternetGateway(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t)

	// Give the router a valid FQDN name first, so we can assert on the
	// reverse pointer of the resulting InternetGatewayAddresses.
	updateNameRouterRequest := &cloudscale.RouterUpdateRequest{
		Name: "router.example.com",
	}
	if err := client.Routers.Update(t.Context(), router.UUID, updateNameRouterRequest); err != nil {
		t.Fatalf("Routers.Update returned error: %s", err)
	}

	updateInternetGatewayRouterRequest := &cloudscale.RouterUpdateRequest{
		InternetGateway: true,
	}

	err := client.Routers.Update(t.Context(), router.UUID, updateInternetGatewayRouterRequest)
	if err != nil {
		t.Fatalf("Routers.Update returned error: %s", err)
	}

	updated, err := client.Routers.Get(t.Context(), router.UUID)
	if err != nil {
		t.Fatalf("Routers.Get returned error: %s", err)
	}

	if updateInternetGatewayRouterRequest.InternetGateway != updated.InternetGateway {
		t.Errorf("router.InternetGateway got=%t\nwant=%t", updated.InternetGateway, updateInternetGatewayRouterRequest.InternetGateway)
	}

	// The number of InternetGatewayAddresses should be 1, since InternetGateway was updated to true.
	// This assertion must be updated after the API supports IPv6 for private networks.
	if numInternetGatewayAddresses := len(updated.InternetGatewayAddresses); numInternetGatewayAddresses != 1 {
		t.Errorf("Number of InternetGatewayAddresses got=%d\nwant=%d", numInternetGatewayAddresses, 1)
	}

	// If InternetGateway is true and the Name is a valid FQDN, all InternetGatewayAddresses should have this FQDN as a reverse pointer.
	for index, address := range updated.InternetGatewayAddresses {
		if reversePTR := *address.ReversePTR; reversePTR != updateNameRouterRequest.Name {
			t.Errorf("Router.InternetGatewayAddresses[%d].ReversePTR got=%s\nwant=%s", index, reversePTR, updateNameRouterRequest.Name)
		}
	}
}

func TestIntegrationRouter_Get(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t)

	got, err := client.Routers.Get(t.Context(), router.UUID)
	if err != nil {
		t.Fatalf("Routers.Get returned error: %s", err)
	}

	if uuid := got.UUID; uuid != router.UUID {
		t.Errorf("Router.UUID got=%s\nwant=%s", uuid, router.UUID)
	}

	if h := time.Since(got.CreatedAt).Hours(); !(-1 < h && h < 1) {
		t.Errorf("router.CreatedAt outside of expected range. got=%v", got.CreatedAt)
	}
}

func TestIntegrationRouter_List(t *testing.T) {
	t.Parallel()

	_ = newTestRouter(t)

	routers, err := client.Routers.List(t.Context())
	if err != nil {
		t.Fatalf("Routers.List returned error: %s\n", err)
	}

	if numRouters := len(routers); numRouters < 1 {
		t.Errorf("Routers.List got=%d\nwant>=%d\n", numRouters, 1)
	}
}

func TestIntegrationRouter_WaitFor(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t)

	if _, err := client.Routers.WaitFor(t.Context(), router.UUID, cloudscale.RouterIsActive); err != nil {
		t.Errorf("router not in active state: %v", err)
	}
}

func TestIntegrationRouter_AttachInterface(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t)
	network, subnet := newTestNetworkWithSubnet(t)

	createInterfaceRequest := cloudscale.CreateInterfaceRequest{
		Network: network.UUID,
		Addresses: []cloudscale.CreateAddressRequest{
			{
				Subnet:  subnet.UUID,
				Address: "192.168.99.10",
			},
		},
	}
	routerInterface, err := client.Routers.CreateInterface(t.Context(), router.UUID, createInterfaceRequest)
	if err != nil {
		t.Fatalf("Routers.CreateInterface returned error: %s", err)
	}
	t.Cleanup(func() {
		if err := client.Routers.DeleteInterface(context.Background(), router.UUID, routerInterface.UUID); err != nil {
			t.Errorf("Routers.DeleteInterface returned error: %s", err)
		}
	})

	if routerInterface.UUID == "" {
		t.Error("Routers.CreateInterface returned interface without UUID")
	}
	if networkUUID := routerInterface.Network.UUID; networkUUID != network.UUID {
		t.Errorf("interface.Network.UUID got=%s\nwant=%s", networkUUID, network.UUID)
	}
	if numAddresses := len(routerInterface.Addresses); numAddresses != 1 {
		t.Fatalf("interface Addresses got=%d\nwant=%d", numAddresses, 1)
	}
	if subnetUUID := routerInterface.Addresses[0].Subnet.UUID; subnetUUID != subnet.UUID {
		t.Errorf("interface.Addresses[0].Subnet.UUID got=%s\nwant=%s", subnetUUID, subnet.UUID)
	}
	if addr := routerInterface.Addresses[0].Address; addr != "192.168.99.10" {
		t.Errorf("interface.Addresses[0].Address got=%s\nwant=%s", addr, "192.168.99.10")
	}
}

func TestIntegrationRouter_DetachInterface(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t)
	network, subnet := newTestNetworkWithSubnet(t)

	createInterfaceRequest := cloudscale.CreateInterfaceRequest{
		Network: network.UUID,
		Addresses: []cloudscale.CreateAddressRequest{
			{
				Subnet:  subnet.UUID,
				Address: "192.168.99.10",
			},
		},
	}
	routerInterface, err := client.Routers.CreateInterface(t.Context(), router.UUID, createInterfaceRequest)
	if err != nil {
		t.Fatalf("Routers.CreateInterface returned error: %s", err)
	}

	if err := client.Routers.DeleteInterface(t.Context(), router.UUID, routerInterface.UUID); err != nil {
		t.Errorf("Routers.DeleteInterface returned error: %s", err)
	}

	// Verify the interface is actually gone: the router's Interfaces list
	// must no longer contain it.
	err = checkIfRouterInterfaceIsDeleted(t.Context(), router, routerInterface)
	if err != nil {
		t.Errorf("check if router interface is deleted failed: %v", err)
	}
}

func TestIntegrationRouter_Delete(t *testing.T) {
	t.Parallel()

	// Cannot use newTestRouter here, since it registers a cleanup function which also tries to delete the router
	createRouterRequest := &cloudscale.RouterCreateRequest{
		Name:                 "offline-router",
		InternetGateway:      false,
		ZonalResourceRequest: cloudscale.ZonalResourceRequest{Zone: testZone},
	}

	router, err := client.Routers.Create(t.Context(), createRouterRequest)
	if err != nil {
		t.Fatalf("Routers.Create returned error: %s", err)
	}

	err = client.Routers.Delete(t.Context(), router.UUID)
	if err != nil {
		t.Fatalf("Routers.Delete returned error: %s", err)
	}

	err = checkIfRouterIsDeleted(t.Context(), router)
	if err != nil {
		t.Errorf("check if router is deleted failed: %v", err)
	}
}
