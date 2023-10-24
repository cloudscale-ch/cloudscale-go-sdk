//go:build integration
// +build integration

package integration

import (
	"context"
	"github.com/cloudscale-ch/cloudscale-go-sdk/v4"
	"reflect"
	"testing"
	"time"
)

func TestIntegrationLoadBalancerPoolMember_CRUD(t *testing.T) {
	integrationTest(t)

	lb, err := createLoadBalancer()
	if err != nil {
		t.Fatalf("LoadBalancers.Create returned error %s\n", err)
	}

	waitUntilLB("running", lb.UUID, t)

	pool, err := createPoolOnLB(lb)
	if err != nil {
		t.Fatalf("LoadBalancerPools.Create returned error %s\n", err)
	}

	network, subnet, err := createNetworkAndSubnet()
	if err != nil {
		t.Fatalf("error while creating network and subnet: %s\n", err)
	}

	createLoadBalancerPoolMemberRequest := &cloudscale.LoadBalancerPoolMemberRequest{
		Name:         testRunPrefix,
		Address:      "192.168.42.11",
		ProtocolPort: 80,
		Subnet:       subnet.UUID,
	}

	expected, err := client.LoadBalancerPoolMembers.Create(context.Background(), pool.UUID, createLoadBalancerPoolMemberRequest)
	if err != nil {
		t.Fatalf("LoadBalancerPoolMembers.Create returned error %s\n", err)
	}

	loadBalancerPoolMember, err := client.LoadBalancerPoolMembers.Get(context.Background(), pool.UUID, expected.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerPoolMembers.Get returned error %s\n", err)
	}

	if h := time.Since(loadBalancerPoolMember.CreatedAt).Hours(); !(-1 < h && h < 1) {
		t.Errorf("loadBalancerPoolMember.CreatedAt ourside of expected range. got=%v", loadBalancerPoolMember.CreatedAt)
	}

	if !reflect.DeepEqual(loadBalancerPoolMember, expected) {
		t.Errorf("Error = %#v, expected %#v", loadBalancerPoolMember, expected)
	}

	if memberPoolUUID := loadBalancerPoolMember.Pool.UUID; memberPoolUUID != pool.UUID {
		t.Errorf("poolLbUUID \n got=%#v\nwant=%#v", memberPoolUUID, pool.UUID)
	}

	if lbUUID := loadBalancerPoolMember.LoadBalancer.UUID; lbUUID != lb.UUID {
		t.Errorf("poolLbUUID \n got=%#v\nwant=%#v", lbUUID, lb.UUID)
	}

	loadBalancerPoolMembers, err := client.LoadBalancerPoolMembers.List(context.Background(), pool.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerPoolMembers.List returned error %s\n", err)
	}

	if numLoadBalancerPoolMembers := len(loadBalancerPoolMembers); numLoadBalancerPoolMembers < 1 {
		t.Errorf("LoadBalancerListeners.List \n got=%d\nwant=%d", numLoadBalancerPoolMembers, 1)
	}

	err = client.LoadBalancerPoolMembers.Delete(context.Background(), pool.UUID, expected.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerPoolMembers.Delete returned error %s\n", err)
	}

	err = client.LoadBalancerPools.Delete(context.Background(), pool.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerPools.Delete returned error %s\n", err)
	}

	err = client.LoadBalancers.Delete(context.Background(), lb.UUID)
	if err != nil {
		t.Fatalf("LoadBalancers.Delete returned error %s\n", err)
	}

	err = client.Networks.Delete(context.Background(), network.UUID)
	if err != nil {
		t.Fatalf("Networks.Delete returned error %s\n", err)
	}
}

