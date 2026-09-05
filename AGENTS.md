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

### The API contract

The API is unversioned, and Control D warns that breaking changes ship without notice. A field this exporter decodes can change spelling or disappear between two scrapes, and no request parameter pins the old shape.

- **The reference describes the API; a live response defines it.** `/organizations/organization` documents the analytics host as a required `stats_endpoint`, while `internal/controld/organization.go` decodes `statsEndpoint`.
- **A documented schema is not a complete one.** `client_count`, the only field `controld_endpoint_clients_total` reads, appears nowhere in the `/devices` schema. Neither do `clients`, `ip_total`, `last_activity` or `ctrld`, and the reference already announces the removal of two of them.
- **Every absence guard is load-bearing.** A collector that withholds its family on a decode failure is what keeps a vanished field out of Prometheus as a fabricated zero.

### Authentication and cost

Authentication failures answer `400`, never `401` or `403`. A missing token and an invalid one both return `{"success": false, "error": {"code": 40001}}`. The first three digits of the error code restate the HTTP status, so the status alone cannot separate a revoked key from a malformed request.

- **Two endpoints need no token.** `/network` and `/services/categories` declare `security: []` and answer 200 unauthenticated, so they keep publishing after a key is revoked and cannot witness a dead token.
- **A sub-organization is read by impersonation.** The parent token repeats the same request under `X-Force-Org-Id: <PK>`, which Control D documents for the Profiles scope alone, so each sub-organization costs a full round trip.
- **Rate limits are neither documented nor signalled.** No response carries `X-RateLimit-*` or `Retry-After`, so the scrape interval is the only throttle.

### Analytics are opt-in

DNS logging is a per-endpoint setting, not an account-wide one: `stats` on an endpoint is `0 = off`, `1 = basic` and `2 = full`, and a new endpoint starts at `0`. An account with real traffic and logging off reports nothing, so an absent `controld_stats_last_queries_count` says the operator never enabled analytics rather than that no query was resolved.

- **Basic is enough for this metric.** Level `1` stores counts of blocks, redirects and bypasses without the queries themselves, which is exactly what the verdict time series reads.
- **The verdict is a DNS action sent as an integer.** Control D names three — blocked, bypassed and redirected — which `internal/collector/stats.go` maps from `0`, `1` and `3`. The enum is undocumented, so `2` and any future action land in the `unknown` bucket rather than being lost.
- **The analytics host follows data residency, not the API.** Enabling analytics asks for a storage region, and `statsEndpoint` names the host that region resolves to. Personal mode has no organization to read it from and hardcodes `america`, which is wrong for an account storing elsewhere.
- **That host does not speak the JSON envelope.** `https://<statsEndpoint>.analytics.controld.com` is absent from the reference and answers a removed report path with a plain-text `404 page not found`. The client checks the HTTP status before it decodes for that reason.
- **Retention differs by grain.** Raw query logs live one month and aggregate statistics one year, which bounds any backfill but not this exporter, since it reads the newest one-minute bucket alone.

### Anycast decides which node answers

Control D serves DNS from anycast prefixes — `76.76.2.0/24`, `76.76.10.0/24` and `2606:1a40::/48` — so BGP, not the client, picks the point of presence that answers a query. `/network` is therefore a global status board rather than a statement about the path an endpoint takes, and the `iata_code` label names a node the operator's resolvers may never reach.

- **`current_pop` is the only node a scrape proves reachable.** The same body names the node that served the API call, which is the exporter's own path and not the resolvers'.
- **`-1` is not down.** A live snapshot returns `api` and `dns` at `1` everywhere and `pxy` at `-1` on the majority of nodes, where the transparent proxy is not offered. Reading `!= 1` as unhealthy fires on every proxy-less node, which is why the codes pass through uninterpreted.

### An endpoint is a resolver, not a host

Control D calls a resolver an endpoint and maps it to a physical device by convention alone, so a count of endpoints counts policy attachment points rather than machines.

- **Secure DNS carries identity; legacy DNS cannot.** A DoH URL or a DoT hostname embeds the resolver ID, whereas a legacy resolver is a plain UDP 53 address pair shared by every client, so Control D identifies a legacy client by source IP.
- **That asymmetry is why an endpoint tracks IPs.** `learn_ip` and the authorized-IP list exist to give legacy queries an identity, and a restricted endpoint refuses an unknown source before it can be learned.

### Payloads carry more than the metrics need

A device response holds client hostnames, MAC addresses and IP addresses, and an organization response holds contact names, emails and an Okta client secret. `--log.level debug` prints that body verbatim, so it belongs in a lab and never in a log, and `CTRLD_API_KEY` belongs in neither.
