# Contributing

Thank you for your interest in contributing to the controld-exporter.

Please follow **[the shared contribution guide](https://github.com/umatare5/.github/blob/main/CONTRIBUTING.md)**, which covers:

- **Development** – the tools to install and the order the pre-commit hooks run in.
- **Command** – the `make` targets and what each one does.
- **Testing** – test placement, mutation checks and what a fixture must carry.
- **Documentation** – page ownership, pinned headings and the verbatim `--help` transcript.
- **Release** – the three files a release touches and what a push to `main` triggers.
- **Pull Requests** – the branch, commit and changelog steps, and what never enters a commit.

This page specifies what is particular to this one.

## Development

These points are where this repository departs from the shared defaults.

- **Do not assume every check runs.** Four are path-filtered: govulncheck, markdownlint, Link Check and actionlint.
- **Keep coverage at 80 percent.** The gate fails the build below it, and the README badge carries the measured figure.
- **Stamp the version with `make build`.** The `--help` transcript reads `dev` until the build stamps it.

If the API key belongs to an organization:

- **Set `--controld.business-mode`.** It changes what collectors read rather than turning them on or off.

## Testing

Go tests cover the client and the collectors, and the example rules are checked separately.

- **Serve a measured reply.** The fixtures live in [`internal/controld/testdata/`](internal/controld/testdata).
- **Redirect `baseURL` with `controld.NewClientWithBaseURL`.** A test outside the package has no other seam.
- **Redact before a reply enters the tree.** Replace the keys, URLs and names, and keep the counts a metric reads.
- **Check rules with `promtool`.** No Go test covers the alerting expressions, so CI lints and asserts them there.

The `Prometheus Rules` job runs these three commands.

```bash
promtool check rules --lint all --lint-fatal examples/prometheus_alert_rules.yml
promtool test rules examples/prometheus_alert_rules_test.yml
promtool check config --lint all --lint-fatal examples/prometheus.yml
```
