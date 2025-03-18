//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"github.com/cloudscale-ch/cloudscale-go-sdk/v5"
	"regexp"
	"sync"
	"testing"
	"time"
)

func TestIntegrationNetwork_CRUD(t *testing.T) {
	integrationTest(t)

	createNetworkRequest := &cloudscale.NetworkCreateRequest{
		Name: testRunPrefix,
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

	if h := time.Since(network.CreatedAt).Hours(); !(-1 < h && h < 1) {
		t.Errorf("network.CreatedAt ourside of expected range. got=%v", network.CreatedAt)
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
		Name:                 testRunPrefix,
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
		Name:                 testRunPrefix,
		AutoCreateIPV4Subnet: &autoCreateSubnet,
	}
	network, err := client.Networks.Create(context.TODO(), createNetworkRequest)
	if err != nil {
		t.Fatalf("Networks.Create returned error %s\n", err)
	}

	createSubnetRequest := &cloudscale.SubnetCreateRequest{
		Network: network.UUID,
		CIDR:    "192.168.42.0/24",
	}
	subnet, err := client.Subnets.Create(context.TODO(), createSubnetRequest)
	if err != nil {
		t.Fatalf("Subnets.Create returned error %s\n", err)
	}

	cases := []struct {
		name                string
		in                  *[]cloudscale.InterfaceRequest
		expectedNumNetworks int
		expectedIP          string
	}{
		{"Attach by network UUID", &[]cloudscale.InterfaceRequest{
			{
				Network: network.UUID,
			},
		}, 1, `192\.168\.42\.[0-9]*`},
		{"Attach by subnet UUID", &[]cloudscale.InterfaceRequest{
			{
				Addresses: &[]cloudscale.AddressRequest{
					{
						Subnet: subnet.UUID,
					},
				},
			},
		}, 1, `192\.168\.42\.[0-9]*`},
		{"Attach by subnet UUID with predefined IP", &[]cloudscale.InterfaceRequest{
			{
				Addresses: &[]cloudscale.AddressRequest{
					{
						Subnet:  subnet.UUID,
						Address: "192.168.42.242",
					},
				},
			},
		}, 1, `192\.168\.42\.242`},
		{"Attach by network UUID without IP (Layer 2)", &[]cloudscale.InterfaceRequest{
			{
				Network: "public",
			},
			{
				Network:   network.UUID,
				Addresses: &[]cloudscale.AddressRequest{},
			},
		}, 2, ""},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			createServerRequest := &cloudscale.ServerRequest{
				Name:         testRunPrefix,
				Flavor:       "flex-4-2",
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
			_, err = client.Servers.WaitFor(
				context.Background(),
				server.UUID,
				serverRunningCondition,
			)
			if err != nil {
				t.Fatalf("Servers.WaitFor returned error %s\n", err)
			}

			if numNetworks := len(server.Interfaces); numNetworks != tt.expectedNumNetworks {
				t.Errorf("Attatched to number of Networks\ngot=%#v\nwant=%#v", numNetworks, tt.expectedNumNetworks)
			}
			lastNetworkInterface := server.Interfaces[len(server.Interfaces)-1]
			if lastNetworkInterface.Network.UUID != network.UUID {
				t.Errorf("Attatched to wrong Network\ngot=%#v\nwant=%#v", lastNetworkInterface.Network.UUID, network.UUID)
			}
			if tt.expectedIP != "" {
				re := regexp.MustCompile(tt.expectedIP)
				if !re.Match([]byte(lastNetworkInterface.Addresses[0].Address)) {
					t.Errorf("Expected IP regex does not match\ngot=%#v\nwant=%#v", lastNetworkInterface.Addresses[0].Address, tt.expectedIP)
				}
			} else {
				if len(lastNetworkInterface.Addresses) != 0 {
					t.Errorf("Expected no IP addresses\ngot=%#v", len(lastNetworkInterface.Addresses))
				}
			}

			// this is required especially for the 'without IP' case.
			time.Sleep(20 * time.Second)
			err = client.Servers.Delete(context.Background(), server.UUID)
			if err != nil {
				t.Fatalf("Servers.Delete returned error %s\n", err)
			}
		})
	}

	err = client.Networks.Delete(context.Background(), network.UUID)
	if err != nil {
		t.Fatalf("Networks.Delete returned error %s\n", err)
	}
}

