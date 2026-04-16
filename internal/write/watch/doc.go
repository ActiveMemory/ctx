//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package watch provides terminal output for the stdin
// watch command (ctx watch).
//
// Watch monitors stdin for context-update tags and
// applies them to context files in real time. Output
// functions cover the full lifecycle from startup
// through per-update results.
//
// # Startup
//
// [Started] confirms the watch loop began and is
// reading from stdin. [DryRun] prints a notice
// that updates will be previewed but not applied.
// [StopHint] shows the Ctrl+C hint for exiting.
//
// # Per-Update Results
//
// [DryRunPreview] shows what would be applied,
// printing the update type and content.
// [ApplySuccess] confirms an update was applied
// successfully. [ApplyFailed] reports a failure
// with the update type and error.
//
// # Visual Structure
//
// [Separator] prints a blank line between updates
// for visual clarity in the output stream.
//
// # Cleanup
//
// [CloseLogError] reports a log file close error
// during shutdown, printing the underlying error.
package watch
