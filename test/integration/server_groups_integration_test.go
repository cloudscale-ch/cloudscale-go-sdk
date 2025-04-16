//go:build integration
// +build integration

package integration

import (
	"context"
	"sync"
	"testing"

	"github.com/cloudscale-ch/cloudscale-go-sdk/v6"
)

func TestIntegrationServerGroup_CRUD(t *testing.T) {
	integrationTest(t)

	expected, err := createServerGroup(t)
	if err != nil {
		t.Fatalf("ServerGroups.Create returned error %s\n", err)
	}

	serverGroup, err := client.ServerGroups.Get(context.Background(), expected.UUID)
	if err != nil {
		t.Fatalf("ServerGroups.Get returned error %s\n", err)
	}

	if uuid := serverGroup.UUID; uuid != expected.UUID {
		t.Errorf("ServerGroup.UUID got=%s\nwant=%s", uuid, expected.UUID)
	}

	err = client.ServerGroups.Delete(context.Background(), serverGroup.UUID)
	if err != nil {
		t.Fatalf("ServerGroups.Get returned error %s\n", err)
	}

}

func createServerGroup(t *testing.T) (*cloudscale.ServerGroup, error) {
	createRequest := &cloudscale.ServerGroupRequest{
		Name: testRunPrefix + "-group",
		Type: "anti-affinity",
	}

	return client.ServerGroups.Create(context.Background(), createRequest)
}

func TestIntegrationServerGroup_MultiSite(t *testing.T) {
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
		go createServerGroupInZoneAndAssert(t, zone, &wg)
	}

	wg.Wait()
}

func createServerGroupInZoneAndAssert(t *testing.T, zone cloudscale.Zone, wg *sync.WaitGroup) {
	defer wg.Done()

	createServerGroupRequest := &cloudscale.ServerGroupRequest{
		Name: "Yellow Submarine",
		Type: "anti-affinity",
	}

	createServerGroupRequest.Zone = zone.Slug

	serverGroup, err := client.ServerGroups.Create(context.TODO(), createServerGroupRequest)
	if err != nil {
		t.Fatalf("ServerGroups.Create returned error %s\n", err)
	}

	if serverGroup.Zone != zone {
		t.Errorf("ServerGroup in wrong Zone\n got=%#v\nwant=%#v", serverGroup.Zone, zone)
	}

	err = client.ServerGroups.Delete(context.Background(), serverGroup.UUID)
	if err != nil {
		t.Errorf("ServerGroups.Delete returned error %s\n", err)
	}
}
