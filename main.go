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
		fmt.Printf("google-analytics-proxy: %v", err)
		os.Exit(1)
	}
}

func mainCmd() error {
	log.Printf("joshdk/google-analytics-proxy version %s", version)

	// listenAddress is the host and port that the proxy will listen on.
	// See net.Dial for details of the address format.
	// Example: "localhost:8080" "0.0.0.0:8080" ":8080"
	var listenAddress = os.Getenv("LISTEN")

	// upstreamEndpoint is the address of the upstream service to be
	// proxied.
	// Example: "https://example.com" "http://:80"
	var upstreamEndpoint = os.Getenv("UPSTREAM_ENDPOINT")

	// upstreamHostname optionally is the hostname to used when proxying
	// requests to the upstream. Used for hostname based routing. If empty, the
	// value of $GOOGLE_ANALYTICS_PROPERTY_NAME will be used.
	// Example: "example.com"
	var upstreamHostname = os.Getenv("UPSTREAM_HOSTNAME")

	// googleAnalyticsTrackingID is the tracking id for the Google
	// Analytics property that you want to track pageview events for. This
	// can be found in your Google Analytics dashboard.
	// Example: "UA-123456789-1"
	var googleAnalyticsTrackingID = os.Getenv("GOOGLE_ANALYTICS_TRACKING_ID")

	// googleAnalyticsPropertyName is the name for the Google Analytics
	// property that you want to track pageview events for. This can be
	// found in your Google Analytics dashboard. Will be used as the upstream
	// hostname in proxied requests if $UPSTREAM_HOSTNAME is empty.
	// Example: "example.com"
	var googleAnalyticsPropertyName = os.Getenv("GOOGLE_ANALYTICS_PROPERTY_NAME")

	// googleAnalyticsDryRun can optionally be used to disable reporting
	// pageview events with Google Analytics. See strconv.ParseBool() for
	// acceptable values.
	// Example: "true"
	var googleAnalyticsDryRun = os.Getenv("GOOGLE_ANALYTICS_DRY_RUN")

	// Parse the upstream endpoint address to ensure that it's valid.
	upstreamURL, err := url.Parse(upstreamEndpoint)
	if err != nil {
		return err
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
	log.Printf("listening on %s", listenAddress)
	return http.ListenAndServe(listenAddress, tracker)
}
