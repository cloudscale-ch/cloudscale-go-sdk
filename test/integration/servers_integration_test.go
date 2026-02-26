//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cloudscale-ch/cloudscale-go-sdk/v7"
)

const DefaultImageSlug = "debian-11"

func integrationTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test")
	}
}

func createServer(t *testing.T, createRequest *cloudscale.ServerRequest) (*cloudscale.Server, error) {
	server, err := client.Servers.Create(context.Background(), createRequest)
	if err != nil {
		return nil, err
	}

	server, err = client.Servers.WaitFor(
		context.Background(),
		server.UUID,
		cloudscale.ServerIsRunning,
	)
	if err != nil {
		t.Fatalf("Servers.WaitFor returned error %s\n", err)
	}

	return server, nil
}

func TestIntegrationServer_CRUD(t *testing.T) {
	integrationTest(t)

	serverRequest := getDefaultServerRequest()
	serverRequest.SSHKeys = []string{}
	serverRequest.Password = randomNotVerySecurePassword(10)

	expected, err := createServer(t, &serverRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}

	server, err := client.Servers.Get(context.Background(), expected.UUID)
	if err != nil {
		t.Fatalf("Servers.Get returned error %s\n", err)
	}

	if uuid := server.UUID; uuid != expected.UUID {
		t.Errorf("Server.UUID got=%s\nwant=%s", uuid, expected.UUID)
	}

	if h := time.Since(server.CreatedAt).Hours(); !(-1 < h && h < 1) {
		t.Errorf("server.CreatedAt ourside of expected range. got=%v", server.CreatedAt)
	}

	if server.Image.Slug != DefaultImageSlug {
		t.Errorf("Server.Image.Slug got=%s, want=%s", server.Image.Slug, DefaultImageSlug)
	}

	const expectedValue = "debian"
	if !strings.Contains(strings.ToLower(server.Image.Name), expectedValue) {
		t.Errorf("Server.Image.Name got=%s, want to contain '%s'", server.Image.Name, expectedValue)
	}
	if !strings.Contains(strings.ToLower(server.Image.OperatingSystem), expectedValue) {
		t.Errorf("Server.Image.OperatingSystem got=%s, want to contain '%s'", server.Image.OperatingSystem, expectedValue)
	}
	if !strings.Contains(strings.ToLower(server.Image.DefaultUsername), expectedValue) {
		t.Errorf("Server.Image.DefaultUsername got=%s, want to contain '%s'", server.Image.DefaultUsername, expectedValue)
	}

	servers, err := client.Servers.List(context.Background())
	if err != nil {
		t.Fatalf("Servers.List returned error %s\n", err)
	}

	if numServers := len(servers); numServers < 1 {
		t.Errorf("Server.List got=%d\nwant=%d\n", numServers, 1)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}

}

func TestIntegrationServer_UpdateStatus(t *testing.T) {
	integrationTest(t)

	request := getDefaultServerRequest()
	server, err := createServer(t, &request)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}

	// Stop a server
	req := &cloudscale.ServerUpdateRequest{
		Status: cloudscale.ServerStopped,
	}
	err = client.Servers.Update(context.Background(), server.UUID, req)
	if err != nil {
		t.Fatalf("Servers.Update returned error %s\n", err)
	}
	_, err = client.Servers.WaitFor(
		context.Background(),
		server.UUID,
		cloudscale.ServerIsStopped,
	)
	if err != nil {
		t.Fatalf("Servers.WaitFor returned error %s\n", err)
	}

	// Start a server
	req.Status = cloudscale.ServerRunning
	err = client.Servers.Update(context.Background(), server.UUID, req)
	if err != nil {
		t.Fatalf("Servers.Update returned error %s\n", err)
	}
	_, err = client.Servers.WaitFor(
		context.Background(),
		server.UUID,
		cloudscale.ServerIsRunning,
	)
	if err != nil {
		t.Fatalf("Servers.WaitFor returned error %s\n", err)
	}

	// Reboot a server
	req.Status = cloudscale.ServerRebooted
	err = client.Servers.Update(context.Background(), server.UUID, req)
	if err != nil {
		t.Fatalf("Servers.Update returned error %s\n", err)
	}
	_, err = client.Servers.WaitFor(
		context.Background(),
		server.UUID,
		cloudscale.ServerIsRunning,
	)
	if err != nil {
		t.Fatalf("Servers.WaitFor returned error %s\n", err)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}
}

func getDefaultServerRequest() cloudscale.ServerRequest {
	return cloudscale.ServerRequest{
		Name:         testRunPrefix,
		Zone:         testZone,
		Flavor:       "flex-4-2",
		Image:        DefaultImageSlug,
		VolumeSizeGB: 10,
		SSHKeys: []string{
			"ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFEepRNW5hDct4AdJ8oYsb4lNP5E9XY5fnz3ZvgNCEv7m48+bhUjJXUPuamWix3zigp2lgJHC6SChI/okJ41GUY=",
		},
	}
}

