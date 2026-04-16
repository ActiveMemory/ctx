//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package freshness defines constants for the ctx
// freshness checking subsystem, which detects context
// files that have not been reviewed or updated for an
// extended period.
//
// Context files lose value when they go stale. The
// freshness subsystem runs as a hook-time check: it
// compares each tracked file's modification time
// against a threshold and emits a nudge when files
// need attention.
//
// # Threshold
//
// StaleThreshold (approximately 182 days / 6 months)
// is the duration after which a tracked context file
// is considered stale. Files older than this threshold
// are flagged for review in hook output, prompting the
// user or agent to verify that the content is still
// accurate.
//
// # Throttle
//
// ThrottleID ("freshness-checked") is the state file
// name used for daily throttling. Once a freshness
// check runs, this marker prevents repeat checks for
// the rest of the day, keeping hook output clean
// during active sessions.
//
// # Template Variables
//
// VarStaleFiles is the template variable key injected
// into freshness hook payloads. It contains the list
// of stale file paths so the hook template can render
// a human-readable nudge.
//
// # Why Centralized
//
// The freshness threshold, throttle ID, and template
// variable are referenced by the hook runner, the
// freshness scanner, and the drift detection subsystem.
// Centralizing them here prevents divergence and makes
// it easy to tune the staleness window project-wide.
package freshness
