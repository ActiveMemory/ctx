//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sync provides terminal output for the context
// sync command (ctx sync).
//
// The sync command has two modes: dependency/config
// analysis and MEMORY.md mirror synchronization.
// Output functions cover both workflows.
//
// # Dependency Analysis
//
// [AllClear] prints the all-clear message when the
// context is fully in sync and no actions are
// needed.
//
// [Header] prints the analysis heading with an
// optional dry-run notice. [Action] prints a
// numbered sync action item with its type label,
// description, and optional suggestion text.
// [Summary] closes the analysis with the total
// action count, adjusting wording for dry-run.
//
// # Memory Mirror
//
// [DryRun] prints the dry-run plan block with the
// source path, mirror path, and drift status.
// [Result] prints the full sync result: optional
// archive notice, synced confirmation, source
// path, and line counts with new-content delta.
//
// [ErrAutoMemoryNotActive] prints an informational
// stderr message when automatic memory source
// discovery fails.
package sync
