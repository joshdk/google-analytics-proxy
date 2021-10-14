# The certs stage is used to obtain a current set of CA certificates.
FROM alpine:3.14 as certs

# hadolint ignore=DL3018
RUN apk add --no-cache \
    ca-certificates

# The builder build stage compiles the Go code into a static binary.
FROM golang:1.16-alpine as builder

ARG VERSION=development

WORKDIR /go/src/github.com/joshdk/google-analytics-proxy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -o /bin/google-analytics-proxy \
    -ldflags "-s -w -X main.version=$VERSION" \
    -trimpath \
    main.go

# The upx build stage uses upx to compress the binary.
FROM alpine:3.14 as upx

RUN wget https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz \
 && tar -xf upx-3.96-amd64_linux.tar.xz \
 && install upx-3.96-amd64_linux/upx /bin/upx \
 && rm -rf upx*

COPY --from=builder /bin/google-analytics-proxy /bin/google-analytics-proxy

RUN upx --best --ultra-brute /bin/google-analytics-proxy

# The final build stage copies in the final binary.
FROM scratch

ARG CREATED
ARG REVISION
ARG VERSION

MAINTAINER Josh Komoroske <github.com/joshdk>

# Standard OCI image labels.
# See: https://github.com/opencontainers/image-spec/blob/v1.0.1/annotations.md#pre-defined-annotation-keys
LABEL org.opencontainers.image.created="$CREATED"
LABEL org.opencontainers.image.authors="Josh Komoroske <github.com/joshdk>"
LABEL org.opencontainers.image.url="https://github.com/joshdk/google-analytics-proxy"
LABEL org.opencontainers.image.documentation="https://github.com/joshdk/google-analytics-proxy/blob/master/README.md"
LABEL org.opencontainers.image.source="https://github.com/joshdk/google-analytics-proxy"
LABEL org.opencontainers.image.version="$VERSION"
LABEL org.opencontainers.image.revision="$REVISION"
LABEL org.opencontainers.image.vendor="Josh Komoroske <github.com/joshdk>"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.ref.name="ghcr.io/joshdk/google-analytics-proxy:$VERSION"
LABEL org.opencontainers.image.title="Google Analytics Proxy"
LABEL org.opencontainers.image.description="Transparent HTTP proxy for tracking pageviews with Google Analytics"

COPY LICENSE.txt /
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY README.md /
COPY --from=upx /bin/google-analytics-proxy /bin/google-analytics-proxy

# Switch to a non-root user. Arbitrarily, use the same uid/gid as the "nobody"
# user from Alpine.
USER 65534:65534

EXPOSE 8080/tcp
EXPOSE 8443/tcp

ENTRYPOINT ["/bin/google-analytics-proxy"]
