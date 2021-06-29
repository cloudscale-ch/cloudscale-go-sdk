package cloudscale

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestCustomImage_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/custom-images/11111111-1864-4608-853a-0771b6885a3a", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{"uuid": "11111111-1864-4608-853a-0771b6885a3a"}`)
	})

	objectUser, err := client.CustomImages.Get(ctx, "11111111-1864-4608-853a-0771b6885a3a")
	if err != nil {
		t.Errorf("CustomImage.Get returned error: %v", err)
	}

	expected := &CustomImage{UUID: "11111111-1864-4608-853a-0771b6885a3a"}
	if !reflect.DeepEqual(objectUser, expected) {
		t.Errorf("CustomImage.Get\n got=%#v\nwant=%#v", objectUser, expected)
	}
}

func TestCustomImage_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/custom-images/11111111-1864-4608-853a-0771b6885a3a", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodDelete)
	})

	err := client.CustomImages.Delete(ctx, "11111111-1864-4608-853a-0771b6885a3a")
	if err != nil {
		t.Errorf("CustomImage.Delete returned error: %v", err)
	}
}

func TestCustomImage_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/custom-images", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `[{"uuid": "11111111-1864-4608-853a-0771b6885a3a"}]`)
	})

	customImages, err := client.CustomImages.List(ctx)
	if err != nil {
		t.Errorf("CustomImage.List returned error: %v", err)
	}

	expected := []CustomImage{{UUID: "11111111-1864-4608-853a-0771b6885a3a"}}
	if !reflect.DeepEqual(customImages, expected) {
		t.Errorf("CustomImage.List\n got=%#v\nwant=%#v", customImages, expected)
	}

}

func TestCustomImage_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/v1/custom-images/11111111-1864-4608-853a-0771b6885a3a", func(w http.ResponseWriter, r *http.Request) {
		testHTTPMethod(t, r, http.MethodPatch)
	})

	userID := "11111111-1864-4608-853a-0771b6885a3a"

	req := &CustomImageRequest{
		Name: "new_name",
	}
	err := client.CustomImages.Update(context.TODO(), userID, req)
	if err != nil {
		t.Errorf("CustomImage.Update returned error: %v", err)
	}
}
