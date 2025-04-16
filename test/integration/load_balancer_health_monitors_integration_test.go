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

func TestIntegrationLoadBalancerHealthMonitor_CRUD(t *testing.T) {
	integrationTest(t)

	lb, err := createLoadBalancer()
	if err != nil {
		t.Fatalf("LoadBalancers.Create returned error %s\n", err)
	}

	waitUntilLB(lb.UUID, t)

	pool, err := createPoolOnLB(lb)
	if err != nil {
		t.Fatalf("LoadBalancerPools.Create returned error %s\n", err)
	}

	createLoadBalancerHealthMonitorRequest := &cloudscale.LoadBalancerHealthMonitorRequest{
		Pool:     pool.UUID,
		DelayS:   20,
		TimeoutS: 15,
		Type:     "tcp",
	}

	expected, err := client.LoadBalancerHealthMonitors.Create(context.Background(), createLoadBalancerHealthMonitorRequest)
	if err != nil {
		t.Fatalf("LoadBalancerHealthMonitors.Create returned error %s\n", err)
	}

	loadBalancerHealthMonitor, err := client.LoadBalancerHealthMonitors.Get(context.Background(), expected.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerHealthMonitors.Get returned error %s\n", err)
	}

	if h := time.Since(loadBalancerHealthMonitor.CreatedAt).Hours(); !(-1 < h && h < 1) {
		t.Errorf("loadBalancerHealthMonitor.CreatedAt ourside of expected range. got=%v", loadBalancerHealthMonitor.CreatedAt)
	}

	if !reflect.DeepEqual(loadBalancerHealthMonitor, expected) {
		t.Errorf("Error = %#v, expected %#v", loadBalancerHealthMonitor, expected)
	}

	if poolLbUUID := loadBalancerHealthMonitor.Pool.UUID; poolLbUUID != pool.UUID {
		t.Errorf("poolLbUUID \n got=%#v\nwant=%#v", poolLbUUID, pool.UUID)
	}

	if lbUUID := loadBalancerHealthMonitor.LoadBalancer.UUID; lbUUID != lb.UUID {
		t.Errorf("lbUUID \n got=%#v\nwant=%#v", lbUUID, lb.UUID)
	}

	loadBalancerHealthMonitors, err := client.LoadBalancerHealthMonitors.List(context.Background())
	if err != nil {
		t.Fatalf("LoadBalancerHealthMonitors.List returned error %s\n", err)
	}

	if numLoadbalancerHealthMonitors := len(loadBalancerHealthMonitors); numLoadbalancerHealthMonitors < 1 {
		t.Errorf("LoadBalancerHealthMonitors.List \n got=%d\nwant=%d", numLoadbalancerHealthMonitors, 1)
	}

	err = client.LoadBalancerHealthMonitors.Delete(context.Background(), expected.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerHealthMonitors.Delete returned error %s\n", err)
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

func TestIntegrationLoadBalancerHealthMonitor_Update(t *testing.T) {
	integrationTest(t)

	lb, err := createLoadBalancer()
	if err != nil {
		t.Fatalf("LoadBalancers.Create returned error %s\n", err)
	}

	waitUntilLB(lb.UUID, t)

	pool, err := createPoolOnLB(lb)
	if err != nil {
		t.Fatalf("LoadBalancerPools.Create returned error %s\n", err)
	}

	createLoadBalancerHealthMonitorRequest := &cloudscale.LoadBalancerHealthMonitorRequest{
		Pool:        pool.UUID,
		DelayS:      20,
		TimeoutS:    15,
		Type:        "tcp",
		UpThreshold: 10,
	}

	healthMonitor, err := client.LoadBalancerHealthMonitors.Create(context.Background(), createLoadBalancerHealthMonitorRequest)
	if err != nil {
		t.Fatalf("LoadBalancerHealthMonitors.Create returned error %s\n", err)
	}

	newMaxRetries := 5
	updateRequest := &cloudscale.LoadBalancerHealthMonitorRequest{
		UpThreshold: newMaxRetries,
	}

	uuid := healthMonitor.UUID
	err = client.LoadBalancerHealthMonitors.Update(context.Background(), uuid, updateRequest)
	if err != nil {
		t.Fatalf("LoadBalancerHealthMonitors.Update returned error %s\n", err)
	}

	updated, err := client.LoadBalancerHealthMonitors.Get(context.Background(), uuid)
	if err != nil {
		t.Fatalf("LoadBalancerHealthMonitors.Get returned error %s\n", err)
	}

	if name := updated.UpThreshold; name != newMaxRetries {
		t.Errorf("updated.Name \n got=%#v\nwant=%#v", name, newMaxRetries)
	}

	err = client.LoadBalancerHealthMonitors.Delete(context.Background(), updated.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerHealthMonitors.Delete returned error %s\n", err)
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

func TestIntegrationLoadBalancerHealthMonitor_HTTP_Update(t *testing.T) {
	integrationTest(t)

	lb, err := createLoadBalancer()
	if err != nil {
		t.Fatalf("LoadBalancers.Create returned error %s\n", err)
	}

	waitUntilLB(lb.UUID, t)

	pool, err := createPoolOnLB(lb)
	if err != nil {
		t.Fatalf("LoadBalancerPools.Create returned error %s\n", err)
	}

	hostName := "hostName"
	createLoadBalancerHealthMonitorRequest := &cloudscale.LoadBalancerHealthMonitorRequest{
		Pool:        pool.UUID,
		DelayS:      20,
		TimeoutS:    15,
		Type:        "http",
		UpThreshold: 10,
		HTTP: &cloudscale.LoadBalancerHealthMonitorHTTPRequest{
			Host: &hostName,
		},
	}

	healthMonitor, err := client.LoadBalancerHealthMonitors.Create(context.Background(), createLoadBalancerHealthMonitorRequest)
	if err != nil {
		t.Fatalf("LoadBalancerHealthMonitors.Create returned error %s\n", err)
	}

	http := healthMonitor.HTTP
	expectedHTTP := cloudscale.LoadBalancerHealthMonitorHTTP{
		ExpectedCodes: []string{"200"},
		Method:        "GET",
		UrlPath:       "/",
		Version:       "1.1",
		Host:          &hostName,
	}

	if !reflect.DeepEqual(http, &expectedHTTP) {
		t.Errorf("healthMonitor.HTTP \n got=%#v\nwant=%#v", http, &expectedHTTP)
	}

	httpRequest := cloudscale.LoadBalancerHealthMonitorHTTPRequest{
		ExpectedCodes: []string{"201", "200"},
	}
	updateRequest := &cloudscale.LoadBalancerHealthMonitorRequest{
		HTTP: &httpRequest,
	}

	uuid := healthMonitor.UUID
	err = client.LoadBalancerHealthMonitors.Update(context.Background(), uuid, updateRequest)
	if err != nil {
		t.Fatalf("LoadBalancerHealthMonitors.Update returned error %s\n", err)
	}

	updated, err := client.LoadBalancerHealthMonitors.Get(context.Background(), uuid)
	if err != nil {
		t.Fatalf("LoadBalancerHealthMonitors.Get returned error %s\n", err)
	}

	expectedUpdatedHTTP := cloudscale.LoadBalancerHealthMonitorHTTP{
		ExpectedCodes: []string{"201", "200"},
		Method:        "GET",
		UrlPath:       "/",
		Version:       "1.1",
		Host:          &hostName,
	}
	updatedHttp := updated.HTTP
	if !reflect.DeepEqual(updatedHttp, &expectedUpdatedHTTP) {
		t.Errorf("updated.HTTP \n got=%#v\nwant=%#v", updatedHttp, &expectedUpdatedHTTP)
	}

	err = client.LoadBalancerHealthMonitors.Delete(context.Background(), updated.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerHealthMonitors.Delete returned error %s\n", err)
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
