package cloudscale

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

// metadata implements a client for the cloudscale.ch's OpenStack
// metadata API. This API allows a server to inspect information about itself,
// like its server ID.
//
// Documentation for the API is available at:
//
//	https://www.cloudscale.ch/en/api/v1

const (
	maxErrMsgLen = 128 // arbitrary max length for error messages

	defaultTimeout = 2 * time.Second
	defaultPath    = "/openstack/2017-02-22/"
)

var (
	defaultMetadataBaseURL = func() *url.URL {
		u, err := url.Parse("http://169.254.169.254")
		if err != nil {
			panic(err)
		}
		return u
	}()
)

// MetadataClient to interact with cloudscale.ch's OpenStack metadata API, from inside
// a server.
type MetadataClient struct {
	client  *http.Client
	BaseURL *url.URL
}

// NewMetadataClient creates a client for the metadata API.
func NewMetadataClient(httpClient *http.Client) *MetadataClient {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}

	client := &MetadataClient{
		client:  httpClient,
		BaseURL: defaultMetadataBaseURL,
	}
	return client
}

// GetMetadata contains the entire contents of a OpenStack's metadata.
// This method is unique because it returns all of the
// metadata at once, instead of individual metadata items.
func (c *MetadataClient) GetMetadata(ctx context.Context) (*Metadata, error) {
	metadata := new(Metadata)
	err := c.getResource(ctx, "meta_data.json", func(r io.Reader) error {
		return json.NewDecoder(r).Decode(metadata)
	})
	return metadata, err
}

// GetServerID returns the Server's unique identifier. This is
// automatically generated upon Server creation.
func (c *MetadataClient) GetServerID(ctx context.Context) (string, error) {
	metadata, err := c.GetMetadata(ctx)
	if err != nil {
		return "", err
	}
	if metadata.Meta.CloudscaleUUID == "" {
		return "", errors.New("CloudscaleUUID not defined in metadata")
	}
	return metadata.Meta.CloudscaleUUID, nil
}

// GetRawUserData returns the user data that was provided by the user
// during Server creation. User data for cloudscale.ch is a YAML
// Script that is used for cloud-init.
func (c *MetadataClient) GetRawUserData(ctx context.Context) (string, error) {
	var userdata string
	err := c.getResource(ctx, "user_data", func(r io.Reader) error {
		userdataraw, err := io.ReadAll(r)
		userdata = string(userdataraw)
		return err
	})
	return userdata, err
}

func (c *MetadataClient) getResource(ctx context.Context, resource string, decoder func(r io.Reader) error) error {
	url := c.resolve(defaultPath, resource)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return c.makeError(resp)
	}
	return decoder(resp.Body)
}

func (c *MetadataClient) makeError(resp *http.Response) error {
	body, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrMsgLen))
	if len(body) >= maxErrMsgLen {
		body = append(body[:maxErrMsgLen], []byte("... (elided)")...)
	} else if len(body) == 0 {
		body = []byte(resp.Status)
	}
	return fmt.Errorf("unexpected response from metadata API, status %d: %s",
		resp.StatusCode, string(body))
}

func (c *MetadataClient) resolve(basePath string, resource ...string) string {
	dupe := *c.BaseURL
	dupe.Path = path.Join(append([]string{basePath}, resource...)...)
	return dupe.String()
}
