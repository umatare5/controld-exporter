# Documentation

Reference pages for controld-exporter. The [README](../README.md) covers getting the exporter running and scraped; these pages carry the metric catalogue and the behaviour every collector shares.

| Page                              | Focus                                  |
| :-------------------------------- | :------------------------------------- |
| [Collectors](collectors.md)       | The metric families and their labels   |
| [Configuration](configuration.md) | Flags and defaults, as `--help` prints |

## Technical information

### Scrape path

A scrape is served from Control D directly rather than from a cached snapshot, so its latency is the API's and its data is as fresh as the request.

- **One collector per scrape** — the handler builds a new registry and collector each time, so nothing survives between scrapes.
- **Sequential** — the seven collectors run in order on one goroutine, and a slow endpoint delays every collector after it.
- **Unbounded per request** — the API client carries no timeout and no retry, so a connection that never answers holds the scrape open.
- **Overlap is not prevented** — a second scrape arriving during the first opens its own set of requests against the same key.
- **Interval** — 60s with `scrape_timeout: 50s` is what [`examples/prometheus.yml`](../examples/prometheus.yml) sets, for those two reasons together.

> [!IMPORTANT]
> The server's one-minute write timeout fails the response to Prometheus without cancelling the requests behind it, so a hung Control D connection leaks the work rather than ending it. Keep `scrape_timeout` below that timeout so Prometheus gives up first and records the failure.

### Absence

A Control D call that fails withholds every series behind it — never `0`, and never the previous scrape's value.

- **Failure is total per collector** — one non-2xx status, one JSON error or one `"success": false` body drops that collector's whole family.
- **Empty is not zero** — an account with no payment, device or profile publishes no series for it rather than a count of nothing.
- **Staleness closes the gap** — Prometheus marks a series stale after the scrape that stops carrying it, so a dashboard shows a break rather than a flat line.
- **Only the log names the cause** — the status and the endpoint are logged at `error`, and no series records that a collector failed.

> [!NOTE]
> A scrape whose every collector failed still answers 200 with an empty body, so the target's own `up` stays 1 and cannot distinguish a healthy account with nothing configured from an API key that has been revoked. `absent()` over a series the account always has is what closes that gap, which is what `ControlDMetricsMissing` in [`examples/prometheus_alert_rules.yml`](../examples/prometheus_alert_rules.yml) does.

### Counter semantics

One series is declared a counter, and it does not behave as one: `controld_stats_last_queries_count` carries the newest one-minute bucket of the query report, so it rises and falls with traffic.

- **`rate()` reads it as a reset** — every scrape whose bucket is smaller than the last looks like a counter restart.
- **Ratios are safe** — dividing one verdict by the sum over all of them uses the raw values and needs no range.
- **Every other series is a gauge** — a configuration count or a status code, which Control D restates in full on each request.

### Account scope

`--controld.business-mode` decides which Control D account the same API key is read against, and the `orgId` label records that decision on every series that carries it.

| Mode                       | Reads                             | `orgId` holds     |
| :------------------------- | :-------------------------------- | :---------------- |
| Personal, the default      | The account the key belongs to    | `000000000`       |
| `--controld.business-mode` | The organization and each sub-org | The Control D key |

Business mode reads a sub-organization by repeating the same request under an `X-Force-Org-Id` header, so the request count grows with the sub-organization count. The API key needs organization scope for any of it to succeed.

### Exporter health

The exporter publishes no series about itself: no scrape duration, no error counter, and no `up`-style gauge.

- **No runtime metrics either** — the Go and process collectors are registered on a registry no handler serves, so `/metrics` carries `controld_` series alone.
- **Failure is read from absence** — an alert on a series the account always has is the only signal a scrape produced nothing.
- **The log carries the diagnosis** — the failing endpoint and its status are logged at `error`, and `--log.level debug` adds the request URI and the decoded body.

### Dashboards

[`examples/control-d-exporter-dashboard.json`](../examples/control-d-exporter-dashboard.json) is a Grafana schema covering both modes. Its `$orgID`, `$profileName` and `$queryType` variables are populated from the label values the exporter publishes, so a personal-mode target offers `000000000` alone.

- **Currencies are hard-coded** — the billing panels name `USD` and `JPY`, so an account settling in another currency needs the expression edited.
- **The query panel uses `increase()`** — which the counter semantics above make unreliable, and the ratio form in the alert rules is what to replace it with.
