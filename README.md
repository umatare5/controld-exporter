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
- ⚙️ **Configuration Drift**: Filter, rule and option counts per profile, so an edit is visible as a step
- 💳 **Billing Visibility**: Payment status, refund status and the next billing instant
- 🏢 **Organization Scope**: Members, users, routers and profiles across an organization and its sub-organizations

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

## Flags

`controld-exporter --help` prints every flag, and [`docs/help.md`](docs/help.md) carries the same list.

| Flag                       | Effect                                             |
| :------------------------- | :------------------------------------------------- |
| `--controld.api-key`       | The Control D token, or `CTRLD_API_KEY`            |
| `--controld.business-mode` | Read the organization instead of the account       |
| `--web.listen-port`        | The port the HTTP server binds, `10034` by default |
| `--log.level`              | `debug` adds the request URI and the decoded body  |

> [!IMPORTANT]
> The exporter starts in personal mode, where `orgId` reads `000000000` on every series that carries it and the `controld_organization_*` families are absent. `--controld.business-mode` needs a token with organization scope; without one the exporter terminates on the first scrape.

## Environment Variables

This exporter reads one environment variable:

| Environment Variable | Description                    |
| :------------------- | :----------------------------- |
| `CTRLD_API_KEY`      | Control D API token (required) |

## Endpoints

The exporter serves two endpoints:

- `/` — landing page, which prints the telemetry path when reached at <http://localhost:10034/>
- `/metrics` — metrics endpoint, configurable via `--web.telemetry-path`

Nothing is cached between requests, so a scrape costs one Control D API call per collector and its latency is the API's. See [`docs/README.md`](docs/README.md) for the scrape path and the timeouts around it.

## Metrics

Every series is namespaced `controld_`. The series a dashboard usually starts from:

| Metric                                             | Type    | Description                             |
| :------------------------------------------------- | :------ | :-------------------------------------- |
| `controld_network_health_code`                     | Gauge   | Service status of one point of presence |
| `controld_billing_status`                          | Gauge   | Transaction status of one payment       |
| `controld_billing_subscription_nextbill_timestamp` | Gauge   | Next billing instant, in Unix seconds   |
| `controld_endpoint_clients_total`                  | Gauge   | Clients counted against one device      |
| `controld_profile_rules_total`                     | Gauge   | Rules on one profile                    |
| `controld_service_categories_total`                | Gauge   | Services in one category                |
| `controld_stats_last_queries_count`                | Counter | DNS queries of one verdict              |
| `controld_organization_users_total`                | Gauge   | Users of the organization               |

See [`docs/collectors.md`](docs/collectors.md) for all metrics.

### Exporter Health Metrics

The exporter publishes no series about itself, so a failed scrape shows as missing series.

- **No exporter series** — no scrape duration, no error counter, no `up`-style gauge.
- **No runtime series** — the Go and process collectors sit on a registry no handler serves.
- **Absence is the signal** — a failing collector withholds its family instead of publishing `0`.

> [!IMPORTANT]
> A scrape whose collectors all failed still answers 200 with an empty body, so the target's own `up` stays 1. Alert on the absence of a series the account always has, as [`examples/prometheus_alert_rules.yml`](examples/prometheus_alert_rules.yml) does with `absent()`.

> [!NOTE]
> The failing endpoint and its status are logged at `error`, and `--log.level debug` adds the request URI and the decoded body. See [`docs/README.md`](docs/README.md) for the absence rules each collector follows.

## Use Cases

### Job Configuration Example

Add the job from [`examples/prometheus.yml`](examples/prometheus.yml) to your Prometheus configuration.

### Alerting Rules Configuration Example

Add the rules from [`examples/prometheus_alert_rules.yml`](examples/prometheus_alert_rules.yml) to your configuration.

### Grafana Dashboard

Import [`examples/control-d-exporter-dashboard.json`](examples/control-d-exporter-dashboard.json) to add the dashboard.

![Control D Exporter Dashboard](examples/control-d-exporter-dashboard.png)

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for the make targets, the build and the release procedure.

## Acknowledgement

I launched this project with the help of **GitHub Copilot Coding Assistant**, and I am grateful to the global developer community for their contributions to open source projects and public repositories.

## Licence

MIT. The binary statically links Apache-2.0, MIT and BSD 3-Clause dependencies, whose notices are reproduced in [`NOTICE`](NOTICE) and shipped alongside [`LICENSE`](LICENSE) in every release archive and container image.
