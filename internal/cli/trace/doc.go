//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package trace implements the ctx trace command for
// commit context tracing.
//
// The trace command links git commits to the context
// that motivated them. It captures which decisions,
// tasks, and learnings were active when a commit was
// made, creating an audit trail from code changes back
// to their rationale.
//
// # Subcommands
//
//   - show: display trace entries for recent commits,
//     showing linked context items (default subcommand)
//   - collect: capture the current context snapshot and
//     associate it with the latest commit
//   - file: show trace entries that reference a specific
//     file path
//   - hook: post-commit hook entry point that
//     automatically collects trace data
//   - tag: annotate a trace entry with additional
//     metadata tags
//
// # Subpackages
//
//	cmd/show: trace display and formatting
//	cmd/collect: context snapshot collection
//	cmd/file: file-scoped trace lookup
//	cmd/hook: post-commit automation
//	cmd/tag: trace annotation
//	core: trace storage, git integration, and
//	  context linking
package trace
