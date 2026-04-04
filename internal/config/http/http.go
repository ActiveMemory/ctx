//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package http provides constants for HTTP operations (webhook
// notifications, outbound requests).
//
// Provides MIME types, timeouts, and URL masking constants.
// Import as config/http.
package http

// MIME type constants.
const (
	// MimeJSON is the Content-Type for JSON payloads.
	MimeJSON = "application/json"
)

// Timeout constants (in seconds).
const (
	// WebhookTimeout is the HTTP client timeout for webhook delivery.
	WebhookTimeout = 5
)

// URL scheme constants for validation.
const (
	// SchemeHTTP is the plain HTTP URL scheme.
	SchemeHTTP = "http"
	// SchemeHTTPS is the TLS-secured HTTP URL scheme.
	SchemeHTTPS = "https"
	// SchemeFile is the local file URL scheme.
	SchemeFile = "file"
)

// Full scheme prefix constants for URL prefix matching.
const (
	// PrefixHTTP is the full HTTP scheme prefix.
	PrefixHTTP = SchemeHTTP + "://"
	// PrefixHTTPS is the full HTTPS scheme prefix.
	PrefixHTTPS = SchemeHTTPS + "://"
	// PrefixFile is the full file scheme prefix.
	PrefixFile = SchemeFile + "://"
)

// URL constants.
const (
	// PathSep is the URL path separator.
	PathSep = '/'

	// PathSepStr is the string form of PathSep.
	PathSepStr = "/"
)

// URL masking constants for safe display of webhook URLs.
const (
	// MaskAfterSlash is the number of slashes (scheme://host/) after
	// which the URL path is replaced with MaskSuffix.
	MaskAfterSlash = 3

	// MaskMaxLen is the maximum visible characters when no third
	// slash is found.
	MaskMaxLen = 20

	// MaskSuffix is appended to the visible portion of a masked URL.
	MaskSuffix = "***"
)
