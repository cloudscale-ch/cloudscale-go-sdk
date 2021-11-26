// +build integration

package integration

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/cloudscale-ch/cloudscale-go-sdk"
)

const DefaultImageSlug = "debian-9"

func integrationTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test")
	}
}

func createServer(t *testing.T, createRequest *cloudscale.ServerRequest) (*cloudscale.Server, error) {
	server, err := client.Servers.Create(context.Background(), createRequest)
	if err == nil {
		waitUntil(cloudscale.ServerRunning, server.UUID, t)
	}
	return server, err
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
	s := waitUntil(cloudscale.ServerStopped, server.UUID, t)
	if status := s.Status; status != cloudscale.ServerStopped {
		t.Errorf("Server.Update got=%s\nwant=%s\n", status, cloudscale.ServerStopped)
	}

	// Start a server
	req.Status = cloudscale.ServerRunning
	err = client.Servers.Update(context.Background(), server.UUID, req)
	if err != nil {
		t.Fatalf("Servers.Update returned error %s\n", err)
	}
	s = waitUntil(cloudscale.ServerRunning, server.UUID, t)
	if status := s.Status; status != cloudscale.ServerRunning {
		t.Errorf("Server.Update got=%s\nwant=%s\n", status, cloudscale.ServerRunning)
	}

	// Reboot a server
	req.Status = cloudscale.ServerRebooted
	err = client.Servers.Update(context.Background(), server.UUID, req)
	if err != nil {
		t.Fatalf("Servers.Update returned error %s\n", err)
	}
	s = waitUntil(cloudscale.ServerRunning, server.UUID, t)
	if status := s.Status; status != cloudscale.ServerRunning {
		t.Errorf("Server.Update got=%s\nwant=%s\n", status, cloudscale.ServerRunning)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}
}

func getDefaultServerRequest() cloudscale.ServerRequest {
	return cloudscale.ServerRequest{
		Name:         testRunPrefix,
		Flavor:       "flex-2",
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
	waitUntil("stopped", server.UUID, t)

	multiUpdateRequest := &cloudscale.ServerUpdateRequest{
		Flavor: "flex-4",
		Name:   "bar",
	}
	err = client.Servers.Update(context.TODO(), server.UUID, multiUpdateRequest)
	// This shouldn't work.
	if err == nil {
		t.Error("Expected an error when updating multiple volume attributes\n")
	} else {
		expected := "To keep changes atomic"
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

	const scaleFlavor = "flex-4"
	// Try to scale.
	scaleRequest := &cloudscale.ServerUpdateRequest{Flavor: scaleFlavor}
	err = client.Servers.Update(context.TODO(), server.UUID, scaleRequest)
	if err != nil {
		t.Errorf("Servers.Update failed %s\n", err)
	}

	getServer, err := client.Servers.Get(context.TODO(), server.UUID)
	if err == nil {
		if getServer.Flavor.Slug != scaleFlavor {
			t.Errorf("Scaling failed, could not scale, is at %s\n", getServer.Flavor.Slug)
		}
	} else {
		t.Errorf("Servers.Get returned error %s\n", err)
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
	s := waitUntil("stopped", server.UUID, t)
	if status := s.Status; status != cloudscale.ServerStopped {
		t.Errorf("Server.Stop got=%s\nwant=%s\n", status, cloudscale.ServerStopped)
	}

	// Start a server
	err = client.Servers.Start(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Start returned error %s\n", err)
	}
	s = waitUntil("running", server.UUID, t)
	if status := s.Status; status != cloudscale.ServerRunning {
		t.Errorf("Server.Start got=%s\nwant=%s\n", status, cloudscale.ServerRunning)
	}

	// reboot server
	err = client.Servers.Reboot(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Reboot returned error %s\n", err)
	}
	s = waitUntil("running", server.UUID, t)
	if status := s.Status; status != cloudscale.ServerRunning {
		t.Errorf("Server.Reboot got=%s\nwant=%s\n", status, cloudscale.ServerRunning)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}
}

func TestIntegrationServer_MultipleVolumes(t *testing.T) {
	integrationTest(t)

	request := getDefaultServerRequest()
	request.Volumes = &([]cloudscale.Volume{
		{SizeGB: 3, Type: "ssd"},
		{SizeGB: 100, Type: "bulk"},
	})

	server, err := createServer(t, &request)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}

	expected := []cloudscale.VolumeStub{
		{Type: "ssd", DevicePath: "", SizeGB: 10, UUID: ""},
		{Type: "ssd", DevicePath: "", SizeGB: 3, UUID: ""},
		{Type: "bulk", DevicePath: "", SizeGB: 100, UUID: ""},
	}
	if !reflect.DeepEqual(server.Volumes, expected) {
		t.Errorf("Volumes response\n got=%#v\nwant=%#v", server.Volumes, expected)
	}

	// Wait a bit until the volumes are actually allocated and have UUID's see
	// https://www.cloudscale.ch/en/api/v1#volumes-create
	time.Sleep(5 * time.Second)
	server, err = client.Servers.Get(context.TODO(), server.UUID)
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

func waitUntil(status string, uuid string, t *testing.T) *cloudscale.Server {
	// An operation that may fail.
	server := new(cloudscale.Server)
	operation := func() error {
		s, err := client.Servers.Get(context.Background(), uuid)
		if err != nil {
			return err
		}

		if s.Status != status {
			return errors.New("Status not reached")
		}
		server = s
		return nil
	}

	err := backoff.Retry(operation, backoff.NewExponentialBackOff())
	if err != nil {
		t.Fatalf("Error while waiting for status change %s\n", err)
	}
	return server
}

func waitUntilListed(t *testing.T, server *cloudscale.Server) {
	operation := func() error {
		servers, err := client.Servers.List(context.Background())
		if err != nil {
			return err
		}
		for _, s := range servers {
			if s.UUID == server.UUID {
				return nil
			}
		}
		return errors.New("not contained in server listing")
	}

	err := backoff.Retry(operation, backoff.NewExponentialBackOff())
	if err != nil {
		t.Fatalf("Error while waiting for server to be listed %s\n", err)
	}
}
