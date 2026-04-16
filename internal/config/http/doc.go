//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package http centralizes HTTP-related constants used
// by the ctx CLI for webhook delivery, URL validation,
// and URL display masking.
//
// ctx supports webhook notifications: when context
// events occur, the CLI delivers JSON payloads to
// user-configured URLs. This package provides the MIME
// types, timeouts, scheme identifiers, and masking
// parameters that govern that delivery pipeline and
// URL handling throughout the codebase.
//
// # MIME Types
//
//   - MimeJSON ("application/json") -- the Content-Type
//     header set on webhook POST requests
//
// # Timeouts
//
//   - WebhookTimeout (5 seconds) -- the HTTP client
//     timeout for webhook delivery. Keeps the CLI
//     responsive even when endpoints are slow.
//
// # URL Scheme Constants
//
// Scheme identifiers for URL validation:
//
//   - SchemeHTTP ("http") -- plain HTTP
//   - SchemeHTTPS ("https") -- TLS-secured HTTP
//   - SchemeFile ("file") -- local file URLs
//
// Full scheme prefix strings for prefix matching:
//
//   - PrefixHTTP ("http://")
//   - PrefixHTTPS ("https://")
//   - PrefixFile ("file://")
//
// # URL Path Helpers
//
//   - PathSep ('/') -- URL path separator as a byte
//   - PathSepStr ("/") -- URL path separator as a
//     string
//
// # URL Masking
//
// When displaying webhook URLs in logs or status
// output, ctx masks the path portion to avoid leaking
// tokens or secrets embedded in URLs:
//
//   - MaskAfterSlash (3) -- number of slashes after
//     which the path is replaced (scheme://host/...)
//   - MaskMaxLen (20) -- max visible characters when
//     no third slash is found
//   - MaskSuffix ("***") -- appended to the visible
//     portion of a masked URL
//
// # Why Centralized
//
// HTTP constants are referenced by webhook delivery,
// URL validation in .ctxrc parsing, hub client
// connections, and CLI display code. A single source
// of truth prevents inconsistencies in timeouts or
// scheme handling across these subsystems.
package http
