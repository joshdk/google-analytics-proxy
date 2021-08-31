// Copyright Josh Komoroske. All rights reserved.
// Use of this source code is governed by the MIT license,
// a copy of which can be found in the LICENSE.txt file.

package analytics

import (
	"net/http"

	"github.com/gofrs/uuid"
)

const twoYearsInSeconds = 63072000

func getCookie(request *http.Request, name string) (string, *http.Cookie) {
	// Check if the named cookie was provided along with the request and return
	// the value it contains.
	if cookie, _ := request.Cookie(name); cookie != nil {
		return cookie.Value, nil
	}

	// The named cookie was not provided, so make a new one.
	cookie := &http.Cookie{
		Name: name,

		// A UUID (version 4) for identifying the current client.
		// See: https://developers.google.com/analytics/devguides/collection/protocol/v1/parameters#cid
		Value: uuid.Must(uuid.NewV4()).String(),

		// Set the cookie expiration to two years (in seconds).
		// See: https://developers.google.com/analytics/devguides/collection/analyticsjs/cookies-user-id#configuring_cookie_field_settings
		MaxAge: twoYearsInSeconds,
	}

	return cookie.Value, cookie
}
