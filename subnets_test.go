package cloudscale

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestSubnets_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/subnets/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc"}`)
	})

	subnet, err := client.Subnets.Get(ctx, "cfde831a-4e87-4a75-960f-89b0148aa2cc")
	if err != nil {
		t.Errorf("Subnets.Get returned error: %v", err)
	}

	expected := &Subnet{UUID: "cfde831a-4e87-4a75-960f-89b0148aa2cc"}
	if !reflect.DeepEqual(subnet, expected) {
		t.Errorf("Subnets.Get\n got=%#v\nwant=%#v", subnet, expected)
	}
}

func TestSubnets_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/subnets", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[{"uuid": "47cec963-fcd2-482f-bdb6-24461b2d47b1"}]`)
	})

	subnets, err := client.Subnets.List(ctx)
	if err != nil {
		t.Errorf("Subnets.List returned error: %v", err)
	}

	expected := []Subnet{{UUID: "47cec963-fcd2-482f-bdb6-24461b2d47b1"}}
	if !reflect.DeepEqual(subnets, expected) {
		t.Errorf("Subnets.List\n got=%#v\nwant=%#v", subnets, expected)
	}

}
