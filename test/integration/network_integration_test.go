// +build integration

package integration

import (
	"context"
	"github.com/cloudscale-ch/cloudscale-go-sdk"
	"regexp"
	"sync"
	"testing"
	"time"
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

	autoCreateSubnet := false
	createNetworkRequest := &cloudscale.NetworkCreateRequest{
		Name: networkBaseName,
		AutoCreateIPV4Subnet: &autoCreateSubnet,
	}
	network, err := client.Networks.Create(context.TODO(), createNetworkRequest)
	if err != nil {
		t.Fatalf("Networks.Create returned error %s\n", err)
	}

	createSubnetRequest := &cloudscale.SubnetCreateRequest{
		Network: network.UUID,
		CIDR: "192.168.42.0/24",
	}
	subnet, err := client.Subnets.Create(context.TODO(), createSubnetRequest)
	if err != nil {
		t.Fatalf("Subnets.Create returned error %s\n", err)
	}

	cases := []struct {
		name       string
		in         *[]cloudscale.InterfaceRequest
		expectedIP string
	}{
		{"Attach by network UUID", &[]cloudscale.InterfaceRequest{
			{
				Network: network.UUID,
			},
		}, `192\.168\.42\.[0-9]*`},
		{"Attach by subnet UUID", &[]cloudscale.InterfaceRequest{
			{
				Addresses: &[]cloudscale.AddressRequest{
					{
						Subnet: subnet.UUID,
					},
				},
			},
		}, `192\.168\.42\.[0-9]*`},
		{"Attach by subnet UUID with predefined IP", &[]cloudscale.InterfaceRequest{
			{
				Addresses: &[]cloudscale.AddressRequest{
					{
						Subnet:  subnet.UUID,
						Address: "192.168.42.242",
					},
				},
			},
		}, `192\.168\.42\.242`},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			createServerRequest := &cloudscale.ServerRequest{
				Name:         "go-sdk-integration-test-network",
				Flavor:       "flex-2",
				Image:        DefaultImageSlug,
				VolumeSizeGB: 10,
				Interfaces:   tt.in,
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
			singleInterface := server.Interfaces[0];
			if singleInterface.Network.UUID != network.UUID {
				t.Errorf("Attatched to wrong Network\ngot=%#v\nwant=%#v", singleInterface.Network.UUID, network.UUID)
			}
			re := regexp.MustCompile(tt.expectedIP)
			if !re.Match([]byte(singleInterface.Addresses[0].Address)) {
				t.Errorf("Expected IP regex does not match\ngot=%#v\nwant=%#v", singleInterface.Addresses[0].Address, tt.expectedIP)
			}

			err = client.Servers.Delete(context.Background(), server.UUID)
			if err != nil {
				t.Fatalf("Servers.Delete returned error %s\n", err)
			}
		})
	}

	// sending the next request immediately can cause errors, since the port cleanup process is still ongoing
	time.Sleep(5 * time.Second)
	err = client.Networks.Delete(context.Background(), network.UUID)
	if err != nil {
		t.Fatalf("Networks.Delete returned error %s\n", err)
	}
}

func TestIntegrationNetwork_Reattach(t *testing.T) {
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

	createSubnetRequest := &cloudscale.SubnetCreateRequest{
		Network: network.UUID,
		CIDR:    "192.168.77.0/24",
	}
	subnet, err := client.Subnets.Create(context.TODO(), createSubnetRequest)
	if err != nil {
		t.Fatalf("Subnets.Create returned error %s\n", err)
	}

	interfaces := []cloudscale.InterfaceRequest{
		{Network: "public"},
	}
	createServerRequest := &cloudscale.ServerRequest{
		Name:         "go-sdk-integration-test-network",
		Flavor:       "flex-2",
		Image:        DefaultImageSlug,
		VolumeSizeGB: 10,
		Interfaces:   &interfaces,
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

	addresses := []cloudscale.AddressRequest{{
		Subnet:  subnet.UUID,
		Address: "192.168.77.77",
	}}
	interfaces = append(interfaces, cloudscale.InterfaceRequest{
		Addresses: &addresses,
	})
	updateRequest := cloudscale.ServerUpdateRequest{
		Interfaces: &interfaces,
	}
	err = client.Servers.Update(context.Background(), server.UUID, &updateRequest)
	if err != nil {
		t.Fatalf("Servers.Update returned error %s\n", err)
	}

	updatedServer, err := client.Servers.Get(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Get returned error %s\n", err)
	}
	if numNetworks := len(updatedServer.Interfaces); numNetworks != 2 {
		t.Errorf("Attatched to number of Networks\ngot=%#v\nwant=%#v", numNetworks, 2)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}

	// sending the next request immediately can cause errors, since the port cleanup process is still ongoing
	time.Sleep(5 * time.Second)
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
