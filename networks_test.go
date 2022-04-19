package cloudscale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNetworks_Create(t *testing.T) {
	setup()
	defer teardown()

	networkRequest := &NetworkCreateRequest{
		Name:       "netzli",
	}

	mux.HandleFunc("/v1/networks", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"name":           "netzli",
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

	network, err := client.Networks.Create(ctx, networkRequest)
	if err != nil {
		t.Errorf("Networks.Create returned error: %v", err)
	}

	if id := network.UUID; id != "42cec963-fcd2-482f-bdb6-24461b2d47b1" {
		t.Errorf("expected id '%s', received '%s'", network.UUID, id)
	}

}

func TestNetworks_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/networks/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc", "created_at": "2019-05-27T16:45:32.241824Z"}`)
	})

	network, err := client.Networks.Get(ctx, "cfde831a-4e87-4a75-960f-89b0148aa2cc")
	if err != nil {
		t.Errorf("Networks.Get returned error: %v", err)
	}

	expected := &Network{UUID: "cfde831a-4e87-4a75-960f-89b0148aa2cc", CreatedAt: time.Date(2019, time.Month(5), 27, 16, 45, 32, 241824000, time.UTC)}
	if !reflect.DeepEqual(network, expected) {
		t.Errorf("Networks.Get\n got=%#v\nwant=%#v", network, expected)
	}
}

func TestNetworks_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/networks/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)
	})

	err := client.Networks.Delete(ctx, "cfde831a-4e87-4a75-960f-89b0148aa2cc")
	if err != nil {
		t.Errorf("Networks.Delete returned error: %v", err)
	}
}

func TestNetworks_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/networks", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[{"uuid": "47cec963-fcd2-482f-bdb6-24461b2d47b1"}]`)
	})

	networks, err := client.Networks.List(ctx)
	if err != nil {
		t.Errorf("Networks.List returned error: %v", err)
	}

	expected := []Network{{UUID: "47cec963-fcd2-482f-bdb6-24461b2d47b1"}}
	if !reflect.DeepEqual(networks, expected) {
		t.Errorf("Networks.List\n got=%#v\nwant=%#v", networks, expected)
	}

}
