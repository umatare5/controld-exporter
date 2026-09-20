# Collectors

Every collector runs on each scrape and no flag turns one off. `--controld.business-mode` changes what most of them read rather than whether they run: personal mode reads the account the key belongs to, business mode the organization beneath it.

A collector that cannot reach Control D withholds its series for that scrape — the [absence rules](README.md#absence) carry what that looks like in a query, and how wide one failed call reaches.

## Metrics

| Collector      | Metric                                             | Type  | Description                             |
| :------------- | :------------------------------------------------- | :---- | :-------------------------------------- |
| `billing`      | `controld_billing_status`                          | Gauge | Transaction status of one payment       |
| `billing`      | `controld_billing_refunded`                        | Gauge | Refund status of one payment            |
| `billing`      | `controld_billing_subscription_amount_total`       | Gauge | Amount of one payment, per currency     |
| `billing`      | `controld_billing_subscription_nextbill_timestamp` | Gauge | Next billing instant, in Unix seconds   |
| `endpoint`     | `controld_endpoint_clients_total`                  | Gauge | Clients counted against one device      |
| `network`      | `controld_network_health_code`                     | Gauge | Service status of one point of presence |
| `profile`      | `controld_profile_preset_filters_total`            | Gauge | Preset filters on one profile           |
| `profile`      | `controld_profile_content_filters_total`           | Gauge | Content filters on one profile          |
| `profile`      | `controld_profile_ip_filters_total`                | Gauge | IP filters on one profile               |
| `profile`      | `controld_profile_rules_total`                     | Gauge | Rules on one profile                    |
| `profile`      | `controld_profile_services_total`                  | Gauge | Service filters on one profile          |
| `profile`      | `controld_profile_groups_total`                    | Gauge | Group filters on one profile            |
| `profile`      | `controld_profile_enabled_option_total`            | Gauge | Enabled options on one profile          |
| `service`      | `controld_service_categories_total`                | Gauge | Services in one category                |
| `organization` | `controld_organization_members_total`              | Gauge | Members of the organization             |
| `organization` | `controld_organization_profiles_total`             | Gauge | Profiles of the organization            |
| `organization` | `controld_organization_users_total`                | Gauge | Users of the organization               |
| `organization` | `controld_organization_routers_total`              | Gauge | Routers of the organization             |
| `organization` | `controld_organization_sub_orgs_total`             | Gauge | Sub-organizations beneath it            |
| `organization` | `controld_sub_organization_members_total`          | Gauge | Members of one sub-organization         |
| `organization` | `controld_sub_organization_profiles_total`         | Gauge | Profiles of one sub-organization        |
| `organization` | `controld_sub_organization_users_total`            | Gauge | Users of one sub-organization           |
| `organization` | `controld_sub_organization_routers_total`          | Gauge | Routers of one sub-organization         |

## Labels

No label is shared across every family. The billing series key on the payment and the network series on the point of presence, while the rest key on a Control D object, plus the account scope it was read under.

| Label                      | Description                                               |
| :------------------------- | :-------------------------------------------------------- |
| `id`                       | The payment's or subscription's Control D primary key     |
| `currency`                 | The ISO code the amount beside it is denominated in       |
| `name`                     | The object's own name, or the category's key on `service` |
| `orgId`                    | The account scope the series was read under               |
| `city_name`/`country_name` | Where Control D places the point of presence              |
| `iata_code`                | The airport code Control D identifies that node by        |
| `service_name`             | `api`, `dns` or `proxy`, one series each per node         |

**`name`**

The field is fixed per family rather than chosen per series. The device, profile and organization families carry the name an operator gave the object, so renaming one in the Control D dashboard ends the old series and opens a new one. `controld_service_categories_total` carries the category's primary key instead, although the endpoint supplies a name beside it.

Control D does not require a device or profile name to be unique, and the exporter separates these series by name and `orgId` alone. Two objects sharing both produce one series rather than two: the registry keeps whichever the API listed first and drops the other without a log line, so a count silently goes missing.

**`orgId`**

Personal mode fills it with `000000000`, a value no Control D organization holds, so a dashboard written against it survives being pointed at a business account. Business mode fills it with the organization's own primary key on the series read for the account, and with a sub-organization's key on the series read for that sub-organization. It never carries the API key, which travels in the `Authorization` header alone.

## Specifications

Each entry carries what the series' HELP text and the shared rules in [Documentation](README.md#technical-information) do not.

**the four `controld_billing_*` series**

they read the account's own payment history, which the organization endpoints do not scope, so business mode publishes them under the payment's `id` alone and carries no `orgId` to separate them by.

- `controld_billing_status` and `controld_billing_refunded` carry the `tx_status` and `tx_refunded` integers unchanged, so the meaning of a non-zero value is Control D's rather than this exporter's.
- The history is unbounded upstream: every payment the account ever made keeps its own series, so the family grows by one `id` per billing period and never shrinks.
- `controld_billing_subscription_nextbill_timestamp` comes from the subscription list rather than the payment list, so its `id` values name subscriptions and join to nothing in the other three.
- `controld_billing_subscription_amount_total` is built from the payment list despite its name and its HELP text, both of which say subscription.

> [!IMPORTANT]
> `controld_billing_subscription_amount_total` publishes two series per payment. One carries the `amount` field under `currency="USD"`, the other `currency_amount` under the payment's own currency. An account billed in USD produces that label set twice, and the registry keeps the first and drops the second without failing the scrape, so the family reports `amount` alone.

**`controld_endpoint_clients_total`**

it counts the clients Control D attributes to one device, keyed by the device's name, so a device renamed in the dashboard ends one series and opens another with the count carried over.

**`controld_network_health_code`**

the value is the `api`, `dns` and `pxy` integer each node publishes, passed through without interpretation, so a code this exporter has never seen reaches Prometheus as readily as a familiar one.

- The series describe Control D's own infrastructure rather than the account, so they are identical for every exporter reading the same region and duplicate across targets.
- The node list is whatever `/network` returns at scrape time, so a point of presence withdrawn upstream stops publishing rather than reading unhealthy.
- `-1` means the service is not offered rather than down, and `proxy` reads it on most nodes, so a rule on `!= 1` fires on all of them. `NetworkServiceDown` in [`examples/prometheus_alert_rules.yml`](../examples/prometheus_alert_rules.yml) tests `== 0` for that reason, and so reports an outage but never an absent proxy.
- The endpoint also names the node that served the call, which is the only one a scrape proves reachable, and the exporter publishes no series for it.

**the seven `controld_profile_*` series**

they count what each profile has configured rather than what it matched, so they move when an operator edits a profile and stay flat under any amount of traffic.

**the nine `controld_organization_*` and `controld_sub_organization_*` series**

they need `--controld.business-mode` and an API key belonging to an organization, and neither is published in personal mode at all — a personal-mode dashboard shows no data rather than zeros.

- The four `controld_sub_organization_*` families come from one sub-organization listing the collector reads in memory, so they cost one request however many sub-organizations exist. The per-sub-organization cost is in the device, profile, service and query-report calls, which repeat once each under `X-Force-Org-Id`.
- Both organization responses are fetched once per scrape and dropped after it, and only a success is shared, so a failing call repeats for each of the five collectors that read it.

> [!WARNING]
> This collector runs first and hands its response to three others, so a failed `/organizations/organization` call costs more than the nine families here. The endpoint, profile and service collectors read the same response and skip with it. The scrape still answers 200, `up` stays 1, and only the log names the endpoint and its status.
