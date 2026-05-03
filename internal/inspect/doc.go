//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package inspect provides general-purpose string
// predicates and position-tracking utilities used
// throughout the ctx codebase for text scanning.
//
// # Position Helpers (pos.go)
//
// Functions that advance a cursor through a string,
// handling both LF and CRLF line endings:
//
//   - [SkipNewline] advances past a newline character
//     (CRLF or LF) at the current position. Returns
//     the position unchanged if no newline is present.
//   - [SkipWhitespace] advances past any sequence of
//     spaces, tabs, and newlines.
//   - [FindNewline] returns the index of the first
//     newline in a string, or -1 if none exists.
//
// # String Predicates (predicate.go)
//
// Boolean checks and index lookups for common text
// patterns:
//
//   - [EndsWithNewline] reports whether a string ends
//     with a newline (CRLF or LF).
//   - [Contains] reports whether a substring exists
//     and returns its index.
//   - [ContainsNewLine] checks for any newline and
//     returns its index.
//   - [ContainsEndComment] checks for a comment close
//     marker and returns its index.
//
// # Design
//
// All functions are pure and safe for concurrent use.
// They use token and marker constants from the config
// packages rather than hardcoded literals, ensuring
// consistency with the rest of the codebase.
package inspect
