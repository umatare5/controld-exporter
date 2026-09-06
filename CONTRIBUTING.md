# Contributing

The [shared contribution guide](https://github.com/umatare5/.github/blob/main/CONTRIBUTING.md) covers what every exporter shares. This page carries the rest.

## Development

CI runs Format and Lint, Test and Build, Coverage, Prometheus Rules, markdownlint, Link Check, actionlint, CodeQL and govulncheck on every pull request.

## Testing

- **No test exists yet** — the tree carries no `*_test.go`, so `make test-unit` reports zero tests.
- **The threshold is 0 percent** — CI holds it there until the first test lands.
- **The example rules carry the coverage** — `promtool` lints them and runs their assertions in CI.
- **That is the only check** — no other automated test covers the alerting expressions.

Three commands reproduce the `Prometheus Rules` job locally.

```bash
promtool check rules --lint all --lint-fatal examples/prometheus_alert_rules.yml
promtool test rules examples/prometheus_alert_rules_test.yml
promtool check config --lint all --lint-fatal examples/prometheus.yml
```

## Code Style

No `--collector.<name>` flag exists here, because every collector runs on each scrape and `--controld.business-mode` changes what each one reads rather than whether it runs.

A collector that cannot reach Control D returns without describing a metric, so its whole family is absent for that scrape and no path publishes a `0` standing for a failed call.

## Documentation

Every fact has one page that owns it, and the other pages link to it rather than restating it.

| Page                 | Owns                                 |
| :------------------- | :----------------------------------- |
| `README.md`          | What it is, how to run and scrape it |
| `docs/README.md`     | The rules every collector obeys      |
| `docs/collectors.md` | The metric catalogue and the labels  |
| `docs/help.md`       | The verbatim `--help` transcript     |

> [!NOTE]
> `CHANGELOG.md` carries one section per release, each with a `### Metrics` and a `### Flags` subsection reading `None.` where that release changed neither, so a reader learns the surface held rather than inferring it from silence.

## Release

The `VERSION:` line in the [`docs/help.md`](docs/help.md) transcript reads `dev` rather than a release number, because the version is stamped at link time and the transcript comes from a locally built binary. The shared procedure's third step therefore has nothing to edit here, and a release pull request carries `CHANGELOG.md` and `VERSION` alone.
