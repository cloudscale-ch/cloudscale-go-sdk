package instrumentation

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/cloudscale-ch/cloudscale-go-sdk/v9"
)

func newTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()
	s := httptest.NewServer(handler)
	t.Cleanup(s.Close)
	return s
}

func statusHandler(status int) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
	}
}

func newInstrumentedClient(reg prometheus.Registerer, tracer trace.Tracer) *http.Client {
	return &http.Client{
		Transport: InstrumentedTransport(http.DefaultTransport, Options{
			Subsystem:          "cloudscale",
			PrometheusRegistry: reg,
			Tracer:             tracer,
		}),
	}
}

func doRequest(t *testing.T, c *http.Client, method, url, opPath string) (*http.Response, error) {
	t.Helper()
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	if opPath != "" {
		req = req.WithContext(cloudscale.WithOperationPath(req.Context(), opPath))
	}
	return c.Do(req)
}

func metricFamily(t *testing.T, reg *prometheus.Registry, name string) *dto.MetricFamily {
	t.Helper()
	families, err := reg.Gather()
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range families {
		if f.GetName() == name {
			return f
		}
	}
	t.Fatalf("metric %s not found", name)
	return nil
}

func labelsOf(m *dto.Metric) map[string]string {
	out := map[string]string{}
	for _, l := range m.Label {
		out[l.GetName()] = l.GetValue()
	}
	return out
}

func attrsOf(span sdktrace.ReadOnlySpan) map[string]attribute.Value {
	out := map[string]attribute.Value{}
	for _, kv := range span.Attributes() {
		out[string(kv.Key)] = kv.Value
	}
	return out
}

func newRecordingTracer(t *testing.T) (trace.Tracer, *tracetest.SpanRecorder) {
	t.Helper()
	sr := tracetest.NewSpanRecorder()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(sr))
	return tp.Tracer("test"), sr
}

func mustDoRequest(t *testing.T, c *http.Client, method, url, opPath string) {
	t.Helper()
	resp, err := doRequest(t, c, method, url, opPath)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}

func singleSpan(t *testing.T, sr *tracetest.SpanRecorder) sdktrace.ReadOnlySpan {
	t.Helper()
	spans := sr.Ended()
	if len(spans) != 1 {
		t.Fatalf("expected 1 ended span, got %d", len(spans))
	}
	return spans[0]
}

func TestInstrumentedTransport_StatusLabels(t *testing.T) {
	cases := []struct {
		name           string
		handler        http.HandlerFunc // nil => transport error (server closed before request)
		method         string
		opPath         string
		expectedStatus string
	}{
		{
			name:           "success",
			handler:        statusHandler(http.StatusOK),
			method:         http.MethodGet,
			opPath:         "v1/servers/:id",
			expectedStatus: "200",
		},
		{
			name:           "server error",
			handler:        statusHandler(http.StatusInternalServerError),
			method:         http.MethodPost,
			opPath:         "v1/servers",
			expectedStatus: "500",
		},
		{
			name:           "transport error",
			handler:        nil,
			method:         http.MethodGet,
			opPath:         "v1/servers",
			expectedStatus: "error",
		},
		{
			name:           "no operation path falls back to unknown",
			handler:        statusHandler(http.StatusOK),
			method:         http.MethodGet,
			opPath:         "",
			expectedStatus: "200",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			reg := prometheus.NewRegistry()
			client := newInstrumentedClient(reg, nil)

			var url string
			if tc.handler != nil {
				url = newTestServer(t, tc.handler).URL + "/" + tc.opPath
			} else {
				s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
				s.Close()
				client.Timeout = 100 * time.Millisecond
				url = s.URL + "/" + tc.opPath
			}

			resp, err := doRequest(t, client, tc.method, url, tc.opPath)
			if tc.expectedStatus == "error" {
				if err == nil {
					t.Fatal("expected error")
				}
			} else {
				if err != nil {
					t.Fatal(err)
				}
				resp.Body.Close()
			}

			wantEndpoint := tc.opPath
			if wantEndpoint == "" {
				wantEndpoint = "unknown"
			}

			total := metricFamily(t, reg, "cloudscale_requests_total")
			if len(total.Metric) != 1 {
				t.Fatalf("expected 1 metric, got %d", len(total.Metric))
			}
			labels := labelsOf(total.Metric[0])
			if labels["method"] != tc.method {
				t.Errorf("method=%s, want %s", labels["method"], tc.method)
			}
			if labels["endpoint"] != wantEndpoint {
				t.Errorf("endpoint=%s, want %s", labels["endpoint"], wantEndpoint)
			}
			if labels["status"] != tc.expectedStatus {
				t.Errorf("status=%s, want %s", labels["status"], tc.expectedStatus)
			}

			// All three metric families should be populated for every case.
			if tc.expectedStatus == "200" {
				duration := metricFamily(t, reg, "cloudscale_request_duration_seconds")
				if got := duration.Metric[0].Histogram.GetSampleCount(); got != 1 {
					t.Errorf("histogram SampleCount=%d, want 1", got)
				}
				inFlight := metricFamily(t, reg, "cloudscale_in_flight_requests")
				if got := inFlight.Metric[0].Gauge.GetValue(); got != 0 {
					t.Errorf("in_flight_requests=%v, want 0", got)
				}
			}
		})
	}
}