func TestIntegrationNetwork_Reattach(t *testing.T) {
	integrationTest(t)

	autoCreateSubnet := false
	createNetworkRequest := &cloudscale.NetworkCreateRequest{
		Name:                 testRunPrefix,
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
		Name:         fmt.Sprintf("%s-network", testRunPrefix),
		Flavor:       "flex-4-2",
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
	_, err = client.Servers.WaitFor(
		context.Background(),
		server.UUID,
		serverRunningCondition,
	)
	if err != nil {
		t.Fatalf("Servers.WaitFor returned error %s\n", err)
	}

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

	err = client.Networks.Delete(context.Background(), network.UUID)
	if err != nil {
		t.Fatalf("Networks.Delete returned error %s\n", err)
	}
}

func TestIntegrationNetwork_Reorder(t *testing.T) {
	integrationTest(t)

	autoCreateSubnet := false
	createNetworkRequest := &cloudscale.NetworkCreateRequest{
		Name:                 testRunPrefix,
		AutoCreateIPV4Subnet: &autoCreateSubnet,
	}
	network, err := client.Networks.Create(context.TODO(), createNetworkRequest)
	if err != nil {
		t.Fatalf("Networks.Create returned error %s\n", err)
	}

	createSubnetRequest := &cloudscale.SubnetCreateRequest{
		Network: network.UUID,
		CIDR:    "192.168.177.0/24",
	}
	subnet, err := client.Subnets.Create(context.TODO(), createSubnetRequest)
	if err != nil {
		t.Fatalf("Subnets.Create returned error %s\n", err)
	}

	interfaces := []cloudscale.InterfaceRequest{
		{Network: "public"},
		{Network: network.UUID},
	}
	createServerRequest := &cloudscale.ServerRequest{
		Name:         fmt.Sprintf("%s-network", testRunPrefix),
		Flavor:       "flex-4-2",
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
	_, err = client.Servers.WaitFor(
		context.Background(),
		server.UUID,
		serverRunningCondition,
	)
	if err != nil {
		t.Fatalf("Servers.WaitFor returned error %s\n", err)
	}
	if numNetworks := len(server.Interfaces); numNetworks != 2 {
		t.Errorf("Attatched to number of Networks\ngot=%#v\nwant=%#v", numNetworks, 2)
	}

	if subnetUUID := server.Interfaces[1].Addresses[0].Subnet.UUID; subnetUUID != subnet.UUID {
		t.Errorf("Unexpected subnet UUID on second interface\ngot=%#v\nwant=%#v", subnetUUID, subnet.UUID)
	}

	for i, j := 0, len(interfaces)-1; i < j; i, j = i+1, j-1 {
		interfaces[i], interfaces[j] = interfaces[j], interfaces[i]
	}

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
	if subnetUUID := updatedServer.Interfaces[0].Addresses[0].Subnet.UUID; subnetUUID != subnet.UUID {
		t.Errorf("Unexpected subnet UUID on first interface\ngot=%#v\nwant=%#v", subnetUUID, subnet.UUID)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}

	err = client.Networks.Delete(context.Background(), network.UUID)
	if err != nil {
		t.Fatalf("Networks.Delete returned error %s\n", err)
	}
}

func TestIntegrationNetwork_Update(t *testing.T) {

	createNetworkRequest := &cloudscale.NetworkCreateRequest{
		Name: testRunPrefix,
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
		Name: testRunPrefix,
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
