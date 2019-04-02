package cloudscale

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestServers_Create(t *testing.T) {
	setup()
	defer teardown()

	serverRequest := &ServerRequest{
		Name:         "mysql",
		Flavor:       "flex-4",
		Image:        "debian",
		VolumeSizeGB: 50,
		SSHKeys:      []string{"key"},
	}

	mux.HandleFunc("/v1/servers", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"name":           "mysql",
			"flavor":         "flex-4",
			"image":          "debian",
			"volume_size_gb": float64(50),
			"ssh_keys":       []interface{}{"key"},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{"uuid": "47cec963-fcd2-482f-bdb6-24461b2d47b1"}`)
	})

	server, err := client.Servers.Create(ctx, serverRequest)
	if err != nil {
		t.Errorf("Servers.Create returned error: %v", err)
	}

	if id := server.UUID; id != "47cec963-fcd2-482f-bdb6-24461b2d47b1" {
		t.Errorf("expected id '%s', received '%s'", server.UUID, id)
	}

}

func TestServers_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/servers/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"}`)
	})

	server, err := client.Servers.Get(ctx, "cfde831a-4e87-4a75-960f-89b0148aa2cc")
	if err != nil {
		t.Errorf("Server.Get returned error: %v", err)
	}

	expected := &Server{UUID: "cfde831a-4e87-4a75-960f-89b0148aa2cc"}
	if !reflect.DeepEqual(server, expected) {
		t.Errorf("Servers.Get\n got=%#v\nwant=%#v", server, expected)
	}
}

func TestServers_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/servers/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)
	})

	err := client.Servers.Delete(ctx, "cfde831a-4e87-4a75-960f-89b0148aa2cc")
	if err != nil {
		t.Errorf("Serveers.Delete returned error: %v", err)
	}
}

func TestServers_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/servers", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[{"uuid": "47cec963-fcd2-482f-bdb6-24461b2d47b1"}]`)
	})

	servers, err := client.Servers.List(ctx)
	if err != nil {
		t.Errorf("Servers.List returned error: %v", err)
	}

	expected := []Server{{UUID: "47cec963-fcd2-482f-bdb6-24461b2d47b1"}}
	if !reflect.DeepEqual(servers, expected) {
		t.Errorf("Servers.List\n got=%#v\nwant=%#v", servers, expected)
	}

}

func TestServers_Reboot(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/servers/47cec963-fcd2-482f-bdb6-24461b2d47b1/reboot", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
	})

	err := client.Servers.Reboot(ctx, "47cec963-fcd2-482f-bdb6-24461b2d47b1")

	if err != nil {
		t.Errorf("Serveers.Reboot returned error: %v", err)
	}
}

func TestServers_Start(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/servers/47cec963-fcd2-482f-bdb6-24461b2d47b1/start", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
	})

	err := client.Servers.Start(ctx, "47cec963-fcd2-482f-bdb6-24461b2d47b1")

	if err != nil {
		t.Errorf("Serveers.Start returned error: %v", err)
	}
}

func TestServers_Stop(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/servers/47cec963-fcd2-482f-bdb6-24461b2d47b1/stop", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
	})

	err := client.Servers.Stop(ctx, "47cec963-fcd2-482f-bdb6-24461b2d47b1")

	if err != nil {
		t.Errorf("Servers.Stop returned error: %v", err)
	}
}

func TestServers_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/servers/47cec963-fcd2-482f-bdb6-24461b2d47b1/stop", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
	})
	mux.HandleFunc("/v1/servers/47cec963-fcd2-482f-bdb6-24461b2d47b1/start", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
	})
	mux.HandleFunc("/v1/servers/47cec963-fcd2-482f-bdb6-24461b2d47b1/reboot", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)
	})

	serverID := "47cec963-fcd2-482f-bdb6-24461b2d47b1"

	err := client.Servers.Update(context.TODO(), serverID, ServerRunning)
	if err != nil {
		t.Errorf("Servers.Update returned error: %v", err)
	}

	err = client.Servers.Update(context.TODO(), serverID, ServerStopped)
	if err != nil {
		t.Errorf("Servers.Update returned error: %v", err)
	}

	err = client.Servers.Update(context.TODO(), serverID, ServerRebooted)
	if err != nil {
		t.Errorf("Servers.Update returned error: %v", err)
	}

	err = client.Servers.Update(context.TODO(), serverID, "Not an actual action")
	if err == nil {
		t.Errorf("Servers.Update returned error: %v", err)
	}

}
