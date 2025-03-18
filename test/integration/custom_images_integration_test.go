//go:build integration
// +build integration

package integration

import (
	"context"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff/v5"
	"github.com/cloudscale-ch/cloudscale-go-sdk/v5"
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

	customImageImport = waitForImport("success", expected.UUID, t)

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

	customImageImport := waitForImport("failed", expected.UUID, t)

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

	customImageImport = waitForImport("success", customImageImport.UUID, t)

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
	}

	customImageImport, err := client.CustomImageImports.Create(context.Background(), createCustomImageImportRequest)
	if err != nil {
		t.Fatalf("CustomImageImports.Create returned error %s\n", err)
	}

	customImageImport = waitForImport("success", customImageImport.UUID, t)

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
		serverRunningCondition,
	)
	if err != nil {
		t.Fatalf("Servers.WaitFor returned error %s\n", err)
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

func waitForImport(status string, uuid string, t *testing.T) *cloudscale.CustomImageImport {
	// An operation that may fail.
	operation := func() (*cloudscale.CustomImageImport, error) {
		i, err := client.CustomImageImports.Get(context.Background(), uuid)
		if err != nil {
			return nil, err
		}

		if i.Status != status {
			return nil, errors.New(fmt.Sprintf("Import status is: %v", i.Status))
		}
		return i, nil
	}

	result, err := backoff.Retry(context.TODO(), operation,
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
		backoff.WithMaxTries(10),
	)
	if err != nil {
		t.Fatalf("Error while waiting for status=%s change %s\n", status, err)
	}
	return result
}
