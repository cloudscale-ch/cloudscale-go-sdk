// +build integration

package integration

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/cloudscale-ch/cloudscale"
)

func TestIntegrationFloatingIP_CRUD(t *testing.T) {
	integrationTest(t)

	createServerRequest := &cloudscale.ServerRequest{
		Name:         "db-master",
		Flavor:       "flex-2",
		Image:        "debian-8",
		VolumeSizeGB: 10,
		SSHKeys: []string{
			"ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFEepRNW5hDct4AdJ8oYsb4lNP5E9XY5fnz3ZvgNCEv7m48+bhUjJXUPuamWix3zigp2lgJHC6SChI/okJ41GUY=",
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

	ip := strings.Split(expectedIP.Network, "/")[0]
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
