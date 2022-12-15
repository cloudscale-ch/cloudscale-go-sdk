//go:build integration
// +build integration

package integration

import (
	"context"
	"github.com/cloudscale-ch/cloudscale-go-sdk/v2"
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

	createLoadBalancerPoolMemberRequest := &cloudscale.LoadBalancerPoolMemberRequest{
		Name:         testRunPrefix,
		Address:      "5.102.144.111",
		ProtocolPort: 80,
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

	createLoadBalancerPoolMemberRequest := &cloudscale.LoadBalancerPoolMemberRequest{
		Name:         testRunPrefix,
		Address:      "5.102.144.111",
		ProtocolPort: 80,
	}

	poolMember, err := client.LoadBalancerPoolMembers.Create(context.Background(), pool.UUID, createLoadBalancerPoolMemberRequest)
	if err != nil {
		t.Fatalf("LoadBalancerPoolMembers.Create returned error %s\n", err)
	}

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
}
