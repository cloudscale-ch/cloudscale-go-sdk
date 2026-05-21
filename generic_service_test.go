package cloudscale

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

// pathCaptureTransport records OperationPath from each request's context and
// returns a minimal successful JSON response so client.Do can decode without
// touching the network.
type pathCaptureTransport struct {
	path string
}

func (p *pathCaptureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	p.path = OperationPath(req.Context())
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader("{}")),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}, nil
}

func newCapturingClient() (*Client, *pathCaptureTransport) {
	pcap := &pathCaptureTransport{}
	c := NewClient(&http.Client{Transport: pcap})
	return c, pcap
}

func TestGenericServiceOperations_OperationPathPrecedence(t *testing.T) {
	cases := []struct {
		name     string
		preset   string // empty means caller did not set WithOperationPath
		op       func(ctx context.Context, g GenericServiceOperations[Region, struct{}, struct{}]) error
		wantPath string
	}{
		{
			name: "Get: no preset → default g.path+/:id",
			op: func(ctx context.Context, g GenericServiceOperations[Region, struct{}, struct{}]) error {
				_, err := g.Get(ctx, "abc")
				return err
			},
			wantPath: "v1/test/:id",
		},
		{
			name:   "Get: caller preset wins",
			preset: "custom/get/path",
			op: func(ctx context.Context, g GenericServiceOperations[Region, struct{}, struct{}]) error {
				_, err := g.Get(ctx, "abc")
				return err
			},
			wantPath: "custom/get/path",
		},
		{
			name: "Create: no preset → default g.path",
			op: func(ctx context.Context, g GenericServiceOperations[Region, struct{}, struct{}]) error {
				_, err := g.Create(ctx, &struct{}{})
				return err
			},
			wantPath: "v1/test",
		},
		{
			name:   "Create: caller preset wins",
			preset: "custom/create/path",
			op: func(ctx context.Context, g GenericServiceOperations[Region, struct{}, struct{}]) error {
				_, err := g.Create(ctx, &struct{}{})
				return err
			},
			wantPath: "custom/create/path",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c, pcap := newCapturingClient()
			g := GenericServiceOperations[Region, struct{}, struct{}]{client: c, path: "v1/test"}

			ctx := t.Context()
			if tc.preset != "" {
				ctx = WithOperationPath(ctx, tc.preset)
			}

			if err := tc.op(ctx, g); err != nil {
				t.Fatalf("operation failed: %v", err)
			}
			if pcap.path != tc.wantPath {
				t.Errorf("captured OperationPath=%q, want %q", pcap.path, tc.wantPath)
			}
		})
	}
}