func TestIntegrationServer_UpdateRest(t *testing.T) {
	integrationTest(t)

	serverRequest := getDefaultServerRequest()
	server, err := createServer(t, &serverRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}
	// We need to stop the server in order to scale
	err = client.Servers.Stop(context.Background(), server.UUID)
	if err != nil {
		t.Errorf("Servers.Stop returned error %s\n", err)
	}
	_, err = client.Servers.WaitFor(
		context.Background(),
		server.UUID,
		cloudscale.ServerIsStopped,
	)
	if err != nil {
		t.Fatalf("Servers.WaitFor returned error %s\n", err)
	}

	multiUpdateRequest := &cloudscale.ServerUpdateRequest{
		Flavor: "flex-4-4",
		Name:   "bar",
	}
	err = client.Servers.Update(context.TODO(), server.UUID, multiUpdateRequest)
	// This shouldn't work.
	if err == nil {
		t.Error("Expected an error when updating multiple volume attributes\n")
	} else {
		expected := "Only one attribute"
		err, ok := err.(*cloudscale.ErrorResponse)
		if !ok {
			t.Errorf("Couldn't cast %s\n", err)
		}
		if err.StatusCode != 400 {
			t.Errorf("Expected bad request and not %d\n", err.StatusCode)
		}
		if !strings.Contains(err.Error(), expected) {
			t.Errorf("Expected \"%s\" not \"%s\"\n", expected, err.Error())
		}
	}

	const scaleFlavor = "flex-4-4"
	// Try to scale.
	scaleRequest := &cloudscale.ServerUpdateRequest{Flavor: scaleFlavor}
	err = client.Servers.Update(context.TODO(), server.UUID, scaleRequest)
	if err != nil {
		t.Errorf("Servers.Update failed %s\n", err)
	}

	getServer, err := client.Servers.WaitFor(
		context.Background(),
		server.UUID,
		cloudscale.ServerIsStopped,
	)
	if err != nil {
		t.Fatalf("Servers.WaitFor returned error %s\n", err)
	}

	if getServer.Flavor.Slug != scaleFlavor {
		t.Errorf("Scaling failed, could not scale, is at %s\n", getServer.Flavor.Slug)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}
}

func TestIntegrationServer_Actions(t *testing.T) {
	integrationTest(t)

	serverRequest := getDefaultServerRequest()
	server, err := createServer(t, &serverRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}

	// Stop a server
	err = client.Servers.Stop(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Stop returned error %s\n", err)
	}
	_, err = client.Servers.WaitFor(
		context.Background(),
		server.UUID,
		cloudscale.ServerIsStopped,
	)
	if err != nil {
		t.Fatalf("Servers.WaitFor returned error %s\n", err)
	}

	// Start a server
	err = client.Servers.Start(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Start returned error %s\n", err)
	}
	_, err = client.Servers.WaitFor(
		context.Background(),
		server.UUID,
		cloudscale.ServerIsRunning,
	)
	if err != nil {
		t.Fatalf("Servers.WaitFor returned error %s\n", err)
	}

	// reboot server
	err = client.Servers.Reboot(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Reboot returned error %s\n", err)
	}
	_, err = client.Servers.WaitFor(
		context.Background(),
		server.UUID,
		cloudscale.ServerIsRunning,
	)
	if err != nil {
		t.Fatalf("Servers.WaitFor returned error %s\n", err)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}
}

func TestIntegrationServer_MultipleVolumes(t *testing.T) {
	integrationTest(t)

	request := getDefaultServerRequest()
	request.Volumes = &([]cloudscale.ServerVolumeRequest{
		{SizeGB: 3, Type: "ssd"},
		{SizeGB: 100, Type: "bulk"},
	})

	server, err := createServer(t, &request)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}

	// Wait until the volumes are actually allocated and have UUIDs
	// https://www.cloudscale.ch/en/api/v1#volumes-create
	server, err = client.Servers.WaitFor(
		context.TODO(),
		server.UUID,
		func(s *cloudscale.Server) (bool, error) {
			for i, volume := range s.Volumes {
				if len(volume.UUID) <= 1 {
					return false, fmt.Errorf("volume at index %d has an invalid or unassigned UUID", i)
				}
			}
			return true, nil
		},
	)
	if err != nil {
		t.Fatalf("Servers.WaitFor returned error %s\n", err)
	}

	// Ignore UUIDs in this comparison
	actual := make([]cloudscale.VolumeStub, len(server.Volumes))
	copy(actual, server.Volumes)
	for i := range actual {
		actual[i].UUID = ""
	}
	expected := []cloudscale.VolumeStub{
		{Type: "ssd", DevicePath: "", SizeGB: 10, UUID: ""},
		{Type: "ssd", DevicePath: "", SizeGB: 3, UUID: ""},
		{Type: "bulk", DevicePath: "", SizeGB: 100, UUID: ""},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Volumes response\n got=%#v\nwant=%#v", actual, expected)
	}

	// delete all volume, except the root volume
	for _, volume := range server.Volumes[1:] {
		volumeUUID := volume.UUID
		if len(volumeUUID) <= 1 {
			t.Errorf("Volume does not seem to have a valid UUID got=%#v", volumeUUID)
		}

		err = client.Volumes.Delete(context.Background(), volumeUUID)
		if err != nil {
			t.Fatalf("Volumes.Delete returned error %s\n", err)
		}
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Get returned error %s\n", err)
	}
}

func TestIntegrationServer_MultiSite(t *testing.T) {
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
		go createServerInZoneAndAssert(t, zone, &wg)
	}

	wg.Wait()
}

func createServerInZoneAndAssert(t *testing.T, zone cloudscale.Zone, wg *sync.WaitGroup) {
	defer wg.Done()

	createRequest := getDefaultServerRequest()
	createRequest.Zone = zone.Slug
	server, err := createServer(t, &createRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}

	if server.Zone != zone {
		t.Errorf("Server in wrong Zone\n got=%#v\nwant=%#v", server.Zone, zone)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Errorf("Servers.Delete returned error %s\n", err)
	}
}
