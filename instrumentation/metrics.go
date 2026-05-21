package instrumentation

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/cloudscale-ch/cloudscale-go-sdk/v9"
)

// metricsCollector holds the Prometheus metrics for cloudscale API calls.
type metricsCollector struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	inFlight        prometheus.Gauge
}

// registerOrReuseCounterVec tries to register a CounterVec and, if it is
// already registered, returns the existing one.
func registerOrReuseCounterVec(reg prometheus.Registerer, opts prometheus.CounterOpts, labelNames []string) *prometheus.CounterVec {
	cv := prometheus.NewCounterVec(opts, labelNames)
	if err := reg.Register(cv); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			return are.ExistingCollector.(*prometheus.CounterVec)
		}
		panic(err)
	}
	return cv
}

// registerOrReuseHistogramVec tries to register a HistogramVec and, if it is
// already registered, returns the existing one.
func registerOrReuseHistogramVec(reg prometheus.Registerer, opts prometheus.HistogramOpts, labelNames []string) *prometheus.HistogramVec {
	hv := prometheus.NewHistogramVec(opts, labelNames)
	if err := reg.Register(hv); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			return are.ExistingCollector.(*prometheus.HistogramVec)
		}
		panic(err)
	}
	return hv
}

// registerOrReuseGauge tries to register a Gauge and, if it is already
// registered, returns the existing one.
func registerOrReuseGauge(reg prometheus.Registerer, opts prometheus.GaugeOpts) prometheus.Gauge {
	g := prometheus.NewGauge(opts)
	if err := reg.Register(g); err != nil {
		var are prometheus.AlreadyRegisteredError
		if errors.As(err, &are) {
			return are.ExistingCollector.(prometheus.Gauge)
		}
		panic(err)
	}
	return g
}

// NewMetricsCollector creates a metricsCollector backed by the given registry.
func NewMetricsCollector(reg prometheus.Registerer, subsystem string) *metricsCollector {
	return &metricsCollector{
		requestsTotal: registerOrReuseCounterVec(reg, prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      "requests_total",
			Help:      "Total API requests by endpoint and response code.",
		}, []string{"method", "endpoint", "status"}),
		requestDuration: registerOrReuseHistogramVec(reg, prometheus.HistogramOpts{
			Subsystem: subsystem,
			Name:      "request_duration_seconds",
			Help:      "Request latency distribution.",
			Buckets:   []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		}, []string{"method", "endpoint"}),
		inFlight: registerOrReuseGauge(reg, prometheus.GaugeOpts{
			Subsystem: subsystem,
			Name:      "in_flight_requests",
			Help:      "Number of requests currently in flight.",
		}),
	}
}

// metricsTransport wraps an http.RoundTripper with Prometheus metrics.
type metricsTransport struct {
	next    http.RoundTripper
	metrics *metricsCollector
}

func (t *metricsTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.metrics.inFlight.Inc()
	defer t.metrics.inFlight.Dec()

	start := time.Now()
	resp, err := t.next.RoundTrip(req)
	duration := time.Since(start)

	endpoint := cloudscale.OperationPath(req.Context())
	if endpoint == "" {
		endpoint = "unknown"
	}
	method := req.Method
	status := "error"

	if resp != nil {
		status = strconv.Itoa(resp.StatusCode)
	}

	t.metrics.requestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
	t.metrics.requestsTotal.WithLabelValues(method, endpoint, status).Inc()

	return resp, err
}
