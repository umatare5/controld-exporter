<div align="center">

  <img alt="controld-exporter" src="docs/assets/logo.png" width="115px" />

  <h1>controld-exporter</h1>

  <p>A third-party Prometheus Exporter for Control D.</p>

  <p>
    <img alt="GitHub Tag" src="https://img.shields.io/github/v/tag/umatare5/controld-exporter?label=Latest%20version" />
    <a href="https://github.com/umatare5/controld-exporter/actions/workflows/go-test-build.yml"><img alt="Test and Build" src="https://github.com/umatare5/controld-exporter/actions/workflows/go-test-build.yml/badge.svg?branch=main" /></a>
    <a href="https://github.com/umatare5/controld-exporter/actions/workflows/go-vulncheck.yml"><img alt="govulncheck" src="https://github.com/umatare5/controld-exporter/actions/workflows/go-vulncheck.yml/badge.svg?branch=main" /></a><br>
    <a href="https://pkg.go.dev/github.com/umatare5/controld-exporter@main"><img alt="Go Reference" src="https://pkg.go.dev/badge/umatare5/controld-exporter.svg" /></a>
    <a href="./LICENSE"><img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" /></a>
  </p>

</div>

## Overview

This exporter reads the [Control D](https://controld.com/) API and publishes the account state as Prometheus metrics.

- 💚 **Service Health**: Control D's own DNS, API and proxy status, per point of presence
- ⚙️ **Configuration Drift**: Filter, rule and option counts per profile, so an edit shows as a step
- 💳 **Billing Visibility**: Payment status, refund status and the next billing instant
- 🏢 **Organization Scope**: Members, users, routers and profiles, per organization and sub-org

> [!IMPORTANT]
> The exporter needs a Control D API token, which is issued from the account dashboard. See the [Control D Getting Started guide](https://docs.controld.com/reference/get-started) for how to register and create one.

## Quick Start

### 1. Set the API token

```bash
export CTRLD_API_KEY="your-control-d-api-token"
```

### 2. Run the exporter with Docker

```bash
docker run -p 10034:10034 -e CTRLD_API_KEY ghcr.io/umatare5/controld-exporter
```

### 3. Scrape it

```bash
curl -s http://localhost:10034/metrics | head
```

> [!TIP]
> If you prefer using binaries, download them from the [release page](https://github.com/umatare5/controld-exporter/releases).
>
> **Supported Platforms:** `linux_amd64`, `linux_arm64`, `darwin_amd64`, `darwin_arm64` and `windows_amd64`

## Collectors

Every collector runs on each scrape, and no flag turns one off.

| Collector      | Publishes                                           |
| :------------- | :-------------------------------------------------- |
| `billing`      | Amounts, status, refunds and the next billing date  |
| `endpoint`     | Clients counted against each device                 |
| `network`      | Service status per point of presence                |
| `profile`      | Filter, rule and option counts per profile          |
| `service`      | Services in each category                           |
| `organization` | Members, users, routers and profiles, business mode |

> [!NOTE]
> `--controld.business-mode` changes what a collector reads, not whether it runs. See [Collectors](docs/collectors.md).

## Flags

`controld-exporter --help` prints every flag, and [`docs/help.md`](docs/help.md) carries the same list with notes.

- `--controld.api-key` is the only required flag, and `CTRLD_API_KEY` fills it instead.
- The variable keeps the token off the process table. See [Help](docs/help.md#notes).
- `--controld.business-mode` widens the account scope of every series that carries `orgId`.
- `--web.*` set the bind address, the port and the telemetry path.
- `--log.level debug` adds the request URI and the response body of every call.

> [!IMPORTANT]
> The exporter starts in personal mode, where `orgId` reads `000000000` on every series that carries it and the `controld_organization_*` families are absent. `--controld.business-mode` needs a token with organization scope; without one the scrape answers 200 carrying the billing and network families alone. See [Account Scope](docs/README.md#account-scope).

## Endpoints

The exporter registers two paths:

- `/` — landing page, which prints the telemetry path when reached at <http://localhost:10034/>
- `/metrics` — metrics endpoint, configurable via `--web.telemetry-path`

Nothing is cached between scrapes, so a scrape's cost grows with the account rather than with the exporter. See [`docs/README.md`](docs/README.md) for the request count, the timeouts around it, and the paths that fall through to the landing page.

## Metrics

Every series is namespaced `controld_`, and the catalogue lives in `docs/`:

| Page                                 | Covers                                       |
| :----------------------------------- | :------------------------------------------- |
| **[Collectors](docs/collectors.md)** | The six collectors, their metrics and labels |
| **[Help](docs/help.md)**             | Flags and defaults, as `--help` prints       |

The series a dashboard usually starts from:

| Collector      | Metric                              | Type  | Description                        |
| :------------- | :---------------------------------- | :---- | :--------------------------------- |
| `network`      | `controld_network_health_code`      | Gauge | Status of one service at one node  |
| `billing`      | `controld_billing_status`           | Gauge | Transaction status of one payment  |
| `endpoint`     | `controld_endpoint_clients_total`   | Gauge | Clients counted against one device |
| `profile`      | `controld_profile_rules_total`      | Gauge | Rules on one profile               |
| `organization` | `controld_organization_users_total` | Gauge | Users of the organization          |

> [!NOTE]
> See [`docs/README.md`](docs/README.md) for the absence and account-scope rules every collector shares.

> [!IMPORTANT]
> `/network` and `/services/categories` need no token, so a revoked key leaves their families publishing and the scrape answering 200. `ControlDMetricsMissing` in [`examples/prometheus_alert_rules.yml`](examples/prometheus_alert_rules.yml) therefore reads a token-gated family beside the token-free one, because absence over either alone misses what the other catches.

### Exporter Health Metrics

The exporter publishes no series about itself, so a failed scrape shows as missing series.

- **No exporter series** — no scrape duration, no error counter, no `up`-style gauge.
- **No runtime series** — the Go and process collectors sit on a registry no handler serves.
- **Absence is the signal** — a failing collector withholds its family instead of publishing `0`.

> [!NOTE]
> See [Exporter Health](docs/README.md#exporter-health) for what a failed scrape looks like, and [Absence](docs/README.md#absence) for the rules each collector follows.

## Examples

### Prometheus Configuration

#### Job Configuration Example

Add the job from [`examples/prometheus.yml`](examples/prometheus.yml) to your Prometheus configuration.

#### Alerting Rules Configuration Example

Add the rules from [`examples/prometheus_alert_rules.yml`](examples/prometheus_alert_rules.yml) to your configuration.

### Grafana Dashboard

Import [`examples/control-d-exporter-dashboard.json`](examples/control-d-exporter-dashboard.json) to add the dashboard.

![Control D Exporter Dashboard](examples/control-d-exporter-dashboard.png)

> [!NOTE]
> The billing panels name `USD` and `JPY`, so an account settling in another currency needs them edited. See [Dashboards](docs/README.md#dashboards).

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for the development setup, what CI runs on a pull request, the tests, the code style, the documentation conventions and the release process.

## License

MIT. The binary statically links Apache-2.0, MIT and BSD 3-Clause dependencies, whose notices are reproduced in [`NOTICE`](NOTICE) and shipped alongside [`LICENSE`](LICENSE) in every release archive and container image.
