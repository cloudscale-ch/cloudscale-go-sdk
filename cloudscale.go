package cloudscale

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	libraryVersion = "1.0"
	defaultBaseURL = "https://api.cloudscale.ch/"
	userAgent      = "cloudscale/" + libraryVersion
	mediaType      = "application/json"

	headerRateLimit     = "RateLimit-Limit"
	headerRateRemaining = "RateLimit-Remaining"
	headerRateReset     = "RateLimit-Reset"
)

// Client manages communication with DigitalOcean V2 API.
type Client struct {
	// HTTP client used to communicate with the DO API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// User agent for client
	UserAgent string
}

// NewClient returns a new DigitalOcean API client.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}

	return c
}

func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

type ErrorResponse struct {
	Message map[string]string
}

func (r *ErrorResponse) Error() string {
	err := ""
	for key, value := range r.Message {
		err = fmt.Sprintf("%s: %s", key, value)
	}
	return err
}
