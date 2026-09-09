package cloudscale

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestRouters_Create(t *testing.T) {
	setup()
	defer teardown()

	routerRequest := &RouterCreateRequest{
		Name:            "gw",
		InternetGateway: true,
	}

	mux.HandleFunc("/v1/routers", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		expected := map[string]any{
			"name":             "gw",
			"internet_gateway": true,
		}

		var v map[string]any
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		_, _ = fmt.Fprint(w, `{"uuid": "42cec963-fcd2-482f-bdb6-24461b2d47b1"}`)
	})

	router, err := client.Routers.Create(ctx, routerRequest)
	if err != nil {
		t.Errorf("Routers.Create returned error: %v", err)
	}

	if id := router.UUID; id != "42cec963-fcd2-482f-bdb6-24461b2d47b1" {
		t.Errorf("expected id '42cec963-fcd2-482f-bdb6-24461b2d47b1', received '%s'", id)
	}
}

func TestRouters_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/routers/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPatch)
	})

	routerID := "cfde831a-4e87-4a75-960f-89b0148aa2cc"

	req := &RouterUpdateRequest{
		Name: "new-router-name",
	}
	err := client.Routers.Update(context.TODO(), routerID, req)
	if err != nil {
		t.Errorf("ObjectsUser.Update returned error: %v", err)
	}
}

