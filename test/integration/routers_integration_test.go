//go:build integration

package integration

import (
	"fmt"
	"testing"
	"time"

	"github.com/cloudscale-ch/cloudscale-go-sdk/v10"
)

func testCreateRouter(t *testing.T) cloudscale.Router {
	t.Helper()

	createRouterRequest := &cloudscale.RouterCreateRequest{
		Name:                 fmt.Sprintf("%s-%s", testRunPrefix, "offline-router"),
		InternetGateway:      false,
		ZonalResourceRequest: cloudscale.ZonalResourceRequest{Zone: testZone},
	}

	initialRouter, err := client.Routers.Create(t.Context(), createRouterRequest)
	if err != nil {
		t.Fatalf("Routers.Create returned error %s", err)
	}

	if initialRouter.Name != createRouterRequest.Name {
		t.Errorf("Router.Name got=%s\nwant=%s", initialRouter.Name, createRouterRequest.Name)
	}

	if initialRouter.InternetGateway != createRouterRequest.InternetGateway {
		t.Errorf("Router.InternetGateway got=%t\nwant=%t", initialRouter.InternetGateway, createRouterRequest.InternetGateway)
	}

	// The number of InternetGatewayAddresses should be 0, since InternetGateway is initially set to false
	if numInternetGatewayAddresses := len(initialRouter.InternetGatewayAddresses); numInternetGatewayAddresses != 0 {
		t.Errorf("Number of InternetGatewayAddresses got=%d\nwant=%d", numInternetGatewayAddresses, 0)
	}

	if _, err := client.Routers.WaitFor(t.Context(), initialRouter.UUID, cloudscale.RouterIsActive); err != nil {
		t.Errorf("router not in active state: %s", err)
	}

	if h := time.Since(initialRouter.CreatedAt).Hours(); !(-1 < h && h < 1) {
		t.Errorf("router.CreatedAt outside of expected range. got=%s", initialRouter.CreatedAt)
	}

	return *initialRouter
}

func testUpdateRouter(t *testing.T, initialRouter cloudscale.Router) cloudscale.Router {
	t.Helper()

	updateNameRouterRequest := &cloudscale.RouterUpdateRequest{
		Name: fmt.Sprintf("%s.%s", testRunPrefix, "router.example.com"),
	}

	err := client.Routers.Update(t.Context(), initialRouter.UUID, updateNameRouterRequest)
	if err != nil {
		t.Fatalf("Routers.Update returned error %s", err)
	}

	updateInternetGatewayRouterRequest := &cloudscale.RouterUpdateRequest{
		InternetGateway: true,
	}

	err = client.Routers.Update(t.Context(), initialRouter.UUID, updateInternetGatewayRouterRequest)
	if err != nil {
		t.Fatalf("Routers.Update returned error %s", err)
	}

	updatedRouter, err := client.Routers.Get(t.Context(), initialRouter.UUID)
	if err != nil {
		t.Fatalf("Routers.Get returned error %s", err)
	}

	if updatedRouter.Name != updateNameRouterRequest.Name {
		t.Errorf("router.Name got=%s\nwant=%s", updatedRouter.Name, updateNameRouterRequest.Name)
	}
	if updatedRouter.InternetGateway != updateInternetGatewayRouterRequest.InternetGateway {
		t.Errorf("router.InternetGateway got=%t\nwant=%t", updatedRouter.InternetGateway, updateInternetGatewayRouterRequest.InternetGateway)
	}
	if uuid := updatedRouter.UUID; uuid != initialRouter.UUID {
		t.Errorf("Router.UUID got=%s\nwant=%s", uuid, initialRouter.UUID)
	}
	// The number of InternetGatewayAddresses should be 1, since InternetGateway was updated to true.
	// This assertion must be updated after the API supports IPv6 for private networks.
	if numInternetGatewayAddresses := len(updatedRouter.InternetGatewayAddresses); numInternetGatewayAddresses != 1 {
		t.Errorf("Number of InternetGatewayAddresses got=%d\nwant=%d", numInternetGatewayAddresses, 1)
	}

	// If InternetGateway is true and the Name is a valid FQDN, all InternetGatewayAddresses should have this FQDN as a reverse pointer.
	for index, address := range updatedRouter.InternetGatewayAddresses {
		reversePTR := ""
		if address.ReversePTR != nil {
			reversePTR = *address.ReversePTR
		}

		if reversePTR != updateNameRouterRequest.Name {
			t.Errorf("Router.InternetGatewayAddresses[%d].ReversePTR got=%s\nwant=%s", index, reversePTR, updateNameRouterRequest.Name)
		}
	}

	return *updatedRouter
}

