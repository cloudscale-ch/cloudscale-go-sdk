# cloudscale.ch Go API SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/cloudscale-ch/cloudscale-go-sdk.svg)](https://pkg.go.dev/github.com/cloudscale-ch/cloudscale-go-sdk/v9)
[![Tests](https://github.com/cloudscale-ch/cloudscale-go-sdk/actions/workflows/test.yaml/badge.svg)](https://github.com/cloudscale-ch/cloudscale-go-sdk/actions/workflows/test.yaml)

If you want to manage your cloudscale.ch server resources with Go, you are at
the right place.

## Getting Started

To use the `cloudscale-go-sdk` for managing your cloudscale.ch resources, follow these steps:

1. **Install the SDK**:\
   Run the following command to install the SDK:

   ```console
   go mod init example.com/m
   go get github.com/cloudscale-ch/cloudscale-go-sdk/v9
   ```

1. **Create a File**:\
   Save the following code into a file, for example, `main.go`.

   ```go
   package main

   import (
       "context"
       "fmt"
       "github.com/cenkalti/backoff/v5"
       "github.com/cloudscale-ch/cloudscale-go-sdk/v9"
       "golang.org/x/oauth2"
       "log"
       "os"
       "time"
   )

   func main() {
       // Read the API token from the environment variable
       apiToken := os.Getenv("CLOUDSCALE_API_TOKEN")
       if apiToken == "" {
           log.Fatalf("CLOUDSCALE_API_TOKEN environment variable is not set")
       }

       // Create a new client
       tc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
           &oauth2.Token{AccessToken: apiToken},
       ))
       client := cloudscale.NewClient(tc)

       // Define server configuration
       createRequest := &cloudscale.ServerRequest{
           Name:    "example-server",
           Flavor:  "flex-8-2",
           Image:   "debian-11",
           Zone:    "rma1",
           SSHKeys: []string{"<KEY>"},
       }

       // Create a server
       server, err := client.Servers.Create(context.Background(), createRequest)
       if err != nil {
           log.Fatalf("Error creating server: %v", err)
       }

       fmt.Printf("Creating server with UUID: %s\n", server.UUID)

       // Wait for the server to be in "running" state
       server, err = client.Servers.WaitFor(
           context.Background(),
           server.UUID,
           cloudscale.ServerIsRunning, // can be replaced with custom condition funcs
           // optionally, pass any option that github.com/cenkalti/backoff/v5 supports
           backoff.WithNotify(func(err error, duration time.Duration) {
               fmt.Printf("Retrying after error: %v, waiting for %v\n", err, duration)
           }),
       )
       if err != nil {
           log.Fatalf("Error waiting for server to start: %v\n", err)
       }

       fmt.Printf("Server is now running with status: %s\n", server.Status)
   }
   ```

1. **Run the File**:\
   Make sure the `CLOUDSCALE_API_TOKEN` environment variable is set. Then, run the file:

   ```bash
   export CLOUDSCALE_API_TOKEN="your_api_token_here"
   go run main.go
   ```

That's it! The code will create a server and leverage the `WaitFor` helper to wait until the server status changes to `running`. For more advanced options, check the [documentation](https://pkg.go.dev/github.com/cloudscale-ch/cloudscale-go-sdk/v9).

## Instrumentation

The SDK ships a transport wrapper in
`github.com/cloudscale-ch/cloudscale-go-sdk/v9/instrumentation` that adds
Prometheus metrics and/or OpenTelemetry spans to every API call. Both signals
are independent — set only the fields you need on `Options`, and leaving both
unset returns the transport unchanged.

Wrap the transport on the client returned by `oauth2.NewClient` before handing
it to `cloudscale.NewClient`:

```go
import (
    "context"
    "os"

    "github.com/cloudscale-ch/cloudscale-go-sdk/v9"
    "github.com/cloudscale-ch/cloudscale-go-sdk/v9/instrumentation"
    "github.com/prometheus/client_golang/prometheus"
    "go.opentelemetry.io/otel"
    "golang.org/x/oauth2"
)

apiToken := os.Getenv("CLOUDSCALE_API_TOKEN")

reg := prometheus.NewRegistry()
tracer := otel.Tracer("cloudscale-go-sdk")

tc := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
    &oauth2.Token{AccessToken: apiToken},
))
tc.Transport = instrumentation.InstrumentedTransport(tc.Transport, instrumentation.Options{
    PrometheusRegistry: reg,
    Tracer:             tracer,
})

client := cloudscale.NewClient(tc)
```

For metrics-only, set just `PrometheusRegistry`; for tracing-only, set just
`Tracer`. The `Subsystem` field overrides the default `cloudscale` metric
prefix when you need to share a registry across collectors.

What is recorded:

- `cloudscale_requests_total{method, endpoint, status}` — counter of API
  requests by HTTP method, path template, and status code.
- `cloudscale_request_duration_seconds{method, endpoint}` — request latency
  histogram.
- `cloudscale_in_flight_requests` — gauge of concurrent in-flight requests.
- Spans named `{METHOD} {endpoint}` (e.g. `GET v1/servers/:id`), or just
  `{METHOD}` when no path template is set. Attributes follow the OpenTelemetry
  HTTP semantic conventions (`http.request.method`, `url.full`,
  `http.response.status_code`) plus a `cloudscale.endpoint` attribute with the
  path template.
- W3C trace context is injected into outbound request headers using the
  globally configured propagator. Call
  `otel.SetTextMapPropagator(propagation.TraceContext{})` at startup to enable
  propagation (it is a no-op by default).

## Testing

The test directory contains integration tests, aside from the unit tests in the
root directory. While the unit tests suite runs very quickly because they
don't make any network calls, this can take some time to run.

### test/integration

This folder contains tests for every type of operation in the cloudscale.ch API
and runs tests against it.

Since the tests are run against live data, there is a higher chance of false
positives and test failures due to network issues, data changes, etc.

Run all integration tests using:

```
CLOUDSCALE_API_TOKEN="HELPIMTRAPPEDINATOKENGENERATOR" make integration
```

To run only the CRUD tests for a quick smoke test:

```
CLOUDSCALE_API_TOKEN="HELPIMTRAPPEDINATOKENGENERATOR" make integration-short
```

There's a possibility to specify the `CLOUDSCALE_API_URL` environment variable to
change the default url of https://api.cloudscale.ch, but you can almost certainly
use the default.

If you want to give params to `go test`, you can use something like this:

```
TESTARGS='-run FloatingIP' make integration
```

Some test default to "rma1" for testing. To override this, you can set the following variable

```
INTEGRATION_TEST_ZONE="lpg1"  make integration
```

## Releasing

To create a new release, please do the following:

- Merge all feature branches into `main`/`master` branch
- Create a release branch from `main`/`master` branch
- Run `make NEW_VERSION=v1.x.x bump-version`
  - For a new major release: follow [these instructions](https://go.dev/doc/modules/major-version)
  - For a new major release: update the `pkg.go.dev` refercenes in this file (multiple!).
- Commit changes
- Open a merge request for the release branch and after code review merge the release branch into master
- Create a [new release](https://github.com/cloudscale-ch/cloudscale-go-sdk/releases/new) on GitHub.
