package instrumentation

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
)

// Options configures the instrumented transport.
type Options struct {
	// Subsystem prefix for metrics (default: "cloudscale").
	Subsystem string

	// Registry to register metrics on. If nil, metrics are disabled.
	PrometheusRegistry prometheus.Registerer

	// Tracer to use for spans. If nil, tracing is disabled.
	Tracer trace.Tracer
}

func (o Options) subsystem() string {
	if o.Subsystem != "" {
		return o.Subsystem
	}
	return "cloudscale"
}

// InstrumentedTransport wraps next with optional metrics and/or tracing.
// If next is nil, http.DefaultTransport is used. With an empty Options value,
// the (possibly defaulted) transport is returned unchanged.
func InstrumentedTransport(next http.RoundTripper, opts Options) http.RoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}
	if opts.PrometheusRegistry == nil && opts.Tracer == nil {
		return next
	}

	transport := next

	if opts.PrometheusRegistry != nil {
		metrics := NewMetricsCollector(opts.PrometheusRegistry, opts.subsystem())
		transport = &metricsTransport{
			next:    transport,
			metrics: metrics,
		}
	}

	if opts.Tracer != nil {
		transport = &tracingTransport{
			next:   transport,
			tracer: opts.Tracer,
		}
	}

	return transport
}
