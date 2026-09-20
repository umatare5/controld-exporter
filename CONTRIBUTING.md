# Contributing

The [shared contribution guide](https://github.com/umatare5/.github/blob/main/CONTRIBUTING.md) covers what every exporter shares. This page carries the rest.

## Development

CI decides what to run from the branch and the paths a pull request touches, so a Go change never waits on the link check and a Markdown-only change never waits on govulncheck.

- **Always** — Format and Lint, Test and Build, Coverage, Prometheus Rules and CodeQL, into `main`.
- **Go changes** — govulncheck, when `go.mod`, `go.sum`, a `.go` file or its own workflow moves.
- **Markdown changes** — markdownlint and Link Check, and on any branch rather than `main` alone.
- **Workflow changes** — actionlint, also on any branch.

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

No `--collector.<name>` flag exists here, because every collector runs on each scrape and `--controld.business-mode` changes what three of the six read rather than whether they run. The organization collector is the one whose output the flag decides, because it emits nothing at all in personal mode.

A collector that cannot reach Control D returns without emitting a sample, so its whole family is absent for that scrape and no path publishes a `0` standing for a failed call. The organization collector does not publish a `0` either. It reads the fetch error before it builds anything, so a failed fetch withholds its families and leaves the collectors behind it to publish.

## Documentation

Every fact has one page that owns it, and the other pages link to it rather than restating it.

| Page                 | Owns                                 |
| :------------------- | :----------------------------------- |
| `README.md`          | What it is, how to run and scrape it |
| `docs/README.md`     | The rules every collector obeys      |
| `docs/collectors.md` | The metric catalogue and the labels  |
| `docs/help.md`       | The verbatim `--help` transcript     |

> [!NOTE]
> `CHANGELOG.md` carries one section per release, and every section since v1.2.0 adds a `### Metrics` and a `### Flags` subsection reading `None.` where that release changed neither, so a reader learns the surface held rather than inferring it from silence.

## Release

The `VERSION:` line in the [`docs/help.md`](docs/help.md) transcript reads `dev` rather than a release number, because the transcript comes from a build that stamps nothing, while `make build` and a release both stamp it. The shared procedure's third step therefore has nothing to edit here, and a release pull request need carry no more than `CHANGELOG.md` and `VERSION`.
