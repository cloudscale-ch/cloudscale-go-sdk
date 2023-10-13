//go:build integration
// +build integration

package integration

import (
	"context"
	"github.com/cloudscale-ch/cloudscale-go-sdk/v3"
	"reflect"
	"testing"
	"time"
)

func TestIntegrationLoadBalancerListener_CRUD(t *testing.T) {
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

	createLoadBalancerListenerRequest := &cloudscale.LoadBalancerListenerRequest{
		Name:         testRunPrefix,
		Pool:         pool.UUID,
		Protocol:     "tcp",
		ProtocolPort: 80,
	}

	expected, err := client.LoadBalancerListeners.Create(context.Background(), createLoadBalancerListenerRequest)
	if err != nil {
		t.Fatalf("LoadBalancerListeners.Create returned error %s\n", err)
	}

	loadBalancerListener, err := client.LoadBalancerListeners.Get(context.Background(), expected.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerListeners.Get returned error %s\n", err)
	}

	if h := time.Since(loadBalancerListener.CreatedAt).Hours(); !(-1 < h && h < 1) {
		t.Errorf("loadBalancerListener.CreatedAt ourside of expected range. got=%v", loadBalancerListener.CreatedAt)
	}

	if !reflect.DeepEqual(loadBalancerListener, expected) {
		t.Errorf("Error = %#v, expected %#v", loadBalancerListener, expected)
	}

	if poolLbUUID := loadBalancerListener.Pool.UUID; poolLbUUID != pool.UUID {
		t.Errorf("poolLbUUID \n got=%#v\nwant=%#v", poolLbUUID, pool.UUID)
	}

	if lbUUID := loadBalancerListener.LoadBalancer.UUID; lbUUID != lb.UUID {
		t.Errorf("lbUUID \n got=%#v\nwant=%#v", lbUUID, lb.UUID)
	}

	loadBalancerListeners, err := client.LoadBalancerListeners.List(context.Background())
	if err != nil {
		t.Fatalf("LoadBalancerListeners.List returned error %s\n", err)
	}

	if numLoadbalancerListeners := len(loadBalancerListeners); numLoadbalancerListeners < 1 {
		t.Errorf("LoadBalancerListeners.List \n got=%d\nwant=%d", numLoadbalancerListeners, 1)
	}

	err = client.LoadBalancerListeners.Delete(context.Background(), expected.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerListeners.Delete returned error %s\n", err)
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

func TestIntegrationLoadBalancerListener_Update(t *testing.T) {
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

	createLoadBalancerListenerRequest := &cloudscale.LoadBalancerListenerRequest{
		Name:         testRunPrefix,
		Pool:         pool.UUID,
		Protocol:     "tcp",
		ProtocolPort: 80,
	}

	listener, err := client.LoadBalancerListeners.Create(context.Background(), createLoadBalancerListenerRequest)
	if err != nil {
		t.Fatalf("LoadBalancerListeners.Create returned error %s\n", err)
	}

	// update name
	newName := testRunPrefix + "-renamed"
	updateRequest := &cloudscale.LoadBalancerListenerRequest{
		Name: newName,
	}

	uuid := listener.UUID
	err = client.LoadBalancerListeners.Update(context.Background(), uuid, updateRequest)
	if err != nil {
		t.Fatalf("LoadBalancerListeners.Update returned error %s\n", err)
	}

	updated, err := client.LoadBalancerListeners.Get(context.Background(), uuid)
	if err != nil {
		t.Fatalf("LoadBalancerListeners.Get returned error %s\n", err)
	}

	if name := updated.Name; name != newName {
		t.Errorf("updated.Name \n got=%s\nwant=%s", name, newName)
	}

	// update allowed ciders
	updatedAllowedCIDRs := []string{"10.0.0.0/24"}
	updateRequest2 := &cloudscale.LoadBalancerListenerRequest{
		AllowedCIDRs: updatedAllowedCIDRs,
	}

	err = client.LoadBalancerListeners.Update(context.Background(), uuid, updateRequest2)
	if err != nil {
		t.Fatalf("LoadBalancerListeners.Update returned error %s\n", err)
	}

	updated2, err := client.LoadBalancerListeners.Get(context.Background(), uuid)
	if err != nil {
		t.Fatalf("LoadBalancerListeners.Get returned error %s\n", err)
	}

	if allowedCIDRs := updated2.AllowedCIDRs; !reflect.DeepEqual(allowedCIDRs, updatedAllowedCIDRs) {
		t.Errorf("updated2.AllowedCIDRs \n got=%s\nwant=%s", allowedCIDRs, updatedAllowedCIDRs)
	}

	err = client.LoadBalancerListeners.Delete(context.Background(), updated.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerListeners.Delete returned error %s\n", err)
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

func createPoolOnLB(lb *cloudscale.LoadBalancer) (*cloudscale.LoadBalancerPool, error) {
	createLoadBalancerPoolRequest := &cloudscale.LoadBalancerPoolRequest{
		Name:         testRunPrefix,
		Algorithm:    "round_robin",
		Protocol:     "tcp",
		LoadBalancer: lb.UUID,
	}

	return client.LoadBalancerPools.Create(context.TODO(), createLoadBalancerPoolRequest)
}
