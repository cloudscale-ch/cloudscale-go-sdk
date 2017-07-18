package cloudscale

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"
)

var (
	ctx = context.TODO()
)

func TestNewClient(t *testing.T) {
	c := NewClient(nil)

	if c.BaseURL == nil || c.BaseURL.String() != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, expected %v", c.BaseURL, defaultBaseURL)
	}

	if c.UserAgent != userAgent {
		t.Errorf("NewClient UserAgent = %v, expected %v", c.UserAgent, userAgent)
	}

}

func TestNewRequest(t *testing.T) {
	c := NewClient(nil)

	inURL, outURL := "/foo", defaultBaseURL+"foo"
	inBody, outBody := &ServerRequest{Name: "l"},
		`{"name":"l","flavor":"","image":"","volume_size_gb":"","ssh_keys":null}`+
			"\n"
	req, _ := c.NewRequest(ctx, http.MethodGet, inURL, inBody)

	// test relative URL was expanded
	if req.URL.String() != outURL {
		t.Errorf("NewRequest(%v) URL = %v, expected %v", inURL, req.URL, outURL)
	}

	// test body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if string(body) != outBody {
		t.Errorf("NewRequest(%v)Body = %v, expected %v", inBody, string(body), outBody)
	}

	// test default user-agent is attached to the request
	userAgent := req.Header.Get("User-Agent")
	if c.UserAgent != userAgent {
		t.Errorf("NewRequest() User-Agent = %v, expected %v", userAgent, c.UserAgent)
	}
}

func TestErrorResponse_Error(t *testing.T) {
	err := ErrorResponse{Message: map[string]string{"name": "This field may not be blank."}}
	if err.Error() == "" {
		t.Errorf("Expected non-empty ErrorResponse.Error()")
	}
}
