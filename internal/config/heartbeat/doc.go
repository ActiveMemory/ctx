//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package heartbeat defines state file prefixes,
// filenames, and template variable keys for the ctx
// heartbeat subsystem.
//
// Heartbeats are periodic signals emitted during an
// AI session to track prompt counts, context window
// usage, and whether context files have been modified.
// The heartbeat hook uses per-session state files to
// persist counters between invocations and a log file
// for auditing.
//
// # State File Prefixes
//
//   - CounterPrefix ("heartbeat-") -- prefix for
//     per-session prompt counter files. The full
//     filename is CounterPrefix + session ID. Each
//     file stores the running prompt count for one
//     session.
//   - MtimePrefix ("heartbeat-mtime-") -- prefix for
//     per-session context modification tracking. The
//     file stores the last-seen mtime of context
//     files so the hook can detect mid-session edits.
//   - LogFile ("heartbeat.log") -- the log file for
//     heartbeat events, stored in .context/state/.
//
// # Template Variables
//
// These keys appear in heartbeat hook payloads and
// are used by the template engine to render heartbeat
// data:
//
//   - VarPromptCount -- number of prompts in the
//     current session
//   - VarSessionID -- the active session identifier
//   - VarContextModified -- boolean flag indicating
//     whether context files changed since last check
//   - VarTokens -- current token count
//   - VarContextWindow -- configured context window
//     size
//   - VarUsagePct -- percentage of context window
//     currently consumed
//
// # Why Centralized
//
// Heartbeat constants are shared between the heartbeat
// hook, the context-size monitor, the state file
// manager, and the doctor check that validates event
// logging. Centralizing them here avoids string drift
// and import cycles across these subsystems.
package heartbeat
