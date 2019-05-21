package cloudscale

import (
	"fmt"
	"net/http"
	"testing"
)

func TestServerId(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/openstack/2017-02-22/meta_data.json", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprintf(w, `{"meta": {"cloudscale_uuid": "foobar"}}`)
	})

	serverID, err := metadataClient.GetServerID()
	if err != nil {
		t.Errorf("GetServerID returned error: %v", err)
	}

	if serverID != "foobar" {
		t.Errorf("expected id 'foobar', received '%s'", serverID)
	}

}

func TestRawUserData(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/openstack/2017-02-22/user_data", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `abcdef`)
	})

	userData, err := metadataClient.GetRawUserData()
	if err != nil {
		t.Errorf("Server.Get returned error: %v", err)
	}

	if userData != "abcdef" {
		t.Errorf("expected 'abcdef', received '%s'", userData)
	}
}
