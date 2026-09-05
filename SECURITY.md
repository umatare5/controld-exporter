# Security Policy

## Supported Versions

Only the most recent tagged release carries fixes — reproduce a finding against it before reporting.

## Reporting a Vulnerability

Report privately through [GitHub Security Advisories](https://github.com/umatare5/controld-exporter/security/advisories/new). **Please do not report a vulnerability through a public GitHub issue or a pull request.**

The response is best effort, with no promised window. The advisory goes out once the fix ships, carries a CVE request, and credits the reporter unless they ask otherwise.

## What to Include

**Redact these first.** None of them belongs in a report.

- The Control D API token, from a log line, a process listing or a container definition
- An organization or sub-organization primary key, which the `orgId` label carries
- A device or profile name, which the `name` label carries verbatim from the account

Then include the following:

- **Affected versions** (required): The `controld-exporter` release, and whether `--controld.business-mode` was set
- **Reproduction steps** (required): The flags and environment variables, and the endpoint the exporter was reading
- **Output** (required): The `/metrics` body or the log lines, with every value above removed
- **Impact assessment** (required): The exploit scenario, and what it reaches
- **Suggested fix** (optional): Proposed remediation, if any
- **Disclosure status** (required): Whether it is shared elsewhere, and your plan for sharing it

## Scope

In scope:

- The API token reaching a log line at any level, or the landing page, or the `/metrics` body
- The API token reaching the process table other than through `--controld.api-key`, whose cost is documented
- Certificate verification weakened on the path to the Control D API, which no flag is meant to relax
- A request reaching an organization or sub-organization the configured mode did not name
- The published container image

Out of scope:

- `/metrics` served unauthenticated over plain HTTP, which is the exporter's documented posture — put it behind a controlled path
- Account data appearing in `/metrics` labels, which is what the exporter exists to publish
- A dependency advisory with no path reachable from `./cmd` — show the path, or a `govulncheck` finding
- A Control D service or API defect, which belongs to Control D rather than to this third-party exporter
- An operator's own configuration, which [`docs/help.md`](docs/help.md) covers
