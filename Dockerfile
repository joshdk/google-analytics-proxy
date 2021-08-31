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

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=upx /bin/google-analytics-proxy /bin/google-analytics-proxy

ENTRYPOINT ["/bin/google-analytics-proxy"]
