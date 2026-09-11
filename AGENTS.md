# Repository Instructions

> [!IMPORTANT]
> Read [`README.md`](README.md) for product overview, flags, metrics, and operator usage.

## Tech Stack

- Go 1.27+ (see [`go.mod`](go.mod))
- [`prometheus/client_golang`](https://github.com/prometheus/client_golang) v1.24+ — metric registration and HTTP handler
- [`urfave/cli/v3`](https://github.com/urfave/cli) v3.11+ — CLI flags and application lifecycle
- [`sirupsen/logrus`](https://github.com/sirupsen/logrus) — structured logging
- [`goreleaser`](https://goreleaser.com/) v2 — cross-platform release builds (see [`.goreleaser.yml`](.goreleaser.yml))

## Repository Structure

- `cmd/` — Entry point (`main.go`); calls `internal/cli` for app setup
- `internal/cli/` — CLI flags and defaults (`0.0.0.0:10034`), the `CTRLD_API_KEY` source, app wiring
- `internal/config/` — flag reads and API-key validation
- `internal/server/` — HTTP server serving `/metrics` and the landing page
- `internal/collector/` — one `prometheus.Collector` fans out to the billing, endpoint, network, profile, service, stats and organization collectors. A new one is a method there, not a registration
- `internal/controld/` — Control D API client and response types
- `internal/log/` — logrus setup
- `docs/` — `README.md` owns the shared rules, `collectors.md` the catalogue, `help.md` the flags
- `examples/` — Prometheus scrape config, alert rules with tests, and a Grafana dashboard

## Setup and Commands

- `make help` — List every target, which is also the default goal
- `make build` — Build the binary into `tmp/controld-exporter`
- `make lint` — `golangci-lint run` + `go mod tidy`
- `make test-unit` — Run unit tests via `gotestsum` with coverage
- `make test-unit-coverage` — Generate HTML report at `coverage/report.html`
- `make clean` — Remove build artifacts and `.bak*` files
- `make image` — Build the Docker image (`$USER/controld-exporter`)
- `make pre-commit-install` / `pre-commit-test` / `pre-commit-uninstall` — Manage the `no-commit-to-main`, `golangci-lint`, `actionlint`, `gitleaks`, and `markdownlint-cli2` hooks (see [`.pre-commit-config.yaml`](.pre-commit-config.yaml))
- The `markdownlint-cli2` hook runs with `--fix`, so it rewrites Markdown in place

## Code Style

- Linting and formatting are enforced by `golangci-lint` (see [`.golangci.yml`](.golangci.yml)).
- Break a metric name, help string, type or label only in a release that bumps the major.
- Call the unit in `internal/collector` a collector, never a module.
- Keep Control D API logic in `internal/controld` so collectors remain thin and testable.
- Comments record only what the code cannot say, and never address the reader.

## Testing

- Run `make lint` and `make test-unit` before committing.
- Place tests next to code under test (`*_test.go`); the repository has no unit tests yet.
- [`CONTRIBUTING.md`](CONTRIBUTING.md) carries what CI runs and what covers the example rules.

## Commits and PRs

- Use [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `chore(deps):`, etc.).
- Sign off commits with `Signed-off-by:` (DCO).
- Open PRs against `main`. [`CONTRIBUTING.md`](CONTRIBUTING.md#development) carries what CI runs.
- Call out metrics or flag changes, because they move operator dashboards and alerts.

## Domain Knowledge

A claim about Control D is written only after a live response carried it, and the reference is read as a description rather than a contract. See [`CONTRIBUTING.md`](CONTRIBUTING.md#documentation) for which page owns which fact.

### The API Contract

The API is unversioned and Control D ships breaking changes without notice, so a field this exporter decodes can change spelling or disappear between two scrapes.

- **A live response defines the shape** — `/organizations/organization` documents `stats_endpoint` while the client decodes `statsEndpoint`. The `/devices` schema omits `client_count` outright, and marks `last_activity` and `clients` for deletion without listing them either.
- **A vanished field publishes as `0`** — only a transport, status, envelope or type error withholds a family, so a field the API stops sending reaches Prometheus as a measurement. See [Absence](docs/README.md#absence).

### Authentication and Cost

A missing or invalid token answers `400` with error `40001`, and a token without access to an endpoint answers `403` with `40301`. The first three digits restate the HTTP status, so the code carries a reason the status does not.

- **Two endpoints need no token** — `/network` and `/services/categories` declare `security: []` and answer a revoked key, which is why absence over their families cannot detect one. See [Absence](docs/README.md#absence).
- **A sub-organization is read by impersonation** — the parent token repeats the device, profile, category and report calls under `X-Force-Org-Id`. Cost grows with the account, not the collector list. See [Account Scope](docs/README.md#account-scope).
- **The scrape interval is the only throttle** — no response carries `X-RateLimit-*` or `Retry-After`, and nothing prevents a second scrape from overlapping the first. See [Scrape Path](docs/README.md#scrape-path).

### Analytics

DNS logging is a per-endpoint setting rather than an account-wide one, and a new endpoint starts with it off, so an account with real traffic can report nothing.

- **Absence has more than one reading** — a failed report call withholds `controld_stats_last_queries_count` and logs at `error`. A business-mode organization failure loses the whole scrape before the stats collector runs. See [Absence](docs/README.md#absence).
- **Both enums are undocumented** — `stats` reads `0` as off, `1` as basic and `2` as full, and a verdict code this exporter predates folds into `unknown`. See [Labels](docs/collectors.md#labels).
- **The analytics host follows data residency** — the organization response names the region label that fronts `analytics.controld.com`, while personal mode hardcodes `america`. That host answers a plain-text `404` rather than the JSON envelope, so the client checks the status before it decodes.

### Anycast

Control D serves DNS from anycast prefixes (`76.76.2.0/24`, `76.76.10.0/24` and `2606:1a40::/48`), so BGP rather than the client picks the point of presence that answers a query.

- **`/network` is a status board, not a path statement** — `iata_code` names a node the operator's resolvers may never reach, and only the node that served the call is proven reachable. See [Specifications](docs/collectors.md#specifications).
- **`-1` is not down** — it means the service is not offered at that node, which `proxy` reports on most of them, so a rule on `!= 1` fires across the fleet. See [Specifications](docs/collectors.md#specifications).

### Endpoints

Control D calls a resolver an endpoint and maps it to a physical device by convention alone, so a count of endpoints counts policy attachment points rather than machines.

- **Secure DNS carries identity, legacy DNS cannot** — a DoH URL or DoT hostname embeds the resolver ID while a legacy resolver is a bare UDP 53 pair. That is why an endpoint tracks `learn_ip` and an authorized-IP list.

### Payload Sensitivity

A device response holds client hostnames, MAC addresses and IP addresses, and an organization response holds contact names, emails and an Okta client secret. `--log.level debug` writes that body as received, so it belongs in a lab and never in a log. See [`SECURITY.md`](SECURITY.md#exposure).
