# Security Policy

Please follow **[the shared security policy](https://github.com/umatare5/.github/blob/main/SECURITY.md)**, which covers:

- **Supported Versions** – only the latest release carries fixes, so reproduce against it.
- **Reporting a Vulnerability** – the private advisory path and what the response promises.
- **What to Include** – the credentials and addresses to redact, and the fields to send.
- **Exposure** – the unauthenticated surface and the container the image ships.
- **Out of Scope** – findings that belong to the monitored system or to an operator's own configuration.

This page specifies what is particular to this one.

## What to Include

Reproduction needs **the collector flags** in force, **whether `--controld.business-mode` was set**.

Redact these before reporting, in addition to the credentials the shared policy names.

- **Tokens** – The Control D API key, from a flag, an environment variable, a process listing or a log line
- **Identifiers** – An organization or sub-organization primary key, which the `orgId` label carries
- **Billing** – A payment or subscription primary key, which the `id` label carries on the billing families
- **Names** – A device, profile, organization or sub-organization name, which `name` carries verbatim

## Exposure

This exporter holds one credential, **the Control D API key**.

- **Variables** – `CTRLD_API_KEY` is the safer source, as `/proc/<pid>/environ` opens to the owner and root alone.
- **Flags** – `--controld.api-key` puts the token on the command line, reaching a process table every account on the host reads.
- **Labels** – The `orgId` label exposes organization primary keys; `name` carries object names or service category keys verbatim.
- **Debug** – `--log.level debug` writes the request URI and every response body as received.
- **Payload** – Those bodies carry client hostnames, MAC addresses, contact emails and SSO credentials.
- **Handling** – The log then holds the whole account, so it is handled as the token is.

## Endpoints

No route authenticates, so the network path the listener sits on is the whole access control. See also [Endpoints](docs/architecture.md#endpoints).

- **Listener** – The server binds a unified listener for all internal routes.
- **Disclosure** – Both the metrics and the names ride an unauthenticated `/metrics`.
- **Landing page** – The `/` route returns HTML with the telemetry path.
- **Isolation** – No operator choice reaches the landing page; it names the listen address and the scrape path alone.

## Ingress Paths

The exporter exposes a listening socket for incoming HTTP scrapes.

- **Reach** – `--web.listen-address` configures the bound interface, controlling which networks can access the exposed telemetry.
- **No allowlist** – The exporter filters no sender, which leaves the packet filter or authenticating proxy to enforce it.
- **Restriction** – Put a packet filter or an authenticating proxy in front to restrict access.

## Egress Paths

The exporter opens outbound connections only to the Control D API.

- **Host** – Every API call strictly egresses to `https://api.controld.com`.
- **Credential** – The client strictly sends the token in the `Authorization` header, keeping it out of URLs and proxy logs.
- **Cost** – The API cost scales inherently with the account tier and the `--controld.business-mode` flag.
- **Scaling** – Business mode repeats the device, profile and category calls for each sub-organization.

## Out of Scope

- **Origin** – A defect in the Control D service or API belongs to Control D rather than to this client.
- **Credit** – Credit or billing exhaustion, which the operator's requirements dictate.
