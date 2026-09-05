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

- **The reference describes, a live response defines** — `/organizations/organization` requires `stats_endpoint`, while the client decodes `statsEndpoint`.
- **A documented schema is not a complete one** — `client_count`, the field `controld_endpoint_clients_total` reads, is absent from the `/devices` schema.
- **Removal arrives announced at best** — that page marks `last_activity` and `clients` for deletion, and neither is in the schema either.
- **Every absence guard is load-bearing** — a collector that withholds its family on a decode failure keeps a vanished field out of Prometheus as a zero.

### Authentication and cost

Authentication failures answer `400`, never `401` or `403`. A missing token and an invalid one both return `{"success": false, "error": {"code": 40001}}`. The first three digits of the error code restate the HTTP status, so the status alone cannot separate a revoked key from a malformed request.

- **Two endpoints need no token** — `/network` and `/services/categories` declare `security: []`, so they keep publishing after a key is revoked.
- **A sub-organization is read by impersonation** — the parent token repeats the request under `X-Force-Org-Id: <PK>`, so each costs a full round trip.
- **Rate limits are neither documented nor signalled** — no response carries `X-RateLimit-*` or `Retry-After`, so the scrape interval is the only throttle.

### Analytics

DNS logging is a per-endpoint setting rather than an account-wide one, and a new endpoint starts with it off. An account with real traffic and logging off reports nothing, so an absent `controld_stats_last_queries_count` says the operator never enabled analytics rather than that no query was resolved.

- **Three levels, set per endpoint** — `stats` is `0 = off`, `1 = basic` and `2 = full`, and `1` already stores the verdict counts this exporter reads.
- **The verdict is a DNS action sent as an integer** — Control D names blocked, bypassed and redirected, which the collector maps from `0`, `1` and `3`.
- **The enum is undocumented** — `2` and any action added later land in the `unknown` bucket rather than being dropped.
- **The analytics host follows data residency** — `statsEndpoint` names the host of the region an account chose, while personal mode hardcodes `america`.
- **That host does not speak the JSON envelope** — its removed report path answers a plain-text `404 page not found`, so the client checks the status first.
- **Retention differs by grain** — raw query logs live one month and aggregate statistics one year, which bounds a backfill but not this exporter.

### Anycast

Control D serves DNS from the anycast prefixes `76.76.2.0/24`, `76.76.10.0/24` and `2606:1a40::/48`, so BGP rather than the client picks the point of presence that answers a query. `/network` is therefore a global status board rather than a statement about the path an endpoint takes, and the `iata_code` label names a node the operator's resolvers may never reach.

- **`current_pop` is the only node a scrape proves reachable** — it names the node that served the API call, which is the exporter's path, not the resolvers'.
- **`-1` is not down** — `pxy` reads `-1` on 23 of 37 nodes, where no transparent proxy is offered, so an alert on `!= 1` fires on every one of them.

### Endpoints

Control D calls a resolver an endpoint and maps it to a physical device by convention alone, so a count of endpoints counts policy attachment points rather than machines.

- **Secure DNS carries identity, legacy DNS cannot** — a DoH URL or DoT hostname embeds the resolver ID, while a legacy resolver is a bare UDP 53 pair.
- **That asymmetry is why an endpoint tracks IPs** — `learn_ip` and the authorized-IP list give a legacy query the identity its transport lacks.

### Payload sensitivity

A device response holds client hostnames, MAC addresses and IP addresses, and an organization response holds contact names, emails and an Okta client secret. `--log.level debug` prints that body verbatim, so it belongs in a lab and never in a log, and `CTRLD_API_KEY` belongs in neither.
