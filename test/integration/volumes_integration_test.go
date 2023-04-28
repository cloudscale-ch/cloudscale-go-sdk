//go:build integration
// +build integration

package integration

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cloudscale-ch/cloudscale-go-sdk/v3"
)

func TestIntegrationVolume_CreateAttached(t *testing.T) {
	integrationTest(t)

	createServerRequest := &cloudscale.ServerRequest{
		Name:         testRunPrefix,
		Flavor:       "flex-4-2",
		Image:        DefaultImageSlug,
		VolumeSizeGB: 10,
		SSHKeys: []string{
			pubKey,
		},
	}

	server, err := client.Servers.Create(context.Background(), createServerRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}

	waitUntil("running", server.UUID, t)

	createVolumeRequest := &cloudscale.VolumeRequest{
		Name:        testRunPrefix,
		SizeGB:      50,
		ServerUUIDs: &[]string{server.UUID},
	}

	volume, err := client.Volumes.Create(context.TODO(), createVolumeRequest)
	if err != nil {
		t.Fatalf("Volumes.Create returned error %s\n", err)
	}

	if h := time.Since(volume.CreatedAt).Hours(); !(-1 < h && h < 1) {
		t.Errorf("volume.CreatedAt ourside of expected range. got=%v", volume.CreatedAt)
	}

	time.Sleep(3 * time.Second)
	detachVolumeRequest := &cloudscale.VolumeRequest{
		ServerUUIDs: &[]string{},
	}
	err = client.Volumes.Update(context.TODO(), volume.UUID, detachVolumeRequest)
	if err != nil {
		t.Errorf("Volumes.Update returned error %s\n", err)
	}
	attachVolumeRequest := &cloudscale.VolumeRequest{
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
	createVolumeRequest := &cloudscale.VolumeRequest{
		Name:   testRunPrefix,
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

	inList := false
	for _, listVolume := range volumes {
		if listVolume.UUID == volume.UUID {
			inList = true
		}
	}
	if !inList {
		t.Errorf("Volume %s not found\n", volume.UUID)
	}

	multiUpdateVolumeRequest := &cloudscale.VolumeRequest{
		SizeGB: 50,
		Name:   testRunPrefix + "Foo",
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
	scaleVolumeRequest := &cloudscale.VolumeRequest{SizeGB: scaleSize}
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

func TestIntegrationVolume_AttachToNewServer(t *testing.T) {
	createVolumeRequest := &cloudscale.VolumeRequest{
		Name:   testRunPrefix,
		SizeGB: 50,
	}

	volume, err := client.Volumes.Create(context.TODO(), createVolumeRequest)
	if err != nil {
		t.Fatalf("Volumes.Create returned error %s\n", err)
	}

	createServerRequest := &cloudscale.ServerRequest{
		Name:         testRunPrefix,
		Flavor:       "flex-4-2",
		Image:        DefaultImageSlug,
		VolumeSizeGB: 10,
		SSHKeys: []string{
			pubKey,
		},
	}

	server, err := client.Servers.Create(context.Background(), createServerRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}

	waitUntil("running", server.UUID, t)

	volumeAttachRequest := &cloudscale.VolumeRequest{
		ServerUUIDs: &[]string{server.UUID},
	}

	err = client.Volumes.Update(context.Background(), volume.UUID, volumeAttachRequest)
	if err != nil {
		t.Fatalf("Volumes.Update returned error %s\n", err)
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

func TestIntegrationVolume_ListByName(t *testing.T) {
	volumeName := testRunPrefix + "-name-test"
	createVolumeRequest := &cloudscale.VolumeRequest{
		Name:   volumeName,
		SizeGB: 5,
	}

	volume, err := client.Volumes.Create(context.TODO(), createVolumeRequest)
	if err != nil {
		t.Fatalf("Volumes.Create returned error %s\n", err)
	}

	volumes, err := client.Volumes.List(context.Background())
	if err != nil {
		t.Fatalf("Volumes.List returned error %s\n", err)
	}
	if actual := len(volumes); actual <= 0 {
		t.Errorf("Expected at lest one volume, got: %#v", actual)
	}

	volumes, err = client.Volumes.List(context.Background(), cloudscale.WithNameFilter(volumeName))
	if err != nil {
		t.Fatalf("Volumes.List returned error %s\n", err)
	}
	if actual := len(volumes); actual != 1 {
		t.Errorf("Expected at exactly one volume, got: %#v", volumes)
	}

	volumes, err = client.Volumes.List(context.Background(), cloudscale.WithNameFilter("reykjavik"))
	if err != nil {
		t.Fatalf("Volumes.List returned error %s\n", err)
	}
	if actual := len(volumes); actual != 0 {
		t.Errorf("Expected no volumes, got: %#v", volumes)
	}

	err = client.Volumes.Delete(context.Background(), volume.UUID)
	if err != nil {
		t.Fatalf("Volumes.Delete returned error %s\n", err)
	}
}

func TestIntegrationVolume_MultiSite(t *testing.T) {
	integrationTest(t)

	allZones, err := getAllZones()
	if err != nil {
		t.Fatalf("getAllRegions returned error %s\n", err)
	}

	if len(allZones) <= 1 {
		t.Skip("Skipping MultiSite test.")
	}

	var wg sync.WaitGroup

	for _, zone := range allZones {
		wg.Add(1)
		go createVolumeInZoneAndAssert(t, zone, &wg)
	}

	wg.Wait()
}

func createVolumeInZoneAndAssert(t *testing.T, zone cloudscale.Zone, wg *sync.WaitGroup) {
	defer wg.Done()

	createVolumeRequest := &cloudscale.VolumeRequest{
		Name:   testRunPrefix,
		SizeGB: 50,
	}

	createVolumeRequest.Zone = zone.Slug

	volume, err := client.Volumes.Create(context.TODO(), createVolumeRequest)
	if err != nil {
		t.Fatalf("Volumes.Create returned error %s\n", err)
	}

	if volume.Zone != zone {
		t.Errorf("Volume in wrong Zone\n got=%#v\nwant=%#v", volume.Zone, zone)
	}

	err = client.Volumes.Delete(context.Background(), volume.UUID)
	if err != nil {
		t.Errorf("Volumes.Delete returned error %s\n", err)
	}
}
