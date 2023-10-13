package cloudscale

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestIntegrationLoadBalancerListener_GetWithPool(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/load-balancers/listeners/754c3797-57de-4fd2-a5c9-97efa2a0c466", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
            "href": "https://lab-api.cloudscale.ch/v1/load-balancers/listeners/754c3797-57de-4fd2-a5c9-97efa2a0c466",
            "uuid": "754c3797-57de-4fd2-a5c9-97efa2a0c466",
            "name": "web-lb1-listener",
            "created_at": "2023-10-13T13:28:38.672592Z",
            "pool": {
              "href": "https://lab-api.cloudscale.ch/v1/load-balancers/pools/0aa55841-d151-4e65-b63c-3282bc06cac0",
              "uuid": "0aa55841-d151-4e65-b63c-3282bc06cac0",
              "name": "web-lb1-pool"
            },
            "load_balancer": {
              "href": "https://lab-api.cloudscale.ch/v1/load-balancers/bc7b04c9-04ee-471f-b719-26f29c767f6c",
              "uuid": "bc7b04c9-04ee-471f-b719-26f29c767f6c",
              "name": "web-lb1"
            },
            "protocol": "tcp",
            "protocol_port": 80,
            "allowed_cidrs": [],
            "timeout_client_data_ms": 50000,
            "timeout_member_connect_ms": 5000,
            "timeout_member_data_ms": 50000,
            "tags": {}
          }`)
	})

	listener, err := client.LoadBalancerListeners.Get(ctx, "754c3797-57de-4fd2-a5c9-97efa2a0c466")
	if err != nil {
		t.Errorf("LoadBalancerListeners.Get returned error: %v", err)
	}

	expected := &LoadBalancerListener{
		TaggedResource: TaggedResource{
			Tags: nil,
		},
		HREF: "https://lab-api.cloudscale.ch/v1/load-balancers/listeners/754c3797-57de-4fd2-a5c9-97efa2a0c466",
		UUID: "754c3797-57de-4fd2-a5c9-97efa2a0c466",
		Name: "web-lb1-listener",
		Pool: &LoadBalancerPoolStub{
			HREF: "https://lab-api.cloudscale.ch/v1/load-balancers/pools/0aa55841-d151-4e65-b63c-3282bc06cac0",
			UUID: "0aa55841-d151-4e65-b63c-3282bc06cac0",
			Name: "web-lb1-pool",
		},
		LoadBalancer: LoadBalancerStub{
			HREF: "https://lab-api.cloudscale.ch/v1/load-balancers/bc7b04c9-04ee-471f-b719-26f29c767f6c",
			UUID: "bc7b04c9-04ee-471f-b719-26f29c767f6c",
			Name: "web-lb1",
		},
		Protocol:               "tcp",
		ProtocolPort:           80,
		AllowedCIDRs:           []string{},
		TimeoutClientDataMS:    50000,
		TimeoutMemberConnectMS: 5000,
		TimeoutMemberDataMS:    50000,
		CreatedAt:              time.Date(2023, time.Month(10), 13, 13, 28, 38, 672592000, time.UTC),
	}
	expected.Tags = TagMap{}

	if !reflect.DeepEqual(listener, expected) {
		t.Errorf("LoadBalancerListeners.Get\n got=%#v\nwant=%#v", listener, expected)
	}
}

func TestIntegrationLoadBalancerListener_GetWithoutPool(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/load-balancers/listeners/3d6ca1f4-5aea-41f5-b724-0f3054b60e85", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
            "href": "https://lab-api.cloudscale.ch/v1/load-balancers/listeners/3d6ca1f4-5aea-41f5-b724-0f3054b60e85",
            "uuid": "3d6ca1f4-5aea-41f5-b724-0f3054b60e85",
            "name": "web-lb1-listener-without-pool",
            "created_at": "2023-10-13T13:28:38.672592Z",
            "pool": null,
            "load_balancer": {
              "href": "https://lab-api.cloudscale.ch/v1/load-balancers/bc7b04c9-04ee-471f-b719-26f29c767f6c",
              "uuid": "bc7b04c9-04ee-471f-b719-26f29c767f6c",
              "name": "web-lb1"
            },
            "protocol": "brieftaube",
            "protocol_port": 80,
            "allowed_cidrs": [],
            "timeout_client_data_ms": 18180,
            "timeout_member_connect_ms": 1818,
            "timeout_member_data_ms": 18180,
            "tags": {}
          }`)
	})

	listener, err := client.LoadBalancerListeners.Get(ctx, "3d6ca1f4-5aea-41f5-b724-0f3054b60e85")
	if err != nil {
		t.Errorf("LoadBalancerListeners.Get returned error: %v", err)
	}

	expected := &LoadBalancerListener{
		TaggedResource: TaggedResource{
			Tags: nil,
		},
		HREF: "https://lab-api.cloudscale.ch/v1/load-balancers/listeners/3d6ca1f4-5aea-41f5-b724-0f3054b60e85",
		UUID: "3d6ca1f4-5aea-41f5-b724-0f3054b60e85",
		Name: "web-lb1-listener-without-pool",
		Pool: nil,
		LoadBalancer: LoadBalancerStub{
			HREF: "https://lab-api.cloudscale.ch/v1/load-balancers/bc7b04c9-04ee-471f-b719-26f29c767f6c",
			UUID: "bc7b04c9-04ee-471f-b719-26f29c767f6c",
			Name: "web-lb1",
		},
		Protocol:               "brieftaube",
		ProtocolPort:           80,
		AllowedCIDRs:           []string{},
		TimeoutClientDataMS:    18180,
		TimeoutMemberConnectMS: 1818,
		TimeoutMemberDataMS:    18180,
		CreatedAt:              time.Date(2023, time.Month(10), 13, 13, 28, 38, 672592000, time.UTC),
	}
	expected.Tags = TagMap{}

	if !reflect.DeepEqual(listener, expected) {
		t.Errorf("LoadBalancerListeners.Get\n got=%#v\nwant=%#v", listener, expected)
	}
}
