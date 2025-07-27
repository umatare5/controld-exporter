# Dockerfile for controld-exporter
# Multi-stage build for optimized container size

# Build stage
FROM golang:1.24-alpine AS builder

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates git

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o controld-exporter ./cmd/main.go

# Final stage - minimal container
FROM scratch

# Copy ca-certificates for HTTPS requests to controld-exporter controllers
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary from builder stage
COPY --from=builder /app/controld-exporter /controld-exporter

# Create a non-root user (using numeric ID for scratch image)
USER 65534:65534

# Set the entrypoint
ENTRYPOINT ["/controld-exporter"]

# Default command shows help
CMD ["--help"]

# Metadata
LABEL org.opencontainers.image.title="controld-exporter"
LABEL org.opencontainers.image.description="Prometheus Exporter for Control D"
LABEL org.opencontainers.image.vendor="umatare5"
LABEL org.opencontainers.image.source="https://github.com/umatare5/controld-exporter"
LABEL org.opencontainers.image.documentation="https://github.com/umatare5/controld-exporter/blob/main/README.md"
