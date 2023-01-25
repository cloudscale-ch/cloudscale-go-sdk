package cloudscale

import (
	"time"
)

const loadBalancerHealthMonitorBasePath = "v1/load-balancers/health-monitors"

type LoadBalancerHealthMonitor struct {
	TaggedResource
	// Just use omitempty everywhere. This makes it easy to use restful. Errors
	// will be coming from the API if something is disabled.
	HREF           string                              `json:"href,omitempty"`
	UUID           string                              `json:"uuid,omitempty"`
	Pool           LoadBalancerPoolStub                `json:"pool,omitempty"`
	Delay          int                                 `json:"delay,omitempty"`
	Timeout        int                                 `json:"timeout,omitempty"`
	MaxRetries     int                                 `json:"max_retries,omitempty"`
	MaxRetriesDown int                                 `json:"max_retries_down,omitempty"`
	Type           string                              `json:"type,omitempty"`
	HTTP           LoadBalancerHealthMonitorHTTPOption `json:"http,omitempty"`
	CreatedAt      time.Time                           `json:"created_at,omitempty"`
}

type LoadBalancerHealthMonitorHTTPOption struct {
	ExpectedCodes []string `json:"expected_codes,omitempty"`
	Method        string   `json:"method,omitempty"`
	UrlPath       string   `json:"url_path,omitempty"`
	Version       string   `json:"version,omitempty"`
	DomainName    *string  `json:"domain_name,omitempty"`
}

type LoadBalancerHealthMonitorRequest struct {
	TaggedResourceRequest
	Pool           string                               `json:"pool,omitempty"`
	Delay          int                                  `json:"delay,omitempty"`
	Timeout        int                                  `json:"timeout,omitempty"`
	MaxRetries     int                                  `json:"max_retries,omitempty"`
	MaxRetriesDown int                                  `json:"max_retries_down,omitempty"`
	Type           string                               `json:"type,omitempty"`
	HTTP           LoadBalancerHealthMonitorHTTPRequest `json:"http,omitempty"`
}

type LoadBalancerHealthMonitorHTTPRequest struct {
	ExpectedCodes []string `json:"expected_codes,omitempty"`
	Method        string   `json:"method,omitempty"`
	UrlPath       string   `json:"url_path,omitempty"`
	Version       string   `json:"version,omitempty"`
	DomainName    *string  `json:"domain_name,omitempty"`
}

type LoadBalancerHealthMonitorService interface {
	GenericCreateService[LoadBalancerHealthMonitor, LoadBalancerHealthMonitorRequest, LoadBalancerHealthMonitorRequest]
	GenericGetService[LoadBalancerHealthMonitor, LoadBalancerHealthMonitorRequest, LoadBalancerHealthMonitorRequest]
	GenericListService[LoadBalancerHealthMonitor, LoadBalancerHealthMonitorRequest, LoadBalancerHealthMonitorRequest]
	GenericUpdateService[LoadBalancerHealthMonitor, LoadBalancerHealthMonitorRequest, LoadBalancerHealthMonitorRequest]
	GenericDeleteService[LoadBalancerHealthMonitor, LoadBalancerHealthMonitorRequest, LoadBalancerHealthMonitorRequest]
}
