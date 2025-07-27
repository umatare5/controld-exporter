# Dockerfile for controld-exporter
# Optimized for GoReleaser builds

FROM scratch

# Copy ca-certificates for HTTPS requests to controld-exporter controllers
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the pre-built binary from GoReleaser
COPY controld-exporter /controld-exporter

# Create a non-root user (using numeric ID for scratch image)
USER 65534:65534

# Set the entrypoint
ENTRYPOINT ["/controld-exporter"]

# Default command shows help
CMD ["--help"]

# Metadata labels will be set by GoReleaser
