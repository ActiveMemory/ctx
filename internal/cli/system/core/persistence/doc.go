//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package persistence tracks edit counter state for the
// persistence nudge system. It decides when to remind the
// agent to save context based on the number of prompts
// since the last nudge or context file update.
//
// # State Management
//
// [State] holds three counters: the total edit/write call
// count, the prompt number of the last nudge, and the
// Unix timestamp of the last TASKS.md modification.
//
// [ReadState] parses a key=value state file from disk
// into a State struct. Returns ok=false when the file
// does not exist or cannot be read.
//
// [WriteState] serializes a State struct back to disk
// using the same key=value format.
//
// # Nudge Scheduling
//
// [NudgeNeeded] determines whether a persistence nudge
// should fire. It uses a two-tier interval strategy:
//
//   - Early range: between PersistenceEarlyMin and
//     PersistenceEarlyMax prompts, nudge every
//     PersistenceEarlyInterval prompts since last nudge
//   - Late range: beyond PersistenceEarlyMax prompts,
//     nudge every PersistenceLateInterval prompts
//
// This ensures frequent reminders during the critical
// early phase of a session, tapering to less frequent
// nudges as the session progresses.
package persistence
