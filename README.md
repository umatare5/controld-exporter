<div align="center">

  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="docs/assets/logo.png" width="115px" />
    <source media="(prefers-color-scheme: light)" srcset="docs/assets/logo.png" width="115px" />
    <img alt="controld-exporter" src="docs/assets/logo.png" width="115px" />
  </picture>

  <h1>controld-exporter</h1>

  <p>A third-party Prometheus Exporter for Control D.</p>

  <p>
    <img alt="GitHub Tag" src="https://img.shields.io/github/v/tag/umatare5/controld-exporter?label=Latest%20version" />
    <a href="https://github.com/umatare5/controld-exporter/actions/workflows/go-test-build.yml"><img alt="Test and Build" src="https://github.com/umatare5/controld-exporter/actions/workflows/go-test-build.yml/badge.svg?branch=main" /></a>
    <a href="https://github.com/umatare5/controld-exporter/actions/workflows/go-vulncheck.yml"><img alt="govulncheck" src="https://github.com/umatare5/controld-exporter/actions/workflows/go-vulncheck.yml/badge.svg?branch=main" /></a>
    <a href="./LICENSE"><img alt="License: MIT" src="https://img.shields.io/badge/License-MIT-yellow.svg" /></a>
  </p>

</div>

## Overview

This exporter allows a Prometheus instance to monitor service health, profiles, billing on [Control D](https://controld.com/).

- 💚 **Service Health**: Watch health checks for DNS, API, and proxy services.
- 💰️ **Billing Visibility**: Watch billing status, refund status and the next billing instant.
- ⚙️ **Configuration Audit**: Track changes in predefined and custom settings with trend visualization.
- 🏢 **Organization Support**: Fetch members, profiles, routers and users for the organization.

> [!NOTE]
>
> Control D is a subscription service, so this exporter needs a token from a paid account once the free trial ends. For the plans and their limits, refer to [Control D - Personal Plans](https://controld.com/plans) and [Control D - Business Pricing](https://controld.com/pricing).

## Installation

This exporter supports both container images and OS-specific binaries installations.

**A. Using Container**

```bash
docker pull ghcr.io/umatare5/controld-exporter
```

**B. Using OS-Specific binaries**

Download from [Releases](https://github.com/umatare5/controld-exporter/releases). `linux_(amd64|arm64)`, `darwin_(amd64|arm64)` and `windows_amd64` are supported.

## Quick Start

This exporter needs an API key. See **[Control D Getting Started Guide](https://docs.controld.com/reference/get-started)** to get an API key first.

**1. Set the API key**

```bash
export CTRLD_API_KEY="your-control-d-api-token"
```

**2. Run the exporter with Docker**

```bash
docker run -p 10034:10034 -e CTRLD_API_KEY ghcr.io/umatare5/controld-exporter:v1.2.1
```

**3. Scrape it**

```bash
curl -s http://localhost:10034/metrics
```

> [!TIP]
>
> See [Metrics](#metrics) for available metrics, and [Prometheus Configuration](#prometheus-configuration) for the job and the alerting rules.

## Flags

The exporter supports the following command-line flags:

```text
NAME:
   controld-exporter - A Prometheus exporter for metrics from the Control D

USAGE:
   controld-exporter [options...]

VERSION:
   1.2.1

GLOBAL OPTIONS:
   --web.listen-address string             Address to bind the HTTP server to. (default: "0.0.0.0")
   --web.listen-port int                   Port number to bind the HTTP server to. (default: 10034)
   --web.telemetry-path string, -p string  Path for the metrics endpoint. (default: "/metrics")
   --controld.api-key string, -k string    API key for authenticating with the Control D API. [$CTRLD_API_KEY]
   --controld.business-mode                Enable the metrics collection available in the business subscription.
   --log.level string                      Set the logging level. One of: [debug, info, warn, error] (default: "info")
   --help, -h                              show help
   --version, -v                           print the version
```

## Endpoints

The exporter serves two endpoints. See [Endpoints](docs/architecture.md#endpoints) for the details.

| Path       | Detail                                             |
| :--------- | :------------------------------------------------- |
| `/`        | Landing page, reached at <http://localhost:10034/> |
| `/metrics` | Metrics endpoint, moved by `--web.telemetry-path`  |

## Metrics

This exporter exposes metrics for the Control D API state.

### Collector Metrics

The following table lists the metrics this exporter publishes. The `controld_organization_*` and `controld_sub_organization_*` families need `--controld.business-mode`. See **Appendix** below the table for more details.

| Metric                                             | Type  | Description                             |
| :------------------------------------------------- | :---- | :-------------------------------------- |
| `controld_billing_status`                          | Gauge | Transaction status of one payment       |
| `controld_billing_refunded`                        | Gauge | Refund status of one payment            |
| `controld_billing_subscription_amount_total`       | Gauge | Amount of one payment, per currency     |
| `controld_billing_subscription_nextbill_timestamp` | Gauge | Next billing instant, in Unix seconds   |
| `controld_endpoint_clients_total`                  | Gauge | Clients counted against one device      |
| `controld_network_health_code`                     | Gauge | Service status of one point of presence |
| `controld_profile_preset_filters_total`            | Gauge | Preset filters on one profile           |
| `controld_profile_content_filters_total`           | Gauge | Content filters on one profile          |
| `controld_profile_ip_filters_total`                | Gauge | IP filters on one profile               |
| `controld_profile_rules_total`                     | Gauge | Rules on one profile                    |
| `controld_profile_services_total`                  | Gauge | Service filters on one profile          |
| `controld_profile_groups_total`                    | Gauge | Group filters on one profile            |
| `controld_profile_enabled_option_total`            | Gauge | Enabled options on one profile          |
| `controld_service_categories_total`                | Gauge | Services in one category                |
| `controld_organization_members_total`              | Gauge | Members of the organization             |
| `controld_organization_profiles_total`             | Gauge | Profiles of the organization            |
| `controld_organization_users_total`                | Gauge | Users of the organization               |
| `controld_organization_routers_total`              | Gauge | Routers of the organization             |
| `controld_organization_sub_orgs_total`             | Gauge | Sub-organizations beneath it            |
| `controld_sub_organization_members_total`          | Gauge | Members of one sub-organization         |
| `controld_sub_organization_profiles_total`         | Gauge | Profiles of one sub-organization        |
| `controld_sub_organization_users_total`            | Gauge | Users of one sub-organization           |
| `controld_sub_organization_routers_total`          | Gauge | Routers of one sub-organization         |

<details><summary><b>Appendix - Collector Metrics Details</b></summary><p>

#### About the Metrics

**the four `controld_billing_*` series**: they read the account's own payment history, which the organization endpoints do not scope. Business mode therefore publishes them under the payment's `id` alone, with no `orgId` to separate them by.

**`controld_endpoint_clients_total`**: it counts the clients Control D attributes to one device, keyed by the device's name. A device renamed in the dashboard ends one series and opens another with the count carried over.

**`controld_network_health_code`**: the value is the `api`, `dns` and `pxy` integer each node publishes, passed through without interpretation. A code this exporter has never seen reaches Prometheus as readily as a familiar one.

**the seven `controld_profile_*` series**: they count what each profile has configured rather than what it matched, so they move when an operator edits a profile and stay flat under any amount of traffic.

**the nine `controld_organization_*` and `controld_sub_organization_*` series**: they need `--controld.business-mode` and an API key belonging to an organization. Neither is published in personal mode, so a personal-mode dashboard shows no data rather than zeros.

#### About the Labels

| Label                      | Description                                               |
| :------------------------- | :-------------------------------------------------------- |
| `id`                       | The payment's or subscription's Control D primary key     |
| `currency`                 | The ISO code the amount beside it is denominated in       |
| `name`                     | The object's own name, or the category's key on `service` |
| `orgId`                    | The account scope the series was read under               |
| `city_name`/`country_name` | Where Control D places the point of presence              |
| `iata_code`                | The airport code Control D identifies that node by        |
| `service_name`             | `api`, `dns` or `proxy`, one series each per node         |

**`name`**: The field is fixed per family rather than chosen per series. The device, profile and organization families carry the name an operator gave the object, so renaming one in the Control D dashboard ends the old series and opens a new one. `controld_service_categories_total` carries the category's primary key instead, although the endpoint supplies a name beside it.

**`orgId`**: Personal mode fills it with `000000000`, a value no Control D organization holds, so a dashboard written against it survives being pointed at a business account. Business mode fills it with the organization's own primary key on the series read for the account, and with a sub-organization's key on the series read for that sub-organization. It never carries the API key, which travels in the `Authorization` header alone.

</p></details>

### Exporter Health Metrics

The exporter publishes no series about itself, so a failed scrape shows as missing series.

## Examples

### Exporter Configuration

By default, the exporter runs in personal mode:

```bash
$ CTRLD_API_KEY="your-control-d-api-token" ./controld-exporter
time="2026-01-01T00:00:00+09:00" level=info msg="Starting the personal mode exporter on port 10034."
```

To run it for the organizations, activate the business mode with `--controld.business-mode`.

### Prometheus Configuration

There are several Prometheus configuration examples provided below:

- **Example Job:** Add from [`examples/prometheus.yml`](./examples/prometheus.yml) to your Prometheus.
- **Example Alerting Rules:** Add from [`examples/prometheus_alert_rules.yml`](./examples/prometheus_alert_rules.yml) to your Prometheus.

### Grafana Configuration

Import [`examples/control-d-exporter-dashboard.json`](./examples/control-d-exporter-dashboard.json) and visualize the metrics.

<picture>
  <img alt="Control D Exporter Dashboard" src="examples/control-d-exporter-dashboard.png">
</picture>

## Documentation

The reference pages under [`docs/`](docs/) carry the behaviour behind the metrics above.

- **[Architecture](docs/architecture.md)** – the scrape path, the absence rules and others.

## Contributing

See [`CONTRIBUTING.md`](CONTRIBUTING.md) for the development setup, test conventions and others.

## License

MIT. The binary statically links Apache-2.0, MIT and BSD 3-Clause dependencies, whose notices are reproduced in [`NOTICE`](NOTICE) and shipped alongside [`LICENSE`](LICENSE) in every release archive and container image.
