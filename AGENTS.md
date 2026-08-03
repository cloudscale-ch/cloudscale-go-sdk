# AGENTS.md

This file is a supplement for AI agents working on the cloudscale SDK (cloudscale-go-sdk).

## What this is

A Go client library for the cloudscale.ch API, imported as
`github.com/cloudscale-ch/cloudscale-go-sdk/v9`. It is source-only: there is no `main` package and nothing to build or
ship, so consumers just `go get` it.

## What to run after a change

| Change                                                      | Command                                                       |
|-------------------------------------------------------------|---------------------------------------------------------------|
| Edited any `.go` file                                       | `make fmt`, then `make vet`                                   |
| Before committing                                           | `make lint` (or `make lint-fix` to apply fixes automatically) |
| Changed logic                                               | `make test`, unit tests, run with `-race`, no network calls   |
| Touched a service and want to check it against the real API | `make integration` or `make integration-short`                |

`make vet` also vets the integration tests under the `integration` build tag, so run it even when you only touched files
in `test/integration`.

## Where things live

| Path                      | Contents                                                                                                  |
|---------------------------|-----------------------------------------------------------------------------------------------------------|
| root package `cloudscale` | one file per API resource (`servers.go`, `servers.go`, `volumes.go`, …) plus its `*_test.go`              |
| `cloudscale.go`           | the `Client`, service wiring in `NewClient`, `NewRequest`/`Do`, error handling, and `libraryVersion`      |
| `generic_service.go`      | generic CRUD operations and `WaitFor`, shared by most services                                            |
| `ctxutil.go`              | `WithOperationPath` / `OperationPath`, the endpoint template carried on the context for metrics and spans |
| `instrumentation/`        | optional transport wrapper adding Prometheus metrics and OpenTelemetry spans                              |
| `test/integration/`       | tests that hit the live API, behind the `//go:build integration` tag                                      |

## Adding or changing a service

Most services are just an instance of the generics in `generic_service.go`. Look at `servers.go`
next to `cloudscale.go` for a full example, including the case where a resource needs extra methods on top of the
generic ones. The usual steps:

1. Define the resource struct and its request structs with `json` tags. Embed `ZonalResource` /
   `TaggedResource` (and the matching `ZonalResourceRequest` / `TaggedResourceRequest`) when the resource is zonal or
   taggable.
2. Declare a `XService` interface built from the `Generic…Service[...]` interfaces, listing only the operations the API
   actually supports.
3. Wire it up in `NewClient` in `cloudscale.go`. Use a bare
   `GenericServiceOperations[Resource, CreateRequest, UpdateRequest]{client, path}` when the generic operations are
   enough, or a wrapper struct that embeds it when you need custom methods (see `ServerServiceOperations` and its
   `CreateInterface` / `DeleteInterface`).
4. For custom endpoints, set the operation path before making the request, e.g.
   `ctx = WithOperationPath(ctx, serverBasePath+"/:id/reboot")`, so metrics and spans get a stable template instead
   of a URL with a UUID baked in.
5. If callers need to wait on a status, add status constants and a condition function in the style of `XIsRunning`
   for use with `WaitFor`.

## Conventions

- Linting is golangci-lint v2, configured in `.golangci.yml`. Imports are grouped with `goimports`
  using the local prefix `github.com/cloudscale-ch/cloudscale-go-sdk`.
- Unit tests are table-driven and use an `httptest` mock server. They must not make real network calls, which is why
  `make test` stays fast.
- Every request method takes a `context.Context` as its first argument. API errors come back as
  `*ErrorResponse`.

## Integration tests

These run against a real cloudscale account, so they cost money and create real resources. Resources are named
`go-sdk-<random>` and cleaned up automatically after each test.

- `CLOUDSCALE_API_TOKEN` — required, the API token to test against.
- `CLOUDSCALE_API_URL` — optional, defaults to `https://api.cloudscale.ch`.
- `INTEGRATION_TEST_ZONE` — optional, defaults to `rma1`.

To run a subset, pass `go test` arguments through `TESTARGS`, e.g.
`TESTARGS='-run FloatingIP' make integration`.

## Releasing

Releases are automated: bump the version with `make NEW_VERSION=vX.Y.Z bump-version`, merge, then push a signed tag. The
`Release` workflow in `.github/workflows/release.yml` verifies the tag and publishes the GitHub release. The README's
"Releasing" section has the full procedure.

## References

- `README.md` — usage, instrumentation, testing, and release steps.
- API docs: <https://pkg.go.dev/github.com/cloudscale-ch/cloudscale-go-sdk/v9>
