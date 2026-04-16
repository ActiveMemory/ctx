//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package http defines the typed error constructors
// for HTTP client safety checks. These errors fire
// when the CLI validates URLs before making outbound
// requests, enforcing scheme restrictions and
// redirect limits.
//
// # Domain
//
// Three constructors cover the entire surface:
//
//   - [UnsafeURLScheme]: a URL uses a scheme
//     other than http or https. This is a security
//     boundary that prevents file://, ftp://, or
//     other protocol handlers from being invoked.
//   - [ParseURL]: a URL string failed to parse.
//     Wraps the underlying url.Parse error.
//   - [TooManyRedirects]: an HTTP response chain
//     exceeded the configured redirect limit.
//
// # Wrapping Strategy
//
// [ParseURL] wraps its cause with fmt.Errorf %w.
// The other two return plain errors because they
// represent policy violations, not IO failures.
// All user-facing text is resolved through
// [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package http
