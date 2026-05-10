//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package validate provides input-validation helpers
// that ctx uses at filesystem, security, and CLI
// boundaries.
//
// # Path Validation
//
//   - [Boundary] checks that a directory resolves to
//     a path within the current working directory.
//     Resolves symlinks in both paths so traversal
//     via symlinked parents is caught. On Windows,
//     comparisons are case-insensitive to handle
//     NTFS path normalization. Returns a typed error
//     from [internal/err/context] if the path escapes
//     the project root.
//   - [Symlinks] checks whether a directory or any
//     of its immediate children are symlinks. Returns
//     a typed error describing the first symlink
//     found. Non-existent directories are not an
//     error (let the caller handle that).
//
// # CLI Body-Flag Validation
//
//   - [BodyFlags] reads each named flag from a cobra
//     command and rejects empty, whitespace-only, or
//     placeholder values (TBD, see chat, n/a, etc.).
//     Pure function: does not mutate the command.
//     Called from a noun-level command's own PreRunE
//     so the wiring is visible at the call site.
//   - [RejectPlaceholder] is the per-value primitive
//     used by [BodyFlags]; exported for tests and
//     ad-hoc reuse. Placeholder values live in
//     [internal/config/validate].
//
// # Design Philosophy
//
// Unlike [internal/sanitize] (which transforms bad
// input into safe values), this package rejects bad
// input outright. Unlike [internal/io] (which guards
// against system directory access), this package
// guards against project-boundary escapes,
// symlink-based traversal, and missing or
// placeholder CLI body fields.
//
// # Concurrency
//
// All functions are pure and safe for concurrent
// use. They rely on os.Getwd, filepath.Abs, and
// filepath.EvalSymlinks, which are themselves
// goroutine-safe.
package validate
