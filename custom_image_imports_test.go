package cloudscale

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestCustomImageImport_Create(t *testing.T) {
	setup()
	defer teardown()

	CustomImageImportRequest := &CustomImageImportRequest{
		Name: "Test Image",
		TaggedResourceRequest: TaggedResourceRequest{
			TagMap{
				"tag":   "one",
				"other": "tag",
			},
		},
	}

	mux.HandleFunc("/v1/custom-images/import", func(w http.ResponseWriter, r *http.Request) {
		expected := map[string]interface{}{
			"name": "Test Image",
			"tags": map[string]interface{}{
				"tag":   "one",
				"other": "tag",
			},
		}

		var v map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&v)
		if err != nil {
			t.Fatalf("decode json: %v", err)
		}

		if !reflect.DeepEqual(v, expected) {
			t.Errorf("Request body\n got=%#v\nwant=%#v", v, expected)
		}
		jsonStr := `{
  						"href": "https://api.cloudscale.ch/v1/custom-images/import/11111111-1864-4608-853a-0771b6885a3a",
  						"uuid": "11111111-1864-4608-853a-0771b6885a3a",
  						"custom_image": {
  						    "href": "https://api.cloudscale.ch/v1/custom-images/11111111-1864-4608-853a-0771b6885a3a",
  						    "uuid": "11111111-1864-4608-853a-0771b6885a3a",
  						    "name": "my-foo"
  						},
  						"url": "https://example.com/foo.raw",
  						"status": "in_progress",
  						"error_message": "",
  						"tags": {}
					}`
		fmt.Fprintf(w, jsonStr)
	})

	customImageImport, err := client.CustomImageImports.Create(ctx, CustomImageImportRequest)
	if err != nil {
		t.Errorf("CustomImageImport.Create returned error: %v", err)
		return
	}

	expectedUUID := "11111111-1864-4608-853a-0771b6885a3a"
	if id := customImageImport.UUID; id != expectedUUID {
		t.Errorf("expected id '%s', received '%s'", expectedUUID, id)
	}
}

func TestCustomImageImport_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/custom-images/import/6fe39134bf4178747eebc429f82cfafdd08891d4279d0d899bc4012db1db6a15", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"href": "https://api.cloudscale.ch/v1/custom-images/import/11111111-1864-4608-853a-0771b6885a3a",
			"uuid": "11111111-1864-4608-853a-0771b6885a3a",
			"custom_image": {
			    "href": "https://api.cloudscale.ch/v1/custom-images/11111111-1864-4608-853a-0771b6885a3a",
			    "uuid": "11111111-1864-4608-853a-0771b6885a3a",
			    "name": "my-foo"
			},
			"url": "https://example.com/foo.raw",
			"status": "in_progress",
			"error_message": "",
			"tags": {}
	}`)
	})

	objectUser, err := client.CustomImageImports.Get(ctx, "6fe39134bf4178747eebc429f82cfafdd08891d4279d0d899bc4012db1db6a15")
	if err != nil {
		t.Errorf("CustomImageImport.Get returned error: %v", err)
	}

	expected := &CustomImageImport{
		HREF: "https://api.cloudscale.ch/v1/custom-images/import/11111111-1864-4608-853a-0771b6885a3a",
		UUID: "11111111-1864-4608-853a-0771b6885a3a",
		CustomImage: CustomImageStub{
			HREF: "https://api.cloudscale.ch/v1/custom-images/11111111-1864-4608-853a-0771b6885a3a",
			UUID: "11111111-1864-4608-853a-0771b6885a3a",
			Name: "my-foo",
		},
		URL:          "https://example.com/foo.raw",
		Status:       "in_progress",
		ErrorMessage: "",
	}
	expected.Tags = TagMap{}

	if !reflect.DeepEqual(objectUser, expected) {
		t.Errorf("CustomImageImport.Get\n got=%#v\nwant=%#v", objectUser, expected)
	}
}

func TestCustomImageImport_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/custom-images/import", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[{
			"href": "https://api.cloudscale.ch/v1/custom-images/import/11111111-1864-4608-853a-0771b6885a3a",
			"uuid": "11111111-1864-4608-853a-0771b6885a3a",
			"custom_image": {
			    "href": "https://api.cloudscale.ch/v1/custom-images/11111111-1864-4608-853a-0771b6885a3a",
			    "uuid": "11111111-1864-4608-853a-0771b6885a3a",
			    "name": "my-foo"
			},
			"url": "https://example.com/foo.raw",
			"status": "in_progress",
			"error_message": "",
			"tags": {}
}]`)
	})

	objectUsers, err := client.CustomImageImports.List(ctx)
	if err != nil {
		t.Errorf("CustomImageImport.List returned error: %v", err)
	}

	expectedImage := CustomImageImport{
		HREF: "https://api.cloudscale.ch/v1/custom-images/import/11111111-1864-4608-853a-0771b6885a3a",
		UUID: "11111111-1864-4608-853a-0771b6885a3a",
		CustomImage: CustomImageStub{
			HREF: "https://api.cloudscale.ch/v1/custom-images/11111111-1864-4608-853a-0771b6885a3a",
			UUID: "11111111-1864-4608-853a-0771b6885a3a",
			Name: "my-foo",
		},
		URL:          "https://example.com/foo.raw",
		Status:       "in_progress",
		ErrorMessage: "",
	}
	expectedImage.Tags = TagMap{}
	expected := []CustomImageImport{
		expectedImage,
	}
	if !reflect.DeepEqual(objectUsers, expected) {
		t.Errorf("CustomImageImport.List\n got=%#v\nwant=%#v", objectUsers, expected)
	}

}
