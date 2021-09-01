![Google Analytics Proxy](https://user-images.githubusercontent.com/307183/131599246-f0516c3b-2f91-43ab-9db9-52e5180c59ad.png)

[![Actions][github-actions-badge]][github-actions-link]
[![License][license-badge]][license-link]
[![Releases][github-release-badge]][github-release-link]

# Google Analytics Proxy

📊 Transparent proxy for tracking page views with Google Analytics

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

Browsing to [https://localhost:8080](https://localhost:8080) afterwards will display the proxied upstream.

## License

This code is distributed under the [MIT License][license-link], see [LICENSE.txt][license-file] for more information.

[github-actions-badge]:  https://github.com/joshdk/google-analytics-proxy/workflows/Build/badge.svg
[github-actions-link]:   https://github.com/joshdk/google-analytics-proxy/actions
[github-master-link]:    https://github.com/joshdk/google-analytics-proxy/tree/master
[github-release-badge]:  https://img.shields.io/github/release/joshdk/google-analytics-proxy/all.svg
[github-release-link]:   https://github.com/joshdk/google-analytics-proxy/releases
[license-badge]:         https://img.shields.io/badge/license-MIT-green.svg
[license-file]:          https://github.com/joshdk/google-analytics-proxy/blob/master/LICENSE.txt
[license-link]:          https://opensource.org/licenses/MIT
