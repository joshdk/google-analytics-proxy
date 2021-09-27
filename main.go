// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.

package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"

	"github.com/joshdk/google-analytics-proxy/analytics"
)

// version is used to hold the version string. Is replaced at go build time
// with -ldflags.
var version = "development"

func main() {
	if err := mainCmd(); err != nil {
		fmt.Println("joshdk/google-analytics-proxy:", err) //nolint:forbidigo
		os.Exit(1)
	}
}

func mainCmd() error {
	log.Printf("joshdk/google-analytics-proxy version %s", version)

	// listenAddress is the host and port that the proxy will listen on.
	// See net.Dial for details of the address format.
	// Example: "localhost:8080" "0.0.0.0:8080" ":8080"
	listenAddress := os.Getenv("LISTEN")

	// tlsCertFile is optionally the path to a TLS certificate file, used for
	// listening and serving HTTPS connections. Must always be configured with
	// tlsKeyFile.
	// Example: "/path/to/tls.pem"
	tlsCertFile := os.Getenv("TLS_CERT_PATH")

	// tlsKeyFile is optionally the path to a TLS private key file, used for
	// listening and serving HTTPS connections. Must always be configured with
	// tlsCertFile.
	// Example: "/path/to/tls.key"
	tlsKeyFile := os.Getenv("TLS_KEY_PATH")

	// upstreamEndpoint is the address of the upstream service to be
	// proxied.
	// Example: "https://example.com" "http://:80"
	upstreamEndpoint := os.Getenv("UPSTREAM_ENDPOINT")

	// upstreamHostname optionally is the hostname to used when proxying
	// requests to the upstream. Used for hostname based routing. If empty, the
	// value of $GOOGLE_ANALYTICS_PROPERTY_NAME will be used.
	// Example: "example.com"
	upstreamHostname := os.Getenv("UPSTREAM_HOSTNAME")

	// googleAnalyticsTrackingID is the tracking id for the Google
	// Analytics property that you want to track pageview events for. This
	// can be found in your Google Analytics dashboard.
	// Example: "UA-123456789-1"
	googleAnalyticsTrackingID := os.Getenv("GOOGLE_ANALYTICS_TRACKING_ID")

	// googleAnalyticsPropertyName is the name for the Google Analytics
	// property that you want to track pageview events for. This can be
	// found in your Google Analytics dashboard. Will be used as the upstream
	// hostname in proxied requests if $UPSTREAM_HOSTNAME is empty.
	// Example: "example.com"
	googleAnalyticsPropertyName := os.Getenv("GOOGLE_ANALYTICS_PROPERTY_NAME")

	// googleAnalyticsDryRun can optionally be used to disable reporting
	// pageview events with Google Analytics. See strconv.ParseBool() for
	// acceptable values.
	// Example: "true"
	googleAnalyticsDryRun := os.Getenv("GOOGLE_ANALYTICS_DRY_RUN")

	// Validate that the required settings are not empty.
	switch {
	case googleAnalyticsTrackingID == "":
		return fmt.Errorf("GOOGLE_ANALYTICS_TRACKING_ID was not provided")
	case googleAnalyticsPropertyName == "":
		return fmt.Errorf("GOOGLE_ANALYTICS_PROPERTY_NAME was not provided")
	case upstreamEndpoint == "":
		return fmt.Errorf("UPSTREAM_ENDPOINT was not provided")
	}

	// Validate the TLS settings, and set sane defaults.
	switch {
	// Validate HTTP listen mode.
	case tlsCertFile == "" && tlsKeyFile == "":
		if listenAddress == "" {
			// Set a default listen address if none was given.
			listenAddress = "0.0.0.0:8080"
		}
	// Validate HTTPS listen mode.
	case tlsCertFile != "" && tlsKeyFile != "":
		if listenAddress == "" {
			// Set a default listen address if none was given.
			listenAddress = "0.0.0.0:8443"
		}
	default:
		// HTTPS listen mode was only partially (mis)configured.
		return fmt.Errorf("TLS_CERT_PATH and TLS_KEY_PATH were not both provided")
	}

	// Parse the upstream endpoint address to ensure that it's valid.
	upstreamURL, err := url.Parse(upstreamEndpoint)
	if err != nil {
		return err
	}

	// Use the property name for the upstream hostname, if one was not
	// explicitly given.
	if upstreamHostname == "" {
		upstreamHostname = googleAnalyticsPropertyName
	}

	// Create a reverse proxy HTTP handler for our upstream. This handler is
	// responsible for relaying all downstream client requests to the upstream
	// service, and the upstream service responses back to the downstream
	// client.
	log.Printf("proxying traffic to %s (%s)", upstreamEndpoint, upstreamHostname)
	proxy := httputil.NewSingleHostReverseProxy(upstreamURL)

	// Modify the original proxy director function, only updating the request
	// hostname so that any hostname base routing that is performed by the
	// upstream service continues to work correctly.
	original := proxy.Director
	proxy.Director = func(request *http.Request) {
		request.Host = upstreamHostname
		original(request)
	}

	// Parse the Google Analytics dry run value. Intentionally ignore all
	// errors and default to false.
	googleAnalyticsDryRunBool, _ := strconv.ParseBool(googleAnalyticsDryRun)

	// Create a tracker for sending pageviews to Google Analytics.
	if !googleAnalyticsDryRunBool {
		log.Printf("tracking analytics for %s (%s)", googleAnalyticsTrackingID, googleAnalyticsPropertyName)
	} else {
		log.Printf("skipping analytics for %s (%s)", googleAnalyticsTrackingID, googleAnalyticsPropertyName)
	}
	tracker := &analytics.Tracker{
		TrackingID:   googleAnalyticsTrackingID,
		PropertyName: googleAnalyticsPropertyName,
		DryRun:       googleAnalyticsDryRunBool,
		Handler:      proxy,
	}

	// Start the server and listen for incoming requests!
	if tlsCertFile != "" && tlsKeyFile != "" {
		log.Printf("serving HTTPS on %s", listenAddress)
		return http.ListenAndServeTLS(listenAddress, tlsCertFile, tlsKeyFile, tracker)
	}
	log.Printf("serving HTTP on %s", listenAddress)
	return http.ListenAndServe(listenAddress, tracker)
}
