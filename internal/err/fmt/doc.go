//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package fmt defines the typed error constructors
// for the `ctx fmt` command, which normalizes
// whitespace, heading levels, and list markers
// across context files (TASKS.md, DECISIONS.md,
// LEARNINGS.md, CONVENTIONS.md).
//
// # Domain
//
// Errors fall into three categories:
//
//   - **Missing context**: the context directory
//     does not exist or contains no files.
//     Constructors: [NoContextDir], [NoFiles].
//   - **File IO**: a context file could not be
//     read or written during formatting.
//     Constructors: [FileRead], [FileWrite].
//   - **Check mode**: the formatter ran in
//     --check mode and found files that need
//     formatting. Constructor: [NeedsFormatting].
//
// # Wrapping Strategy
//
// IO constructors ([FileRead], [FileWrite]) wrap
// their cause with fmt.Errorf %w so callers can
// errors.Is against system errors. Pure validation
// constructors return plain errors. All user-facing
// text is resolved through
// [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package fmt
