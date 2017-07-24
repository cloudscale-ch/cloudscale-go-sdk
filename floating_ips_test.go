package cloudscale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestFloatingIPs_Create(t *testing.T) {
	setup()
	defer teardown()

	floatingIPrequest := &FloatingIPCreateRequest{
		IPVersion: 6,
		Server:    "47cec963-fcd2-482f-bdb6-24461b2d47b1",
	}

	mux.HandleFunc("/v1/floating-ips", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"ip_version": float64(6),
			"server":     "47cec963-fcd2-482f-bdb6-24461b2d47b1",
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		fmt.Fprintf(w, `{"network": "2001:db8::cafe/128"}`)
	})

	floatingIP, err := client.FloatingIPs.Create(ctx, floatingIPrequest)
	if err != nil {
		t.Errorf("FloatingIPs.Create returned error: %v", err)
	}

	if network := floatingIP.Network; network != "2001:db8::cafe/128" {
		t.Errorf("expected network '%s', received '%s'", floatingIP.Network, network)
	}

}

func TestFloatingIPs_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/floating-ips/192.0.2.123", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"network": "192.0.2.123/32"}`)
	})

	floatingIP, err := client.FloatingIPs.Get(ctx, "192.0.2.123")
	if err != nil {
		t.Errorf("FloatingIPs.Get returned error: %v", err)
	}

	expected := &FloatingIP{Network: "192.0.2.123/32"}
	if !reflect.DeepEqual(floatingIP, expected) {
		t.Errorf("FloatingIPs.Get\n got=%#v\nwant=%#v", floatingIP, expected)
	}
}

func TestFloatingIPs_Update(t *testing.T) {
	setup()
	defer teardown()

	updateRequest := &FloatingIPUpdateRequest{
		Server: "47777777-fcd2-482f-bdb6-24461b2d47b1",
	}

	mux.HandleFunc("/v1/floating-ips/192.0.2.123", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"server": "47777777-fcd2-482f-bdb6-24461b2d47b1",
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body = %#v, expected %#v", v, expected)
		}

		fmt.Fprintf(w, `{"network": "192.0.2.123/32"}`)
	})

	floatingIP, err := client.FloatingIPs.Update(ctx, "192.0.2.123", updateRequest)
	if err != nil {
		t.Errorf("FloatingIps.Update returned error: %v", err)
	} else {
		if network := floatingIP.Network; network != "192.0.2.123/32" {
			t.Errorf("expected network '%s', received '%s'", "192.0.2.123/32", network)
		}
	}
}

func TestFloatingIPs_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/floating-ips/192.0.2.123", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)
	})

	err := client.FloatingIPs.Delete(ctx, "192.0.2.123")
	if err != nil {
		t.Errorf("FloatingIPs.Delete returned error: %v", err)
	}
}

func TestFloatingIPs_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/floating-ips", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[{"network": "192.0.2.123/32"}]`)
	})

	servers, err := client.FloatingIPs.List(ctx)
	if err != nil {
		t.Errorf("Servers.List returned error: %v", err)
	}

	expected := []FloatingIP{{Network: "192.0.2.123/32"}}
	if !reflect.DeepEqual(servers, expected) {
		t.Errorf("FloatingIPs.List\n got=%#v\nwant=%#v", servers, expected)
	}

}
