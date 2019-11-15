// +build integration

package integration

import (
	"context"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/cloudscale-ch/cloudscale-go-sdk"
	"sync"
	"testing"
)

const networkBaseName = "go-sdk-integration-test-network"

func DeleteNetworkWithRetry(network *cloudscale.Network) error {
	operation := func() error {
		err := client.Networks.Delete(context.Background(), network.UUID)
		if err != nil {
			msg := fmt.Sprintf("Networks.Delete returned error %s\n", err)
			return errors.New(msg)
		}
		return nil
	}
	err := backoff.Retry(operation, backoff.NewExponentialBackOff())
	if err != nil {
		msg := fmt.Sprintf("Retries exeeded: %s\n", err)
		return errors.New(msg)
	}
	return nil
}

func TestIntegrationNetwork_CRUD(t *testing.T) {
	integrationTest(t)

	createNetworkRequest := &cloudscale.NetworkCreateRequest{
		Name: networkBaseName,
	}

	expected, err := client.Networks.Create(context.TODO(), createNetworkRequest)
	if err != nil {
		t.Fatalf("Networks.Create returned error %s\n", err)
	}

	network, err := client.Networks.Get(context.Background(), expected.UUID)
	if err != nil {
		t.Fatalf("Networks.Get returned error %s\n", err)
	}

	if uuid := network.UUID; uuid != expected.UUID {
		t.Errorf("Network.UUID got=%s\nwant=%s", uuid, expected.UUID)
	}

	expectedNumberOfSubnets := 1
	if numberOfSubnets := len(network.Subnets); numberOfSubnets != expectedNumberOfSubnets {
		t.Errorf("Number of Subnets got=%#v\nwant=%#v", numberOfSubnets, expectedNumberOfSubnets)
	}

	networks, err := client.Networks.List(context.Background())
	if err != nil {
		t.Fatalf("Networks.List returned error %s\n", err)
	}

	if numNetworks := len(networks); numNetworks == 0 {
		t.Errorf("Network.List got=%d\nwant=%d\n", numNetworks, 1)
	}

	err = client.Networks.Delete(context.Background(), network.UUID)
	if err != nil {
		t.Fatalf("Networks.Delete returned error %s\n", err)
	}
}

func TestIntegrationNetwork_CreateWithoutSubnet(t *testing.T) {
	integrationTest(t)

	autoCreateSubnet := false
	createNetworkRequest := &cloudscale.NetworkCreateRequest{
		Name:                 networkBaseName,
		AutoCreateIPV4Subnet: &autoCreateSubnet,
	}

	network, err := client.Networks.Create(context.TODO(), createNetworkRequest)
	if err != nil {
		t.Fatalf("Networks.Create returned error %s\n", err)
	}

	expectedNumberOfSubnets := 0
	if numberOfSubnets := len(network.Subnets); numberOfSubnets != expectedNumberOfSubnets {
		t.Errorf("Number of Subnets got=%#v\nwant=%#v", numberOfSubnets, expectedNumberOfSubnets)
	}

	err = client.Networks.Delete(context.Background(), network.UUID)
	if err != nil {
		t.Fatalf("Networks.Delete returned error %s\n", err)
	}
}

func TestIntegrationNetwork_CreateAttached(t *testing.T) {
	integrationTest(t)

	createNetworkRequest := &cloudscale.NetworkCreateRequest{
		Name: networkBaseName,
	}

	network, err := client.Networks.Create(context.TODO(), createNetworkRequest)
	if err != nil {
		t.Fatalf("Networks.Create returned error %s\n", err)
	}

	interfaceRequests := []cloudscale.InterfaceRequest{{
		Network: network.UUID,
	},}
	createServerRequest := &cloudscale.ServerRequest{
		Name:         "go-sdk-integration-test-network",
		Flavor:       "flex-2",
		Image:        DefaultImageSlug,
		VolumeSizeGB: 10,
		Interfaces:   &interfaceRequests,
		SSHKeys: []string{
			pubKey,
		},
	}

	server, err := client.Servers.Create(context.Background(), createServerRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}
	waitUntil("running", server.UUID, t)

	if numNetworks := len(server.Interfaces); numNetworks != 1 {
		t.Errorf("Attatched to number of Networks\ngot=%#v\nwant=%#v", numNetworks, 1)
	}
	if singleInterface := server.Interfaces[0]; singleInterface.Network.UUID != network.UUID {
		t.Errorf("Attatched to wrong Network\ngot=%#v\nwant=%#v", singleInterface.Network.UUID, network.UUID)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}

	err = DeleteNetworkWithRetry(network)
	if err != nil {
		t.Fatalf("Could not delete network: %s\n", err)
	}
}

