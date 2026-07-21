//go:build integration

package integration

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/cloudscale-ch/cloudscale-go-sdk/v9"
)

func TestIntegrationRouter_CR_D(t *testing.T) {
	t.Parallel()

	createRouterRequest := &cloudscale.RouterCreateRequest{
		Name:                 testRunPrefix,
		InternetGateway:      true,
		ZonalResourceRequest: cloudscale.ZonalResourceRequest{Zone: testZone},
	}

	expected, err := client.Routers.Create(t.Context(), createRouterRequest)
	if err != nil {
		t.Fatalf("Routers.Create returned error %s", err)
	}

	router, err := client.Routers.Get(t.Context(), expected.UUID)
	if err != nil {
		t.Fatalf("Routers.Get returned error %s", err)
	}

	if uuid := router.UUID; uuid != expected.UUID {
		t.Errorf("Router.UUID got=%s\nwant=%s", uuid, expected.UUID)
	}

	if h := time.Since(router.CreatedAt).Hours(); !(-1 < h && h < 1) {
		t.Errorf("router.CreatedAt outside of expected range. got=%v", router.CreatedAt)
	}

	if !router.InternetGateway {
		t.Errorf("router.InternetGateway got=%v\nwant=%v", router.InternetGateway, true)
	}

	if _, err := client.Routers.WaitFor(t.Context(), router.UUID, cloudscale.RouterIsActive); err != nil {
		t.Errorf("router not in active state: %v", err)
	}

	routers, err := client.Routers.List(t.Context())
	if err != nil {
		t.Fatalf("Routers.List returned error %s\n", err)
	}

	if numRouters := len(routers); numRouters < 1 {
		t.Errorf("Routers.List got=%d\nwant>=%d\n", numRouters, 1)
	}

	// Set up a network with a subnet so we can attach an interface to the router.
	createNetworkRequest := &cloudscale.NetworkCreateRequest{
		Name:                 testRunPrefix,
		AutoCreateIPV4Subnet: new(false),
		ZonalResourceRequest: cloudscale.ZonalResourceRequest{Zone: testZone},
	}
	network, err := client.Networks.Create(t.Context(), createNetworkRequest)
	if err != nil {
		t.Fatalf("Networks.Create returned error %s", err)
	}

	createSubnetRequest := &cloudscale.SubnetCreateRequest{
		Network: network.UUID,
		CIDR:    "192.168.99.0/24",
	}
	subnet, err := client.Subnets.Create(t.Context(), createSubnetRequest)
	if err != nil {
		t.Fatalf("Subnets.Create returned error %s", err)
	}

	createInterfaceRequest := cloudscale.CreateInterfaceRequest{
		Network: network.UUID,
		Addresses: []cloudscale.CreateAddressRequest{
			{
				Subnet:  subnet.UUID,
				Address: "192.168.99.10",
			},
		},
	}
	iface, err := client.Routers.CreateInterface(t.Context(), router.UUID, createInterfaceRequest)
	if err != nil {
		t.Fatalf("Routers.CreateInterface returned error %s", err)
	}

	if iface.UUID == "" {
		t.Error("Routers.CreateInterface returned interface without UUID")
	}
	if networkUUID := iface.Network.UUID; networkUUID != network.UUID {
		t.Errorf("interface.Network.UUID got=%s\nwant=%s", networkUUID, network.UUID)
	}
	if numAddresses := len(iface.Addresses); numAddresses != 1 {
		t.Fatalf("interface Addresses got=%d\nwant=%d", numAddresses, 1)
	}
	if subnetUUID := iface.Addresses[0].Subnet.UUID; subnetUUID != subnet.UUID {
		t.Errorf("interface.Addresses[0].Subnet.UUID got=%s\nwant=%s", subnetUUID, subnet.UUID)
	}
	if addr := iface.Addresses[0].Address; addr != "192.168.99.10" {
		t.Errorf("interface.Addresses[0].Address got=%s\nwant=%s", addr, "192.168.99.10")
	}

	// Clean up: delete interface
	if err := client.Routers.DeleteInterface(t.Context(), router.UUID, iface.UUID); err != nil {
		t.Errorf("Routers.DeleteInterface returned error: %s", err)
	}

	// Verify the interface is actually gone before deleting the router: the router's
	// Interfaces list must no longer contain it.
	err = waitForDeleted(t.Context(), func() (exists bool, err error) {
		r, err := client.Routers.Get(t.Context(), router.UUID)
		if err != nil {
			return true, err
		}
		for _, i := range r.Interfaces {
			if i.UUID == iface.UUID {
				t.Logf("interface %q still attached to router %q", iface.UUID, router.UUID)
				return true, nil
			}
		}
		return false, nil
	})
	if err != nil {
		t.Errorf("waiting for interface delete failed: %v", err)
	}

	// Clean up: delete router
	err = client.Routers.Delete(t.Context(), router.UUID)
	if err != nil {
		t.Fatalf("Routers.Delete returned error %s", err)
	}

	err = waitForDeleted(t.Context(), func() (exists bool, err error) {
		r, err := client.Routers.Get(t.Context(), router.UUID)
		if err != nil {
			if cerr, ok := errors.AsType[*cloudscale.ErrorResponse](err); ok && cerr.StatusCode == http.StatusNotFound {
				return false, nil
			}
			return true, err
		}
		t.Logf("router %q still exists with status %q", r.UUID, r.Status)
		return true, nil
	})
	if err != nil {
		t.Errorf("waiting for router delete failed: %v", err)
	}

	// Clean up: delete network
	err = client.Networks.Delete(t.Context(), network.UUID)
	if err != nil {
		t.Fatalf("Networks.Delete returned error %s", err)
	}
}
