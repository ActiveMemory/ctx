//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package mcp implements a Model Context Protocol (MCP) server for ctx.
//
// MCP is a standard protocol (JSON-RPC 2.0 over stdin/stdout) that allows
// AI tools to discover and consume context from external sources. This
// package exposes ctx's context files as MCP resources and ctx commands
// as MCP tools, enabling any MCP-compatible AI tool (Claude Desktop,
// Cursor, Windsurf, VS Code Copilot, etc.) to access project context
// without tool-specific integrations.
//
// # Architecture
//
//	AI Tool → stdin → MCP Server → ctx internals
//	AI Tool ← stdout ← MCP Server ← ctx internals
//
// The server communicates via JSON-RPC 2.0 over stdin/stdout.
//
// # Resources
//
// Resources expose context files as read-only content:
//
//	ctx://context/tasks         → TASKS.md
//	ctx://context/decisions     → DECISIONS.md
//	ctx://context/conventions   → CONVENTIONS.md
//	ctx://context/constitution  → CONSTITUTION.md
//	ctx://context/architecture  → ARCHITECTURE.md
//	ctx://context/learnings     → LEARNINGS.md
//	ctx://context/glossary      → GLOSSARY.md
//	ctx://context/agent         → All files assembled in read order
//
// # Tools
//
// Tools expose ctx commands as callable operations:
//
//	ctx_status                → Context health summary
//	ctx_add                   → Add a task, decision, learning, or convention
//	ctx_complete              → Mark a task as done
//	ctx_drift                 → Detect stale or invalid context
//	ctx_recall                → Query past session history
//	ctx_watch_update          → Apply structured context updates to files
//	ctx_compact               → Move completed tasks to archive
//	ctx_next                  → Get the next pending task
//	ctx_check_task_completion → Nudge when a recent action may complete a task
//	ctx_session_event         → Signal session start/end lifecycle
//	ctx_remind                → List active reminders
//
// # Prompts
//
// Prompts provide pre-built templates for common workflows:
//
//	ctx-session-start  → Load full context at session start
//	ctx-add-decision   → Format an architectural decision entry
//	ctx-add-learning   → Format a learning entry
//	ctx-reflect        → Guide end-of-session reflection
//	ctx-checkpoint     → Report session statistics
//
// # Usage
//
//	server := mcp.NewServer(contextDir, version)
//	server.Serve()  // blocks, reads stdin, writes stdout
//
// # Design Invariants
//
// This implementation preserves all six ctx design invariants:
//
//   - Markdown-on-filesystem: all state remains in .context/ files
//   - Zero runtime dependencies: no external services required
//   - Deterministic assembly: same files + budget = same output
//   - Human authority: tools propose changes through file writes
//   - Local-first: no network required for core operation
//   - No telemetry: no data leaves the local machine
//
// # File Organization
//
//   - prompts.go: MCP prompt definitions and handlers
//   - protocol.go: JSON-RPC 2.0 message types for MCP
//   - resources.go: MCP resource definitions, handlers, and subscription poller
//   - server.go: Server lifecycle, dispatch, and I/O handling
//   - session.go: Per-session advisory state tracking
//   - tools.go: MCP tool definitions and handlers
//   - types.go: Server struct definition
package mcp
