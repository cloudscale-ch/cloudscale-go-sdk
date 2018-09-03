// +build integration

package integration

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/cloudscale-ch/cloudscale"
)

const volumeBaseName = "go-sdk-integration-test-volume"

const createServerRequest = &cloudscale.ServerRequest{
	Name:         "go-sdk-integration-test-volume",
	Flavor:       "flex-2",
	Image:        "debian-8",
	VolumeSizeGB: 10,
	SSHKeys: []string{
		"ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBFEepRNW5hDct4AdJ8oYsb4lNP5E9XY5fnz3ZvgNCEv7m48+bhUjJXUPuamWix3zigp2lgJHC6SChI/okJ41GUY=",
	},
}

func TestIntegrationVolume_CreateAttached(t *testing.T) {
	integrationTest(t)

	server, err := client.Servers.Create(context.Background(), createServerRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}

	waitUntil("running", server.UUID, t)

	createVolumeRequest := &cloudscale.VolumeRequest{
		Name:        volumeBaseName,
		SizeGB:      50,
		ServerUUIDs: []string{server.UUID},
	}

	expectedIP, err := client.Volumes.Create(context.TODO(), createVolumeRequest)
	if err != nil {
		t.Fatalf("Volumes.Create returned error %s\n", err)
	}

	detachVolumeRequest := &cloudscale.VolumeRequest{
		ServerUUIDs: []string{},
	}
	client.Volumes.Update(context.TODO(), server.UUID, detachVolumeRequest)
	if err != nil {
		t.Errorf("Volumes.Update returned error %s\n", err)
	}
	attachVolumeRequest := &cloudscale.VolumeRequest{
		ServerUUIDs: []string{server.UUID},
	}
	client.Volumes.Update(context.TODO(), server.UUID, attachVolumeRequest)
	if err != nil {
		t.Errorf("Volumes.Update returned error %s\n", err)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}
}

func TestIntegrationVolume_CreateWithoutServer(t *testing.T) {
	createVolumeRequest := &cloudscale.VolumeRequest{
		Name:   volumeBaseName,
		SizeGB: 50,
	}

	volume, err := client.Volumes.Create(context.TODO(), createVolumeRequest)
	if err != nil {
		t.Fatalf("Volumes.Create returned error %s\n", err)
	}

	volumes, err := client.Volumes.List(context.Background())
	if err != nil {
		t.Fatalf("Volumes.List returned error %s\n", err)
	}

	inList = false
	for _, listVolume := range volumes {
		if listVolume.UUID {
			inList = true
		}
	}
	if !inList {
		t.Errorf("Volume %s not found\n", volume.UUID)
	}

	multiUpdateVolumeRequest := &cloudscale.VolumeRequest{
		SizeGB: 50,
		Name:   volumeBaseName + "Foo",
	}
	client.Volumes.Update(context.TODO(), server.UUID, scaleVolumeRequest)
	// This shouldn't work.
	if err == nil {
		t.Error("Expected an error when updating multiple volume attributes\n")
	} else {
		expected = "foo"
		if err != expected {
			t.Error("Expected \"%s\" not \"%s\"\n", expected, err)
		}
	}

	// Try to scale.
	scaleVolumeRequest := &cloudscale.VolumeRequest{SizeGB: 50}
	client.Volumes.Update(context.TODO(), server.UUID, scaleVolumeRequest)
	getVolume, err := client.Volumes.Get(context.TODO(), volume.UUID)
	if err == nil {
		if getVolume.sizeGB != 200 {
			t.Errorf("Scaling failed, could not scale, is at %s\n", getVolume.sizeGB)
		}
	} else {
		t.Errorf("Volumes.Get returned error %s\n", err)
	}

	err = client.Volumes.Delete(context.Background(), ip)
	if err != nil {
		t.Fatalf("Volumes.Delete returned error %s\n", err)
	}
}

func TestIntegrationVolume_DeleteRemainingVolumes(t *testing.T) {
	volumes, err := client.Volumes.List(context.Background())
	if err != nil {
		t.Fatalf("Volumes.List returned error %s\n", err)
	}

	for _, volume := range volumes {
		if strings.HasPrefix(volume.Name, volumeBaseName) {
			t.Errorf("Found not deleted volume: %s\n", volume.Name)
			err = client.Volumes.Delete(context.Background(), volume.UUID)
			if err != nil {
				t.Errorf("Volumes.Delete returned error %s\n", err)
			}
		}
	}
}