func TestInstrumentedTransport_Concurrent(t *testing.T) {
	const n = 10

	reg := prometheus.NewRegistry()

	// `arrived` and `release` are used for blocking requests and counting in flight requests metric:
	// 1. Each request is added to `arrived` once, so we can wait until all requests have been done.
	// 2. in flight requests are counted and should equal `n`
	// 3. release channel is closed, so requests can finish
	// 4. we wait on the WaitGroup which signals that all requests have finished before checking the other metrics

	arrived := make(chan struct{}, n)
	release := make(chan struct{})
	server := newTestServer(t, func(w http.ResponseWriter, _ *http.Request) {
		arrived <- struct{}{}
		<-release
		w.WriteHeader(http.StatusOK)
	})

	client := newInstrumentedClient(reg, nil)

	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			resp, err := doRequest(t, client, http.MethodGet, server.URL+"/v1/servers", "v1/servers")
			if err != nil {
				t.Errorf("request failed: %v", err)
				return
			}
			resp.Body.Close()
		}()
	}

	for i := 0; i < n; i++ {
		<-arrived
	}

	// Peak: all n requests are blocked in the handler.
	inFlight := metricFamily(t, reg, "cloudscale_in_flight_requests")
	if got := inFlight.Metric[0].Gauge.GetValue(); got != float64(n) {
		t.Errorf("peak in_flight=%v, want %d", got, n)
	}

	close(release)
	wg.Wait()

	inFlight = metricFamily(t, reg, "cloudscale_in_flight_requests")
	if got := inFlight.Metric[0].Gauge.GetValue(); got != 0 {
		t.Errorf("post-completion in_flight=%v, want 0", got)
	}

	total := metricFamily(t, reg, "cloudscale_requests_total")
	if got := total.Metric[0].Counter.GetValue(); got != float64(n) {
		t.Errorf("requests_total=%v, want %d", got, n)
	}

	duration := metricFamily(t, reg, "cloudscale_request_duration_seconds")
	if got := duration.Metric[0].Histogram.GetSampleCount(); got != uint64(n) {
		t.Errorf("histogram SampleCount=%d, want %d", got, n)
	}
}

func TestInstrumentedTransport_NilOptions(t *testing.T) {
	transport := InstrumentedTransport(http.DefaultTransport, Options{})
	if transport != http.DefaultTransport {
		t.Fatal("expected default transport when options are empty")
	}
}

func TestInstrumentedTransport_MetricsOnly(t *testing.T) {
	reg := prometheus.NewRegistry()
	server := newTestServer(t, statusHandler(http.StatusTeapot))
	client := newInstrumentedClient(reg, nil)

	mustDoRequest(t, client, http.MethodGet, server.URL+"/v1/flavors", "v1/flavors")

	for _, name := range []string{
		"cloudscale_requests_total",
		"cloudscale_request_duration_seconds",
		"cloudscale_in_flight_requests",
	} {
		// verifies metric is available
		metricFamily(t, reg, name)
	}
}

func TestInstrumentedTransport_DefaultSubsystem(t *testing.T) {
	reg := prometheus.NewRegistry()
	server := newTestServer(t, statusHandler(http.StatusOK))
	client := &http.Client{
		Transport: InstrumentedTransport(http.DefaultTransport, Options{
			// Subsystem intentionally omitted — should default to "cloudscale".
			PrometheusRegistry: reg,
		}),
	}

	mustDoRequest(t, client, http.MethodGet, server.URL+"/v1/flavors", "v1/flavors")

	metricFamily(t, reg, "cloudscale_requests_total")
}

func TestInstrumentedTransport_SharedRegistry(t *testing.T) {
	reg := prometheus.NewRegistry()
	server := newTestServer(t, statusHandler(http.StatusOK))

	// Two independent transports sharing the same registry must reuse the same
	// underlying collectors — both increments should land on a single series.
	for _, c := range []*http.Client{newInstrumentedClient(reg, nil), newInstrumentedClient(reg, nil)} {
		mustDoRequest(t, c, http.MethodGet, server.URL+"/v1/servers", "v1/servers")
	}

	total := metricFamily(t, reg, "cloudscale_requests_total")
	if len(total.Metric) != 1 {
		t.Fatalf("expected 1 series (shared collector), got %d", len(total.Metric))
	}
	if got := total.Metric[0].Counter.GetValue(); got != 2 {
		t.Errorf("counter=%v, want 2", got)
	}
}

