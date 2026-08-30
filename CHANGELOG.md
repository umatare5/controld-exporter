# Changelog

Notable changes to the metric surface, one section per release — a short preamble, the breaking change where the release has one, then the metric changes and the flag changes.

This changelog starts at v1.1.0; earlier releases are described by their [release notes](https://github.com/umatare5/controld-exporter/releases) alone.

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

The binaries build with Go 1.27 and pinned `CGO_ENABLED=0` on every platform. The image declares port `10034/tcp` and carries the third-party licence notices, and release archives add `CHANGELOG.md`, `SECURITY.md`, and `NOTICE`.

### Metrics

None.

### Flags

None.

## [v1.1.0]

This release takes dependency updates only. No metric, label, flag or HELP string changes.

[v1.2.1]: https://github.com/umatare5/controld-exporter/releases/tag/v1.2.1
[v1.2.0]: https://github.com/umatare5/controld-exporter/releases/tag/v1.2.0
[v1.1.0]: https://github.com/umatare5/controld-exporter/releases/tag/v1.1.0
