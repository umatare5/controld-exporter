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
- **Do not lean on the coverage gate.** No test exists yet, so CI passes at 0 percent coverage until the first one lands.
- **Stamp the version with `make build`.** The `--help` transcript reads `dev` until the build stamps it.

If the API key belongs to an organization:

- **Set `--controld.business-mode`.** It changes what collectors read rather than turning them on or off.

## Testing

Testing here rests on the example rules rather than on Go tests.

- **Check rules with `promtool`.** No Go test covers the alerting expressions, so CI lints and asserts them there.
- **Redirect `baseURL` to an `httptest` server.** It is unexported, so only a test inside `internal/controld` can rewrite it.
- **Keep a device or organization reply out of the tree.** It carries client hostnames, MAC addresses and SSO credentials.

The `Prometheus Rules` job runs these three commands.

```bash
promtool check rules --lint all --lint-fatal examples/prometheus_alert_rules.yml
promtool test rules examples/prometheus_alert_rules_test.yml
promtool check config --lint all --lint-fatal examples/prometheus.yml
```
