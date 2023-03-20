//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"github.com/cloudscale-ch/cloudscale-go-sdk/v2"
	"reflect"
	"testing"
)

func getInitialTags() cloudscale.TagMap {
	return cloudscale.TagMap{
		fmt.Sprintf("yin-%s", testRunPrefix): "yang",
	}
}

func getInitialTagsKeyOnly() cloudscale.TagMap {
	return cloudscale.TagMap{
		fmt.Sprintf("yin-%s", testRunPrefix): "",
	}
}

func getNewTags() cloudscale.TagMap {
	return cloudscale.TagMap{
		fmt.Sprintf("yab-%s", testRunPrefix): "yum",
	}
}

func getEmptyTags() cloudscale.TagMap {
	return cloudscale.TagMap{}
}

func getNewTagsKeyOnly() cloudscale.TagMap {
	return cloudscale.TagMap{
		fmt.Sprintf("yab-%s", testRunPrefix): "",
	}
}

func TestIntegrationTags_Server(t *testing.T) {
	integrationTest(t)

	createRequest := getDefaultServerRequest()
	initialTags := getInitialTags()
	createRequest.Tags = &initialTags

	server, err := createServer(t, &createRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}

	getResult, err := client.Servers.Get(context.Background(), server.UUID)
	if err != nil {
		t.Errorf("Servers.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult.Tags, initialTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	updateRequest := cloudscale.ServerUpdateRequest{}
	newTags := getNewTags()
	updateRequest.Tags = &newTags

	err = client.Servers.Update(context.Background(), server.UUID, &updateRequest)
	if err != nil {
		t.Errorf("Servers.Update returned error: %v", err)
	}
	getResult2, err := client.Servers.Get(context.Background(), server.UUID)
	if err != nil {
		t.Errorf("Servers.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult2.Tags, newTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	// test querying with tags
	initialTagsKeyOnly := getInitialTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.Servers.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("Servers.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	newTagsKeyOnly := getNewTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.Servers.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("Servers.List returned error %s\n", err)
		}
		if len(res) != 1 {
			t.Errorf("Expected exactly one result when filter with %#v, got: %#v", tags, len(res))
		}
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}
}

func TestIntegrationTags_Volume(t *testing.T) {
	integrationTest(t)

	createRequest := cloudscale.VolumeRequest{
		Name:   testRunPrefix,
		SizeGB: 3,
	}
	initialTags := getInitialTags()
	createRequest.Tags = &initialTags

	volume, err := client.Volumes.Create(context.Background(), &createRequest)
	if err != nil {
		t.Fatalf("Volumes.Create returned error %s\n", err)
	}

	getResult, err := client.Volumes.Get(context.Background(), volume.UUID)
	if err != nil {
		t.Errorf("Volumes.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult.Tags, initialTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	updateRequest := cloudscale.VolumeRequest{}
	newTags := getNewTags()
	updateRequest.Tags = &newTags

	err = client.Volumes.Update(context.Background(), volume.UUID, &updateRequest)
	if err != nil {
		t.Errorf("Volumes.Update returned error: %v", err)
	}
	getResult2, err := client.Volumes.Get(context.Background(), volume.UUID)
	if err != nil {
		t.Errorf("Volumes.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult2.Tags, newTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	// test querying with tags
	initialTagsKeyOnly := getInitialTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.Volumes.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("Volumes.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	newTagsKeyOnly := getNewTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.Volumes.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("Volumes.List returned error %s\n", err)
		}
		if len(res) != 1 {
			t.Errorf("Expected exactly one result when filter with %#v, got: %#v", tags, len(res))
		}
	}

	err = client.Volumes.Delete(context.Background(), volume.UUID)
	if err != nil {
		t.Fatalf("Volumes.Delete returned error %s\n", err)
	}
}

func TestIntegrationTags_FloatingIP(t *testing.T) {
	integrationTest(t)

	serverCreateRequest := getDefaultServerRequest()

	server, err := createServer(t, &serverCreateRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}

	createRequest := cloudscale.FloatingIPCreateRequest{
		IPVersion: 6,
		Server:    server.UUID,
	}
	initialTags := getInitialTags()
	createRequest.Tags = &initialTags

	floatingIP, err := client.FloatingIPs.Create(context.Background(), &createRequest)
	if err != nil {
		t.Fatalf("FloatingIPs.Create returned error %s\n", err)
	}

	getResult, err := client.FloatingIPs.Get(context.Background(), floatingIP.IP())
	if err != nil {
		t.Errorf("FloatingIPs.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult.Tags, initialTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	updateRequest := cloudscale.FloatingIPUpdateRequest{}
	newTags := getNewTags()
	updateRequest.Tags = &newTags

	err = client.FloatingIPs.Update(context.Background(), floatingIP.IP(), &updateRequest)
	if err != nil {
		t.Errorf("FloatingIPs.Update returned error: %v", err)
	}
	getResult2, err := client.FloatingIPs.Get(context.Background(), floatingIP.IP())
	if err != nil {
		t.Errorf("FloatingIPs.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult2.Tags, newTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	// test querying with tags
	initialTagsKeyOnly := getInitialTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.FloatingIPs.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("FloatingIPs.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	newTagsKeyOnly := getNewTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.FloatingIPs.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("FloatingIPs.List returned error %s\n", err)
		}
		if len(res) != 1 {
			t.Errorf("Expected exactly one result when filter with %#v, got: %#v", tags, len(res))
		}
	}

	err = client.FloatingIPs.Delete(context.Background(), floatingIP.IP())
	if err != nil {
		t.Fatalf("FloatingIPs.Delete returned error %s\n", err)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}
}

func TestIntegrationTags_ObjectsUser(t *testing.T) {
	integrationTest(t)

	createRequest := cloudscale.ObjectsUserRequest{
		DisplayName: testRunPrefix,
	}
	initialTags := getInitialTags()
	createRequest.Tags = &initialTags

	objectsUser, err := client.ObjectsUsers.Create(context.Background(), &createRequest)
	if err != nil {
		t.Fatalf("ObjectsUsers.Create returned error %s\n", err)
	}

	getResult, err := client.ObjectsUsers.Get(context.Background(), objectsUser.ID)
	if err != nil {
		t.Errorf("ObjectsUsers.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult.Tags, initialTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	updateRequest := cloudscale.ObjectsUserRequest{}
	newTags := getNewTags()
	updateRequest.Tags = &newTags

	err = client.ObjectsUsers.Update(context.Background(), objectsUser.ID, &updateRequest)
	if err != nil {
		t.Errorf("ObjectsUsers.Update returned error: %v", err)
	}
	getResult2, err := client.ObjectsUsers.Get(context.Background(), objectsUser.ID)
	if err != nil {
		t.Errorf("ObjectsUsers.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult2.Tags, newTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	// test querying with tags
	initialTagsKeyOnly := getInitialTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.ObjectsUsers.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("ObjectsUsers.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	newTagsKeyOnly := getNewTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.ObjectsUsers.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("ObjectsUsers.List returned error %s\n", err)
		}
		if len(res) != 1 {
			t.Errorf("Expected exactly one result when filter with %#v, got: %#v", tags, len(res))
		}
	}

	err = client.ObjectsUsers.Delete(context.Background(), objectsUser.ID)
	if err != nil {
		t.Fatalf("ObjectsUsers.Delete returned error %s\n", err)
	}

}

func TestIntegrationTags_Network(t *testing.T) {
	integrationTest(t)

	createRequest := cloudscale.NetworkCreateRequest{
		Name: testRunPrefix,
	}
	initialTags := getInitialTags()
	createRequest.Tags = &initialTags

	network, err := client.Networks.Create(context.Background(), &createRequest)
	if err != nil {
		t.Fatalf("Networks.Create returned error %s\n", err)
	}

	getResult, err := client.Networks.Get(context.Background(), network.UUID)
	if err != nil {
		t.Errorf("Networks.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult.Tags, initialTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	updateRequest := cloudscale.NetworkUpdateRequest{}
	newTags := getNewTags()
	updateRequest.Tags = &newTags

	err = client.Networks.Update(context.Background(), network.UUID, &updateRequest)
	if err != nil {
		t.Errorf("Networks.Update returned error: %v", err)
	}
	getResult2, err := client.Networks.Get(context.Background(), network.UUID)
	if err != nil {
		t.Errorf("Networks.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult2.Tags, newTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	// test querying with tags
	initialTagsKeyOnly := getInitialTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.Networks.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("Networks.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	newTagsKeyOnly := getNewTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.Networks.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("Networks.List returned error %s\n", err)
		}
		if len(res) != 1 {
			t.Errorf("Expected exactly one result when filter with %#v, got: %#v", tags, len(res))
		}
	}

	err = client.Networks.Delete(context.Background(), network.UUID)
	if err != nil {
		t.Fatalf("Networks.Delete returned error %s\n", err)
	}

}

func TestIntegrationTags_Subnet(t *testing.T) {
	integrationTest(t)

	autoCreateSubnet := false
	createNetworkRequest := cloudscale.NetworkCreateRequest{
		Name:                 testRunPrefix,
		AutoCreateIPV4Subnet: &autoCreateSubnet,
	}
	network, err := client.Networks.Create(context.Background(), &createNetworkRequest)
	if err != nil {
		t.Fatalf("Networks.Create returned error %s\n", err)
	}

	createRequest := cloudscale.SubnetCreateRequest{
		CIDR:    "172.16.0.0/14",
		Network: network.UUID,
	}
	initialTags := getInitialTags()
	createRequest.Tags = &initialTags
	subnet, err := client.Subnets.Create(context.Background(), &createRequest)
	if err != nil {
		t.Fatalf("Subnets.Create returned error %s\n", err)
	}

	getResult, err := client.Subnets.Get(context.Background(), subnet.UUID)
	if err != nil {
		t.Errorf("Subnets.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult.Tags, initialTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	updateRequest := cloudscale.SubnetUpdateRequest{}
	newTags := getNewTags()
	updateRequest.Tags = &newTags

	err = client.Subnets.Update(context.Background(), subnet.UUID, &updateRequest)
	if err != nil {
		t.Errorf("Subnets.Update returned error: %v", err)
	}
	getResult2, err := client.Subnets.Get(context.Background(), subnet.UUID)
	if err != nil {
		t.Errorf("Subnets.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult2.Tags, newTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	// test querying with tags
	initialTagsKeyOnly := getInitialTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.Subnets.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("Subnets.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	newTagsKeyOnly := getNewTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.Subnets.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("Subnets.List returned error %s\n", err)
		}
		if len(res) != 1 {
			t.Errorf("Expected exactly one result when filter with %#v, got: %#v", tags, len(res))
		}
	}

	err = client.Subnets.Delete(context.Background(), subnet.UUID)
	if err != nil {
		t.Fatalf("Subnets.Delete returned error %s\n", err)
	}
	err = client.Networks.Delete(context.Background(), network.UUID)
	if err != nil {
		t.Fatalf("Networks.Delete returned error %s\n", err)
	}

}

func TestIntegrationTags_ServerGroup(t *testing.T) {
	integrationTest(t)

	createRequest := cloudscale.ServerGroupRequest{
		Name: testRunPrefix + "-group",
		Type: "anti-affinity",
	}
	initialTags := getInitialTags()
	createRequest.Tags = &initialTags

	serverGroup, err := client.ServerGroups.Create(context.Background(), &createRequest)
	if err != nil {
		t.Fatalf("ServerGroups.Create returned error %s\n", err)
	}

	getResult, err := client.ServerGroups.Get(context.Background(), serverGroup.UUID)
	if err != nil {
		t.Errorf("ServerGroups.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult.Tags, initialTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	updateRequest := cloudscale.ServerGroupRequest{}
	// Test update with empty tags
	emptyTags := getEmptyTags()
	updateRequest.Tags = &emptyTags
	err = client.ServerGroups.Update(context.Background(), serverGroup.UUID, &updateRequest)
	if err != nil {
		t.Errorf("ServerGroups.Update returned error: %v", err)
	}
	getResult2, err := client.ServerGroups.Get(context.Background(), serverGroup.UUID)
	if err != nil {
		t.Errorf("ServerGroups.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult2.Tags, emptyTags) {
		t.Errorf("Tagging failed, could not untag, is at %s\n", getResult.Tags)
	}

	// Test update with some tags again
	newTags := getNewTags()
	updateRequest.Tags = &newTags

	err = client.ServerGroups.Update(context.Background(), serverGroup.UUID, &updateRequest)
	if err != nil {
		t.Errorf("ServerGroups.Update returned error: %v", err)
	}
	getResult2, err = client.ServerGroups.Get(context.Background(), serverGroup.UUID)
	if err != nil {
		t.Errorf("ServerGroups.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult2.Tags, newTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	// test querying with tags
	initialTagsKeyOnly := getInitialTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.ServerGroups.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("ServerGroups.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	newTagsKeyOnly := getNewTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.ServerGroups.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("ServerGroups.List returned error %s\n", err)
		}
		if len(res) != 1 {
			t.Errorf("Expected exactly one result when filter with %#v, got: %#v", tags, len(res))
		}
	}

	err = client.ServerGroups.Delete(context.Background(), serverGroup.UUID)
	if err != nil {
		t.Fatalf("ServerGroups.Delete returned error %s\n", err)
	}

}

func TestIntegrationTags_CustomImage(t *testing.T) {
	integrationTest(t)

	createRequest := cloudscale.CustomImageImportRequest{
		Name:             testRunPrefix,
		URL:              "https://at-images.objects.lpg.cloudscale.ch/alpine",
		UserDataHandling: "extend-cloud-config",
		Zones:            []string{"lpg1"},
		SourceFormat:     "raw",
	}
	initialTags := getInitialTags()
	createRequest.Tags = &initialTags

	customImageImport, err := client.CustomImageImports.Create(context.Background(), &createRequest)
	if err != nil {
		t.Fatalf("CustomImages.Create returned error %s\n", err)
	}
	waitForImport("success", customImageImport.UUID, t)

	imageID := customImageImport.CustomImage.UUID
	getResult, err := client.CustomImages.Get(context.Background(), imageID)
	if err != nil {
		t.Errorf("CustomImages.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult.Tags, initialTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	updateRequest := cloudscale.CustomImageRequest{}
	newTags := getNewTags()
	updateRequest.Tags = &newTags

	err = client.CustomImages.Update(context.Background(), imageID, &updateRequest)
	if err != nil {
		t.Errorf("CustomImages.Update returned error: %v", err)
	}
	getResult2, err := client.CustomImages.Get(context.Background(), imageID)
	if err != nil {
		t.Errorf("CustomImages.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult2.Tags, newTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	// test querying with tags
	initialTagsKeyOnly := getInitialTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.CustomImages.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("CustomImages.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	newTagsKeyOnly := getNewTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.CustomImages.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("CustomImages.List returned error %s\n", err)
		}
		if len(res) != 1 {
			t.Errorf("Expected exactly one result when filter with %#v, got: %#v", tags, len(res))
		}
	}

	err = client.CustomImages.Delete(context.Background(), imageID)
	if err != nil {
		t.Fatalf("CustomImages.Delete returned error %s\n", err)
	}
}

func TestIntegrationTags_LoadBalancerAndRelatedResources(t *testing.T) {
	integrationTest(t)

	createRequest := cloudscale.LoadBalancerRequest{
		Name:   testRunPrefix,
		Flavor: "lb-small",
	}
	createRequest.Zone = "rma1"
	initialTags := getInitialTags()
	createRequest.Tags = &initialTags

	loadBalancer, err := client.LoadBalancers.Create(context.Background(), &createRequest)
	if err != nil {
		t.Fatalf("LoadBalancers.Create returned error %s\n", err)
	}

	waitUntilLB("running", loadBalancer.UUID, t)

	getResult, err := client.LoadBalancers.Get(context.Background(), loadBalancer.UUID)
	if err != nil {
		t.Errorf("LoadBalancers.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult.Tags, initialTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	updateRequest := cloudscale.LoadBalancerRequest{}
	newTags := getNewTags()
	updateRequest.Tags = &newTags

	err = client.LoadBalancers.Update(context.Background(), loadBalancer.UUID, &updateRequest)
	if err != nil {
		t.Errorf("LoadBalancers.Update returned error: %v", err)
	}
	getResult2, err := client.LoadBalancers.Get(context.Background(), loadBalancer.UUID)
	if err != nil {
		t.Errorf("LoadBalancers.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult2.Tags, newTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	// test querying with tags
	initialTagsKeyOnly := getInitialTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.LoadBalancers.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("LoadBalancers.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	newTagsKeyOnly := getNewTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.LoadBalancers.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("LoadBalancers.List returned error %s\n", err)
		}
		if len(res) != 1 {
			t.Errorf("Expected exactly one result when filter with %#v, got: %#v", tags, len(res))
		}
	}

	// call these test cases inline to avoid recreating the load balancer
	testIntegrationTags_LoadBalancerPool(t, loadBalancer)

	err = client.LoadBalancers.Delete(context.Background(), loadBalancer.UUID)
	if err != nil {
		t.Fatalf("LoadBalancers.Delete returned error %s\n", err)
	}

}

func testIntegrationTags_LoadBalancerPool(t *testing.T, b *cloudscale.LoadBalancer) {
	createRequest := cloudscale.LoadBalancerPoolRequest{
		Name:         testRunPrefix,
		Algorithm:    "round_robin",
		Protocol:     "tcp",
		LoadBalancer: b.UUID,
	}
	initialTags := getInitialTags()
	createRequest.Tags = &initialTags

	pool, err := client.LoadBalancerPools.Create(context.Background(), &createRequest)
	if err != nil {
		t.Fatalf("LoadBalancerPools.Create returned error %s\n", err)
	}

	getResult, err := client.LoadBalancerPools.Get(context.Background(), pool.UUID)
	if err != nil {
		t.Errorf("LoadBalancers.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult.Tags, initialTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	updateRequest := cloudscale.LoadBalancerPoolRequest{}
	newTags := getNewTags()
	updateRequest.Tags = &newTags

	err = client.LoadBalancerPools.Update(context.Background(), pool.UUID, &updateRequest)
	if err != nil {
		t.Errorf("LoadBalancers.Update returned error: %v", err)
	}
	getResult2, err := client.LoadBalancerPools.Get(context.Background(), pool.UUID)
	if err != nil {
		t.Errorf("LoadBalancers.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult2.Tags, newTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	// test querying with tags
	initialTagsKeyOnly := getInitialTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.LoadBalancerPools.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("LoadBalancers.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	newTagsKeyOnly := getNewTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.LoadBalancerPools.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("LoadBalancers.List returned error %s\n", err)
		}
		if len(res) != 1 {
			t.Errorf("Expected exactly one result when filter with %#v, got: %#v", tags, len(res))
		}
	}

	// call these test cases inline to avoid recreating the load balancer
	testIntegrationTags_LoadBalancerListner(t, pool)
	testIntegrationTags_LoadBalancerPoolMember(t, pool)
	testIntegrationTags_LoadBalancerHealthMonitor(t, pool)

	err = client.LoadBalancerPools.Delete(context.Background(), pool.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerPools.Delete returned error %s\n", err)
	}
}

func testIntegrationTags_LoadBalancerListner(t *testing.T, p *cloudscale.LoadBalancerPool) {
	createRequest := cloudscale.LoadBalancerListenerRequest{
		Name:         testRunPrefix,
		Pool:         p.UUID,
		Protocol:     "tcp",
		ProtocolPort: 8080,
	}
	initialTags := getInitialTags()
	createRequest.Tags = &initialTags

	listener, err := client.LoadBalancerListeners.Create(context.Background(), &createRequest)
	if err != nil {
		t.Fatalf("LoadBalancerListeners.Create returned error %s\n", err)
	}

	getResult, err := client.LoadBalancerListeners.Get(context.Background(), listener.UUID)
	if err != nil {
		t.Errorf("LoadBalancers.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult.Tags, initialTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	updateRequest := cloudscale.LoadBalancerListenerRequest{}
	newTags := getNewTags()
	updateRequest.Tags = &newTags

	err = client.LoadBalancerListeners.Update(context.Background(), listener.UUID, &updateRequest)
	if err != nil {
		t.Errorf("LoadBalancers.Update returned error: %v", err)
	}
	getResult2, err := client.LoadBalancerListeners.Get(context.Background(), listener.UUID)
	if err != nil {
		t.Errorf("LoadBalancers.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult2.Tags, newTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	// test querying with tags
	initialTagsKeyOnly := getInitialTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.LoadBalancerListeners.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("LoadBalancers.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	newTagsKeyOnly := getNewTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.LoadBalancerListeners.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("LoadBalancers.List returned error %s\n", err)
		}
		if len(res) != 1 {
			t.Errorf("Expected exactly one result when filter with %#v, got: %#v", tags, len(res))
		}
	}

	err = client.LoadBalancerListeners.Delete(context.Background(), listener.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerListeners.Delete returned error %s\n", err)
	}
}

func testIntegrationTags_LoadBalancerPoolMember(t *testing.T, p *cloudscale.LoadBalancerPool) {
	network, subnet, err := createNetworkAndSubnet()
	if err != nil {
		t.Fatalf("error while creating network and subnet: %s\n", err)
	}

	createRequest := cloudscale.LoadBalancerPoolMemberRequest{
		Name:         testRunPrefix,
		ProtocolPort: 8080,
		Address:      "192.168.42.12",
		Subnet:       subnet.UUID,
	}
	initialTags := getInitialTags()
	createRequest.Tags = &initialTags

	member, err := client.LoadBalancerPoolMembers.Create(context.Background(), p.UUID, &createRequest)
	if err != nil {
		t.Fatalf("LoadBalancerPoolMembers.Create returned error %s\n", err)
	}

	getResult, err := client.LoadBalancerPoolMembers.Get(context.Background(), p.UUID, member.UUID)
	if err != nil {
		t.Errorf("LoadBalancers.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult.Tags, initialTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	updateRequest := cloudscale.LoadBalancerPoolMemberRequest{}
	newTags := getNewTags()
	updateRequest.Tags = &newTags

	err = client.LoadBalancerPoolMembers.Update(context.Background(), p.UUID, member.UUID, &updateRequest)
	if err != nil {
		t.Errorf("LoadBalancers.Update returned error: %v", err)
	}
	getResult2, err := client.LoadBalancerPoolMembers.Get(context.Background(), p.UUID, member.UUID)
	if err != nil {
		t.Errorf("LoadBalancers.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult2.Tags, newTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	// test querying with tags
	initialTagsKeyOnly := getInitialTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.LoadBalancerPoolMembers.List(context.Background(), p.UUID, cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("LoadBalancers.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	newTagsKeyOnly := getNewTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.LoadBalancerPoolMembers.List(context.Background(), p.UUID, cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("LoadBalancers.List returned error %s\n", err)
		}
		if len(res) != 1 {
			t.Errorf("Expected exactly one result when filter with %#v, got: %#v", tags, len(res))
		}
	}

	err = client.LoadBalancerPoolMembers.Delete(context.Background(), p.UUID, member.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerPoolMembers.Delete returned error %s\n", err)
	}

	err = client.Networks.Delete(context.Background(), network.UUID)
	if err != nil {
		t.Fatalf("Networks.Delete returned error %s\n", err)
	}
}

func testIntegrationTags_LoadBalancerHealthMonitor(t *testing.T, p *cloudscale.LoadBalancerPool) {
	createRequest := cloudscale.LoadBalancerHealthMonitorRequest{
		Pool:          p.UUID,
		DelayS:        10,
		TimeoutS:      3,
		UpThreshold:   5,
		DownThreshold: 10,
		Type:          "tcp",
	}
	initialTags := getInitialTags()
	createRequest.Tags = &initialTags

	healthMonitor, err := client.LoadBalancerHealthMonitors.Create(context.Background(), &createRequest)
	if err != nil {
		t.Fatalf("LoadBalancerHealthMonitors.Create returned error %s\n", err)
	}

	getResult, err := client.LoadBalancerHealthMonitors.Get(context.Background(), healthMonitor.UUID)
	if err != nil {
		t.Errorf("LoadBalancers.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult.Tags, initialTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	updateRequest := cloudscale.LoadBalancerHealthMonitorRequest{}
	newTags := getNewTags()
	updateRequest.Tags = &newTags

	err = client.LoadBalancerHealthMonitors.Update(context.Background(), healthMonitor.UUID, &updateRequest)
	if err != nil {
		t.Errorf("LoadBalancers.Update returned error: %v", err)
	}
	getResult2, err := client.LoadBalancerHealthMonitors.Get(context.Background(), healthMonitor.UUID)
	if err != nil {
		t.Errorf("LoadBalancers.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult2.Tags, newTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	// test querying with tags
	initialTagsKeyOnly := getInitialTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.LoadBalancerHealthMonitors.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("LoadBalancers.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	newTagsKeyOnly := getNewTagsKeyOnly()
	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.LoadBalancerHealthMonitors.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("LoadBalancers.List returned error %s\n", err)
		}
		if len(res) != 1 {
			t.Errorf("Expected exactly one result when filter with %#v, got: %#v", tags, len(res))
		}
	}

	err = client.LoadBalancerHealthMonitors.Delete(context.Background(), healthMonitor.UUID)
	if err != nil {
		t.Fatalf("LoadBalancerHealthMonitors.Delete returned error %s\n", err)
	}
}
