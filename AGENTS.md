# Repository Instructions

> [!IMPORTANT]
> Read [`README.md`](README.md) for the project overview.

## Tech Stack

- Go 1.27+ (see [`go.mod`](go.mod))
- [`prometheus/client_golang`](https://github.com/prometheus/client_golang) v1.24+ – metric registration and HTTP handler
- [`urfave/cli/v3`](https://github.com/urfave/cli) v3.11+ – CLI flags and application lifecycle
- [`sirupsen/logrus`](https://github.com/sirupsen/logrus) – structured logging
- [`goreleaser`](https://goreleaser.com/) v2 – cross-platform release builds (see [`.goreleaser.yml`](.goreleaser.yml))

## Repository Structure

Read from `cmd/main.go`. Each package is named for what it owns.

- [`cmd/main.go`](cmd/main.go) – application entry point
- [`internal/cli/`](internal/cli) – command-line flags, defaults and app wiring
- [`internal/config/`](internal/config) – flag reads and API key validation
- [`internal/collector/`](internal/collector) – metric descriptions and collection logic
- [`internal/controld/`](internal/controld) – upstream API client and data structures
- [`internal/server/`](internal/server) – HTTP server configuration and routing
- [`internal/log/`](internal/log) – logrus level and formatter setup
- [`docs/`](docs) – reference pages behind the README
- [`scripts/`](scripts) – helper scripts the pre-commit hooks run
- [`examples/`](examples) – Prometheus configuration, alert rules and the Grafana dashboard

## Setup and Commands

Run `make pre-commit-install` first.

- Read [`Makefile`](Makefile) which lists all available make targets and their descriptions.
- Read [`CONTRIBUTING.md`](CONTRIBUTING.md) which provides guidelines for contributing to the project.

## Code Style

Follow [Effective Go](https://go.dev/doc/effective_go) conventions and the software development principles DRY/YAGNI/SRP.

- Keep code simple and readable, avoiding clever tricks that obscure intent.
- Keep minimal for all changes, coding, testing, commenting, and documentation.
- Write simple comments that explain the reasoning behind the code, not just what it does.

## Testing

Follow [`CONTRIBUTING.md`](CONTRIBUTING.md).

- Run `make lint` and `make test-unit` before creating a commit.

## Commits and PRs

Follow [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `chore(deps):`, etc.).

- Run pre-commit and ensure all hooks pass before committing.
- Must Sign off all commits with `Signed-off-by:` (DCO).
- Open PRs against `main`. Create Draft PR as default.

## Domain Knowledge

Learn the constraints outside the exporter, because they decide what a metric can mean.

### About the API

The API is unversioned and ships breaking changes without notice. See [Absence](docs/architecture.md#absence).

- A vanished field publishes as `0` because only a transport or envelope error withholds it.
- Two endpoints need no token, so `/network` and `/services/categories` answer a revoked key.
- Sub-organizations are read by impersonation, repeating calls under `X-Force-Org-Id`.
- The scrape interval is the only throttle, because no response carries `X-RateLimit-*`.
- The error code restates the HTTP status, so `40001`, `40301` and `40401` name the reason.
- Query reporting is out of reach, because the analytics host rejects the API key.

### About the anycast

Control D serves DNS from anycast, so BGP picks the node. See [Metrics](README.md#metrics).

- The anycast prefixes are `76.76.2.0/24`, `76.76.10.0/24` and `2606:1a40::/48`.
- `/network` is a status board, so `iata_code` may name a node the resolvers never reach.
- `-1` means the service is not offered there, which `proxy` reports on most nodes.
