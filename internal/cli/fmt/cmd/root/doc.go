//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx fmt" command that
// formats context files to a consistent line width.
//
// # What It Does
//
// Iterates over the four context files (TASKS.md,
// DECISIONS.md, LEARNINGS.md, CONVENTIONS.md),
// applies word-wrapping to each, and writes back
// only files that changed. In check mode it reports
// which files need formatting without modifying them.
//
// # Flags
//
//   - --width: Target line width in characters.
//     Defaults to [wrap.DefaultWidth] (72).
//   - --check: Dry-run mode. Reports files that
//     would change and exits with code 1 if any
//     need formatting. Useful in CI pipelines.
//
// # Output
//
// In normal mode, prints a summary line like
// "Formatted 2/4 files." In check mode, prints
// each file that needs formatting and exits with
// a non-zero status if any were found.
//
// # Delegation
//
// [Cmd] builds the cobra.Command with --width and
// --check flags. [Run] resolves the context
// directory via [rc.ContextDir], reads each file
// with [io.SafeReadUserFile], wraps content via
// [wrap.ContextFile], and writes back with
// [io.SafeWriteFile].
package root
