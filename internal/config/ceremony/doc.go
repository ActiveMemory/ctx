//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package ceremony defines configuration constants for the
// session ceremony detection system.
//
// ctx encourages two session rituals: /ctx-remember at the
// start (to reload context) and /ctx-wrap-up at the end
// (to persist learnings). The ceremony subsystem scans
// recent journal files to detect whether these rituals
// were performed and nudges the agent when they are
// missing.
//
// # Throttle Control
//
// [ThrottleID] names the daily throttle state file
// ("ceremony-reminded") so the nudge fires at most once
// per day. This prevents the ceremony reminder from
// becoming noisy during long sessions with many hook
// invocations.
//
// # Journal Scanning
//
// The ceremony detector looks backward through recent
// sessions to decide whether a nudge is needed:
//
//   - [JournalLookback] limits the scan to the 3 most
//     recent journal files. Scanning too far back would
//     generate false negatives (old sessions that pre-date
//     the ceremony system).
//   - [RememberCmd] is the command string ("ctx-remember")
//     searched for in journal content.
//   - [WrapUpCmd] is the command string ("ctx-wrap-up")
//     searched for in journal content.
//
// If neither command appears in the lookback window, the
// hook emits a nudge suggesting the agent run the missing
// ceremony.
//
// # Why Centralized
//
// The ceremony hook, the nudge renderer, and the journal
// scanner all need the same command names and lookback
// depth. Centralizing them here keeps the detection
// logic consistent and easy to tune.
package ceremony
