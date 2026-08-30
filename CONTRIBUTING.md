# Contributing

Thank you for considering a contribution.

## Commands

The following commands are available for development and testing:

| Command                                     | Description                            |
| :------------------------------------------ | :------------------------------------- |
| `go build ./...`                            | Compile all packages                   |
| `golangci-lint run`                         | Run the configured linters             |
| `golangci-lint fmt`                         | Apply the configured formatters        |
| `pre-commit install --allow-missing-config` | Install the pre-commit hooks           |
| `pre-commit run --all-files`                | Run every hook across the whole tree   |

Markdown style is enforced by the `markdownlint-cli2` hook that `pre-commit install` wires in, and again in CI. Links are checked in CI only, because that run reaches third-party hosts. Run `lychee .` to reproduce a link failure locally.

## Build

Released container images are built and pushed to `ghcr.io/umatare5/controld-exporter` by GoReleaser, so a local Docker build is not part of the day-to-day workflow. To build the binary locally:

```bash
go build -o tmp/controld-exporter ./cmd
```

## Release

To release a new version, follow these steps:

1. Add the `## [vX.Y.Z]` section to `CHANGELOG.md` above the previous release, and add that version's release link at the foot of the file.
2. Update the version in the `VERSION` file.
3. Submit a pull request with both files.

A push to `main` touching `VERSION` runs the [release workflow](https://github.com/umatare5/controld-exporter/actions/workflows/go-release.yml), which tags the commit and publishes the release in the same run. The workflow has no manual trigger, so there is no step to perform by hand.

## Pull requests

1. [Fork](https://github.com/umatare5/controld-exporter/fork) the repository
2. Create a feature branch
3. Commit your changes
4. Record any change to the metric surface under a `## [vX.Y.Z]` section for the coming version in `CHANGELOG.md`, adding the section if it is not there yet
5. Rebase your local changes against the `main` branch
6. Create a new Pull Request
