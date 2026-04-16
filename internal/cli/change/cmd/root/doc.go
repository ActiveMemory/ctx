//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package root implements the "ctx change" command for
// detecting context and code changes since a reference
// time.
//
// # What It Does
//
// The command shows what changed in the project since
// a given point in time. It scans both context files
// (.context/ directory) and git history, then renders
// a unified summary to stdout.
//
// # Flags
//
//   - --since: Reference time as a duration (e.g.
//     "24h", "7d") or a date string (e.g.
//     "2026-03-01"). When omitted the command falls
//     back to session markers or the event log to
//     find the most recent reference point.
//
// # Output
//
// A human-readable list grouped by change source.
// Context-file changes show which files were modified
// and how many entries were added. Code changes show
// a git-log summary with commit counts and affected
// files.
//
// # Delegation
//
// [Cmd] builds the cobra.Command and binds the
// --since flag. [Run] resolves the reference time
// via [detect.ReferenceTime], scans for context and
// code changes via [scan], and renders through
// [render.List].
package root
