FROM golang:1.17.5-bullseye as builder
COPY src /var/src
RUN cd /var/src && \
    go build

FROM debian:bullseye-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates \
    && update-ca-certificates

COPY --from=builder /var/src/cclashx /usr/local/bin/cclashx