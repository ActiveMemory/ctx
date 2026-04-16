//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package warn provides a centralized stderr warning
// sink for best-effort operations whose errors would
// otherwise be silently discarded.
//
// # The Problem It Solves
//
// Many ctx operations are fire-and-forget: closing
// file handles, removing temporary files, writing
// state markers, and appending to JSONL logs. When
// these fail, the error is not actionable by the
// caller, but silently swallowing it makes debugging
// harder. This package provides a single function
// that formats and emits the warning consistently.
//
// # Public Surface
//
//   - [Warn] formats a message with Printf-style
//     arguments, prefixes it with "ctx: ", appends
//     a newline, and writes to the sink. Sink write
//     failures are silently dropped because there is
//     nowhere else to report them.
//
// # Sink Replacement
//
// The sink variable defaults to os.Stderr. Tests
// replace it with io.Discard to suppress output
// during test runs.
//
// # Usage Pattern
//
// Callers throughout ctx use warn.Warn in place of
// log.Println or fmt.Fprintf(os.Stderr, ...) to keep
// all warning output prefixed and formatted
// identically.
package warn
