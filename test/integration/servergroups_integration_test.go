// +build integration

package integration

import (
	"context"
	"strings"
	"testing"

	cloudscale "github.com/cloudscale-ch/cloudscale-go-sdk"
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
		Name:         serverBaseName + "-group",
		Type:         "anti-affinity",
	}

	return client.ServerGroups.Create(context.Background(), createRequest)
}

func TestIntegrationServerGroup_DeleteRemaining(t *testing.T) {
	serverGroups, err := client.ServerGroups.List(context.Background())
	if err != nil {
		t.Fatalf("Servers.List returned error %s\n", err)
	}

	for _, serverGroup := range serverGroups {
		if strings.HasPrefix(serverGroup.Name, serverBaseName) {
			t.Errorf("Found not deleted serverGroup: %s\n", serverGroup.Name)
			err = client.ServerGroups.Delete(context.Background(), serverGroup.UUID)
			if err != nil {
				t.Errorf("ServerGroups.Delete returned error %s\n", err)
			}
		}
	}
}
