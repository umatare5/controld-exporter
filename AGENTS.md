# Repository Instructions

> [!IMPORTANT]
> Read [README.md](README.md) for product overview, flags, metrics, and operator usage.

## Tech Stack

- Go 1.27+ (see [go.mod](go.mod))
- [`prometheus/client_golang`](https://github.com/prometheus/client_golang) v1.24+ — metric registration and HTTP handler
- [`urfave/cli/v3`](https://github.com/urfave/cli) v3.11+ — CLI flags and application lifecycle
- [`sirupsen/logrus`](https://github.com/sirupsen/logrus) — structured logging
- [`jinzhu/configor`](https://github.com/jinzhu/configor) — environment-variable loading into the config struct
- [`goreleaser`](https://goreleaser.com/) v2 — cross-platform release builds (see [.goreleaser.yml](.goreleaser.yml))

## Repository Structure

- `cmd/` — Entry point (`main.go`); calls `internal/cli` for app setup
- `internal/cli/` — CLI flag definitions and app wiring (urfave/cli/v3)
- `internal/config/` — flag/env parsing, defaults (`0.0.0.0:10034`), and API-key validation
- `internal/server/` — HTTP server serving `/metrics` and the landing page
- `internal/collector/` — billing, endpoint, network, profile, service, stats, and organization collectors; `prometheus.Collector` implementations
- `internal/controld/` — Control D API client and response types
- `internal/log/` — logrus setup
- `examples/` — Prometheus scrape config, alert rules, and a Grafana dashboard

## Setup and Commands

- `make build` — Build the binary into `tmp/controld-exporter`
- `make lint` — `golangci-lint run` + `go mod tidy`
- `make test-unit` — Run unit tests via `gotestsum` with coverage
- `make test-unit-coverage` — Generate HTML report at `coverage/report.html`
- `make clean` — Remove build artifacts and `.bak*` files
- `make image` — Build the Docker image (`$USER/controld-exporter`)
- `make pre-commit-install` / `pre-commit-test` / `pre-commit-uninstall` — Manage the `no-commit-to-main`, `golangci-lint`, `actionlint`, `gitleaks`, and `markdownlint-cli2` hooks (see [.pre-commit-config.yaml](.pre-commit-config.yaml))

## Code Style

- Linting and formatting are enforced by `golangci-lint` in the pre-commit hook (see [.golangci.yml](.golangci.yml)).
- Keep metric names, help strings, types, and labels stable unless a SemVer-signaled breaking change is intentional.
- Keep Control D API logic in `internal/controld` so collectors remain thin and testable.
- Comments record only what the code cannot say, and never address the reader.

## Testing Instructions

- Run `make lint` and `make test-unit` before committing.
- Place tests next to code under test (`*_test.go`); the repository has no unit tests yet.

## Commits and PRs

- Use [Conventional Commits](https://www.conventionalcommits.org/) (`feat:`, `fix:`, `chore(deps):`, etc.).
- Sign off commits with `Signed-off-by:` (DCO).
- Open PRs against `main`. CI runs lint, build, and CodeQL.
- Call out any metrics or flag changes explicitly because they affect operator dashboards and alerts.

## Domain Knowledge

- This repository ships a Prometheus exporter binary and container image for Control D, not a reusable Go SDK.
- The exporter starts in personal mode by default and fills the `orgId` label with `000000000`; `--controld.business-mode` opts in to organization-scoped metrics.
- Never log `CTRLD_API_KEY` or raw upstream payloads that may contain sensitive data; mask secrets in any diagnostic output.
- Favor low-cardinality labels and predictable metric behavior over exhaustive upstream mirroring.
- The exporter listens on `0.0.0.0:10034` and serves metrics at `/metrics` by default; keep the flag surface (`--web.*`, `--controld.*`, `--log.level`) small and stable.
