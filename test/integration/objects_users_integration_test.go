//go:build integration
// +build integration

package integration

import (
	"context"
	"fmt"
	"testing"

	"github.com/cloudscale-ch/cloudscale-go-sdk/v4"
)

func TestIntegrationObjectsUser_CRUD(t *testing.T) {
	integrationTest(t)

	expected, err := createObjectsUser(t)
	if err != nil {
		t.Fatalf("ObjectsUsers.Create returned error %s\n", err)
	}

	objectsUser, err := client.ObjectsUsers.Get(context.Background(), expected.ID)
	if err != nil {
		t.Fatalf("ObjectsUsers.Get returned error %s\n", err)
	}

	if id := objectsUser.ID; id != expected.ID {
		t.Errorf("ObjectsUser.ID got=%s\nwant=%s", id, expected.ID)
	}
	if access_key := objectsUser.Keys[0]["access_key"]; access_key != expected.Keys[0]["access_key"] {
		t.Errorf("ObjectsUser.Keys[0][\"access_key\"] got=%s\nwant=%s", access_key, expected.Keys[0]["access_key"])
	}
	if secret_key := objectsUser.Keys[0]["secret_key"]; secret_key != expected.Keys[0]["secret_key"] {
		t.Errorf("ObjectsUser.Keys[0][\"secret_key\"] got=%s\nwant=%s", secret_key, expected.Keys[0]["secret_key"])
	}

	err = client.ObjectsUsers.Delete(context.Background(), objectsUser.ID)
	if err != nil {
		t.Fatalf("ObjectsUsers.Get returned error %s\n", err)
	}

}

func TestIntegrationObjectsUser_UpdateRest(t *testing.T) {
	integrationTest(t)

	objectsUser, err := createObjectsUser(t)
	if err != nil {
		t.Fatalf("ObjectsUsers.Create returned error %s\n", err)
	}

	// Try to rename.
	renamedName := fmt.Sprintf("%s-renamed", testRunPrefix)
	renameRequest := &cloudscale.ObjectsUserRequest{DisplayName: renamedName}
	err = client.ObjectsUsers.Update(context.TODO(), objectsUser.ID, renameRequest)
	if err != nil {
		t.Errorf("ObjectsUsers.Update failed %s\n", err)
	}

	getObjectsUser, err := client.ObjectsUsers.Get(context.TODO(), objectsUser.ID)
	if err == nil {
		if getObjectsUser.DisplayName != renamedName {
			t.Errorf("Renaming failed, could not rename, is at %s\n", getObjectsUser.DisplayName)
		}
	} else {
		t.Errorf("ObjectsUserRequest.Get returned error %s\n", err)
	}

	err = client.ObjectsUsers.Delete(context.Background(), objectsUser.ID)
	if err != nil {
		t.Fatalf("ObjectsUsers.Get returned error %s\n", err)
	}

}

func createObjectsUser(t *testing.T) (*cloudscale.ObjectsUser, error) {
	createRequest := &cloudscale.ObjectsUserRequest{
		DisplayName: testRunPrefix,
	}

	return client.ObjectsUsers.Create(context.Background(), createRequest)
}
