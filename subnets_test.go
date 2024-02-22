package cloudscale

import (
	"encoding/json"
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

func TestMarshalingOfDNSServersInSubnetUpdateRequest(t *testing.T) {
	testCases := []struct {
		name     string
		request  SubnetUpdateRequest
		expected string
	}{
		{
			name: "one dns server",
			request: SubnetUpdateRequest{
				DNSServers: &[]string{"8.8.8.8"},
			},
			expected: "{\"dns_servers\":[\"8.8.8.8\"]}",
		},
		{
			name: "two dns servers",
			request: SubnetUpdateRequest{
				DNSServers: &[]string{"8.8.8.8", "8.8.4.4"},
			},
			expected: "{\"dns_servers\":[\"8.8.8.8\",\"8.8.4.4\"]}",
		},
		{
			name: "no dns servers",
			request: SubnetUpdateRequest{
				DNSServers: &[]string{},
			},
			expected: "{\"dns_servers\":[]}",
		},
		{
			name: "defaults",
			request: SubnetUpdateRequest{
				DNSServers: &UseCloudscaleDefaults,
			},
			expected: "{\"dns_servers\":null}",
		},
		{
			name: "gateway",
			request: SubnetUpdateRequest{
				GatewayAddress: "192.168.1.1",
			},
			expected: "{\"gateway_address\":\"192.168.1.1\"}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.request)
			if err != nil {
				t.Errorf("Error marshaling JSON: %v", err)
			}
			if actualOutput := string(b); actualOutput != tc.expected {
				t.Errorf("Unexpected JSON output:\nExpected: %s\nActual:   %s", tc.expected, actualOutput)
			}
		})
	}
}

func TestMarshalingOfDNSServersInSubnetSubnetCreateRequest(t *testing.T) {
	testCases := []struct {
		name     string
		request  SubnetCreateRequest
		expected string
	}{
		{
			name: "one dns server",
			request: SubnetCreateRequest{
				DNSServers: &[]string{"8.8.8.8"},
			},
			expected: "{\"dns_servers\":[\"8.8.8.8\"]}",
		},
		{
			name: "two dns servers",
			request: SubnetCreateRequest{
				DNSServers: &[]string{"8.8.8.8", "8.8.4.4"},
			},
			expected: "{\"dns_servers\":[\"8.8.8.8\",\"8.8.4.4\"]}",
		},
		{
			name: "no dns servers",
			request: SubnetCreateRequest{
				DNSServers: &[]string{},
			},
			expected: "{\"dns_servers\":[]}",
		},
		{
			name: "defaults",
			request: SubnetCreateRequest{
				DNSServers: &UseCloudscaleDefaults,
			},
			expected: "{\"dns_servers\":null}",
		},
		{
			name: "gateway",
			request: SubnetCreateRequest{
				GatewayAddress: "192.168.1.1",
			},
			expected: "{\"gateway_address\":\"192.168.1.1\"}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.request)
			if err != nil {
				t.Errorf("Error marshaling JSON: %v", err)
			}
			if actualOutput := string(b); actualOutput != tc.expected {
				t.Errorf("Unexpected JSON output:\nExpected: %s\nActual:   %s", tc.expected, actualOutput)
			}
		})
	}
}