func TestIntegrationLoadBalancerPoolMember_Update(t *testing.T) {
	integrationTest(t)

	lb, err := createLoadBalancer()
	if err != nil {
		t.Fatalf("LoadBalancers.Create returned error %s\n", err)
	}

	waitUntilLB("running", lb.UUID, t)

	pool, err := createPoolOnLB(lb)
	if err != nil {
		t.Fatalf("LoadBalancerPools.Create returned error %s\n", err)
	}

	network, subnet, err := createNetworkAndSubnet()
	if err != nil {
		t.Fatalf("error while creating network and subnet: %s\n", err)
	}

	createLoadBalancerPoolMemberRequest := &cloudscale.LoadBalancerPoolMemberRequest{
		Name:         testRunPrefix,
		Address:      "192.168.42.11",
		ProtocolPort: 80,
		Subnet:       subnet.UUID,
	}

	poolMember, err := client.LoadBalancerPoolMembers.Create(context.Background(), pool.UUID, createLoadBalancerPoolMemberRequest)
	if err != nil {
		t.Fatalf("LoadBalancerPoolMembers.Create returned error %s\n", err)
	}

	// Update Name
	newName := testRunPrefix + "-renamed"
	updateRequest := &cloudscale.LoadBalancerPoolMemberRequest{
		Name: newName,
	}

	uuid := poolMember.UUID
	err = client.LoadBalancerPoolMembers.Update(context.Background(), pool.UUID, uuid, updateRequest)
	if err != nil {
		t.Fatalf("LoadBalancerPoolMembers.Update returned error %s\n", err)
	}

	updated, err := client.LoadBalancerPoolMembers.Get(context.Background(), pool.UUID, uuid)
	if err != nil {
		t.Fatalf("LoadBalancerPoolMembers.Get returned error %s\n", err)
	}

	if name := updated.Name; name != newName {
		t.Errorf("updated.Name \n got=%s\nwant=%s", name, newName)
	}

	// Disable
	newEnabled := false
	updateRequest2 := &cloudscale.LoadBalancerPoolMemberRequest{
		Enabled: &newEnabled,
	}

	err = client.LoadBalancerPoolMembers.Update(context.Background(), pool.UUID, uuid, updateRequest2)
	if err != nil {
		t.Fatalf("LoadBalancerPoolMembers.Update returned error %s\n", err)
	}

	updated2, err := client.LoadBalancerPoolMembers.Get(context.Background(), pool.UUID, uuid)
	if err != nil {
		t.Fatalf("LoadBalancerPoolMembers.Get returned error %s\n", err)
	}

	if enabled := updated2.Enabled; enabled != newEnabled {
		t.Errorf("updated2.Enabled \n got=%t\nwant=%t", enabled, newEnabled)
	}

	err = client.LoadBalancerPoolMembers.Delete(context.Background(), pool.UUID, updated.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerPoolMembers.Delete returned error %s\n", err)
	}

	err = client.LoadBalancerPools.Delete(context.Background(), pool.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerPools.Delete returned error %s\n", err)
	}

	err = client.LoadBalancers.Delete(context.Background(), lb.UUID)
	if err != nil {
		t.Fatalf("LoadBalancers.Delete returned error %s\n", err)
	}

	err = client.Networks.Delete(context.Background(), network.UUID)
	if err != nil {
		t.Fatalf("Networks.Delete returned error %s\n", err)
	}
}

func createNetworkAndSubnet() (*cloudscale.Network, *cloudscale.Subnet, error) {
	autoCreateSubnet := false

	network, err := client.Networks.Create(context.TODO(), &cloudscale.NetworkCreateRequest{
		Name:                 testRunPrefix,
		AutoCreateIPV4Subnet: &autoCreateSubnet,
		ZonalResourceRequest: cloudscale.ZonalResourceRequest{
			Zone: testZone,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	subnets, err := client.Subnets.Create(context.TODO(), &cloudscale.SubnetCreateRequest{
		CIDR:    "192.168.42.0/24",
		Network: network.UUID,
	})
	if err != nil {
		return network, nil, err
	}

	return network, subnets, err
}
