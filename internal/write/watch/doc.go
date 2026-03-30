//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package watch provides terminal output for the stdin watch
// command (ctx watch).
//
// Watch monitors stdin for context-update tags and applies them
// to context files. Output functions cover the lifecycle:
// [Started] confirms the watch loop began, [DryRun] enables
// preview mode, [StopHint] shows how to exit, [DryRunPreview]
// shows what would be applied, [ApplySuccess]/[ApplyFailed]
// report per-update results, and [Separator] visually separates
// updates. [CloseLogError] reports log cleanup failures.
package watch
