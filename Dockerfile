# Dockerfile for controld-exporter

FROM scratch

# dockers_v2 lays the build context out as linux/<arch>/<binary>
ARG TARGETPLATFORM

# Copy ca-certificates for HTTPS requests to the Control D API
COPY --from=alpine:3.24.4@sha256:f59abe656b692d901e184be3445e084786602108599d2a51d229c0f5a48c71bf /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the pre-built binary from GoReleaser
COPY $TARGETPLATFORM/controld-exporter /controld-exporter

# extra_files in .goreleaser.yml is what puts these in the build context
COPY LICENSE NOTICE /

# Create a non-root user (using numeric ID for scratch image)
USER 65534:65534

# Declare the port; publishing it still requires docker run -p
EXPOSE 10034

# Set the entrypoint
ENTRYPOINT ["/controld-exporter"]
