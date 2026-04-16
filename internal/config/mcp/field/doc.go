//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package field defines the JSON property key names that
// appear in MCP tool input schemas. Every tool exposed
// by the ctx MCP server declares an input object whose
// properties are keyed by the constants in this package.
//
// When an AI agent calls a tool it sends a JSON object
// like {"content": "...", "priority": "high"}. The MCP
// handler extracts values by these exact key names, so
// both the schema declaration and the handler must agree
// on the strings. Centralizing them here prevents drift.
//
// # Key Constants
//
//   - [Content]      : the main text body passed to
//     ctx_add and ctx_watch_update.
//   - [Priority]     : task priority ("high",
//     "medium", "low") for ctx_add.
//   - [Section]      : target section in TASKS.md.
//   - [Query]        : search text or task number
//     for ctx_complete and ctx_search.
//   - [RecentAction] : recent action description for
//     the task completion nudge.
//   - [Caller]       : identifies the MCP client
//     (cursor, vscode, etc.).
//   - [Limit], [Since]: pagination and date filters
//     for ctx_journal_source.
//   - [SessionID], [Branch], [Commit]: provenance
//     metadata attached to journal entries.
//   - [Prompt], [Summary]: optional fields for
//     steering file matching and session-end hooks.
//
// # Why These Are Centralized
//
// Tool input schemas are declared in one place and
// parsed in another. If a key is renamed in the schema
// but not in the handler (or vice-versa), the tool
// silently ignores the field. Typed constants make such
// mismatches a compile error.
package field
