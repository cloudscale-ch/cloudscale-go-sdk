// +build integration

package integration

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/cloudscale-ch/cloudscale-go-sdk"
	"golang.org/x/oauth2"
)

var (
	client *cloudscale.Client
)

func TestMain(m *testing.M) {
	// setup tests
	token := os.Getenv("CLOUDSCALE_TOKEN")
	if token == "" {
		log.Fatal("Missing CLOUDSCALE_TOKEN, tests won't run!\n")
	}
	tc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
	client = cloudscale.NewClient(tc)

	// run the tests
	runResult := m.Run()

	log.Printf("Checking for reamaining resources..\n")
	foundResource := false
	foundResource = foundResource || DeleteRemainingServer()
	foundResource = foundResource || DeleteRemainingServerGroups()
	foundResource = foundResource || DeleteRemainingVolumes()
	foundResource = foundResource || DeleteRemainingNetworks()
	foundResource = foundResource || DeleteRemainingObjectsUsers()

	if (foundResource) {
		log.Fatal("Failing due to remaining resource\n")
	}
	os.Exit(runResult)
}

func DeleteRemainingServer() bool {
	foundResource := false

	servers, err := client.Servers.List(context.Background())
	if err != nil {
		log.Fatalf("Servers.List returned error %s\n", err)
	}

	for _, server := range servers {
		if strings.HasPrefix(server.Name, serverBaseName) {
			foundResource = true
			log.Printf("Found not deleted server: %s\n", server.Name)
			err = client.Servers.Delete(context.Background(), server.UUID)
			if err != nil {
				log.Fatalf("Servers.Delete returned error %s\n", err)
			}
		}
	}

	return foundResource
}

func DeleteRemainingServerGroups() bool {
	foundResource := false

	serverGroups, err := client.ServerGroups.List(context.Background())
	if err != nil {
		log.Fatalf("ServerGroups.List returned error %s\n", err)
	}

	for _, serverGroup := range serverGroups {
		if strings.HasPrefix(serverGroup.Name, serverBaseName) {
			foundResource = true
			log.Printf("Found not deleted serverGroup: %s\n", serverGroup.Name)
			err = client.ServerGroups.Delete(context.Background(), serverGroup.UUID)
			if err != nil {
				log.Fatalf("ServerGroups.Delete returned error %s\n", err)
			}
		}
	}

	return foundResource
}

func DeleteRemainingVolumes() bool {
	foundResource := false

	volumes, err := client.Volumes.List(context.Background(), nil)
	if err != nil {
		log.Fatalf("Volumes.List returned error %s\n", err)
	}

	for _, volume := range volumes {
		if strings.HasPrefix(volume.Name, "go-sdk-integration-test") {
			foundResource = true
			log.Printf("Found not deleted volume: %s\n", volume.Name)
			err = client.Volumes.Delete(context.Background(), volume.UUID)
			if err != nil {
				log.Fatalf("Volumes.Delete returned error %s\n", err)
			}
		}
	}

	return foundResource
}

func DeleteRemainingNetworks() bool {
	foundResource := false

	networks, err := client.Networks.List(context.Background())
	if err != nil {
		log.Fatalf("Networks.List returned error %s\n", err)
	}

	for _, network := range networks {
		if strings.HasPrefix(network.Name, "go-sdk-integration-test") {
			foundResource = true
			log.Printf("Found not deleted network: %s\n", network.Name)
			err = client.Networks.Delete(context.Background(), network.UUID)
			if err != nil {
				log.Fatalf("Networks.Delete returned error %s\n", err)
			}
		}
	}

	return foundResource
}

func DeleteRemainingObjectsUsers() bool {
	foundResource := false

	objectsUsers, err := client.ObjectsUsers.List(context.Background())
	if err != nil {
		log.Fatalf("ObjectsUsers.List returned error %s\n", err)
	}

	for _, objectsUser := range objectsUsers {
		if strings.HasPrefix(objectsUser.DisplayName, serverBaseName) {
			foundResource = true
			log.Printf("Found not deleted objectsUser: %s\n", objectsUser.DisplayName)
			err = client.ObjectsUsers.Delete(context.Background(), objectsUser.ID)
			if err != nil {
				log.Fatalf("ObjectsUsers.Delete returned error %s\n", err)
			}
		}
	}

	return foundResource
}
