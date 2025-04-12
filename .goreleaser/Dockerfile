# Build stage
FROM golang:1-alpine AS builder

# Set build arguments for goreleaser
ARG UID=60000
ARG CTRLD_API_KEY
ENV CTRLD_API_KEY=$CTRLD_API_KEY

WORKDIR /app
COPY controld-exporter /bin/

EXPOSE 10016
USER ${UID}
ENTRYPOINT [ "/bin/controld-exporter" ]
