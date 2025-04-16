//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"github.com/cloudscale-ch/cloudscale-go-sdk/v6"
	"io"
	"net/http"
	"testing"
	"time"
)

const testImageURL = "https://at-images.objects.lpg.cloudscale.ch/alpine"

func TestIntegrationCustomImage_CRUD(t *testing.T) {
	integrationTest(t)

	createCustomImageRequest := &cloudscale.CustomImageImportRequest{
		Name:             testRunPrefix,
		URL:              "https://at-images.objects.lpg.cloudscale.ch/alpine",
		UserDataHandling: "extend-cloud-config",
		Zones:            []string{"lpg1", "rma1"},
		SourceFormat:     "raw",
		FirmwareType:     "uefi",
	}

	expected, err := client.CustomImageImports.Create(context.TODO(), createCustomImageRequest)
	if err != nil {
		t.Fatalf("CustomImageImports.Create returned error %s\n", err)
	}

	customImageImport, err := client.CustomImageImports.Get(context.Background(), expected.UUID)
	if err != nil {
		t.Fatalf("CustomImageImports.Get returned error %s\n", err)
	}

	customImageImport, err = client.CustomImageImports.WaitFor(context.Background(), customImageImport.UUID, cloudscale.ImportIsSuccessful)
	if err != nil {
		t.Fatalf("CustomImageImports.WaitFor returned error %s\n", err)
	}

	if uuid := customImageImport.UUID; uuid != expected.UUID {
		t.Errorf("customImageImport.UUID got=%s\nwant=%s", uuid, expected.UUID)
	}

	if name := customImageImport.CustomImage.Name; name != expected.CustomImage.Name {
		t.Errorf("customImageImport.CustomImage.Name got=%s\nwant=%s", name, expected.CustomImage.Name)
	}

	if uuid := customImageImport.CustomImage.UUID; uuid != expected.CustomImage.UUID {
		t.Errorf("customImageImport.CustomImage.UUID got=%s\nwant=%s", uuid, expected.CustomImage.UUID)
	}

	customImages, err := client.CustomImages.List(context.Background())
	if err != nil {
		t.Fatalf("CustomImages.List returned error %s\n", err)
	}

	if numCustomImages := len(customImages); numCustomImages == 0 {
		t.Errorf("CustomImage.List got=%d\nwant=%d\n", numCustomImages, 1)
	}

	if h := time.Since(customImages[0].CreatedAt).Hours(); !(-1 < h && h < 1) {
		t.Errorf("customImages[0].CreatedAt ourside of expected range. got=%v", customImages[0].CreatedAt)
	}

	customImageImports, err := client.CustomImageImports.List(context.Background())
	if err != nil {
		t.Fatalf("CustomImageImports.List returned error %s\n", err)
	}

	if numCustomImageImports := len(customImageImports); numCustomImageImports == 0 {
		t.Errorf("CustomImageImport.List got=%d\nwant=%d\n", numCustomImageImports, 1)
	}

	customImage, err := client.CustomImages.Get(context.Background(), customImageImport.CustomImage.UUID)

	if err != nil {
		t.Fatalf("CustomImages.Get returned error %s\n", err)
	}

	if firmwareType := customImage.FirmwareType; customImage.FirmwareType != "uefi" {
		t.Errorf("customImage.FirmwareType got=%s\nwant=%s", firmwareType, "uefi")
	}

	for _, algo := range []string{"md5", "sha256"} {
		checksumURL := fmt.Sprintf("%s.%s", testImageURL, algo)
		resp, err := http.Get(checksumURL)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatal(fmt.Sprintf("Wrong http status code\n got=%#v\nwant=%#v", resp.Status, http.StatusOK))
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		checksum := string(body)
		if checksum != customImage.Checksums[algo] {
			t.Error(fmt.Sprintf("Checksum does not match\n got=%#v\nwant=%#v", customImage.Checksums[algo], checksum))
		}
	}

	err = client.CustomImages.Delete(context.Background(), customImageImport.CustomImage.UUID)
	if err != nil {
		t.Fatalf("CustomImages.Delete returned error %s\n", err)
	}
}

func TestIntegrationCustomImage_InvalidURL(t *testing.T) {
	integrationTest(t)

	createCustomImageRequest := &cloudscale.CustomImageImportRequest{
		Name:             testRunPrefix,
		URL:              "http://www.cloudscale.ch/this-does-and-will-never-exist",
		UserDataHandling: "extend-cloud-config",
		Zones:            []string{"rma1"},
		SourceFormat:     "raw",
	}

	expected, err := client.CustomImageImports.Create(context.TODO(), createCustomImageRequest)
	if err != nil {
		t.Fatalf("CustomImageImports.Create returned error %s\n", err)
	}

	customImageImport, err := client.CustomImageImports.WaitFor(
		context.Background(),
		expected.UUID,
		func(resource *cloudscale.CustomImageImport) (bool, error) { return resource.Status == "failed", nil },
	)
	if err != nil {
		t.Fatalf("CustomImageImports.WaitFor returned error %s\n", err)
	}
	expectedMessage := "Expected HTTP 200, got HTTP 404"
	if errorMessage := customImageImport.ErrorMessage; errorMessage != expectedMessage {
		t.Errorf("customImageImport.ErrorMessage got=%s\nwant=%s\n", errorMessage, expectedMessage)
	}
}

