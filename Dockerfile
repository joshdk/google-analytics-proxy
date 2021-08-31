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
    -ldflags "-X main.version=$VERSION" \
    main.go

# The final build stage copies in the final binary.
FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /bin/google-analytics-proxy /bin/google-analytics-proxy

ENTRYPOINT ["/bin/google-analytics-proxy"]
