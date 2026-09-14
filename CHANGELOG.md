# Changelog

Notable changes to the metric surface, one section per release — a short preamble, the breaking change where the release has one, then the metric changes and the flag changes.

This changelog starts at v1.1.0; earlier releases are described by their [release notes](https://github.com/umatare5/controld-exporter/releases) alone.

## [Unreleased]

This release keeps a failed or empty Control D response from costing the scrape, and gives the shipped alert rule a family a revoked key can silence. It also links the shared policy pages, narrows the release archives to the two files a redistributed binary needs, and rebuilds the reference pages. No metric, label, flag or HELP string changes.

`SECURITY.md` and `CONTRIBUTING.md` now open with the baseline every exporter under `umatare5` shares and carry only what is specific to this one, so a convention stated once is no longer restated per repository.

Release archives carry `LICENSE` and `NOTICE` alone. The exporter parses none of the files they held — `examples/prometheus*.yml` are a Prometheus server configuration and its rule files — and each is a click away on the page the archive was downloaded from.

`README.md` gains a `## Collectors` section, folds `## Environment Variables` into `## Flags` and delegates every mechanism to the page that owns it. The reference pages under `docs/` now carry the request count, the endpoint behaviour and the per-family facts the README only summarised.

A failed organization fetch no longer costs the whole scrape. The organization collector built its metrics before it read the fetch error, so a non-2xx answer from `/organizations/organization` dereferenced a nil response and panicked. The recovered panic answered `500` with no family and no log line — the collector now reads the error first and leaves the six behind it to publish.

The empty-payload guards never fired, because each tested a concrete response against a type switch matching only `[]any` and `map[string]any`. A `200` carrying an empty query list therefore reached an unguarded index and panicked the same way. Each guard now reads its own slice.

`ControlDMetricsMissing` read `absent(controld_network_health_code)`, which a revoked key cannot trigger because `/network` answers without a token. It now reads that family beside `controld_profile_rules_total`, so a revoked key and a failed `/network` call each fire it.

The contributor pages now carry a claim and a link where they carried a mechanism. `AGENTS.md` keeps its seven sections and rewrites Domain Knowledge around what Control D does, `CONTRIBUTING.md` states which CI jobs a path filter gates, and `SECURITY.md` names the calls each mode makes.

More statements were corrected against the source. A field the API stops sending publishes as `0` rather than being withheld, a token without access to an endpoint answers `403` with `40301`, the listen defaults live in `internal/cli`, and `statsEndpoint` names a region label in front of `analytics.controld.com` rather than a host.

### Metrics

None.

### Flags

None.

## [v1.2.1]

This release reports a failed Control D API call by its HTTP status. No metric, label, flag or HELP string changes.

A request that answers with a non-2xx status is now an error before the body is decoded. `controld_stats_last_queries_count` has been unavailable since Control D removed the analytics endpoint it reads, and that failure logged as `invalid character 'p' after top-level value` because the plain-text `404 page not found` body was decoded as JSON; it now names the status and the endpoint. The metric stays unavailable — only the diagnosis changes.

### Metrics

None.

### Flags

None.

## [v1.2.0]

This release rebuilds the distribution on Go 1.27 and moves container publishing to GoReleaser `dockers_v2`. No metric, label, flag or HELP string changes.

> [!IMPORTANT]
>
> ### BREAKING CHANGE
>
> - Per-arch image tags (`latest-amd64`, `v1.1.0-arm64`, and the other `-amd64`/`-arm64` suffixes) and the standalone `v1` tag are no longer published; the existing ones stay frozen at v1.1.0. Pull the multi-arch tags (`latest`, `vX.Y.Z`, `vX.Y`) instead.
> - `docker run` without arguments now starts the exporter instead of printing help, matching the README quick start.

The binaries build with Go 1.27 and pinned `CGO_ENABLED=0` on every platform. The image declares port `10034/tcp` and carries the third-party license notices, and this release's archives added `CHANGELOG.md`, `SECURITY.md` and `NOTICE` beside the binary.

### Metrics

None.

### Flags

None.

## [v1.1.0]

This release takes dependency updates only. No metric, label, flag or HELP string changes.

[Unreleased]: https://github.com/umatare5/controld-exporter/compare/v1.2.1...main
[v1.2.1]: https://github.com/umatare5/controld-exporter/releases/tag/v1.2.1
[v1.2.0]: https://github.com/umatare5/controld-exporter/releases/tag/v1.2.0
[v1.1.0]: https://github.com/umatare5/controld-exporter/releases/tag/v1.1.0
