//go:build integration
// +build integration

package integration

import (
	"context"
	"errors"
	"github.com/cenkalti/backoff"
	"github.com/cloudscale-ch/cloudscale-go-sdk/v2"
	"reflect"
	"testing"
	"time"
)

func TestIntegrationLoadBalancer_CRUD(t *testing.T) {
	integrationTest(t)

	createLoadBalancerRequest := &cloudscale.LoadBalancerRequest{
		Name:   testRunPrefix,
		Flavor: "lb-flex-4-2",
	}
	createLoadBalancerRequest.Zone = "rma1"

	expected, err := client.LoadBalancers.Create(context.TODO(), createLoadBalancerRequest)
	if err != nil {
		t.Fatalf("LoadBalancers.Create returned error %s\n", err)
	}

	loadBalancer, err := client.LoadBalancers.Get(context.Background(), expected.UUID)
	if err != nil {
		t.Fatalf("LoadBalancers.Get returned error %s\n", err)
	}

	waitUntilLB("running", expected.UUID, t)

	if h := time.Since(loadBalancer.CreatedAt).Hours(); !(-1 < h && h < 1) {
		t.Errorf("loadBalancer.CreatedAt ourside of expected range. got=%v", loadBalancer.CreatedAt)
	}

	if !reflect.DeepEqual(loadBalancer, expected) {
		t.Errorf("Error = %#v, expected %#v", loadBalancer, expected)
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

func waitUntilLB(status string, uuid string, t *testing.T) *cloudscale.LoadBalancer {
	// An operation that may fail.
	loadBalancer := new(cloudscale.LoadBalancer)
	operation := func() error {
		lb, err := client.LoadBalancers.Get(context.Background(), uuid)
		if err != nil {
			return err
		}

		if lb.Status != status {
			return errors.New("Status not reached")
		}
		loadBalancer = lb
		return nil
	}

	err := backoff.Retry(operation, backoff.NewExponentialBackOff())
	if err != nil {
		t.Fatalf("Error while waiting for status change %s\n", err)
	}
	return loadBalancer
}
