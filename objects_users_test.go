package cloudscale

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestObjectsUser_Create(t *testing.T) {
	setup()
	defer teardown()

	ObjectsUserRequest := &ObjectsUserRequest{
		DisplayName: "TestBucket",
		TaggedResourceRequest: TaggedResourceRequest{
			&TagMap{
				"tag":   "one",
				"other": "tag",
			},
		},
	}

	mux.HandleFunc("/v1/objects-users", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"display_name": "TestBucket",
			"tags": map[string]interface{}{
				"tag":   "one",
				"other": "tag",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}
		jsonStr := `{
						"href": "https://api.cloudscale.ch/v1/objects-users/6fe39134bf4178747eebc429f82cfafdd08891d4279d0d899bc4012db1db6a15",
						"id": "6fe39134bf4178747eebc429f82cfafdd08891d4279d0d899bc4012db1db6a15",
						"display_name": "alan",
						"keys": [{
							"access_key": "0ZTAIBKSGYBRHQ09G11W",
							"secret_key": "bn2ufcwbIa0ARLc5CLRSlVaCfFxPHOpHmjKiH34T"
						}],
						"tags": {
							"tag": "one",
							"other": "tag"
						}
					}`
		fmt.Fprintf(w, jsonStr)
	})

	objectsUser, err := client.ObjectsUsers.Create(ctx, ObjectsUserRequest)
	if err != nil {
		t.Errorf("ObjectsUser.Create returned error: %v", err)
		return
	}

	if id := objectsUser.ID; id != "6fe39134bf4178747eebc429f82cfafdd08891d4279d0d899bc4012db1db6a15" {
		t.Errorf("expected id '%s', received '%s'", objectsUser.ID, id)
	}
}

func TestObjectsUser_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/objects-users/6fe39134bf4178747eebc429f82cfafdd08891d4279d0d899bc4012db1db6a15", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"id": "6fe39134bf4178747eebc429f82cfafdd08891d4279d0d899bc4012db1db6a15"}`)
	})

	objectUser, err := client.ObjectsUsers.Get(ctx, "6fe39134bf4178747eebc429f82cfafdd08891d4279d0d899bc4012db1db6a15")
	if err != nil {
		t.Errorf("ObjectsUser.Get returned error: %v", err)
	}

	expected := &ObjectsUser{ID: "6fe39134bf4178747eebc429f82cfafdd08891d4279d0d899bc4012db1db6a15"}
	if !reflect.DeepEqual(objectUser, expected) {
		t.Errorf("ObjectsUser.Get\n got=%#v\nwant=%#v", objectUser, expected)
	}
}

func TestObjectsUser_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/objects-users/6fe39134bf4178747eebc429f82cfafdd08891d4279d0d899bc4012db1db6a15", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)
	})

	err := client.ObjectsUsers.Delete(ctx, "6fe39134bf4178747eebc429f82cfafdd08891d4279d0d899bc4012db1db6a15")
	if err != nil {
		t.Errorf("ObjectsUser.Delete returned error: %v", err)
	}
}

func TestObjectsUser_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/objects-users", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[{"id": "6fe39134bf4178747eebc429f82cfafdd08891d4279d0d899bc4012db1db6a15"}]`)
	})

	objectUsers, err := client.ObjectsUsers.List(ctx)
	if err != nil {
		t.Errorf("ObjectsUser.List returned error: %v", err)
	}

	expected := []ObjectsUser{{ID: "6fe39134bf4178747eebc429f82cfafdd08891d4279d0d899bc4012db1db6a15"}}
	if !reflect.DeepEqual(objectUsers, expected) {
		t.Errorf("ObjectsUser.List\n got=%#v\nwant=%#v", objectUsers, expected)
	}

}

func TestObjectsUser_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/objects-users/6fe39134bf4178747eebc429f82cfafdd08891d4279d0d899bc4012db1db6a15", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPatch)
	})

	userID := "6fe39134bf4178747eebc429f82cfafdd08891d4279d0d899bc4012db1db6a15"

	req := &ObjectsUserRequest{
		DisplayName: "new_name",
	}
	err := client.ObjectsUsers.Update(context.TODO(), userID, req)
	if err != nil {
		t.Errorf("ObjectsUser.Update returned error: %v", err)
	}
}