func TestIntegrationCustomImage_Update(t *testing.T) {

	createCustomImageImportRequest := &cloudscale.CustomImageImportRequest{
		Name:             testRunPrefix,
		URL:              testImageURL,
		UserDataHandling: "extend-cloud-config",
		Zones:            []string{"lpg1", "rma1"},
		SourceFormat:     "raw",
	}

	customImageImport, err := client.CustomImageImports.Create(context.TODO(), createCustomImageImportRequest)
	if err != nil {
		t.Fatalf("CustomImageImports.Create returned error %s\n", err)
	}

	customImageImport, err = client.CustomImageImports.WaitFor(context.Background(), customImageImport.UUID, cloudscale.ImportIsSuccessful)
	if err != nil {
		t.Fatalf("CustomImageImports.WaitFor returned error %s\n", err)
	}

	expectedNewName := fmt.Sprintf("%s-renamed", testRunPrefix)
	updateRequest := &cloudscale.CustomImageRequest{
		Name: expectedNewName,
	}

	err = client.CustomImages.Update(context.Background(), customImageImport.UUID, updateRequest)
	if err != nil {
		t.Fatalf("CustomImageImports.Update returned error %s\n", err)
	}

	updatedCustomImage, err := client.CustomImages.Get(context.Background(), customImageImport.UUID)
	if err != nil {
		t.Fatalf("CustomImageImports.Get returned error %s\n", err)
	}

	if actualName := updatedCustomImage.Name; actualName != expectedNewName {
		t.Errorf("CustomImageImport Name\ngot=%#v\nwant=%#v", updatedCustomImage.Name, expectedNewName)
	}

	err = client.CustomImages.Delete(context.Background(), customImageImport.UUID)
	if err != nil {
		t.Fatalf("CustomImageImports.Delete returned error %s\n", err)
	}
}

func TestIntegrationCustomImage_Boot(t *testing.T) {

	createCustomImageImportRequest := &cloudscale.CustomImageImportRequest{
		Name:             testRunPrefix,
		URL:              testImageURL,
		UserDataHandling: "extend-cloud-config",
		Zones:            []string{"lpg1", "rma1"},
		SourceFormat:     "raw",
		Slug:             "fancy",
	}

	customImageImport, err := client.CustomImageImports.Create(context.Background(), createCustomImageImportRequest)
	if err != nil {
		t.Fatalf("CustomImageImports.Create returned error %s\n", err)
	}

	customImageImport, err = client.CustomImageImports.WaitFor(context.Background(), customImageImport.UUID, cloudscale.ImportIsSuccessful)
	if err != nil {
		t.Fatalf("CustomImageImports.WaitFor returned error %s\n", err)
	}

	createServerRequest := &cloudscale.ServerRequest{
		Name:         testRunPrefix,
		Flavor:       "flex-4-2",
		Image:        customImageImport.CustomImage.UUID,
		VolumeSizeGB: 10,
		SSHKeys: []string{
			pubKey,
		},
	}

	server, err := client.Servers.Create(context.Background(), createServerRequest)
	if err != nil {
		t.Fatalf("Servers.Create returned error %s\n", err)
	}
	_, err = client.Servers.WaitFor(
		context.Background(),
		server.UUID,
		cloudscale.ServerIsRunning,
	)
	if err != nil {
		t.Fatalf("Servers.WaitFor returned error %s\n", err)
	}

	if server.Image.Slug != "custom:fancy" {
		t.Errorf("Server.Image.Slug got=%s, want=%s", server.Image.Slug, "=custom:fancy")
	}
	if server.Image.Name != createCustomImageImportRequest.Name {
		t.Errorf("Server.Image.Name got=%s, want '%s'", server.Image.Name, createCustomImageImportRequest.Name)
	}
	if server.Image.OperatingSystem != "" {
		t.Errorf("Server.Image.OperatingSystem got=%s, want an empty string", server.Image.OperatingSystem)
	}
	if server.Image.DefaultUsername != "" {
		t.Errorf("Server.Image.DefaultUsername got=%s, want an empty string", server.Image.DefaultUsername)
	}

	err = client.Servers.Delete(context.Background(), server.UUID)
	if err != nil {
		t.Fatalf("Servers.Delete returned error %s\n", err)
	}

	err = client.CustomImages.Delete(context.Background(), customImageImport.UUID)
	if err != nil {
		t.Fatalf("CustomImageImports.Delete returned error %s\n", err)
	}
}
