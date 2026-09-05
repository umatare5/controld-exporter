# Repository Instructions

> [!IMPORTANT]
> Read [`README.md`](README.md) for product overview, flags, metrics, and operator usage.

## Tech Stack

- Go 1.27+ (see [`go.mod`](go.mod))
- [`prometheus/client_golang`](https://github.com/prometheus/client_golang) v1.24+ — metric registration and HTTP handler
- [`urfave/cli/v3`](https://github.com/urfave/cli) v3.11+ — CLI flags and application lifecycle
- [`sirupsen/logrus`](https://github.com/sirupsen/logrus) — structured logging
- [`jinzhu/configor`](https://github.com/jinzhu/configor) — environment-variable loading into the config struct
- [`goreleaser`](https://goreleaser.com/) v2 — cross-platform release builds (see [`.goreleaser.yml`](.goreleaser.yml))

## Repository Structure

- `cmd/` — Entry point (`main.go`); calls `internal/cli` for app setup
- `internal/cli/` — CLI flag definitions and app wiring (urfave/cli/v3)
- `internal/config/` — flag/env parsing, defaults (`0.0.0.0:10034`), and API-key validation
- `internal/server/` — HTTP server serving `/metrics` and the landing page
- `internal/collector/` — billing, endpoint, network, profile, service, stats, and organization collectors; `prometheus.Collector` implementations
- `internal/controld/` — Control D API client and response types
- `internal/log/` — logrus setup
- `docs/` — the reference set: `README.md` for the shared rules, `collectors.md` for the catalogue, `help.md` for the flags
- `examples/` — Prometheus scrape config, alert rules with their unit tests, and a Grafana dashboard

## Setup and Commands

- `make build` — Build the binary into `tmp/controld-exporter`
- `make lint` — `golangci-lint run` + `go mod tidy`
- `make test-unit` — Run unit tests via `gotestsum` with coverage
- `make test-unit-coverage` — Generate HTML report at `coverage/report.html`
- `make clean` — Remove build artifacts and `.bak*` files
- `make image` — Build the Docker image (`$USER/controld-exporter`)
- `make pre-commit-install` / `pre-commit-test` / `pre-commit-uninstall` — Manage the `no-commit-to-main`, `golangci-lint`, `actionlint`, `gitleaks`, and `markdownlint-cli2` hooks (see [`.pre-commit-config.yaml`](.pre-commit-config.yaml))

## Code Style

- Linting and formatting are enforced by `golangci-lint` in the pre-commit hook (see [`.golangci.yml`](.golangci.yml)).
- Keep metric names, help strings, types, and labels stable unless a SemVer-signaled breaking change is intentional.
- Keep Control D API logic in `internal/controld` so collectors remain thin and testable.
- Comments record only what the code cannot say, and never address the reader.

## Testing

- Run `make lint` and `make test-unit` before committing.
- Place tests next to code under test (`*_test.go`); the repository has no unit tests yet.

## Commits and PRs

- Use [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `chore(deps):`, etc.).
- Sign off commits with `Signed-off-by:` (DCO).
- Open PRs against `main`. CI runs lint, build, CodeQL, `govulncheck`, and `promtool` over the example rules.
- Call out any metrics or flag changes explicitly because they affect operator dashboards and alerts.

## Domain Knowledge

- **The API is unversioned, and Control D says breaking changes ship without warning.** A field this exporter decodes can change spelling or disappear between two scrapes, and no request parameter pins the old shape. Every absence guard in `internal/collector` is load-bearing rather than defensive.
- **The reference describes the API; only a live response defines it.** `/organizations/organization` documents the analytics host as a required `stats_endpoint`. The field a real response carries is `statsEndpoint`, which is what `internal/controld/organization.go` decodes.
- **Authentication failures answer `400`, never `401` or `403`.** A missing token and an invalid one both return `{"success": false, "error": {"code": 40001}}` under HTTP 400. The first three digits of the error code restate the status, so the status alone cannot separate a revoked key from a malformed request.
- **Two endpoints need no token at all.** `/network` and `/services/categories` declare `security: []` and answer 200 unauthenticated. Those two collectors keep publishing after a key is revoked, so an alert watching only `controld_network_health_code` cannot detect a dead token.
- **A sub-organization is read by impersonation rather than by a path.** The parent token repeats the same `/devices` or `/profiles` request under `X-Force-Org-Id: <PK>`, which Control D documents for the Profiles scope alone. Each sub-organization therefore costs a full round trip.
- **Rate limits are neither documented nor signalled.** No response carries `X-RateLimit-*` or `Retry-After`, and no page states a quota. The scrape interval is the only throttle, and one scrape costs a call per collector plus one per sub-organization.
- **The analytics host is a separate service that does not speak the JSON envelope.** `https://<statsEndpoint>.analytics.controld.com` is a regional alias absent from the reference. Its report path answers a plain-text `404 page not found`, which is why the client checks the HTTP status before it decodes.
- **A payload carries more than the metrics need.** A device response holds client hostnames, MAC addresses and IP addresses, and an organization response holds contact names, emails and an Okta client secret. `--log.level debug` prints that body, so it belongs in a lab and never in a log.
