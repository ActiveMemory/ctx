//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package core contains business logic for the change
// command, which reports what changed in the project
// since a reference point in time.
//
// This package is an umbrella that coordinates three
// subpackages, each handling one stage of the change
// detection pipeline:
//
// # Reference Time Detection (detect/)
//
// The detect subpackage resolves the reference time
// from which changes are measured. It parses human
// timestamps, duration strings, and session identifiers
// to produce a concrete time.Time.
//
// # Filesystem and Git Scanning (scan/)
//
// The scan subpackage queries two data sources:
//
//   - Context file changes: it reads .context/ and
//     returns Markdown files modified after the
//     reference time, sorted by modification time
//     descending.
//   - Code changes: it runs git log to summarize
//     commit count, latest message, affected
//     directories, and contributing authors since the
//     reference time.
//
// # Output Rendering (render/)
//
// The render subpackage formats scan results for two
// audiences:
//
//   - [render.List] produces a full Markdown report for
//     terminal display with headings, file lists, and
//     code summaries.
//   - [render.ChangesForHook] produces a compact
//     single-line summary for hook relay injection
//     into AI tool prompts.
//
// # Data Flow
//
// The cmd/ layer resolves a reference time via detect,
// calls scan.FindContextChanges and
// scan.SummarizeCodeChanges, then passes the results
// to render.List or render.ChangesForHook depending on
// the output mode. The write/ layer handles final
// output to the cobra command.
package core
