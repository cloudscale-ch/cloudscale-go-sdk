package cloudscale

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

const (
	libraryVersion = "v6.0.0"
	defaultBaseURL = "https://api.cloudscale.ch/"
	userAgent      = "cloudscale/" + libraryVersion
	mediaType      = "application/json"
)

// Client manages communication with CloudScale API.
type Client struct {

	// HTTP client used to communicate with the CloudScale API.
	client *http.Client

	// Base URL for API requests.
	BaseURL *url.URL

	// Authentication token
	AuthToken string

	// User agent for client
	UserAgent string

	Regions                    RegionService
	Servers                    ServerService
	Volumes                    VolumeService
	Networks                   NetworkService
	Subnets                    SubnetService
	FloatingIPs                FloatingIPsService
	ServerGroups               ServerGroupService
	ObjectsUsers               ObjectsUsersService
	CustomImages               CustomImageService
	CustomImageImports         CustomImageImportsService
	LoadBalancers              LoadBalancerService
	LoadBalancerPools          LoadBalancerPoolService
	LoadBalancerPoolMembers    LoadBalancerPoolMemberService
	LoadBalancerListeners      LoadBalancerListenerService
	LoadBalancerHealthMonitors LoadBalancerHealthMonitorService
	Metrics                    MetricsService
}

// NewClient returns a new CloudScale API client.
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	// To allow more complicated testing we allow changing the cloudscale.ch
	// URL.
	defaultURL := os.Getenv("CLOUDSCALE_API_URL")

	if defaultURL == "" {
		defaultURL = defaultBaseURL
	}
	baseURL, _ := url.Parse(defaultURL)

	c := &Client{client: httpClient, BaseURL: baseURL, UserAgent: userAgent}
	c.Regions = RegionServiceOperations{client: c}
	c.Servers = ServerServiceOperations{
		GenericServiceOperations: GenericServiceOperations[Server, ServerRequest, ServerUpdateRequest]{
			client: c,
			path:   serverBasePath,
		},
		client: c,
	}
	c.Networks = GenericServiceOperations[Network, NetworkCreateRequest, NetworkUpdateRequest]{
		client: c,
		path:   networkBasePath,
	}
	c.Subnets = GenericServiceOperations[Subnet, SubnetCreateRequest, SubnetUpdateRequest]{
		client: c,
		path:   subnetBasePath,
	}
	c.FloatingIPs = GenericServiceOperations[FloatingIP, FloatingIPCreateRequest, FloatingIPUpdateRequest]{
		client: c,
		path:   floatingIPsBasePath,
	}
	c.Volumes = GenericServiceOperations[Volume, VolumeRequest, VolumeRequest]{
		client: c,
		path:   volumeBasePath,
	}
	c.ServerGroups = GenericServiceOperations[ServerGroup, ServerGroupRequest, ServerGroupRequest]{
		client: c,
		path:   serverGroupsBasePath,
	}
	c.ObjectsUsers = GenericServiceOperations[ObjectsUser, ObjectsUserRequest, ObjectsUserRequest]{
		client: c,
		path:   objectsUsersBasePath,
	}
	c.CustomImages = GenericServiceOperations[CustomImage, CustomImageRequest, CustomImageRequest]{
		client: c,
		path:   customImagesBasePath,
	}
	c.CustomImageImports = GenericServiceOperations[CustomImageImport, CustomImageImportRequest, CustomImageImportRequest]{
		client: c,
		path:   customImageImportsBasePath,
	}
	c.LoadBalancers = GenericServiceOperations[LoadBalancer, LoadBalancerRequest, LoadBalancerRequest]{
		client: c,
		path:   loadBalancerBasePath,
	}
	c.LoadBalancerPools = GenericServiceOperations[LoadBalancerPool, LoadBalancerPoolRequest, LoadBalancerPoolRequest]{
		client: c,
		path:   loadBalancerPoolBasePath,
	}
	c.LoadBalancerPoolMembers = LoadBalancerPoolMemberServiceOperations{
		client: c,
	}
	c.LoadBalancerListeners = GenericServiceOperations[LoadBalancerListener, LoadBalancerListenerRequest, LoadBalancerListenerRequest]{
		client: c,
		path:   loadBalancerListenerBasePath,
	}
	c.LoadBalancerHealthMonitors = GenericServiceOperations[LoadBalancerHealthMonitor, LoadBalancerHealthMonitorRequest, LoadBalancerHealthMonitorRequest]{
		client: c,
		path:   loadBalancerHealthMonitorBasePath,
	}
	c.Metrics = MetricsServiceOperations{client: c}

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

	if len(c.AuthToken) != 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AuthToken))
	}

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) error {

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	err = CheckResponse(resp)
	if err != nil {
		return err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return err
			}
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err != nil {
				return err
			}
		}
	}

	return err
}

func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	data, err := ioutil.ReadAll(r.Body)
	res := map[string]string{}
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, &res)
		if err != nil {
			return err
		}
	}

	return &ErrorResponse{
		StatusCode: r.StatusCode,
		Message:    res,
	}
}

type ErrorResponse struct {
	StatusCode int
	Message    map[string]string
}

func (r *ErrorResponse) Error() string {
	err := ""
	for key, value := range r.Message {
		err = fmt.Sprintf("%s: %s", key, value)
	}
	return err
}

type ListRequestModifier func(r *http.Request)
