package instrumentation

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/cloudscale-ch/cloudscale-go-sdk/v9"
)

// tracingTransport wraps an http.RoundTripper with OpenTelemetry spans.
type tracingTransport struct {
	next   http.RoundTripper
	tracer trace.Tracer
}

func (t *tracingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	endpoint := cloudscale.OperationPath(req.Context())
	spanName := req.Method
	if endpoint != "" {
		spanName = req.Method + " " + endpoint
	}

	// Query parameters on certain endpoints (bucket_name, objects_user_id,
	// tag filters, volume name filters) may carry user-chosen identifiers.
	// We always drop the query string to avoid leaking potentially sensitive data.
	sanitizedURL := *req.URL
	sanitizedURL.RawQuery = ""

	ctx, span := t.tracer.Start(req.Context(), spanName,
		trace.WithAttributes(
			semconv.HTTPRequestMethodKey.String(req.Method),
			semconv.URLFull(sanitizedURL.String()),
		),
	)
	defer span.End()

	if endpoint != "" {
		span.SetAttributes(attribute.String("cloudscale.endpoint", endpoint))
	}

	// Clone so we can inject trace headers without aliasing the caller's
	// header map across retries.
	req = req.Clone(ctx)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := t.next.RoundTrip(req)

	if resp != nil {
		span.SetAttributes(semconv.HTTPResponseStatusCode(resp.StatusCode))
		if resp.StatusCode >= 400 {
			span.SetStatus(codes.Error, http.StatusText(resp.StatusCode))
		}
	}

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
	}

	return resp, err
}
