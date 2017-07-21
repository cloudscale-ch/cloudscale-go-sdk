// +build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/cloudscale-ch/cloudscale"
)

func integrationTest(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping acceptance test")
	}
}

func TestIntegrationServer_CRUD(t *testing.T) {
	integrationTest(t)

	createRequest := &cloudscale.ServerRequest{
		Name:         "db-master",
		Flavor:       "flex-2",
		Image:        "debian-8",
		VolumeSizeGB: 10,
		SSHKeys: []string{
			"ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFEepRNW5hDct4AdJ8oYsb4lNP5E9XY5fnz3ZvgNCEv7m48+bhUjJXUPuamWix3zigp2lgJHC6SChI/okJ41GUY=",
		},
	}

	expected, err := client.Servers.Create(context.Background(), createRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}
	time.Sleep(20 * time.Second)

	server, err := client.Servers.Get(context.Background(), expected.UUID)
	if err != nil {
		t.Fatalf("Servers.Get returned error %s\n", err)
	}

	if uuid := server.UUID; uuid != expected.UUID {
		t.Errorf("Server.UUID \n got=%s\nwant=%s", uuid, expected.UUID)
	}

	servers, err := client.Servers.List(context.Background())
	if err != nil {
		t.Fatalf("Servers.List returned error %s\n", err)
	}

	if numServers := len(servers); numServers < 0 {
		t.Errorf("Server.List \n got=%d\nwant=%d", numServers, 1)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Get returned error %s\n", err)
	}

}

func TestIntegrationServer_Actions(t *testing.T) {
	integrationTest(t)

	createRequest := &cloudscale.ServerRequest{
		Name:         "db-master",
		Flavor:       "flex-2",
		Image:        "debian-8",
		VolumeSizeGB: 10,
		SSHKeys: []string{
			"ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFEepRNW5hDct4AdJ8oYsb4lNP5E9XY5fnz3ZvgNCEv7m48+bhUjJXUPuamWix3zigp2lgJHC6SChI/okJ41GUY=",
		},
	}

	server, err := client.Servers.Create(context.Background(), createRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}
	time.Sleep(20 * time.Second)

	err = client.Servers.Stop(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Stop returned error %s\n", err)
	}
	time.Sleep(30 * time.Second)

	err = client.Servers.Start(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Start returned error %s\n", err)
	}
	time.Sleep(20 * time.Second)

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Get returned error %s\n", err)
	}
}
