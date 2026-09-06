# Security Policy

The [shared policy](https://github.com/umatare5/.github/blob/main/SECURITY.md) covers what every exporter here shares. This page carries the rest.

## What to Include

Redact these before reporting, in addition to the credentials the shared policy names.

- The Control D API token, from a log line, a process listing or a container definition
- An organization or sub-organization primary key, which the `orgId` label carries
- A device or profile name, which the `name` label carries verbatim from the account

Reproduction needs the flags in force, whether `--controld.business-mode` was set, and the endpoint being read, because the same token reads two different account scopes.

## Exposure

This exporter holds one Control D API token, reads the whole account under it and publishes that account's own names, so `/metrics` is a copy of the configuration rather than a measurement of traffic.

- **Token** — `CTRLD_API_KEY` passes it out of the environment, while `--controld.api-key` puts it on the command line, where every account on the host reads it out of `ps`.
- **Labels** — `orgId` carries an organization or sub-organization primary key and `name` carries a device or profile name verbatim, on an unauthenticated endpoint.
- **Debug** — `--log.level debug` writes the request URI and the decoded body of every response, so the log holds the whole account and is handled as the token is.

> [!IMPORTANT]
> The token travels in an `Authorization` header alone, and no path writes it to the landing page, the `/metrics` body or a log line, so a token reaching any of them is a vulnerability.
>
> `--controld.business-mode` fixes the scope every collector reads, and a sub-organization is reached only by repeating a request under an `X-Force-Org-Id` header the organization response named. A request reaching an organization the configured mode did not name is a vulnerability.

## Egress

Nothing leaves the host but the calls one scrape makes, and every one of them carries the token, so the exporter reaches Control D and nothing else.

- **API** — each collector reads `https://api.controld.com` in one or two requests, and business mode adds one per sub-organization, so a scrape's request count grows with the account.
- **Analytics** — the query report goes to `analytics.controld.com` under `america` in personal mode, and under the label the organization response supplies in business mode.
- **Certificates** — the requests go through Go's default client, which carries no TLS settings of its own, so verification is on and no flag relaxes it.

## Out of Scope

- Account data in a `/metrics` label, which is what the exporter exists to publish.
- A Control D service or API defect, which belongs to Control D and not to this exporter.