func TestRouters_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/routers/cfde831a-4e87-4a75-960f-89b0148aa2cc", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `{
            "href": "https://api.cloudscale.ch/v1/routers/cfde831a-4e87-4a75-960f-89b0148aa2cc",
            "uuid": "cfde831a-4e87-4a75-960f-89b0148aa2cc",
            "name": "gw",
            "zone": {"slug": "lpg1"},
            "created_at": "2019-05-27T16:45:32.241824Z",
            "status": "up",
            "internet_gateway": true,
            "internet_gateway_addresses": [
              {
                "address": "203.0.113.1",
                "subnet": {
                  "href": "https://api.cloudscale.ch/v1/subnets/8a04e678-4f1c-4d5f-9e40-8f0eaf1d0e0d",
                  "cidr": "203.0.113.0/24",
                  "uuid": "8a04e678-4f1c-4d5f-9e40-8f0eaf1d0e0d"
                },
                "version": 4,
                "reverse_ptr": "203-0-113-1.cust.example.com"
              },
			  {
                "address": "2001:db8::1",
                "subnet": {
                  "href": "https://api.cloudscale.ch/v1/subnets/9204e678-4f1c-4d5f-9e40-8f0eaf1d0eaa",
                  "cidr": "2001:db8::/32",
                  "uuid": "9204e678-4f1c-4d5f-9e40-8f0eaf1d0eaa"
                },
                "version": 6,
                "reverse_ptr": "203-0-113-1.cust.example.com"
              }
            ],
            "interfaces": [
              {
                "uuid": "1e0c6f9c-9f0d-4d1b-9f0d-8f0eaf1d0e0d",
                "network": {
                  "href": "https://api.cloudscale.ch/v1/networks/7f0eaf1d-0e0d-4d5f-9e40-8f0eaf1d0e0d",
                  "name": "my-network",
                  "uuid": "7f0eaf1d-0e0d-4d5f-9e40-8f0eaf1d0e0d"
                },
                "addresses": [
                  {
                    "address": "10.0.0.1",
                    "subnet": {
                      "href": "https://api.cloudscale.ch/v1/subnets/3d6ca1f4-5aea-41f5-b724-0f3054b60e85",
                      "cidr": "10.0.0.0/24",
                      "uuid": "3d6ca1f4-5aea-41f5-b724-0f3054b60e85"
                    },
                    "version": 4,
                    "reverse_ptr": null
                  }
                ],
                "type": "private",
                "mac_address": "00:00:5e:00:53:ab"
              }
            ]
          }`)
	})

	router, err := client.Routers.Get(ctx, "cfde831a-4e87-4a75-960f-89b0148aa2cc")
	if err != nil {
		t.Errorf("Routers.Get returned error: %v", err)
	}

	expected := &Router{
		ZonalResource: ZonalResource{
			Zone: ZoneStub{Slug: "lpg1"},
		},
		HREF:            "https://api.cloudscale.ch/v1/routers/cfde831a-4e87-4a75-960f-89b0148aa2cc",
		UUID:            "cfde831a-4e87-4a75-960f-89b0148aa2cc",
		Name:            "gw",
		CreatedAt:       time.Date(2019, time.Month(5), 27, 16, 45, 32, 241824000, time.UTC),
		Status:          "up",
		InternetGateway: true,
		InternetGatewayAddresses: []IPAddress{
			{
				Address: "203.0.113.1",
				Subnet: SubnetStub{
					HREF: "https://api.cloudscale.ch/v1/subnets/8a04e678-4f1c-4d5f-9e40-8f0eaf1d0e0d",
					CIDR: "203.0.113.0/24",
					UUID: "8a04e678-4f1c-4d5f-9e40-8f0eaf1d0e0d",
				},
				Version:    4,
				ReversePTR: new("203-0-113-1.cust.example.com"),
			},
			{
				Address: "2001:db8::1",
				Subnet: SubnetStub{
					HREF: "https://api.cloudscale.ch/v1/subnets/9204e678-4f1c-4d5f-9e40-8f0eaf1d0eaa",
					CIDR: "2001:db8::/32",
					UUID: "9204e678-4f1c-4d5f-9e40-8f0eaf1d0eaa",
				},
				Version:    6,
				ReversePTR: new("203-0-113-1.cust.example.com"),
			},
		},
		Interfaces: []RouterInterface{
			{
				UUID: "1e0c6f9c-9f0d-4d1b-9f0d-8f0eaf1d0e0d",
				Network: NetworkStub{
					HREF: "https://api.cloudscale.ch/v1/networks/7f0eaf1d-0e0d-4d5f-9e40-8f0eaf1d0e0d",
					Name: "my-network",
					UUID: "7f0eaf1d-0e0d-4d5f-9e40-8f0eaf1d0e0d",
				},
				Addresses: []IPAddress{
					{
						Address: "10.0.0.1",
						Subnet: SubnetStub{
							HREF: "https://api.cloudscale.ch/v1/subnets/3d6ca1f4-5aea-41f5-b724-0f3054b60e85",
							CIDR: "10.0.0.0/24",
							UUID: "3d6ca1f4-5aea-41f5-b724-0f3054b60e85",
						},
						Version:    4,
						ReversePTR: nil,
					},
				},
				Type:       "private",
				MACAddress: "00:00:5e:00:53:ab",
			},
		},
	}

	if !reflect.DeepEqual(router, expected) {
		t.Errorf("Routers.Get\n got=%#v\nwant=%#v", router, expected)
	}
}

func TestRouters_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/routers", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		_, _ = fmt.Fprint(w, `[{"uuid": "47cec963-fcd2-482f-bdb6-24461b2d47b1"}]`)
	})

	routers, err := client.Routers.List(ctx)
	if err != nil {
		t.Errorf("Routers.List returned error: %v", err)
	}

	expected := []Router{{UUID: "47cec963-fcd2-482f-bdb6-24461b2d47b1"}}
	if !reflect.DeepEqual(routers, expected) {
		t.Errorf("Routers.List\n got=%#v\nwant=%#v", routers, expected)
	}
}

func TestRouters_CreateInterface(t *testing.T) {
	setup()
	defer teardown()

	createReq := CreateInterfaceRequest{
		Network: "7f0eaf1d-0e0d-4d5f-9e40-8f0eaf1d0e0d",
		Addresses: []CreateAddressRequest{
			{
				Subnet:  "3d6ca1f4-5aea-41f5-b724-0f3054b60e85",
				Address: "10.0.0.1",
			},
		},
	}

	mux.HandleFunc("/v1/routers/cfde831a-4e87-4a75-960f-89b0148aa2cc/interfaces", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPost)

		expected := map[string]any{
			"network": "7f0eaf1d-0e0d-4d5f-9e40-8f0eaf1d0e0d",
			"addresses": []any{
				map[string]any{
					"subnet":  "3d6ca1f4-5aea-41f5-b724-0f3054b60e85",
					"address": "10.0.0.1",
				},
			},
		}

		var v map[string]any
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}

		_, _ = fmt.Fprint(w, `{
            "uuid": "1e0c6f9c-9f0d-4d1b-9f0d-8f0eaf1d0e0d",
            "network": {
              "href": "https://api.cloudscale.ch/v1/networks/7f0eaf1d-0e0d-4d5f-9e40-8f0eaf1d0e0d",
              "name": "my-network",
              "uuid": "7f0eaf1d-0e0d-4d5f-9e40-8f0eaf1d0e0d"
            },
            "addresses": [
              {
                "address": "10.0.0.1",
                "subnet": {
                  "href": "https://api.cloudscale.ch/v1/subnets/3d6ca1f4-5aea-41f5-b724-0f3054b60e85",
                  "cidr": "10.0.0.0/24",
                  "uuid": "3d6ca1f4-5aea-41f5-b724-0f3054b60e85"
                },
                "version": 4,
                "reverse_ptr": null
              }
            ],
            "type": "vip",
            "mac_address": "aa:bb:cc:dd:ee:ff"
          }`)
	})

	iface, err := client.Routers.CreateInterface(ctx, "cfde831a-4e87-4a75-960f-89b0148aa2cc", createReq)
	if err != nil {
		t.Errorf("Routers.CreateInterface returned error: %v", err)
	}

	expected := &RouterInterface{
		UUID: "1e0c6f9c-9f0d-4d1b-9f0d-8f0eaf1d0e0d",
		Network: NetworkStub{
			HREF: "https://api.cloudscale.ch/v1/networks/7f0eaf1d-0e0d-4d5f-9e40-8f0eaf1d0e0d",
			Name: "my-network",
			UUID: "7f0eaf1d-0e0d-4d5f-9e40-8f0eaf1d0e0d",
		},
		Addresses: []IPAddress{
			{
				Address: "10.0.0.1",
				Subnet: SubnetStub{
					HREF: "https://api.cloudscale.ch/v1/subnets/3d6ca1f4-5aea-41f5-b724-0f3054b60e85",
					CIDR: "10.0.0.0/24",
					UUID: "3d6ca1f4-5aea-41f5-b724-0f3054b60e85",
				},
				Version:    4,
				ReversePTR: nil,
			},
		},
		Type:       "vip",
		MACAddress: "aa:bb:cc:dd:ee:ff",
	}

	if !reflect.DeepEqual(iface, expected) {
		t.Errorf("Routers.CreateInterface\n got=%#v\nwant=%#v", iface, expected)
	}
}

func TestRouters_DeleteInterface(t *testing.T) {
	setup()
	defer teardown()

	called := false
	mux.HandleFunc("/v1/routers/cfde831a-4e87-4a75-960f-89b0148aa2cc/interfaces/1e0c6f9c-9f0d-4d1b-9f0d-8f0eaf1d0e0d", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)

		called = true
	})

	err := client.Routers.DeleteInterface(ctx, "cfde831a-4e87-4a75-960f-89b0148aa2cc", "1e0c6f9c-9f0d-4d1b-9f0d-8f0eaf1d0e0d")
	if err != nil {
		t.Errorf("Routers.DeleteInterface returned error: %v", err)
	}

	if !called {
		t.Error("expected delete_interface endpoint to be called")
	}
}