func TestIntegrationNetwork_AttachWithoutIP(t *testing.T) {
	integrationTest(t)

	interfaceRequests := []cloudscale.InterfaceRequest{{
		Network: "public",
	},}
	createServerRequest := &cloudscale.ServerRequest{
		Name:         "go-sdk-integration-test-network",
		Flavor:       "flex-2",
		Image:        DefaultImageSlug,
		VolumeSizeGB: 10,
		Interfaces:   &interfaceRequests,
		SSHKeys: []string{
			pubKey,
		},
	}

	server, err := client.Servers.Create(context.Background(), createServerRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}
	waitUntil("running", server.UUID, t)

	if numNetworks := len(server.Interfaces); numNetworks != 1 {
		t.Errorf("Attatched to number of Networks\ngot=%#v\nwant=%#v", numNetworks, 1)
	}

	createNetworkRequest := &cloudscale.NetworkCreateRequest{
		Name: networkBaseName,
	}

	network, err := client.Networks.Create(context.TODO(), createNetworkRequest)
	if err != nil {
		t.Fatalf("Networks.Create returned error %s\n", err)
	}

	interfaceRequests = append(interfaceRequests, cloudscale.InterfaceRequest{
		Network:   network.UUID,
		Addresses: &[]string{},
	})
	updateServerRequest := &cloudscale.ServerUpdateRequest{
		Interfaces: &interfaceRequests,
	}

	err = client.Servers.Update(context.TODO(), server.UUID, updateServerRequest)
	if err != nil {
		t.Errorf("Servers.Update returned error: %v", err)
	}

	server, err = client.Servers.Get(context.Background(), server.UUID)
	if err != nil {
		t.Errorf("Server.Get returned error: %v", err)
	}
	if numNetworks := len(server.Interfaces); numNetworks != 2 {
		t.Errorf("Attatched to number of Networks\ngot=%#v\nwant=%#v", numNetworks, 2)
	}
	second := server.Interfaces[1]
	if second.Network.UUID != network.UUID {
		t.Errorf("Attatched to wrong Network\ngot=%#v\nwant=%#v", second.Network.UUID, network.UUID)
	}

	if addressCount := len(second.Addresses); addressCount > 0 {
		t.Errorf("Expected no addresses\ngot=%#v", addressCount)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}

	err = DeleteNetworkWithRetry(network)
	if err != nil {
		t.Fatalf("Could not delete network: %s\n", err)
	}
}

func TestIntegrationNetwork_Update(t *testing.T) {

	createNetworkRequest := &cloudscale.NetworkCreateRequest{
		Name: networkBaseName,
		MTU:  1500,
	}

	network, err := client.Networks.Create(context.TODO(), createNetworkRequest)
	if err != nil {
		t.Fatalf("Networks.Create returned error %s\n", err)
	}

	expectedNewMTU := 6601
	updateRequest := &cloudscale.NetworkUpdateRequest{
		MTU: expectedNewMTU,
	}

	err = client.Networks.Update(context.Background(), network.UUID, updateRequest)
	if err != nil {
		t.Fatalf("Networks.Update returned error %s\n", err)
	}

	updatedNetwork, err := client.Networks.Get(context.Background(), network.UUID)
	if err != nil {
		t.Fatalf("Networks.Get returned error %s\n", err)
	}

	if actualMTU := updatedNetwork.MTU; actualMTU != expectedNewMTU {
		t.Errorf("Network MTU\ngot=%#v\nwant=%#v", updatedNetwork.MTU, expectedNewMTU)
	}

	err = client.Networks.Delete(context.Background(), network.UUID)
	if err != nil {
		t.Fatalf("Networks.Delete returned error %s\n", err)
	}
}

func TestIntegrationNetwork_MultiSite(t *testing.T) {
	integrationTest(t)

	allZones, err := getAllZones()
	if err != nil {
		t.Fatalf("getAllRegions returned error %s\n", err)
	}

	if len(allZones) <= 1 {
		t.Skip("Skipping MultiSite test.")
	}

	var wg sync.WaitGroup

	for _, zone := range allZones {
		wg.Add(1)
		go createNetworkInZoneAndAssert(t, zone, &wg)
	}

	wg.Wait()
}

func createNetworkInZoneAndAssert(t *testing.T, zone cloudscale.Zone, wg *sync.WaitGroup) {
	defer wg.Done()

	createNetworkRequest := &cloudscale.NetworkCreateRequest{
		Name: networkBaseName,
	}

	createNetworkRequest.Zone = zone.Slug

	network, err := client.Networks.Create(context.TODO(), createNetworkRequest)
	if err != nil {
		t.Fatalf("Networks.Create returned error %s\n", err)
	}

	if network.Zone != zone {
		t.Errorf("Network in wrong Zone\n got=%#v\nwant=%#v", network.Zone, zone)
	}

	err = client.Networks.Delete(context.Background(), network.UUID)
	if err != nil {
		t.Errorf("Networks.Delete returned error %s\n", err)
	}
}
