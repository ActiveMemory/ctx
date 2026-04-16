//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package governance defines threshold constants that
// control when the MCP server emits governance nudges --
// gentle reminders that steer the AI agent toward good
// context hygiene.
//
// Two categories of nudge exist:
//
// # Drift Checks
//
// After [DriftCheckMinCalls] tool invocations without a
// drift check, the server reminds the agent to run
// ctx_drift. Once a drift check has occurred, the
// reminder will not fire again until at least
// [DriftCheckInterval] (15 minutes) has elapsed.
//
//   - [DriftCheckInterval] -- minimum wall-clock time
//     between successive drift reminders.
//   - [DriftCheckMinCalls] -- tool call floor before
//     the first drift reminder fires.
//
// # Persist Nudges
//
// If the agent has invoked [PersistNudgeAfter] (10)
// tools without writing to any context file, the server
// nudges it to persist learnings or decisions. The nudge
// repeats every [PersistNudgeRepeat] (8) additional
// calls until a write occurs.
//
//   - [PersistNudgeAfter]  -- tool call count that
//     triggers the first persist reminder.
//   - [PersistNudgeRepeat] -- interval in tool calls
//     between subsequent persist reminders.
//
// # Why These Are Centralized
//
// Hook handlers evaluate these thresholds on every tool
// call. Keeping them in one package makes tuning
// governance behavior a single-constant change visible
// in code review, rather than a scattered edit across
// multiple handler files.
package governance
