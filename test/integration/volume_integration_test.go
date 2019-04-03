// +build integration

package integration

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/cloudscale-ch/cloudscale"
)

const volumeBaseName = "go-sdk-integration-test-volume"

func TestIntegrationVolume_CreateAttached(t *testing.T) {
	integrationTest(t)

	createServerRequest := &cloudscale.ServerRequest{
		Name:         "go-sdk-integration-test-volume",
		Flavor:       "flex-2",
		Image:        DefaultImageSlug,
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

	createVolumeRequest := &cloudscale.Volume{
		Name:        volumeBaseName,
		SizeGB:      50,
		ServerUUIDs: &[]string{server.UUID},
	}

	volume, err := client.Volumes.Create(context.TODO(), createVolumeRequest)
	if err != nil {
		t.Fatalf("Volumes.Create returned error %s\n", err)
	}

	time.Sleep(3 * time.Second)
	detachVolumeRequest := &cloudscale.Volume{
		ServerUUIDs: &[]string{},
	}
	err = client.Volumes.Update(context.TODO(), volume.UUID, detachVolumeRequest)
	if err != nil {
		t.Errorf("Volumes.Update returned error %s\n", err)
	}
	attachVolumeRequest := &cloudscale.Volume{
		ServerUUIDs: &[]string{server.UUID},
	}

	time.Sleep(3 * time.Second)
	err = client.Volumes.Update(context.TODO(), volume.UUID, attachVolumeRequest)
	if err != nil {
		t.Errorf("Volumes.Update returned error %s\n", err)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}
	err = client.Volumes.Delete(context.Background(), volume.UUID)
	if err != nil {
		t.Fatalf("Volumes.Delete returned error %s\n", err)
	}
}

func TestIntegrationVolume_CreateWithoutServer(t *testing.T) {
	createVolumeRequest := &cloudscale.Volume{
		Name:   volumeBaseName,
		SizeGB: 50,
	}

	volume, err := client.Volumes.Create(context.TODO(), createVolumeRequest)
	if err != nil {
		t.Fatalf("Volumes.Create returned error %s\n", err)
	}

	volumes, err := client.Volumes.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("Volumes.List returned error %s\n", err)
	}

	inList := false
	for _, listVolume := range volumes {
		if listVolume.UUID == volume.UUID {
			inList = true
		}
	}
	if !inList {
		t.Errorf("Volume %s not found\n", volume.UUID)
	}

	multiUpdateVolumeRequest := &cloudscale.Volume{
		SizeGB: 50,
		Name:   volumeBaseName + "Foo",
	}
	err = client.Volumes.Update(context.TODO(), volume.UUID, multiUpdateVolumeRequest)
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

	const scaleSize = 200
	// Try to scale.
	scaleVolumeRequest := &cloudscale.Volume{SizeGB: scaleSize}
	err = client.Volumes.Update(context.TODO(), volume.UUID, scaleVolumeRequest)
	getVolume, err := client.Volumes.Get(context.TODO(), volume.UUID)
	if err == nil {
		if getVolume.SizeGB != scaleSize {
			t.Errorf("Scaling failed, could not scale, is at %d\n", getVolume.SizeGB)
		}
	} else {
		t.Errorf("Volumes.Get returned error %s\n", err)
	}

	err = client.Volumes.Delete(context.Background(), volume.UUID)
	if err != nil {
		t.Fatalf("Volumes.Delete returned error %s\n", err)
	}
}

func TestIntegrationVolume_DeleteRemainingVolumes(t *testing.T) {
	volumes, err := client.Volumes.List(context.Background(), nil)
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
