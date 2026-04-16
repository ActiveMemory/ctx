//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sanitize transforms untrusted input into
// safe values suitable for use as filesystem names.
//
// Unlike validation (which rejects bad input),
// sanitization mutates input to conform to
// constraints. The result is always usable; the
// caller never needs to handle an error.
//
// # Public Surface
//
//   - [Filename] converts an arbitrary topic string
//     into a safe filename component: replaces spaces
//     and special characters with hyphens via
//     [regex.FileNameChar], strips leading and
//     trailing hyphens, converts to lowercase, and
//     limits the result to 50 characters. Returns
//     "session" if the input is empty after cleaning.
//
// # Design
//
// The function is idempotent: sanitizing an already-
// safe string returns it unchanged. It uses config
// constants for the replacement character, max length,
// and default fallback rather than hardcoded literals.
//
// All functions are pure and safe for concurrent use.
package sanitize
