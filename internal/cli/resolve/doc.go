//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package resolve provides shared CLI flag resolution
// helpers used across multiple command packages.
//
// Several ctx commands need to determine which AI tool
// is active (e.g. claude, aider, cursor). Rather than
// duplicating the resolution logic, those packages call
// [Tool] to resolve the identifier from a consistent
// chain of sources.
//
// # Resolution Order
//
// [Tool] checks the following sources in order:
//
//  1. The --tool flag on the cobra command, if explicitly
//     set by the user
//  2. The tool field in .ctxrc, read via [rc.Tool]
//
// If neither source provides a value, an error is
// returned so the caller can surface a clear diagnostic.
//
// [Tool] checks the --tool flag first, then falls back
// to the tool field in .ctxrc, and returns an error if
// neither source provides a value.
package resolve
