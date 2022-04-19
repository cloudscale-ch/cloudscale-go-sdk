package cloudscale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestServersGroups_Create(t *testing.T) {
	setup()
	defer teardown()

	serverGroupRequest := &ServerGroupRequest{
		Name:       "db-servers",
		Type:       "anti-affinity",
	}

	mux.HandleFunc("/v1/server-groups", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"name":           "db-servers",
			"type":           "anti-affinity",
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{"uuid": "42cec963-fcd2-482f-bdb6-24461b2d47b1"}`)
	})

	serverGroup, err := client.ServerGroups.Create(ctx, serverGroupRequest)
	if err != nil {
		t.Errorf("ServerGroups.Create returned error: %v", err)
	}

	if id := serverGroup.UUID; id != "42cec963-fcd2-482f-bdb6-24461b2d47b1" {
		t.Errorf("expected id '%s', received '%s'", serverGroup.UUID, id)
	}

}

func TestServerGroups_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/server-groups/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"}`)
	})

	serverGroup, err := client.ServerGroups.Get(ctx, "cfde831a-4e87-4a75-960f-89b0148aa2cc")
	if err != nil {
		t.Errorf("ServerGroups.Get returned error: %v", err)
	}

	expected := &ServerGroup{UUID: "cfde831a-4e87-4a75-960f-89b0148aa2cc"}
	if !reflect.DeepEqual(serverGroup, expected) {
		t.Errorf("ServerGroups.Get\n got=%#v\nwant=%#v", serverGroup, expected)
	}
}

func TestServerGroups_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/server-groups/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)
	})

	err := client.ServerGroups.Delete(ctx, "cfde831a-4e87-4a75-960f-89b0148aa2cc")
	if err != nil {
		t.Errorf("ServerGroups.Delete returned error: %v", err)
	}
}

func TestServerGroups_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/server-groups", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[{"uuid": "47cec963-fcd2-482f-bdb6-24461b2d47b1"}]`)
	})

	serverGroups, err := client.ServerGroups.List(ctx)
	if err != nil {
		t.Errorf("ServerGroups.List returned error: %v", err)
	}

	expected := []ServerGroup{{UUID: "47cec963-fcd2-482f-bdb6-24461b2d47b1"}}
	if !reflect.DeepEqual(serverGroups, expected) {
		t.Errorf("ServerGroups.List\n got=%#v\nwant=%#v", serverGroups, expected)
	}

}
