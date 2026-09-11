# Help

The `controld-exporter --help` text, transcribed from the binary.

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

`--log.level debug` writes the request URI and the response body of every Control D call to the log. The body is written as received and before the status is checked, so a failed response is logged too, which is why [`SECURITY.md`](../SECURITY.md) treats a debug log as sensitive as the API key itself. The key is not among the headers logged.

`--web.telemetry-path` shares one `http.ServeMux` with the landing page registered at `/`, so setting it to `/` is a duplicate registration and panics at startup rather than replacing the landing page.

The transcript above reads `dev` for `VERSION:` because it comes from a build that stamps nothing. `make build` stamps the contents of [`VERSION`](../VERSION) and a release stamps its tag, so neither answers `dev`.
