# Security Policy

The [shared security policy](https://github.com/umatare5/.github/blob/main/SECURITY.md) covers what every exporter here shares. This page carries the rest.

## What to Include

Redact these before reporting, in addition to the credentials the shared policy names.

- The Control D API token, from a log line, a process listing or a container definition
- An organization or sub-organization primary key, which the `orgId` label carries
- A payment or subscription primary key, which the `id` label carries on the billing families
- A device, profile, organization or sub-organization name, which `name` carries verbatim

Reproduction needs the flags in force, whether `--controld.business-mode` was set, and the endpoint being read, because the same token reads two different account scopes.

## Exposure

This exporter holds one Control D API token, reads the whole account under it and publishes that account's own names, so `/metrics` is largely a copy of the configuration. The exception is `controld_stats_last_queries_count`, which counts the last minute's DNS queries by verdict.

- **Token** — `CTRLD_API_KEY` passes it out of the environment.
- **Command line** — `--controld.api-key` puts it there instead.
- **Process table** — every account on the host reads that line out of `ps`.
- **Labels** — `orgId` carries an organization or sub-organization primary key.
- **Names** — `name` carries an object's name verbatim, or a service category's primary key.
- **Endpoint** — both ride an unauthenticated `/metrics`.
- **Debug** — `--log.level debug` writes the request URI and every response body as received.
- **Log handling** — the log then holds the whole account, and is handled as the token is.

> [!IMPORTANT]
> The token travels in an `Authorization` header alone, and no path writes it to the landing page, the `/metrics` body or a log line, so a token reaching any of them is a vulnerability.
>
> `--controld.api-key` is the one documented way the token reaches the process table, so any path that puts it there without that flag is a vulnerability.
>
> `--controld.business-mode` fixes the scope every collector reads, and a sub-organization is reached only under an `X-Force-Org-Id` header carrying a key the sub-organizations response listed. A request reaching an organization the configured mode did not name is a vulnerability.

## Egress Paths

Nothing leaves the host but the calls one scrape makes, and every one of them carries the token, so the exporter reaches Control D and nothing else.

### API

- **Host** — every call but the query report goes to `https://api.controld.com`.
- **Sub-organizations** — business mode repeats the device, profile and category calls for each.
- **Scale** — a scrape's request count therefore grows with the account.

### Analytics

- **Host** — the query report goes to a region label in front of `analytics.controld.com`.
- **Business mode** — the organization response supplies that label, and every sub-org reuses it.

### Transport

- **Certificates** — the requests use Go's default client, which carries no TLS settings of its own.
- **Verification** — it is therefore on, and no flag relaxes it.

## Out of Scope

- Account data in a `/metrics` label, because that is what the exporter exists to publish.
- A Control D service or API defect, because it belongs to Control D and not to this exporter.
