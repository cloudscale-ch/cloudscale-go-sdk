// +build integration

package integration

import (
	"context"
	"github.com/cloudscale-ch/cloudscale-go-sdk"
	"strings"
	"sync"
	"testing"
)

const networkBaseName = "go-sdk-integration-test-network"

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

	if numNetworks := len(networks); numNetworks < 0 {
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
		Name:   networkBaseName,
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

func TestIntegrationNetwork_DeleteRemainingNetworks(t *testing.T) {
	networks, err := client.Networks.List(context.Background())
	if err != nil {
		t.Fatalf("Networks.List returned error %s\n", err)
	}

	for _, network := range networks {
		if strings.HasPrefix(network.Name, "go-sdk-integration-test") {
			t.Errorf("Found not deleted network: %s\n", network.Name)
			err = client.Networks.Delete(context.Background(), network.UUID)
			if err != nil {
				t.Errorf("Networks.Delete returned error %s\n", err)
			}
		}
	}
}
