# Contributing

Thank you for considering a contribution.

## Commands

The following `make` commands are available for development and testing:

| Command                     | Description                                   |
| :-------------------------- | :-------------------------------------------- |
| `make help`                 | Display available targets and requirements    |
| `make build`                | Build the binary to `./tmp/controld-exporter` |
| `make lint`                 | Run golangci-lint and tidy go.mod             |
| `make test-unit`            | Run unit tests with coverage using gotestsum  |
| `make test-unit-coverage`   | Generate HTML coverage report                 |
| `make clean`                | Remove build artifacts and backup files       |
| `make image`                | Build Docker image                            |
| `make pre-commit-install`   | Install the pre-commit hooks                  |
| `make pre-commit-test`      | Run every hook across the tree                |
| `make pre-commit-uninstall` | Remove the pre-commit hooks                   |

Markdown style is enforced by the `markdownlint-cli2` hook that `make pre-commit-install` wires in, and again in CI. Links are checked in CI only, because that run reaches third-party hosts. Run `lychee .` to reproduce a link failure locally.

## Build

The repository includes a ready to use `Dockerfile`. To build a new Docker image:

```bash
make image
```

This cross-compiles a Linux binary into `./tmp/image`, then builds from that directory because the `Dockerfile` expects the binary at the context root. The image is tagged `$USER/controld-exporter`. Released images are pushed to `ghcr.io/umatare5/controld-exporter` by GoReleaser instead.

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
