[![Google Analytics Proxy][banner-image-link]][github-repo-link]

[![License][license-badge]][license-link]
[![Go Reference][godoc-badge]][godoc-link]
[![Go Report Card][goreportcard-badge]][goreportcard-link]
[![Actions][github-actions-badge]][github-actions-link]
[![Packages][github-packages-badge]][github-packages-link]
[![Releases][github-release-badge]][github-release-link]

# Google Analytics Proxy

📈 Transparent HTTP proxy for tracking pageviews with Google Analytics

## Motivations

There are a number of situations where it is difficult (or impossible) to utilize the traditional Google Analytics tracking scripts.

- Server-side problems
  - What if the page is dynamically generated, or the contents isn't under your direct control?
  - What if the page return an HTTP redirect?
  - What if the page returns non-HTML content?
- Client-side problems:
  - What if the browser blocks Google Analytics?
  - What if the browser has JavaScript disabled?
  - What if the browser doesn't even use JavaScript (like `curl`)?

This application is an option to solve all of these problems.

## What this application is not

- A proxy for `https://www.google-analytics.com/collect`.
- A proxy for serving `analytics.js`.
- A JavaScript library.

## How it works

This application is an HTTP proxy server; It listens for a client HTTP request, forwards the request to an upstream HTTP server, waits for a response from the upstream server, and finally returns that response back to the original client.

From the client's perspective everything is functioning identically to if they were connecting to the upstream directly, but all of their traffic is actually being transparently proxied.

While this is happening, each request and response is used to construct a pageview event, that is then reported to Google Analytics.

Since there is no JavaScript whatsoever, it is not possible to disable Google Analytics reporting. 🚫

Additionally, the upstream HTTP service doesn't need to integrate with (or have any knowledge of) Google Analytics.

## Installing

A [release][github-release-link] version Docker image can be pulled by running:

```shell
docker pull ghcr.io/joshdk/google-analytics-proxy:v0.1.0
```

Or, a [development][github-master-link] version binary can be installed by running:

```shell
go install github.com/joshdk/google-analytics-proxy@master
```

## Configuration

This tool uses several environment variables as configuration.

| Name                              | Purpose                                                             | Example               |
| --------------------------------- | ------------------------------------------------------------------- | --------------------- |
| `$LISTEN`                         | Host and port that the proxy will listen on.                        | `0.0.0.0:8080`        |
| `$TLS_CERT_PATH`                  | (Optional) Path to TLS certificate file.                            | `/path/to/tls.pem`    |
| `$TLS_KEY_PATH`                   | (Optional) Path to TLS private key file.                            | `/path/to/tls.key`    |
| `$UPSTREAM_ENDPOINT`              | Address of the upstream service to be proxied.                      | `https://example.com` |
| `$UPSTREAM_HOSTNAME`              | (Optional) Hostname to used when proxying requests to the upstream. | `example.com`         |
| `$GOOGLE_ANALYTICS_TRACKING_ID`   | Tracking ID for your Google Analytics property                      | `UA-123456789-1`      |
| `$GOOGLE_ANALYTICS_PROPERTY_NAME` | Name of your Google Analytics property.                             | `example.com`         |
| `$GOOGLE_ANALYTICS_DRY_RUN`       | (Optional) Disables Google Analytics reporting.                     | `true`                |

## Usage

To run the Docker image, you can use a command like:

```shell
docker run \
  --rm \
  -p 8080:8080 \
  -e LISTEN=0.0.0.0:8080 \
  -e UPSTREAM_ENDPOINT=https://example.com \
  -e UPSTREAM_HOSTNAME=example.com \
  -e GOOGLE_ANALYTICS_PROPERTY_NAME=example.com \
  -e GOOGLE_ANALYTICS_TRACKING_ID=UA-123456789-1 \
    ghcr.io/joshdk/google-analytics-proxy:v0.1.0
```

Or, to run the local binary, you can use a command like:

```shell
LISTEN=0.0.0.0:8080 \
UPSTREAM_ENDPOINT=https://example.com \
UPSTREAM_HOSTNAME=example.com \
GOOGLE_ANALYTICS_PROPERTY_NAME=example.com \
GOOGLE_ANALYTICS_TRACKING_ID=UA-123456789-1 \
  $GOPATH/bin/google-analytics-proxy
```

### Kubernetes

This application is designed to be deployed into a Kubernetes cluster, ideally as a side-car container inside the same pod as your existing web service.
While not the only configuration option, this technique is beneficial because it enables you to:

- Proxy your service with minimal networking overhead.
- Scale along with your service.
- Be configured in the same deployment spec as your service.

As a demo, there is a [kubernetes][kubernetes-demo-link] directory, containing a [kustomize](https://kustomize.io/) manifest which can be deployed with `kubectl apply -k ./kubernetes`.
Afterwards, you must run `kubectl port-forward svc/demo 8080:8080` in order to expose the service locally.

---

In all cases, browsing to [https://localhost:8080](https://localhost:8080) afterwards will display the proxied upstream.
Realtime pageviews should also appear in your Google Analytics dashboard.

## License

This code is distributed under the [MIT License][license-link], see [LICENSE.txt][license-file] for more information.

[banner-image-link]:      https://user-images.githubusercontent.com/307183/131765571-5303a7f6-42c0-4764-ab5f-0b96ede2fda1.png
[github-actions-badge]:   https://github.com/joshdk/google-analytics-proxy/workflows/Build/badge.svg
[github-actions-link]:    https://github.com/joshdk/google-analytics-proxy/actions
[github-master-link]:     https://github.com/joshdk/google-analytics-proxy/tree/master
[github-packages-badge]:  https://img.shields.io/badge/ghcr.io-images-blue.svg
[github-packages-link]:   https://github.com/joshdk/google-analytics-proxy/pkgs/container/google-analytics-proxy
[github-release-badge]:   https://img.shields.io/github/release/joshdk/google-analytics-proxy/all.svg
[github-release-link]:    https://github.com/joshdk/google-analytics-proxy/releases
[github-repo-link]:       https://github.com/joshdk/google-analytics-proxy
[godoc-badge]:            https://pkg.go.dev/badge/github.com/joshdk/google-analytics-proxy.svg
[godoc-link]:             https://pkg.go.dev/github.com/joshdk/google-analytics-proxy
[goreportcard-badge]:     https://goreportcard.com/badge/github.com/joshdk/google-analytics-proxy
[goreportcard-link]:      https://goreportcard.com/report/github.com/joshdk/google-analytics-proxy
[kubernetes-demo-link]:   https://github.com/joshdk/google-analytics-proxy/tree/master/kubernetes
[license-badge]:          https://img.shields.io/badge/license-MIT-green.svg
[license-file]:           https://github.com/joshdk/google-analytics-proxy/blob/master/LICENSE.txt
[license-link]:           https://opensource.org/licenses/MIT
