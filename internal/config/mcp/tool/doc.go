//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package tool defines the MCP tool registration names
// that map each MCP tool to a ctx CLI subcommand.
//
// When an AI agent calls tools/call, it passes a tool
// name in the request. The MCP server matches this name
// against the constants here to dispatch the request to
// the correct handler. Each constant follows the
// "ctx_<subcommand>" naming convention so the mapping
// between MCP tool and CLI command is transparent.
//
// # Key Constants
//
//   - [Status] ("ctx_status"): returns a summary
//     of all context files and their freshness.
//   - [Add] ("ctx_add"): appends a new entry to a
//     context file (decision, learning, task, etc.).
//   - [Complete] ("ctx_complete"): marks a task as
//     done in TASKS.md.
//   - [Drift] ("ctx_drift"): runs drift detection
//     across context files.
//   - [JournalSource] ("ctx_journal_source"): queries
//     past session journal entries.
//   - [WatchUpdate] ("ctx_watch_update"): writes a
//     structured update to a context file.
//   - [Compact] ("ctx_compact"): compacts completed
//     tasks in TASKS.md.
//   - [Next] ("ctx_next"): suggests the next task
//     to work on.
//   - [CheckTaskCompletion]
//     ("ctx_check_task_completion"): checks whether
//     a recent action completed a pending task.
//   - [SessionEvent] ("ctx_session_event"): records
//     session start/end lifecycle events.
//   - [Remind] ("ctx_remind"): lists active
//     reminders for the current session.
//   - [SteeringGet] ("ctx_steering_get"): retrieves
//     a steering file matched to a prompt.
//   - [Search] ("ctx_search"): full-text search
//     across context files.
//   - [SessionStart] / [SessionEnd]: hooks that run
//     at session boundaries.
//
// # Why These Are Centralized
//
// Tool registration, the dispatch switch, integration
// tests, and governance hooks all reference tool names.
// A constant ensures registration and dispatch always
// agree, and makes the full tool surface discoverable.
package tool
