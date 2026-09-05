# Collectors

Every collector runs on each scrape, in the order below, and no flag turns one off. `--controld.business-mode` changes what each one reads rather than whether it runs: personal mode reads the account the key belongs to, business mode the organization beneath it.

A collector that cannot reach Control D withholds its series for that scrape — the [absence rules](README.md#absence) carry what that looks like in a query.

## Metrics

| Collector      | Metric                                             | Type    | Description                             |
| :------------- | :------------------------------------------------- | :------ | :-------------------------------------- |
| `billing`      | `controld_billing_status`                          | Gauge   | Transaction status of one payment       |
| `billing`      | `controld_billing_refunded`                        | Gauge   | Refund status of one payment            |
| `billing`      | `controld_billing_subscription_amount_total`       | Gauge   | Amount of one payment, per currency     |
| `billing`      | `controld_billing_subscription_nextbill_timestamp` | Gauge   | Next billing instant, in Unix seconds   |
| `endpoint`     | `controld_endpoint_clients_total`                  | Gauge   | Clients counted against one device      |
| `network`      | `controld_network_health_code`                     | Gauge   | Service status of one point of presence |
| `profile`      | `controld_profile_preset_filters_total`            | Gauge   | Preset filters on one profile           |
| `profile`      | `controld_profile_content_filters_total`           | Gauge   | Content filters on one profile          |
| `profile`      | `controld_profile_ip_filters_total`                | Gauge   | IP filters on one profile               |
| `profile`      | `controld_profile_rules_total`                     | Gauge   | Rules on one profile                    |
| `profile`      | `controld_profile_services_total`                  | Gauge   | Service filters on one profile          |
| `profile`      | `controld_profile_groups_total`                    | Gauge   | Group filters on one profile            |
| `profile`      | `controld_profile_enabled_option_total`            | Gauge   | Enabled options on one profile          |
| `service`      | `controld_service_categories_total`                | Gauge   | Services in one category                |
| `stats`        | `controld_stats_last_queries_count`                | Counter | DNS queries of one verdict              |
| `organization` | `controld_organization_members_total`              | Gauge   | Members of the organization             |
| `organization` | `controld_organization_profiles_total`             | Gauge   | Profiles of the organization            |
| `organization` | `controld_organization_users_total`                | Gauge   | Users of the organization               |
| `organization` | `controld_organization_routers_total`              | Gauge   | Routers of the organization             |
| `organization` | `controld_organization_sub_orgs_total`             | Gauge   | Sub-organizations beneath it            |
| `organization` | `controld_sub_organization_members_total`          | Gauge   | Members of one sub-organization         |
| `organization` | `controld_sub_organization_profiles_total`         | Gauge   | Profiles of one sub-organization        |
| `organization` | `controld_sub_organization_users_total`            | Gauge   | Users of one sub-organization           |
| `organization` | `controld_sub_organization_routers_total`          | Gauge   | Routers of one sub-organization         |

## Labels

No label is shared across every family: the billing series key on the payment, the network series on the point of presence, and the rest on a Control D object and the account scope it was read under.

| Label                      | Description                                                  |
| :------------------------- | :----------------------------------------------------------- |
| `id`                       | The payment's or subscription's Control D primary key        |
| `currency`                 | The ISO code the amount beside it is denominated in          |
| `name`                     | The object's own name, or its primary key where unnamed      |
| `orgId`                    | The account scope the series was read under                  |
| `city_name`/`country_name` | Where Control D places the point of presence                 |
| `iata_code`                | The airport code Control D identifies that node by           |
| `service_name`             | `api`, `dns` or `proxy`, one series each per node            |
| `type`                     | The verdict a one-minute bucket of queries was counted under |

**`name`**

A profile and an organization carry the name an operator gave them, so renaming one in the Control D dashboard ends the old series and opens a new one. `controld_service_categories_total` carries the category's primary key here instead, because the categories endpoint publishes no separate display name.

**`orgId`**

Personal mode fills it with `000000000`, a value no Control D organization holds, so a dashboard written against it survives being pointed at a business account. Business mode fills it with the organization's key on the account-wide series, and with the sub-organization's key on everything read through `X-Force-Org-Id`.

**`type`**

The report returns a verdict code rather than a name, and the exporter maps `0` to `blocked`, `1` to `bypassed` and `3` to `redirected`. Every other code folds into `unknown`, so a verdict Control D adds after this release lands there rather than opening a series nobody is alerting on.

## Specifications

Each entry carries what the series' HELP text and the shared rules in [Documentation](README.md#technical-information) do not.

**the four `controld_billing_*` series**

they read the account's own payment history, which the organization endpoints do not scope, so business mode publishes them under the payment's `id` alone and carries no `orgId` to separate them by.

- `controld_billing_status` and `controld_billing_refunded` carry the `tx_status` and `tx_refunded` integers unchanged, so the meaning of a non-zero value is Control D's rather than this exporter's.
- The history is unbounded upstream: every payment the account ever made keeps its own series, so the family grows by one `id` per billing period and never shrinks.
- `controld_billing_subscription_nextbill_timestamp` comes from the subscription list rather than the payment list, so its `id` values name subscriptions and join to nothing in the other three.

> [!IMPORTANT]
> `controld_billing_subscription_amount_total` publishes two series per payment: one labelled `currency="USD"` carrying the `amount` field, and one labelled with the payment's own currency carrying `currency_amount`. An account billed in USD produces the same label set twice, which the registry rejects for the whole scrape, so this family is usable only where the account settles in something other than USD.

**`controld_endpoint_clients_total`**

it counts the clients Control D attributes to one device, keyed by the device's name, so a device renamed in the dashboard ends one series and opens another with the count carried over.

**`controld_network_health_code`**

the value is the `api`, `dns` and `pxy` integer each node publishes, passed through without interpretation, so a code this exporter has never seen reaches Prometheus as readily as a familiar one.

- The series describe Control D's own infrastructure rather than the account, so they are identical for every exporter reading the same region and duplicate across targets.
- The node list is whatever `/network` returns at scrape time, so a point of presence withdrawn upstream stops publishing rather than reading unhealthy.

**the seven `controld_profile_*` series**

they count what each profile has configured rather than what it matched, so they move when an operator edits a profile and stay flat under any amount of traffic.

> [!IMPORTANT]
> `controld_profile_ip_filters_total` publishes the content-filter count rather than the IP-filter count, so it duplicates `controld_profile_content_filters_total` on every profile. Read the IP-filter count from the Control D dashboard until this is corrected.

**`controld_stats_last_queries_count`**

it carries the newest one-minute bucket of the DNS query report rather than a running total, so its value falls whenever traffic falls and `rate()` over it reads as a counter reset. Take ratios from the raw values instead.

- The report is fetched with a start timestamp one minute behind the scrape, so a scrape interval other than 60s either double-counts a bucket or skips one.
- It is declared to Prometheus as a counter, which is what makes the `_count` suffix and the type disagree; the alert rules in [`examples/prometheus_alert_rules.yml`](../examples/prometheus_alert_rules.yml) are written around that.
- Control D withdrew the analytics endpoint this series reads, so it has published nothing since then and the failure appears in the log as a non-2xx status rather than as a zero.

**the nine `controld_organization_*` and `controld_sub_organization_*` series**

they need `--controld.business-mode` and an API key belonging to an organization, and neither is published in personal mode at all — a personal-mode dashboard shows no data rather than zeros.

- Sub-organization series are read one request per sub-organization, so a scrape's duration grows linearly with the number of sub-organizations beneath the account.
- The main organization's response also supplies the analytics hostname the `stats` collector uses, so an organization fetch that fails takes the query counts with it.

> [!WARNING]
> A failed `/organizations/organization` call in business mode reaches the metric-building code with no response to read, which terminates the process rather than skipping the scrape. Keep personal mode until an organization is actually configured.
