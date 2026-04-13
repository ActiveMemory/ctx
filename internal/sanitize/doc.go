//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sanitize transforms untrusted input into safe values.
//
// Unlike validation (which rejects bad input), sanitization mutates
// input to conform to constraints. [Filename] converts arbitrary
// strings into safe filename components, [Content] neutralizes
// Markdown structure injections, [Reflect] truncates and strips
// control characters for error messages, and [SessionID] produces
// path-safe session identifiers.
// Part of the internal subsystem.
package sanitize
