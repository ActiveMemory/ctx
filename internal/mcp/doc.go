//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package mcp implements a Model Context Protocol server
// for ctx.
//
// MCP is a standard protocol (JSON-RPC 2.0 over
// stdin/stdout) that allows AI tools to discover and
// consume context from external sources. This package
// exposes ctx's context files as MCP resources and ctx
// commands as MCP tools, enabling any MCP-compatible AI
// tool (Claude Desktop, Cursor, Windsurf, VS Code
// Copilot, etc.) to access project context without
// tool-specific integrations.
//
// # Architecture
//
//	AI Tool -> stdin  -> MCP Server -> ctx internals
//	AI Tool <- stdout <- MCP Server <- ctx internals
//
// The server communicates via JSON-RPC 2.0 over
// stdin/stdout.
//
// # Resources
//
// Resources expose context files as read-only content:
//
//	ctx://context/tasks         -> TASKS.md
//	ctx://context/decisions     -> DECISIONS.md
//	ctx://context/conventions   -> CONVENTIONS.md
//	ctx://context/constitution  -> CONSTITUTION.md
//	ctx://context/architecture  -> ARCHITECTURE.md
//	ctx://context/learnings     -> LEARNINGS.md
//	ctx://context/glossary      -> GLOSSARY.md
//	ctx://context/agent         -> All files assembled
//
// # Tools
//
// Tools expose ctx commands as callable operations:
//
//	ctx_status             -> Context health summary
//	ctx_add                -> Add a task, decision, etc.
//	ctx_complete           -> Mark a task as done
//	ctx_drift              -> Detect stale context
//	ctx_journal_source     -> Query session history
//	ctx_watch_update       -> Apply structured updates
//	ctx_compact            -> Archive completed tasks
//	ctx_next               -> Get next pending task
//	ctx_check_task_completion -> Nudge on completion
//	ctx_session_event      -> Signal session lifecycle
//	ctx_remind             -> List active reminders
//
// # Prompts
//
// Prompts provide pre-built templates:
//
//	ctx-session-start  -> Load full context
//	ctx-decision-add   -> Format a decision entry
//	ctx-learning-add   -> Format a learning entry
//	ctx-reflect        -> Guide reflection
//	ctx-checkpoint     -> Report session statistics
//
// # Usage
//
//	server := mcp.New(contextDir, version)
//	server.Serve()  // blocks on stdin/stdout
//
// # Design Invariants
//
// This implementation preserves all ctx invariants:
//
//   - Markdown-on-filesystem: state in .context/
//   - Zero runtime dependencies
//   - Deterministic assembly
//   - Human authority
//   - Local-first: no network required
//   - No telemetry
package mcp
