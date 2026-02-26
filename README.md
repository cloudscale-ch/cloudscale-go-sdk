# cloudscale.ch Go API SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/cloudscale-ch/cloudscale-go-sdk.svg)](https://pkg.go.dev/github.com/cloudscale-ch/cloudscale-go-sdk/v7)
[![Tests](https://github.com/cloudscale-ch/cloudscale-go-sdk/actions/workflows/test.yaml/badge.svg)](https://github.com/cloudscale-ch/cloudscale-go-sdk/actions/workflows/test.yaml)

If you want to manage your cloudscale.ch server resources with Go, you are at
the right place.

## Getting Started

To use the `cloudscale-go-sdk` for managing your cloudscale.ch resources, follow these steps:

1. **Install the SDK**:\
   Run the following command to install the SDK:

   ```console
   go mod init example.com/m
   go get github.com/cloudscale-ch/cloudscale-go-sdk/v7
   ```

1. **Create a File**:\
   Save the following code into a file, for example, `main.go`.

   ```go
   package main

   import (
       "context"
       "fmt"
       "github.com/cenkalti/backoff/v5"
       "github.com/cloudscale-ch/cloudscale-go-sdk/v7"
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

That's it! The code will create a server and leverage the `WaitFor` helper to wait until the server status changes to `running`. For more advanced options, check the [documentation](https://pkg.go.dev/github.com/cloudscale-ch/cloudscale-go-sdk/v7).

## Testing

The test directory contains integration tests, aside from the unit tests in the
root directory. While the unit tests suite runs very quickly because they
don't make any network calls, this can take some time to run.

### test/integration

This folder contains tests for every type of operation in the cloudscale.ch API
and runs tests against it.

Since the tests are run against live data, there is a higher chance of false
positives and test failures due to network issues, data changes, etc.

Run the tests using:

```
CLOUDSCALE_API_TOKEN="HELPIMTRAPPEDINATOKENGENERATOR" make integration
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
