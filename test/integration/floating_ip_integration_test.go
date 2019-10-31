// +build integration

package integration

import (
	"context"
	"reflect"
	"sync"
	"testing"

	"github.com/cloudscale-ch/cloudscale-go-sdk"
)

const pubKey string = "ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFEepRNW5hDct4AdJ8oYsb4lNP5E9XY5fnz3ZvgNCEv7m48+bhUjJXUPuamWix3zigp2lgJHC6SChI/okJ41GUY="

func TestIntegrationFloatingIP_CRUD(t *testing.T) {
	integrationTest(t)

	createServerRequest := &cloudscale.ServerRequest{
		Name:         serverBaseName,
		Flavor:       "flex-2",
		Image:        DefaultImageSlug,
		VolumeSizeGB: 10,
		SSHKeys: []string{
			pubKey,
		},
	}

	server, err := client.Servers.Create(context.Background(), createServerRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}

	waitUntil("running", server.UUID, t)

	createFloatingIPRequest := &cloudscale.FloatingIPCreateRequest{
		IPVersion: 4,
		Server:    server.UUID,
	}

	expectedIP, err := client.FloatingIPs.Create(context.TODO(), createFloatingIPRequest)
	if err != nil {
		t.Fatalf("floatingIP.Create returned error %s\n", err)
	}

	if uuid := expectedIP.Server.UUID; uuid != server.UUID {
		t.Errorf("expectedIP.Server.UUID \n got=%s\nwant=%s", uuid, server.UUID)
	}

	ip := expectedIP.IP()
	floatingIP, err := client.FloatingIPs.Get(context.Background(), ip)
	if err != nil {
		t.Fatalf("Servers.Get returned error %s\n", err)
	}

	if !reflect.DeepEqual(floatingIP, expectedIP) {
		t.Errorf("Error = %#v, expected %#v", floatingIP, expectedIP)
	}

	floatingIps, err := client.FloatingIPs.List(context.Background())
	if err != nil {
		t.Fatalf("FloatingIPs.List returned error %s\n", err)
	}

	if numFloatingIps := len(floatingIps); numFloatingIps < 0 {
		t.Errorf("FloatingIPs.List \n got=%d\nwant=%d", numFloatingIps, 1)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}

	err = client.FloatingIPs.Delete(context.Background(), ip)
	if err != nil {
		t.Fatalf("FloatingIPs.Delete returned error %s\n", err)
	}
}

func TestIntegrationFloatingIP_Update(t *testing.T) {
	createServerRequest := &cloudscale.ServerRequest{
		Name:         serverBaseName,
		Flavor:       "flex-2",
		Image:        DefaultImageSlug,
		VolumeSizeGB: 10,
		SSHKeys: []string{
			pubKey,
		},
	}

	server, err := client.Servers.Create(context.Background(), createServerRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}

	createServerRequest2 := &cloudscale.ServerRequest{
		Name:         serverBaseName + "-floating",
		Flavor:       "flex-2",
		Image:        DefaultImageSlug,
		VolumeSizeGB: 10,
		SSHKeys: []string{
			pubKey,
		},
	}

	expected, err := client.Servers.Create(context.Background(), createServerRequest2)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}

	waitUntil("running", server.UUID, t)
	waitUntil("running", expected.UUID, t)

	createFloatingIPRequest := &cloudscale.FloatingIPCreateRequest{
		IPVersion: 4,
		Server:    server.UUID,
	}

	expectedIP, err := client.FloatingIPs.Create(context.TODO(), createFloatingIPRequest)
	if err != nil {
		t.Fatalf("floatingIP.Create returned error %s\n", err)
	}

	updateRequest := &cloudscale.FloatingIPUpdateRequest{
		Server: expected.UUID,
	}

	ip := expectedIP.IP()
	err = client.FloatingIPs.Update(context.Background(), ip, updateRequest)
	if err != nil {
		t.Fatalf("floatingIP.Update returned error %s\n", err)
	}

	floatingIP, err := client.FloatingIPs.Get(context.Background(), ip)
	if err != nil {
		t.Fatalf("floatingIP.Get returned error %s\n", err)
	}

	if uuid := floatingIP.Server.UUID; uuid != expected.UUID {
		t.Errorf("Server UUID \n got=%s\nwant=%s", uuid, expected.UUID)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}
	err = client.Servers.Delete(context.Background(), expected.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}

	err = client.FloatingIPs.Delete(context.Background(), ip)
	if err != nil {
		t.Fatalf("FloatingIPs.Delete returned error %s\n", err)
	}
}

func TestIntegrationFloatingIP_MultiSite(t *testing.T) {
	integrationTest(t)

	allRegions := getAllRegions(t)
	if len(allRegions) <= 1 {
		t.Skip("Skipping MultiSite test.")
	}

	var wg sync.WaitGroup

	for _, region := range allRegions {
		wg.Add(1)
		go createFloatingIPInRegionAndAssert(t, region, &wg)
	}

	wg.Wait()
}

func createFloatingIPInRegionAndAssert(t *testing.T, region cloudscale.Region, wg *sync.WaitGroup) {
	defer wg.Done()

	createServerRequest := &cloudscale.ServerRequest{
		Name:         serverBaseName,
		Flavor:       "flex-2",
		Image:        DefaultImageSlug,
		VolumeSizeGB: 10,
		SSHKeys: []string{
			pubKey,
		},
	}
	createServerRequest.Zone = region.Zones[0].Slug

	server, err := client.Servers.Create(context.Background(), createServerRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}

	waitUntil("running", server.UUID, t)

	createFloatingIPRequest := &cloudscale.FloatingIPCreateRequest{
		IPVersion: 6,
		Server:    server.UUID,
	}

	createFloatingIPRequest.Region = region.Slug

	floatingIP, err := client.FloatingIPs.Create(context.TODO(), createFloatingIPRequest)
	if err != nil {
		t.Fatalf("FloatingIPs.Create returned error %s\n", err)
	}

	if floatingIP.Region.Slug != region.Slug {
		t.Errorf("FloatingIP in wrong Region\n got=%#v\nwant=%#v", floatingIP.Region, region)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}

	err = client.FloatingIPs.Delete(context.Background(), floatingIP.IP())
	if err != nil {
		t.Errorf("FloatingIPs.Delete returned error %s\n", err)
	}
}

func TestIntegrationFloatingIP_PrefixLength(t *testing.T) {
	createServerRequest := &cloudscale.ServerRequest{
		Name:         serverBaseName,
		Flavor:       "flex-2",
		Image:        DefaultImageSlug,
		VolumeSizeGB: 10,
		SSHKeys: []string{
			pubKey,
		},
	}

	server, err := client.Servers.Create(context.Background(), createServerRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}

	waitUntil("running", server.UUID, t)

	createFloatingIPRequest := &cloudscale.FloatingIPCreateRequest{
		IPVersion:    6,
		PrefixLength: 56,
		Server:       server.UUID,
	}

	expectedIP, err := client.FloatingIPs.Create(context.TODO(), createFloatingIPRequest)
	if err != nil {
		t.Fatalf("floatingIP.Create returned error %s\n", err)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}

	err = client.FloatingIPs.Delete(context.Background(), expectedIP.IP())
	if err != nil {
		t.Fatalf("FloatingIPs.Delete returned error %s\n", err)
	}

}
