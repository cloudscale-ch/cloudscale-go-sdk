//go:build integration
// +build integration

package integration

import (
	"context"
	"github.com/cloudscale-ch/cloudscale-go-sdk/v6"
	"reflect"
	"testing"
	"time"
)

func TestIntegrationLoadBalancerPool_CRUD(t *testing.T) {
	integrationTest(t)

	lb, err := createLoadBalancer()
	if err != nil {
		t.Fatalf("LoadBalancers.Create returned error %s\n", err)
	}

	waitUntilLB(lb.UUID, t)

	createLoadBalancerPoolRequest := &cloudscale.LoadBalancerPoolRequest{
		Name:         testRunPrefix,
		Algorithm:    "round_robin",
		Protocol:     "tcp",
		LoadBalancer: lb.UUID,
	}

	expected, err := client.LoadBalancerPools.Create(context.TODO(), createLoadBalancerPoolRequest)
	if err != nil {
		t.Fatalf("LoadBalancerPools.Create returned error %s\n", err)
	}

	loadBalancerPool, err := client.LoadBalancerPools.Get(context.Background(), expected.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerPools.Get returned error %s\n", err)
	}

	if h := time.Since(loadBalancerPool.CreatedAt).Hours(); !(-1 < h && h < 1) {
		t.Errorf("loadBalancerPool.CreatedAt ourside of expected range. got=%v", loadBalancerPool.CreatedAt)
	}

	if !reflect.DeepEqual(loadBalancerPool, expected) {
		t.Errorf("Error = %#v, expected %#v", loadBalancerPool, expected)
	}

	if poolLbUUID := loadBalancerPool.LoadBalancer.UUID; poolLbUUID != lb.UUID {
		t.Errorf("poolLbUUID \n got=%#v\nwant=%#v", poolLbUUID, lb.UUID)
	}

	loadBalancerPools, err := client.LoadBalancerPools.List(context.Background())
	if err != nil {
		t.Fatalf("LoadBalancerPools.List returned error %s\n", err)
	}

	if numLoadBalancerPools := len(loadBalancerPools); numLoadBalancerPools < 1 {
		t.Errorf("LoadBalancerPools.List \n got=%d\nwant=%d", numLoadBalancerPools, 1)
	}

	err = client.LoadBalancerPools.Delete(context.Background(), loadBalancerPool.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerPools.Delete returned error %s\n", err)
	}

	err = client.LoadBalancers.Delete(context.Background(), lb.UUID)
	if err != nil {
		t.Fatalf("LoadBalancers.Delete returned error %s\n", err)
	}
}

func TestIntegrationLoadBalancerPool_Update(t *testing.T) {
	integrationTest(t)

	lb, err := createLoadBalancer()
	if err != nil {
		t.Fatalf("LoadBalancers.Create returned error %s\n", err)
	}

	waitUntilLB(lb.UUID, t)

	createLoadBalancerPoolRequest := &cloudscale.LoadBalancerPoolRequest{
		Name:         testRunPrefix,
		Algorithm:    "round_robin",
		Protocol:     "tcp",
		LoadBalancer: lb.UUID,
	}

	pool, err := client.LoadBalancerPools.Create(context.TODO(), createLoadBalancerPoolRequest)
	if err != nil {
		t.Fatalf("LoadBalancerPools.Create returned error %s\n", err)
	}

	newName := testRunPrefix + "-renamed"
	updateRequest := &cloudscale.LoadBalancerPoolRequest{
		Name: newName,
	}

	uuid := pool.UUID
	err = client.LoadBalancerPools.Update(context.Background(), uuid, updateRequest)
	if err != nil {
		t.Fatalf("LoadBalancerPools.Update returned error %s\n", err)
	}

	updated, err := client.LoadBalancerPools.Get(context.Background(), uuid)
	if err != nil {
		t.Fatalf("LoadBalancerPools.Get returned error %s\n", err)
	}

	if name := updated.Name; name != newName {
		t.Errorf("loadbalancer.Name \n got=%s\nwant=%s", name, newName)
	}

	err = client.LoadBalancerPools.Delete(context.Background(), updated.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerPools.Delete returned error %s\n", err)
	}

	err = client.LoadBalancers.Delete(context.Background(), lb.UUID)
	if err != nil {
		t.Fatalf("LoadBalancers.Delete returned error %s\n", err)
	}
}

func createLoadBalancer() (*cloudscale.LoadBalancer, error) {
	createRequest := &cloudscale.LoadBalancerRequest{
		Name:   testRunPrefix,
		Flavor: "lb-standard",
	}
	createRequest.Zone = testZone

	return client.LoadBalancers.Create(context.Background(), createRequest)
}