func testListRouters(t *testing.T) {
	t.Helper()

	routers, err := client.Routers.List(t.Context())
	if err != nil {
		t.Fatalf("Routers.List returned error %s\n", err)
	}

	if numRouters := len(routers); numRouters < 1 {
		t.Errorf("Routers.List got=%d\nwant>=%d\n", numRouters, 1)
	}
}

func createNetwork(t *testing.T) cloudscale.Network {
	t.Helper()

	createNetworkRequest := &cloudscale.NetworkCreateRequest{
		Name:                 testRunPrefix,
		AutoCreateIPV4Subnet: new(false),
		ZonalResourceRequest: cloudscale.ZonalResourceRequest{Zone: testZone},
	}
	network, err := client.Networks.Create(t.Context(), createNetworkRequest)
	if err != nil {
		t.Fatalf("Networks.Create returned error %s", err)
	}

	return *network
}

func createSubnet(t *testing.T, network cloudscale.Network) cloudscale.Subnet {
	t.Helper()

	createSubnetRequest := &cloudscale.SubnetCreateRequest{
		Network: network.UUID,
		CIDR:    "192.168.99.0/24",
	}
	subnet, err := client.Subnets.Create(t.Context(), createSubnetRequest)
	if err != nil {
		t.Fatalf("Subnets.Create returned error %s", err)
	}

	return *subnet
}

func testCreateRouterInterface(t *testing.T, router cloudscale.Router, subnet cloudscale.Subnet) cloudscale.RouterInterface {
	t.Helper()

	createRouterInterfaceRequest := cloudscale.CreateInterfaceRequest{
		Network: subnet.Network.UUID,
		Addresses: []cloudscale.CreateAddressRequest{
			{
				Subnet:  subnet.UUID,
				Address: "192.168.99.10",
			},
		},
	}

	routerInterface, err := client.Routers.CreateInterface(t.Context(), router.UUID, createRouterInterfaceRequest)
	if err != nil {
		t.Fatalf("Routers.CreateInterface returned error %s", err)
	}

	if routerInterface.UUID == "" {
		t.Error("Routers.CreateInterface returned interface without UUID")
	}
	if networkUUID := routerInterface.Network.UUID; networkUUID != subnet.Network.UUID {
		t.Errorf("interface.Network.UUID got=%s\nwant=%s", networkUUID, subnet.Network.UUID)
	}
	if numAddresses := len(routerInterface.Addresses); numAddresses != 1 {
		t.Errorf("interface Addresses got=%d\nwant=%d", numAddresses, 1)
	}
	if subnetUUID := routerInterface.Addresses[0].Subnet.UUID; subnetUUID != subnet.UUID {
		t.Errorf("interface.Addresses[0].Subnet.UUID got=%s\nwant=%s", subnetUUID, subnet.UUID)
	}
	if addr := routerInterface.Addresses[0].Address; addr != "192.168.99.10" {
		t.Errorf("interface.Addresses[0].Address got=%s\nwant=%s", addr, "192.168.99.10")
	}

	return *routerInterface
}

func testDeleteRouterInterface(t *testing.T, router cloudscale.Router, routerInterface cloudscale.RouterInterface) {
	t.Helper()

	err := client.Routers.DeleteInterface(t.Context(), router.UUID, routerInterface.UUID)
	if err != nil {
		t.Errorf("Routers.DeleteInterface returned error: %s", err)
	}
}

func testDeleteRouter(t *testing.T, router cloudscale.Router) {
	t.Helper()

	err := client.Routers.Delete(t.Context(), router.UUID)
	if err != nil {
		t.Fatalf("Routers.Delete returned error %s", err)
	}
}

func deleteNetwork(t *testing.T, network cloudscale.Network) {
	t.Helper()

	err := client.Networks.Delete(t.Context(), network.UUID)
	if err != nil {
		t.Fatalf("Networks.Delete returned error %s", err)
	}
}

func TestIntegrationRouter_CRUD(t *testing.T) {
	t.Parallel()

	router := testCreateRouter(t)
	router = testUpdateRouter(t, router)
	testListRouters(t)

	// Set up a network with a subnet so we can attach an interface to the router.
	network := createNetwork(t)
	subnet := createSubnet(t, network)

	routerInterface := testCreateRouterInterface(t, router, subnet)

	// Clean up and test router interface deletion
	testDeleteRouterInterface(t, router, routerInterface)

	// Clean up and test router deletion
	testDeleteRouter(t, router)

	// Clean up: delete network (automatically deletes the subnet too)
	deleteNetwork(t, network)
}
