# Build stage
FROM golang:1-alpine AS builder

ARG UID=60000
ARG CTRLD_API_KEY
ENV CTRLD_API_KEY=$CTRLD_API_KEY

# Set the working directory
WORKDIR /tmp/build

# Copy go.mod and go.sum first to leverage Docker layer caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o controld-exporter

# Final stage
FROM alpine:3.21.3

ARG UID=60000
WORKDIR /app
COPY --from=builder /tmp/build/controld-exporter /bin/

EXPOSE 10016
USER ${UID}
ENTRYPOINT [ "/bin/controld-exporter" ]