func TestInstrumentedTransport_CombinedMetricsAndTracing(t *testing.T) {
	reg := prometheus.NewRegistry()
	tracer := noop.NewTracerProvider().Tracer("test")

	server := newTestServer(t, statusHandler(http.StatusOK))
	client := newInstrumentedClient(reg, tracer)

	mustDoRequest(t, client, http.MethodGet, server.URL+"/v1/servers/123", "v1/servers/:id")

	total := metricFamily(t, reg, "cloudscale_requests_total")
	labels := labelsOf(total.Metric[0])
	if labels["endpoint"] != "v1/servers/:id" {
		t.Errorf("endpoint=%s, want v1/servers/:id", labels["endpoint"])
	}
}

func TestInstrumentedTransport_SpanAttributes(t *testing.T) {
	cases := []struct {
		name         string
		opPath       string
		wantSpanName string
		wantEndpoint bool
	}{
		{
			name:         "with operation path",
			opPath:       "v1/servers/:id",
			wantSpanName: "GET v1/servers/:id",
			wantEndpoint: true,
		},
		{
			name:         "without operation path",
			opPath:       "",
			wantSpanName: "GET",
			wantEndpoint: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tracer, sr := newRecordingTracer(t)

			server := newTestServer(t, statusHandler(http.StatusOK))
			client := newInstrumentedClient(nil, tracer)

			reqURL := server.URL + "/some/path?bucket_name=secret&objects_user_id=user-uuid"
			mustDoRequest(t, client, http.MethodGet, reqURL, tc.opPath)

			span := singleSpan(t, sr)

			if span.Name() != tc.wantSpanName {
				t.Errorf("span name=%q, want %q", span.Name(), tc.wantSpanName)
			}

			attrs := attrsOf(span)
			if got := attrs["http.request.method"].AsString(); got != http.MethodGet {
				t.Errorf("http.request.method=%q, want %q", got, http.MethodGet)
			}
			gotURL := attrs["url.full"].AsString()
			if gotURL == "" {
				t.Errorf("url.full not set")
			}
			if strings.ContainsAny(gotURL, "?") {
				t.Errorf("url.full=%q contains a query string; expected query to be redacted", gotURL)
			}
			for _, forbidden := range []string{"bucket_name", "objects_user_id", "secret", "user-uuid"} {
				if strings.Contains(gotURL, forbidden) {
					t.Errorf("url.full=%q contains %q; expected query parameters to be redacted", gotURL, forbidden)
				}
			}
			if got := attrs["http.response.status_code"].AsInt64(); got != http.StatusOK {
				t.Errorf("http.response.status_code=%d, want %d", got, http.StatusOK)
			}

			if tc.wantEndpoint {
				if got := attrs["cloudscale.endpoint"].AsString(); got != tc.opPath {
					t.Errorf("cloudscale.endpoint=%q, want %q", got, tc.opPath)
				}
			} else if _, ok := attrs["cloudscale.endpoint"]; ok {
				t.Errorf("cloudscale.endpoint should be absent when no operation path is set")
			}
		})
	}
}

func TestInstrumentedTransport_URLFullRedactsQuery(t *testing.T) {
	tracer, sr := newRecordingTracer(t)

	server := newTestServer(t, statusHandler(http.StatusOK))
	client := newInstrumentedClient(nil, tracer)

	reqURL := server.URL + "/v1/metrics/buckets?start=2026-01-01&end=2026-01-02&bucket_name=secret&objects_user_id=user-uuid"
	mustDoRequest(t, client, http.MethodGet, reqURL, "v1/metrics/buckets")

	got := attrsOf(singleSpan(t, sr))["url.full"].AsString()
	want := server.URL + "/v1/metrics/buckets"
	if got != want {
		t.Errorf("url.full=%q, want %q", got, want)
	}
}

func TestInstrumentedTransport_Propagation(t *testing.T) {
	prev := otel.GetTextMapPropagator()
	otel.SetTextMapPropagator(propagation.TraceContext{})
	t.Cleanup(func() { otel.SetTextMapPropagator(prev) })

	tracer, _ := newRecordingTracer(t)

	var got string
	server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("traceparent")
		w.WriteHeader(http.StatusOK)
	})

	client := newInstrumentedClient(nil, tracer)

	mustDoRequest(t, client, http.MethodGet, server.URL+"/v1/servers", "v1/servers")

	if got == "" {
		t.Fatal("expected traceparent header to be injected")
	}
}
