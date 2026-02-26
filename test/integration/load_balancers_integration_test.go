//go:build integration
// +build integration

package integration

import (
	"context"
	"github.com/cloudscale-ch/cloudscale-go-sdk/v7"
	"reflect"
	"testing"
	"time"
)

func TestIntegrationLoadBalancer_CRUD(t *testing.T) {
	integrationTest(t)

	createLoadBalancerRequest := &cloudscale.LoadBalancerRequest{
		Name:   testRunPrefix,
		Flavor: "lb-standard",
	}
	createLoadBalancerRequest.Zone = testZone

	expected, err := client.LoadBalancers.Create(context.TODO(), createLoadBalancerRequest)
	if err != nil {
		t.Fatalf("LoadBalancers.Create returned error %s\n", err)
	}

	loadBalancer, err := client.LoadBalancers.Get(context.Background(), expected.UUID)
	if err != nil {
		t.Fatalf("LoadBalancers.Get returned error %s\n", err)
	}

	waitUntilLB(expected.UUID, t)

	if h := time.Since(loadBalancer.CreatedAt).Hours(); !(-1 < h && h < 1) {
		t.Errorf("loadBalancer.CreatedAt ourside of expected range. got=%v", loadBalancer.CreatedAt)
	}

	if !reflect.DeepEqual(loadBalancer, expected) {
		t.Errorf("Error = %#v, expected %#v", loadBalancer, expected)
	}

	if numberOfVIPAddresses := len(loadBalancer.VIPAddresses); numberOfVIPAddresses != 2 {
		t.Errorf("numberOfVIPAddresses \n got=%d\nwant=%d", numberOfVIPAddresses, 2)
	}

	loadBalancers, err := client.LoadBalancers.List(context.Background())
	if err != nil {
		t.Fatalf("LoadBalancers.List returned error %s\n", err)
	}

	if numLoadBalancers := len(loadBalancers); numLoadBalancers < 1 {
		t.Errorf("LoadBalancers.List \n got=%d\nwant=%d", numLoadBalancers, 1)
	}

	err = client.LoadBalancers.Delete(context.Background(), loadBalancer.UUID)
	if err != nil {
		t.Fatalf("LoadBalancers.Delete returned error %s\n", err)
	}
}

func TestIntegrationLoadBalancer_PrivateNetwork(t *testing.T) {
	integrationTest(t)

	autoCreateIPV4Subnet := false
	networkRequest := &cloudscale.NetworkCreateRequest{
		Name:                 testRunPrefix,
		AutoCreateIPV4Subnet: &autoCreateIPV4Subnet,
	}
	networkRequest.Zone = testZone

	network, err := client.Networks.Create(context.Background(), networkRequest)
	if err != nil {
		t.Fatalf("Networks.Create returned error %s\n", err)
	}

	subnetRequest := &cloudscale.SubnetCreateRequest{
		Network: network.UUID,
		CIDR:    "192.168.7.0/24",
	}

	subnet, err := client.Subnets.Create(context.Background(), subnetRequest)
	if err != nil {
		t.Fatalf("Subnets.Create returned error %s\n", err)
	}

	vipAddress := "192.168.7.7"
	createLoadBalancerRequest := &cloudscale.LoadBalancerRequest{
		Name:   testRunPrefix,
		Flavor: "lb-standard",
		VIPAddresses: &[]cloudscale.VIPAddressRequest{
			{
				Address: vipAddress,
				Subnet:  subnet.UUID,
			},
		},
	}
	createLoadBalancerRequest.Zone = testZone

	loadBalancer, err := client.LoadBalancers.Create(context.TODO(), createLoadBalancerRequest)
	if err != nil {
		t.Fatalf("LoadBalancers.Create returned error %s\n", err)
	}

	if numVIPAddresses := len(loadBalancer.VIPAddresses); numVIPAddresses != 1 {
		t.Errorf("numVIPAddresses \n got=%d\nwant=%d", numVIPAddresses, 1)
	}

	if loadBalancerVIPAddress := loadBalancer.VIPAddresses[0].Address; loadBalancerVIPAddress != vipAddress {
		t.Errorf("loadBalancerVIPAddress \n got=%s\nwant=%s", loadBalancerVIPAddress, vipAddress)
	}

	if loadBalancerSubnetUUID := loadBalancer.VIPAddresses[0].Subnet.UUID; loadBalancerSubnetUUID != subnet.UUID {
		t.Errorf("loadBalancerSubnetUUID \n got=%s\nwant=%s", loadBalancerSubnetUUID, subnet.UUID)
	}

	waitUntilLB(loadBalancer.UUID, t)

	err = client.LoadBalancers.Delete(context.Background(), loadBalancer.UUID)
	if err != nil {
		t.Fatalf("LoadBalancers.Delete returned error %s\n", err)
	}

	err = client.Networks.Delete(context.Background(), network.UUID)
	if err != nil {
		t.Fatalf("Networks.Delete returned error %s\n", err)
	}
}

func TestIntegrationLoadBalancer_Update(t *testing.T) {
	integrationTest(t)

	createLoadBalancerRequest := &cloudscale.LoadBalancerRequest{
		Name:   testRunPrefix,
		Flavor: "lb-standard",
	}
	createLoadBalancerRequest.Zone = testZone

	lb, err := client.LoadBalancers.Create(context.TODO(), createLoadBalancerRequest)
	if err != nil {
		t.Fatalf("loadBalancer.Create returned error %s\n", err)
	}

	waitUntilLB(lb.UUID, t)

	newName := testRunPrefix + "-renamed"
	updateRequest := &cloudscale.LoadBalancerRequest{
		Name: newName,
	}

	uuid := lb.UUID
	err = client.LoadBalancers.Update(context.Background(), uuid, updateRequest)
	if err != nil {
		t.Fatalf("LoadBalancers.Update returned error %s\n", err)
	}

	loadBalancer, err := client.LoadBalancers.Get(context.Background(), uuid)
	if err != nil {
		t.Fatalf("LoadBalancers.Get returned error %s\n", err)
	}

	if name := loadBalancer.Name; name != newName {
		t.Errorf("loadbalancer.Name \n got=%s\nwant=%s", name, newName)
	}

	err = client.LoadBalancers.Delete(context.Background(), uuid)
	if err != nil {
		t.Fatalf("LoadBalancers.Delete returned error %s\n", err)
	}
}

func waitUntilLB(uuid string, t *testing.T) *cloudscale.LoadBalancer {
	lb, err := client.LoadBalancers.WaitFor(context.Background(), uuid, cloudscale.LoadBalancerIsRunning)
	if err != nil {
		t.Fatalf("client.LoadBalancers.WaitFor returned error %s\n", err)
	}

	return lb
}
