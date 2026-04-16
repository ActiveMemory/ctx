//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package date defines the typed error constructors
// for date parsing and validation. These errors fire
// when the user supplies a malformed date string to
// flags like --since or --until, or when a date value
// in context metadata fails to parse.
//
// # Domain
//
// Two constructors cover the entire surface:
//
//   - [InvalidValue] -- a standalone date string
//     does not match the expected YYYY-MM-DD format.
//     Used during metadata validation.
//   - [Invalid] -- a date flag value fails to parse.
//     Wraps the underlying time.Parse error and
//     includes the flag name for context.
//
// # Wrapping Strategy
//
// [Invalid] wraps its cause with fmt.Errorf %w so
// callers can inspect the underlying parse error.
// [InvalidValue] returns a plain error because
// there is no system cause to chain. All user-
// facing text is resolved through
// [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package date
