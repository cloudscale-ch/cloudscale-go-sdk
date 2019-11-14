// +build integration

package integration

import (
	"context"
	"github.com/cloudscale-ch/cloudscale-go-sdk"
	"testing"
)


func TestIntegrationSubnet_GetAndList(t *testing.T) {
	integrationTest(t)

	createNetworkRequest := &cloudscale.NetworkCreateRequest{
		Name: networkBaseName,
	}

	network, err := client.Networks.Create(context.TODO(), createNetworkRequest)
	if err != nil {
		t.Fatalf("Networks.Create returned error %s\n", err)
	}

	expectedNumberOfSubnets := 1
	if numberOfSubnets := len(network.Subnets); numberOfSubnets != expectedNumberOfSubnets {
		t.Errorf("Number of Subnets got=%#v\nwant=%#v", numberOfSubnets, expectedNumberOfSubnets)
	}

	subnet, err := client.Subnets.Get(context.Background(), network.Subnets[0].UUID)
	if err != nil {
		t.Fatalf("Subnets.Get returned error %s\n", err)
	}

	if uuid := subnet.UUID; uuid != network.Subnets[0].UUID {
		t.Errorf("Subnet.UUID got=%s\nwant=%s", uuid, network.Subnets[0].UUID)
	}

	subnets, err := client.Subnets.List(context.Background())
	if err != nil {
		t.Fatalf("Subnets.List returned error %s\n", err)
	}

	if numSubnets := len(subnets); numSubnets < 0 {
		t.Errorf("Subnets.List got=%d\nwant=%d\n", numSubnets, 1)
	}

	err = client.Networks.Delete(context.Background(), network.UUID)
	if err != nil {
		t.Fatalf("Networks.Delete returned error %s\n", err)
	}
}
