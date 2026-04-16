//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package doctor centralizes check-name, category, and
// threshold constants consumed by the ctx doctor command.
//
// The doctor subsystem runs a suite of health checks
// against a project's .context/ directory. Each check
// produces a Result whose Name field is one of the
// Check* constants defined here (e.g. CheckContextInit,
// CheckDrift, CheckPluginInstalled). Results are grouped
// by Category for display purposes.
//
// # Check Names
//
// Every constant maps 1-to-1 to a diagnostic routine:
//
//   - CheckContextInit -- verifies .context/ exists
//   - CheckRequiredFiles -- ensures core markdown files
//     (TASKS.md, DECISIONS.md, etc.) are present
//   - CheckCtxrcValidation -- validates .ctxrc syntax
//   - CheckDrift -- delegates to drift detection
//   - CheckPluginInstalled, CheckPluginEnabledGlobal,
//     CheckPluginEnabledLocal -- verify the companion
//     plugin is installed and activated
//   - CheckCompanionConfig -- checks companion tool
//     configuration status
//   - CheckEventLogging -- confirms event log is writable
//   - CheckWebhook -- validates webhook URL and delivery
//   - CheckReminders -- flags pending reminders
//   - CheckTaskCompletion -- warns when too many tasks
//     remain incomplete
//   - CheckContextSize, CheckContextFilePrefix -- measure
//     token usage against the configured budget
//   - CheckRecentEvents -- reviews recent event history
//   - CheckResourceMemory, CheckResourceSwap,
//     CheckResourceDisk, CheckResourceLoad -- system
//     resource health probes
//
// # Categories
//
// Results are bucketed into categories for grouped
// display: Structure, Quality, Plugin, Hooks, State,
// Size, Resources, and Events.
//
// # Thresholds
//
// Numeric thresholds control when a check transitions
// from pass to warning:
//
//   - TaskCompletionWarnPct (80) -- warn when fewer
//     than 80% of tasks are completed
//   - TaskCompletionMinCount (5) -- skip ratio check
//     until at least 5 tasks exist
//   - ContextSizeWarnPct (20) -- warn when context
//     window usage exceeds 20%
//
// # Why Centralized
//
// Keeping these strings and thresholds in a dedicated
// config package prevents import cycles between the
// doctor runner, individual check implementations, and
// CLI rendering code. Any package that produces or
// consumes doctor results imports config/doctor instead
// of depending on the runner.
package doctor
