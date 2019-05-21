// +build integration

package integration

import (
	"context"
	"errors"
	"strings"
	"testing"

	"cloudscale"
	"github.com/cenkalti/backoff"
)

const serverBaseName = "go-sdk-integration-test"
const DefaultImageSlug = "debian-9"

func integrationTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test")
	}
}

func TestIntegrationServer_CRUD(t *testing.T) {
	integrationTest(t)

	expected, err := createServer(t)
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

	if numServers := len(servers); numServers < 0 {
		t.Errorf("Server.List got=%d\nwant=%d\n", numServers, 1)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Get returned error %s\n", err)
	}

}

func TestIntegrationServer_UpdateStatus(t *testing.T) {
	integrationTest(t)

	server, err := createServer(t)
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

func createServer(t *testing.T) (*cloudscale.Server, error) {
	createRequest := &cloudscale.ServerRequest{
		Name:         serverBaseName,
		Flavor:       "flex-2",
		Image:        DefaultImageSlug,
		VolumeSizeGB: 10,
		SSHKeys: []string{
			"ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFEepRNW5hDct4AdJ8oYsb4lNP5E9XY5fnz3ZvgNCEv7m48+bhUjJXUPuamWix3zigp2lgJHC6SChI/okJ41GUY=",
		},
	}

	server, err := client.Servers.Create(context.Background(), createRequest)
	if err == nil {
		waitUntil(cloudscale.ServerRunning, server.UUID, t)
	}
	return server, err
}

func TestIntegrationServer_UpdateRest(t *testing.T) {
	integrationTest(t)

	server, err := createServer(t)
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

	server, err := createServer(t)
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

func TestIntegrationServer_DeleteRemainingServer(t *testing.T) {
	servers, err := client.Servers.List(context.Background())
	if err != nil {
		t.Fatalf("Servers.List returned error %s\n", err)
	}

	for _, server := range servers {
		if strings.HasPrefix(server.Name, serverBaseName) {
			t.Errorf("Found not deleted server: %s\n", server.Name)
			err = client.Servers.Delete(context.Background(), server.UUID)
			if err != nil {
				t.Errorf("Servers.Delete returned error %s\n", err)
			}
		}
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
