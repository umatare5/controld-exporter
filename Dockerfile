# Dockerfile for controld-exporter

FROM scratch

# Copy ca-certificates for HTTPS requests to controld-exporter controllers
COPY --from=alpine:latest@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the pre-built binary from GoReleaser
COPY controld-exporter /controld-exporter

# Create a non-root user (using numeric ID for scratch image)
USER 65534:65534

# Set the entrypoint
ENTRYPOINT ["/controld-exporter"]

# Default command shows help
CMD ["--help"]
