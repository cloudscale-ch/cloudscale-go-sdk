// +build integration

package integration

import (
	"context"
	"github.com/cloudscale-ch/cloudscale-go-sdk"
	"reflect"
	"testing"
)

var initialTags = cloudscale.TagMap{
	"yin": "yang",
}

var initialTagsKeyOnly = cloudscale.TagMap{
	"yin": "",
}

var newTags = cloudscale.TagMap{
	"yab": "yum",
}

var newTagsKeyOnly = cloudscale.TagMap{
	"yab": "",
}

func TestIntegrationTags_Server(t *testing.T) {
	integrationTest(t)

	createRequest := getDefaultServerRequest()
	createRequest.Tags = initialTags

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
	updateRequest.Tags = newTags

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
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.Servers.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("Servers.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.Servers.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("Servers.List returned error %s\n", err)
		}
		if len(res) < 1 {
			t.Errorf("Expected at least one result when filter with %#v, got: %#v", tags, len(res))
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
		Name:   volumeBaseName,
		SizeGB: 3,
	}
	createRequest.Tags = initialTags

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
	updateRequest.Tags = newTags

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
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.Volumes.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("Volumes.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.Volumes.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("Volumes.List returned error %s\n", err)
		}
		if len(res) < 1 {
			t.Errorf("Expected at least one result when filter with %#v, got: %#v", tags, len(res))
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
	createRequest.Tags = initialTags

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
	updateRequest.Tags = newTags

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
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.FloatingIPs.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("FloatingIPs.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.FloatingIPs.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("FloatingIPs.List returned error %s\n", err)
		}
		if len(res) < 1 {
			t.Errorf("Expected at least one result when filter with %#v, got: %#v", tags, len(res))
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
		DisplayName: baseObjectsUserName,
	}
	createRequest.Tags = initialTags

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
	updateRequest.Tags = newTags

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
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.ObjectsUsers.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("ObjectsUsers.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.ObjectsUsers.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("ObjectsUsers.List returned error %s\n", err)
		}
		if len(res) < 1 {
			t.Errorf("Expected at least one result when filter with %#v, got: %#v", tags, len(res))
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
		Name: networkBaseName,
	}
	createRequest.Tags = initialTags

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
	updateRequest.Tags = newTags

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
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.Networks.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("Networks.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.Networks.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("Networks.List returned error %s\n", err)
		}
		if len(res) < 1 {
			t.Errorf("Expected at least one result when filter with %#v, got: %#v", tags, len(res))
		}
	}

	err = client.Networks.Delete(context.Background(), network.UUID)
	if err != nil {
		t.Fatalf("Networks.Delete returned error %s\n", err)
	}

}

func TestIntegrationTags_Subnet(t *testing.T) {
	integrationTest(t)

	createNetworkRequest := cloudscale.NetworkCreateRequest{
		Name: networkBaseName,
	}
	network, err := client.Networks.Create(context.Background(), &createNetworkRequest)
	if err != nil {
		t.Fatalf("Networks.Create returned error %s\n", err)
	}

	createRequest := cloudscale.SubnetCreateRequest{
		CIDR:    "172.16.0.0/14",
		Network: network.UUID,
	}
	createRequest.Tags = initialTags
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

	// test querying with tags
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.Subnets.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("Subnets.List returned error %s\n", err)
		}
		if len(res) < 1 {
			t.Errorf("Expected at least one result when filter with %#v, got: %#v", tags, len(res))
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
		Name: serverBaseName + "-group",
		Type: "anti-affinity",
	}
	createRequest.Tags = initialTags

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
	updateRequest.Tags = newTags

	err = client.ServerGroups.Update(context.Background(), serverGroup.UUID, &updateRequest)
	if err != nil {
		t.Errorf("ServerGroups.Update returned error: %v", err)
	}
	getResult2, err := client.ServerGroups.Get(context.Background(), serverGroup.UUID)
	if err != nil {
		t.Errorf("ServerGroups.Get returned error %s\n", err)
	}
	if !reflect.DeepEqual(getResult2.Tags, newTags) {
		t.Errorf("Tagging failed, could not tag, is at %s\n", getResult.Tags)
	}

	// test querying with tags
	for _, tags := range []cloudscale.TagMap{initialTags, initialTagsKeyOnly} {
		res, err := client.ServerGroups.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("ServerGroups.List returned error %s\n", err)
		}
		if len(res) > 0 {
			t.Errorf("Expected no result when filter with %#v, got: %#v", tags, res)
		}
	}

	for _, tags := range []cloudscale.TagMap{newTags, newTagsKeyOnly} {
		res, err := client.ServerGroups.List(context.Background(), cloudscale.WithTagFilter(tags))
		if err != nil {
			t.Errorf("ServerGroups.List returned error %s\n", err)
		}
		if len(res) < 1 {
			t.Errorf("Expected at least one result when filter with %#v, got: %#v", tags, len(res))
		}
	}

	err = client.ServerGroups.Delete(context.Background(), serverGroup.UUID)
	if err != nil {
		t.Fatalf("ServerGroups.Delete returned error %s\n", err)
	}

}
