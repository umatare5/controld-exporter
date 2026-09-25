# Changelog

Notable changes to the metric surface, one section per release, listing the pull requests that release carries.

## [v1.3.1]

- [#99](https://github.com/umatare5/controld-exporter/pull/99) – test: cover the client and the collectors against measured replies
- [#98](https://github.com/umatare5/controld-exporter/pull/98) – docs: write American English in the changelog and the build config
- [#97](https://github.com/umatare5/controld-exporter/pull/97) – Rebuild on every make build and keep worktrees on make clean
- [#96](https://github.com/umatare5/controld-exporter/pull/96) – Clarify command-line flags in README
- [#95](https://github.com/umatare5/controld-exporter/pull/95) – docs: align README structure with the sibling exporters
- [#94](https://github.com/umatare5/controld-exporter/pull/94) – Warn in every release note that a minor may move the metric surface

> [!IMPORTANT]
>
> **BEHAVIOR CHANGE**
>
> - `controld_billing_status` reports the transaction status the reply carries, where every payment published `0` before. A settled payment reads `1`.
> - The bundled alert rule and the Grafana panel move with it. `BillingStatusFailed` fires on `!= 1`, and the panel maps `1` to OK, so replace both if you deployed the versions v1.3.0 shipped.

## [v1.3.0]

- [#92](https://github.com/umatare5/controld-exporter/pull/92) – Rebuild the reference set around one architecture page
- [#91](https://github.com/umatare5/controld-exporter/pull/91) – Bump umatare5/common to v0.21.1 to fix the CodeQL workflow
- [#90](https://github.com/umatare5/controld-exporter/pull/90) – Remove the stats collector and the query series it published
- [#88](https://github.com/umatare5/controld-exporter/pull/88) – chore(deps): update all patch dependencies
- [#87](https://github.com/umatare5/controld-exporter/pull/87) – Survive what Control D actually returns, and alert when it stops
- [#86](https://github.com/umatare5/controld-exporter/pull/86) – Rewrite the contributor pages as claims that link their owner
- [#85](https://github.com/umatare5/controld-exporter/pull/85) – Split the operator pages by owner and correct them against the code
- [#84](https://github.com/umatare5/controld-exporter/pull/84) – Link the shared baseline and narrow the release archive
- [#83](https://github.com/umatare5/controld-exporter/pull/83) – Stop reading configuration from CONFIGOR\_\* variables
- [#82](https://github.com/umatare5/controld-exporter/pull/82) – Report failed scheduled runs and build a weekly release snapshot
- [#81](https://github.com/umatare5/controld-exporter/pull/81) – Run the build and tests weekly
- [#79](https://github.com/umatare5/controld-exporter/pull/79) – Extend the shared Renovate profile and pin the Alpine tag
- [#78](https://github.com/umatare5/controld-exporter/pull/78) – docs: add a reference set and validate the example rules in CI
- [#77](https://github.com/umatare5/controld-exporter/pull/77) – Update dependency golangci/golangci-lint to v2.13.2
- [#76](https://github.com/umatare5/controld-exporter/pull/76) – Update module github.com/sirupsen/logrus to v1.10.2

> [!IMPORTANT]
>
> **BREAKING CHANGE**
>
> - `controld_stats_last_queries_count` and its `type` label are removed, with the `stats` collector behind them.
> - `controld_profile_ip_filters_total` now reports the IP filter count instead of duplicating the content filter count.
> - The standalone major image tag is published again, reversing the note v1.2.0 carried.

## [v1.2.1]

- [#75](https://github.com/umatare5/controld-exporter/pull/75) – Release v1.2.1
- [#74](https://github.com/umatare5/controld-exporter/pull/74) – Report HTTP failures by status instead of a JSON parse error

## [v1.2.0]

- [#73](https://github.com/umatare5/controld-exporter/pull/73) – Release v1.2.0
- [#72](https://github.com/umatare5/controld-exporter/pull/72) – Migrate GoReleaser to dockers_v2 and ship license notices
- [#71](https://github.com/umatare5/controld-exporter/pull/71) – Refresh the Makefile with the shared development targets
- [#70](https://github.com/umatare5/controld-exporter/pull/70) – Install the shared pre-commit stack and strict docs lint

> [!IMPORTANT]
>
> **BREAKING CHANGE**
>
> - Per-arch image tags (`latest-amd64`, `v1.1.0-arm64`, and the other `-amd64`/`-arm64` suffixes) and the standalone `v1` tag are no longer published; the existing ones stay frozen at v1.1.0. Pull the multi-arch tags (`latest`, `vX.Y.Z`, `vX.Y`) instead.
> - `docker run` without arguments now starts the exporter instead of printing help, matching the README quick start.

## [v1.1.0]

This release takes dependency updates only. No metric, label, flag or HELP string changes.

[v1.3.1]: https://github.com/umatare5/controld-exporter/releases/tag/v1.3.1
[v1.3.0]: https://github.com/umatare5/controld-exporter/releases/tag/v1.3.0
[v1.2.1]: https://github.com/umatare5/controld-exporter/releases/tag/v1.2.1
[v1.2.0]: https://github.com/umatare5/controld-exporter/releases/tag/v1.2.0
[v1.1.0]: https://github.com/umatare5/controld-exporter/releases/tag/v1.1.0
