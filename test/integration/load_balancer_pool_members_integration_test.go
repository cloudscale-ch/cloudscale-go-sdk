//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"github.com/cloudscale-ch/cloudscale-go-sdk/v5"
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

func TestIntegrationServer_LoadBalancer_PrivateNetwork_Port22(t *testing.T) {
	// Ensure integration tests are enabled
	integrationTest(t)

	// Step 1: Create a private network and subnet
	network, subnet, err := createNetworkAndSubnet()
	if err != nil {
		t.Fatalf("Networks and Subnets creation returned error %s\n", err)
	}

	defer func() {
		// Cleanup: Delete the private network
		err = client.Networks.Delete(context.Background(), network.UUID)
		if err != nil {
			t.Fatalf("Networks.Delete returned error %s\n", err)
		}
	}()

	// Step 2: Create a server on the private network
	serverRequest := getDefaultServerRequest()
	serverRequest.Interfaces = &[]cloudscale.InterfaceRequest{{Network: network.UUID}}
	serverRequest.SSHKeys = []string{}
	serverRequest.Password = randomNotVerySecurePassword(10)

	server, err := createServer(t, &serverRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}

	defer func() {
		// Cleanup: Remove the server
		err = client.Servers.Delete(context.Background(), server.UUID)
		if err != nil {
			t.Fatalf("Servers.Delete returned error %s\n", err)
		}
	}()

	privateIP := server.Interfaces[0].Addresses[0].Address

	// Step 3: Create a load balancer
	lbRequest := &cloudscale.LoadBalancerRequest{
		Name:   "test-lb",
		Flavor: "lb-standard",
		ZonalResourceRequest: cloudscale.ZonalResourceRequest{
			Zone: server.Zone.Slug,
		},
	}

	loadBalancer, err := client.LoadBalancers.Create(context.Background(), lbRequest)
	if err != nil {
		t.Fatalf("LoadBalancers.Create returned error %s\n", err)
	}

	defer func() {
		// Cleanup: Remove the load balancer
		err = client.LoadBalancers.Delete(context.Background(), loadBalancer.UUID)
		if err != nil {
			t.Fatalf("LoadBalancers.Delete returned error %s\n", err)
		}
	}()

	// Step 4: Wait for the load balancer to be running
	waitUntilLB("running", loadBalancer.UUID, t)

	// Step 5: Create a load balancer pool
	poolRequest := &cloudscale.LoadBalancerPoolRequest{
		Name:         "test-pool",
		Algorithm:    "round_robin",
		Protocol:     "tcp",
		LoadBalancer: loadBalancer.UUID,
	}

	pool, err := client.LoadBalancerPools.Create(context.Background(), poolRequest)
	if err != nil {
		t.Fatalf("LoadBalancerPools.Create returned error %s\n", err)
	}

	// Step 6: Add the server to the load balancer pool, forwarding traffic to port 22
	memberRequest := &cloudscale.LoadBalancerPoolMemberRequest{
		Name:         "test-member",
		Address:      privateIP,
		ProtocolPort: 22,
		Subnet:       subnet.UUID,
	}

	member, err := client.LoadBalancerPoolMembers.Create(context.Background(), pool.UUID, memberRequest)
	if err != nil {
		t.Fatalf("LoadBalancerPoolMembers.Create returned error %s\n", err)
	}

	// Step 7: Add a listener to the load balancer, using the created pool
	listenerRequest := &cloudscale.LoadBalancerListenerRequest{
		Name:         "test-listener",
		Pool:         pool.UUID, // Associate the listener with the previously created pool
		Protocol:     "tcp",
		ProtocolPort: 22,
	}

	_, err = client.LoadBalancerListeners.Create(context.Background(), listenerRequest)
	if err != nil {
		t.Fatalf("LoadBalancerListeners.Create returned error %s\n", err)
	}

	// Step 8: Add a TCP health monitor to the load balancer pool
	monitorRequest := &cloudscale.LoadBalancerHealthMonitorRequest{
		Type: "tcp",
		Pool: pool.UUID,
	}

	_, err = client.LoadBalancerHealthMonitors.Create(context.Background(), monitorRequest)
	if err != nil {
		t.Fatalf("LoadBalancerHealthMonitors.Create returned error %s\n", err)
	}

	// Define the condition to check for the desired status.
	condition := func(member *cloudscale.LoadBalancerPoolMember) (bool, error) {
		status := "up"
		if member.MonitorStatus == status {
			return true, nil
		}
		return false, fmt.Errorf("waiting for status: %s, current status: %s", status, member.MonitorStatus)
	}

	// Wait for the pool member to reach the desired status.
	_, err = client.LoadBalancerPoolMembers.WaitFor(
		context.Background(),
		pool.UUID,
		member.UUID,
		condition,
	)
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
