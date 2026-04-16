//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package knowledge defines configuration constants
// for the knowledge-check hook, which validates that
// context knowledge files (LEARNINGS.md, DECISIONS.md,
// CONVENTIONS.md) remain healthy.
//
// # What the Hook Does
//
// During a session the knowledge hook inspects the
// context knowledge files for common problems: stale
// references, missing sections, oversized entries,
// and structural drift. When issues are found it
// builds a warning list and injects it into the
// nudge message via the VarFileWarnings template
// variable.
//
// # Throttling
//
// The check runs at most once per day per session.
// ThrottleID ("check-knowledge") names the state
// file that tracks the last-run timestamp, reusing
// the same daily-throttle mechanism as the journal
// and persistence checks.
//
// # Key Constants
//
//   - ThrottleID -- state file name for daily
//     throttle ("check-knowledge").
//   - VarFileWarnings -- template variable key
//     injected into hook messages with the list
//     of file-level warnings.
package knowledge
