package cloudscale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const flavorsResponse = `[
  {
    "slug": "flex-4-1",
    "name": "Flex-4-1",
    "vcpu_count": 1,
    "memory_gb": 4,
	"gpu": null,
    "zones": [
      {
        "slug": "rma1"
      },
      {
        "slug": "lpg1"
      }
    ]
  },
  {
    "slug": "gpu2-128-16-1-200",
    "name": "GPU2-128-16-1-200",
    "vcpu_count": 16,
    "memory_gb": 128,
	"gpu": {
		"count": 1,
		"name": "NVIDIA RTX PRO 6000 Max-Q",
		"vram_per_gpu_gb": 96,
		"total_vram_gb": 96
	},
    "zones": [
      {
        "slug": "rma1"
      },
      {
        "slug": "lpg1"
      }
    ]
  }
]`

func TestFlavors_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/flavors", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, flavorsResponse)
	})

	flavors, err := client.Flavors.List(ctx)
	if err != nil {
		t.Errorf("Flavors.List returned error: %v", err)
	}

	expected := []Flavor{
		{
			Slug:      "flex-4-1",
			Name:      "Flex-4-1",
			VCPUCount: 1,
			MemoryGB:  4,
			GPU:       nil,
			Zones: []ZoneStub{
				{
					Slug: "rma1",
				},
				{
					Slug: "lpg1",
				},
			},
		},
		{
			Slug:      "gpu2-128-16-1-200",
			Name:      "GPU2-128-16-1-200",
			VCPUCount: 16,
			MemoryGB:  128,
			GPU: &FlavorGPU{
				Name:         "NVIDIA RTX PRO 6000 Max-Q",
				Count:        1,
				VRAMPerGPUGB: 96,
			},
			Zones: []ZoneStub{
				{
					Slug: "rma1",
				},
				{
					Slug: "lpg1",
				},
			},
		},
	}
	if !reflect.DeepEqual(flavors, expected) {
		got, _ := json.MarshalIndent(flavors, "", "  ")
		want, _ := json.MarshalIndent(expected, "", "  ")
		t.Errorf("Flavors.List\n got=%s\nwant=%s", got, want)
	}

}
