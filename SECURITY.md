# Security Policy

## Supported versions

Only the latest release carries fixes, and no older tag gets a patch branch. Reproduce a finding against the latest release before reporting it.

## Reporting a vulnerability

Report privately through GitHub Security Advisories, never an issue or a pull request — open the repository's **Security** tab and choose **Report a vulnerability**.

One maintainer works on this in their own time, so no response time is promised. The advisory goes out after the fix ships and credits the reporter unless they ask otherwise.

## What this exporter holds and exposes

This exporter reads the Control D API with one account token and exposes the data as Prometheus metrics. The token's privileges determine what the exporter can read.

- **Credential** — one Control D API token, which makes the exporter exactly as sensitive as that token, so prefer `CTRLD_API_KEY` to `--controld.api-key`.
- **Metrics** — unauthenticated plain HTTP carrying billing status, device names, profile names, and organization identifiers as label values, so keep it on a controlled path.
- **Logs** — `--log.level debug` writes each raw API response to the log unredacted, so treat debug logs as sensitive as the account itself.

## Out of scope

A defect in the Control D service or its API belongs to **Control D** — report it there, not to this third-party exporter.
