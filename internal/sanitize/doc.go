//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sanitize transforms untrusted input into
// safe values suitable for use as filesystem names
// or for writing to context files.
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
//   - [Content] escapes Markdown structural patterns
//     (entry headers, task checkboxes, constitution
//     rules) and strips null bytes from untrusted
//     content before writing to .context/.
//
//   - [StripControl] removes ASCII control characters
//     from a string while preserving tabs and newlines.
//
//   - [Reflect] strips control chars and truncates to
//     a maximum length; used when reflecting untrusted
//     input back in error messages.
//
//   - [SessionID] converts an arbitrary string into a
//     path-safe session identifier: strips null bytes,
//     path traversal sequences, slashes, replaces
//     unsafe characters with hyphens, and truncates to
//     [config/sanitize.MaxSessionIDLen].
//
// # Design
//
// All functions are pure and safe for concurrent use.
package sanitize
