# The builder build stage compiles the Go code into a static binary.
FROM golang:1.16-alpine as builder

WORKDIR /go/src/github.com/joshdk/google-analytics-proxy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -o /bin/google-analytics-proxy \
    main.go

# The final build stage copies in the final binary.
FROM scratch

COPY --from=builder /bin/google-analytics-proxy /bin/google-analytics-proxy

ENTRYPOINT ["/bin/google-analytics-proxy"]
