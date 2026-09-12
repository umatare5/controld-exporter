# Documentation

Reference pages for controld-exporter. The [README](../README.md) covers getting the exporter running and scraped; these pages carry the metric catalogue and the behaviour every collector shares.

| Page                        | Focus                                  |
| :-------------------------- | :------------------------------------- |
| [Collectors](collectors.md) | The metric families and their labels   |
| [Help](help.md)             | Flags and defaults, as `--help` prints |

## Technical Information

### Scrape Path

Every scrape reads Control D rather than a cached snapshot, so its latency is the API's.

- **Nothing survives a scrape** — the handler builds a registry and a collector per request.
- **Sequential** — the seven collectors share one goroutine, so a slow call delays the rest.
- **The cost is the account's** — personal mode makes 7 requests, business mode 9 plus 4 per sub-org.
- **Unbounded per request** — the API client carries no timeout and no retry.
- **Overlap is not prevented** — a second scrape opens its own requests against the same key.
- **The interval is the only throttle** — no response carries `X-RateLimit-*` or `Retry-After`.
- **Interval** — [`examples/prometheus.yml`](../examples/prometheus.yml) sets 60s with `scrape_timeout: 50s`.

> [!IMPORTANT]
> The server's one-minute write timeout fails the response to Prometheus without cancelling the requests behind it, so a hung Control D connection leaks the work rather than ending it. Keep `scrape_timeout` below that timeout so Prometheus gives up first and records the failure.

### Endpoints

Both paths answer on the address `--web.listen-address` and `--web.listen-port` bind, and neither authenticates — [`SECURITY.md`](../SECURITY.md) carries the surface that leaves exposed.

- **`/` is the fallback** — every path no handler claims routes to the landing page, so a mistyped telemetry path answers 200 with HTML rather than 404.
- **`/metrics` moves with the flag** — the path it left then falls through to the landing page.

### Absence

A Control D call that fails withholds the series behind it — never `0`, never the last value.

- **Failure is scoped to the call** — one non-2xx status, one JSON error or one `"success": false` body withholds only the series that call feeds.
- **A collector can fail in part** — billing reads payments and subscriptions independently, and a sub-organization whose request fails is skipped while the others still publish.
- **Empty is not zero** — an account with no payment, device or profile publishes no series at all.
- **A vanished field is zero** — a field the API stops sending decodes to `0` and publishes as `0`, because only a transport, status, envelope or type error withholds a family.
- **Staleness closes the gap** — Prometheus marks a series stale after the scrape that stops carrying it, so a dashboard shows a break rather than a flat line.
- **Only the log names the cause** — no series records that a collector failed.

> [!IMPORTANT]
> `/network` and `/services/categories` need no token, so their families keep publishing after a key is revoked and a scrape still answers 200. `absent(controld_network_health_code)`, which `ControlDMetricsMissing` in [`examples/prometheus_alert_rules.yml`](../examples/prometheus_alert_rules.yml) reads, therefore cannot see a revoked key. Alert on a family the token gates and the account fills, such as `controld_profile_rules_total`, because an empty account is silent too.

> [!WARNING]
> A failed `/organizations/organization` call withholds more than the organization families, because the endpoint, profile, service and stats collectors read that same response and skip with it. Business mode then answers 200 carrying the billing and network families alone, and each of the five re-requests the failed endpoint rather than sharing one failure.

### Counter Semantics

One series is declared a counter, and it does not behave as one: `controld_stats_last_queries_count` carries a one-minute bucket of the query report, so it rises and falls with traffic.

- **`rate()` reads it as a reset** — a bucket smaller than the last looks like a counter restart.
- **Ratios are safe** — one verdict over the sum of all of them needs no range.
- **Every other series is a gauge** — a configuration count or a status code, restated in full.
- **The `_total` suffix is not a counter** — nineteen gauge families carry it.

### Account Scope

`--controld.business-mode` decides how much of a Control D account is read, and the `orgId` label records that decision on every series that carries it.

| Mode                       | Reads                             | `orgId` holds              |
| :------------------------- | :-------------------------------- | :------------------------- |
| Personal, the default      | The account the key belongs to    | `000000000`                |
| `--controld.business-mode` | The organization and each sub-org | The org's or sub-org's key |

Business mode reads a sub-organization by repeating the same request under an `X-Force-Org-Id` header, so the request count grows with the sub-organization count. The API key needs organization scope for any of it to succeed, and it travels in the `Authorization` header alone — no label ever carries it.

### Exporter Health

The exporter publishes no series about itself: no scrape duration, no error counter, no gauge.

- **No runtime metrics either** — the Go and process collectors are registered on a registry no handler serves, so `/metrics` carries `controld_` series alone.
- **Failure is read from absence** — a scrape whose collectors all failed still answers 200 with an empty body, so `up` stays 1 and only an absent family shows it.
- **Both modes read alike** — a failed organization fetch withholds the families it feeds rather than the exposition, because the collector returns before it builds a sample.
- **The log carries the diagnosis** — the failing endpoint and its status are logged at `error`, and `--log.level debug` adds the request URI and the response body as received.

### Dashboards

[`examples/control-d-exporter-dashboard.json`](../examples/control-d-exporter-dashboard.json) is a Grafana schema covering both modes. Its `$orgID`, `$profileName` and `$queryType` variables are populated from the label values the exporter publishes, so a personal-mode target offers `000000000` alone.

- **Currencies are hard-coded** — the billing panels name `USD` and `JPY`, so an account settling in another currency needs the expression edited.
- **The blocking-rate panels use `increase()`** — which the counter semantics above make unreliable, and the ratio form in the alert rules is what to replace it with.
