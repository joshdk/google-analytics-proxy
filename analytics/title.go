// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.

package analytics

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http/httptest"
	"strings"

	"github.com/andybalholm/brotli"
	"golang.org/x/net/html"
)

func getTitle(recorder *httptest.ResponseRecorder) (title string, err error) {
	// The content type that was returned. The <title> tag can only be
	// extracted from an HTML response body.
	contentType := recorder.Header().Get("Content-Type")

	// Only bother parsing HTML if there is HTML to parse.
	if !strings.Contains(contentType, "text/html") {
		return
	}

	// The content encoding that was returned. The response body may be be
	// compressed, so detect the type of compression, if any, and decode the
	// response body accordingly.
	contentEncoding := recorder.Header().Get("Content-Encoding")
	var reader io.Reader = bytes.NewBuffer(recorder.Body.Bytes())

	// How is the request compressed?
	switch contentEncoding {
	case "gzip":
		reader, err = gzip.NewReader(reader)
		if err != nil {
			return
		}
	case "br":
		reader = brotli.NewReader(reader)
	}

	// Extract title from the (now decoded) response body.
	return getTitleFromBody(reader), nil
}

func getTitleFromBody(body io.Reader) string {
	tokenizer := html.NewTokenizer(body)

	// Stop iterating over the document after some "reasonable" point.
	const cutoff = 100

	for iteration := 0; iteration < cutoff; iteration++ {
		switch tokenizer.Next() {
		case html.StartTagToken:
			// Is the current tag a <title> tag?
			if tokenizer.Token().Data != "title" {
				continue
			}

			// Extract the actual title from the <title> tag.
			if tokenizer.Next() == html.TextToken {
				return tokenizer.Token().Data
			}

		case html.EndTagToken:
			// We reached the end of the <head> tag, but never found a <title>.
			if tokenizer.Token().Data == "head" {
				return ""
			}

		case html.ErrorToken:
			// We have finished iterating over the document, or there was an
			// error.
			return ""
		}
	}

	// We were unable to locate a title by the cutoff point.
	return ""
}
