//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package nudge defines configuration constants for
// the nudge subsystem, which reminds the AI agent to
// persist context (decisions, learnings, tasks) before
// a session grows too long without saves.
//
// # Persistence Nudge
//
// The persistence check tracks prompt count per
// session using state files prefixed with
// PersistencePrefix ("persistence-nudge-"). Nudging
// begins after PersistenceEarlyMin (11) prompts and
// repeats at PersistenceEarlyInterval (20) prompts
// during the early window (up to PersistenceEarlyMax,
// 25). After that, nudges fire every
// PersistenceLateInterval (15) prompts.
//
// State is tracked via three keys in the state file:
//   - PersistenceKeyCount: current prompt count
//   - KeyLastNudge: prompt number of last nudge
//   - PersistenceKeyLastMtime: last context mtime
//
// Events are logged to PersistenceLogFile
// ("check-persistence.log").
//
// # Task Completion Nudge
//
// The task completion check uses state files prefixed
// with PrefixTask ("task-nudge-") to track whether
// the agent has been reminded to complete open tasks
// in the current session.
//
// # Pause Behavior
//
// When nudges are paused, PauseTurnThreshold (5)
// controls whether the pause message shows a simple
// label or includes the turn count.
//
// # Template Variables
//
// VarPromptCount and VarSinceNudge are injected into
// hook message templates so the nudge text can report
// how many prompts have elapsed and how many since
// the last reminder.
package nudge
