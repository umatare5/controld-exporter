# Dockerfile for controld-exporter

FROM scratch

# Copy ca-certificates for HTTPS requests to controld-exporter controllers
COPY --from=alpine:latest@sha256:865b95f46d98cf867a156fe4a135ad3fe50d2056aa3f25ed31662dff6da4eb62 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the pre-built binary from GoReleaser
COPY controld-exporter /controld-exporter

# Create a non-root user (using numeric ID for scratch image)
USER 65534:65534

# Set the entrypoint
ENTRYPOINT ["/controld-exporter"]

# Default command shows help
CMD ["--help"]
