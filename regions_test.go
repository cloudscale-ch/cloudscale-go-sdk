package cloudscale

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

const regionsResponse = `[
  {
    "slug": "frn",
    "zones": []
  },
  {
    "slug": "usr",
    "zones": [
      {
        "slug": "usr1"
      }
    ]
  }
]`

func TestRegions_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/regions", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, regionsResponse)
	})

	regions, err := client.Regions.List(ctx)
	if err != nil {
		t.Errorf("Regions.List returned error: %v", err)
	}

	expected := []Region{
		{
			Slug:  "frn",
			Zones: []Zone{},
		},
		{
			Slug: "usr",
			Zones: []Zone{
				{Slug: "usr1"},
			},
		},
	}
	if !reflect.DeepEqual(regions, expected) {
		t.Errorf("Regions.List\n got=%#v\nwant=%#v", regions, expected)
	}

}
