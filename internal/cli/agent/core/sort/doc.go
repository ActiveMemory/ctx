//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package sort provides the file read-order for context
// assembly in the agent command.
//
// # ReadOrder
//
// [ReadOrder] returns context file paths in the priority
// sequence defined by cfgCtx.ReadOrder. The order places
// high-priority files first: constitution, then tasks,
// conventions, architecture, decisions, learnings,
// glossary, and playbook.
//
// The function iterates cfgCtx.ReadOrder, looks up each
// file name in the loaded Context, and includes it only
// when the file exists and is non-empty. Paths are
// returned as full paths by joining ctx.Dir with the
// file name.
//
// # Filtering
//
// Empty files are excluded so the agent packet does not
// waste budget on placeholder files. The emptiness check
// uses the IsEmpty field on entity.ContextFile, which is
// set during context loading.
//
// # Data Flow
//
// The budget subpackage calls ReadOrder to determine
// which files to process and in what sequence. The
// returned paths feed into the section assembly loop
// where each file's content is extracted, scored, and
// truncated to fit the token budget.
package sort
