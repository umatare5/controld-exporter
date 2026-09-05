# Help

The `controld-exporter --help` text, transcribed from the binary.

## controld-exporter

```text
NAME:
   controld-exporter - A Prometheus exporter for metrics from the Control D

USAGE:
   controld-exporter [options...]

VERSION:
   dev

GLOBAL OPTIONS:
   --web.listen-address string             Address to bind the HTTP server to. (default: "0.0.0.0")
   --web.listen-port int                   Port number to bind the HTTP server to. (default: 10034)
   --web.telemetry-path string, -p string  Path for the metrics endpoint. (default: "/metrics")
   --controld.api-key string, -k string    API key for authenticating with the Control D API. [$CTRLD_API_KEY]
   --controld.business-mode                Enable the metrics collection available in the business subscription.
   --log.level string                      Set the logging level. One of: [debug, info, warn, error] (default: "info")
   --help, -h                              show help
   --version, -v                           print the version
```

## Notes

`--controld.api-key` is the only required flag, and `CTRLD_API_KEY` fills it where the flag is absent. Startup stops before the listener opens when both are empty, so a misconfigured exporter fails loudly rather than serving an empty `/metrics`.

`--controld.business-mode` decides which Control D account scope every collector reads — [Collectors](collectors.md#specifications) carries what each mode publishes and what fills `orgId` without it.

`--log.level debug` writes the full request URI and the decoded body of every Control D response to the log, which is why [`SECURITY.md`](../SECURITY.md) treats a debug log as sensitive as the API key itself.

`--web.telemetry-path` shares one `http.ServeMux` with the landing page registered at `/`, so setting it to `/` is a duplicate registration and panics at startup rather than replacing the landing page.

The binary reports `dev` for `--version` unless the version is stamped at link time, so a locally built exporter and a release archive of the same commit answer differently.
