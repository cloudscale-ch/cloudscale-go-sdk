//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cloudscale-ch/cloudscale-go-sdk/v3"
	"golang.org/x/oauth2"
)

var (
	client        *cloudscale.Client
	testRunPrefix string
	testZone      string
)

func TestMain(m *testing.M) {
	// setup tests
	rand.Seed(time.Now().UnixNano())
	testRunPrefix = fmt.Sprintf("go-sdk-%d", rand.Intn(100000))

	token := os.Getenv("CLOUDSCALE_API_TOKEN")
	if token == "" {
		log.Fatal("Missing CLOUDSCALE_API_TOKEN, tests won't run!\n")
	}
	tc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
	client = cloudscale.NewClient(tc)

	testZone = os.Getenv("INTEGRATION_TEST_ZONE")
	if testZone == "" {
		testZone = "rma1"
	}

	// run the tests
	runResult := m.Run()

	log.Printf("Checking for leftover resources..\n")
	foundResource := false
	foundResource = foundResource || DeleteRemainingServer()
	foundResource = foundResource || DeleteRemainingServerGroups()
	foundResource = foundResource || DeleteRemainingVolumes()
	foundResource = foundResource || DeleteRemainingSubnets()
	foundResource = foundResource || DeleteRemainingNetworks()
	foundResource = foundResource || DeleteRemainingObjectsUsers()
	foundResource = foundResource || DeleteRemainingCustomImages()
	foundResource = foundResource || DeleteRemainingLoadBalancers()

	if foundResource {
		log.Fatal("Failing due to leftover resource\n")
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
		if strings.HasPrefix(server.Name, testRunPrefix) {
			foundResource = true
			log.Printf("Found not deleted server: %s (%s)\n", server.Name, server.UUID)
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
		if strings.HasPrefix(serverGroup.Name, testRunPrefix) {
			foundResource = true
			log.Printf("Found not deleted serverGroup: %s (%s)\n", serverGroup.Name, serverGroup.UUID)
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

	volumes, err := client.Volumes.List(context.Background())
	if err != nil {
		log.Fatalf("Volumes.List returned error %s\n", err)
	}

	for _, volume := range volumes {
		if strings.HasPrefix(volume.Name, testRunPrefix) {
			foundResource = true
			log.Printf("Found not deleted volume: %s (%s)\n", volume.Name, volume.UUID)
			err = client.Volumes.Delete(context.Background(), volume.UUID)
			if err != nil {
				log.Fatalf("Volumes.Delete returned error %s\n", err)
			}
		}
	}

	return foundResource
}

func DeleteRemainingSubnets() bool {
	foundResource := false

	subnets, err := client.Subnets.List(context.Background())
	if err != nil {
		log.Fatalf("Subnets.List returned error %s\n", err)
	}

	for _, subnet := range subnets {
		if strings.HasPrefix(subnet.Network.Name, testRunPrefix) {
			foundResource = true
			log.Printf("Found not deleted subnet: %s (%s) on network %s (%s)\n", subnet.CIDR, subnet.UUID, subnet.Network.Name, subnet.Network.UUID)
			err = client.Subnets.Delete(context.Background(), subnet.UUID)
			if err != nil {
				log.Fatalf("Subnets.Delete returned error %s\n", err)
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
		if strings.HasPrefix(network.Name, testRunPrefix) {
			foundResource = true
			log.Printf("Found not deleted network: %s (%s)\n", network.Name, network.UUID)
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
		if strings.HasPrefix(objectsUser.DisplayName, testRunPrefix) {
			foundResource = true
			log.Printf("Found not deleted objectsUser: %s (%s)\n", objectsUser.DisplayName, objectsUser.ID)
			err = client.ObjectsUsers.Delete(context.Background(), objectsUser.ID)
			if err != nil {
				log.Fatalf("ObjectsUsers.Delete returned error %s\n", err)
			}
		}
	}

	return foundResource
}

func DeleteRemainingCustomImages() bool {
	foundResource := false

	customImages, err := client.CustomImages.List(context.Background())
	if err != nil {
		log.Fatalf("CustomImages.List returned error %s\n", err)
	}

	for _, customImage := range customImages {
		if strings.HasPrefix(customImage.Name, testRunPrefix) {
			foundResource = true
			log.Printf("Found not deleted customImage: %s (%s)\n", customImage.Name, customImage.UUID)
			err = client.CustomImages.Delete(context.Background(), customImage.UUID)
			if err != nil {
				log.Fatalf("CustomImages.Delete returned error %s\n", err)
			}
		}
	}

	return foundResource
}

func DeleteRemainingLoadBalancers() bool {
	foundResource := false

	loadBalancers, err := client.LoadBalancers.List(context.Background())
	if err != nil {
		log.Fatalf("LoadBalancers.List returned error %s\n", err)
	}

	for _, loadBalancer := range loadBalancers {
		if strings.HasPrefix(loadBalancer.Name, testRunPrefix) {
			foundResource = true
			log.Printf("Found not deleted loadBalancer: %s (%s)\n", loadBalancer.Name, loadBalancer.UUID)
			err = client.LoadBalancers.Delete(context.Background(), loadBalancer.UUID)
			if err != nil {
				log.Fatalf("LoadBalancers.Delete returned error %s\n", err)
			}
		}
	}

	return foundResource
}
