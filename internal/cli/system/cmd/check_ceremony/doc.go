//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package check_ceremony implements the
// **`ctx system check-ceremony`** hidden hook, which
// nudges adoption of session ceremonies when they are
// missing from recent journal entries.
//
// # What It Does
//
// The hook scans recent journal files for evidence of
// two session ceremonies:
//
//   - **/ctx-remember**: the start-of-session
//     context recall ceremony.
//   - **/ctx-wrap-up**: the end-of-session
//     persistence ceremony.
//
// When either ceremony is absent from the lookback
// window, the hook emits a nudge message encouraging
// the user to adopt the missing ceremony. If both are
// present, the hook is silent.
//
// # Input
//
// A JSON hook envelope on stdin with session metadata.
//
// # Output
//
// On missing ceremony: a nudge message identifying
// which ceremony is absent. On both present or
// throttled: no output.
//
// # Throttling
//
// The hook is throttled to fire at most once per day
// using a marker file in the state directory.
//
// # Delegation
//
// [Cmd] builds the hidden cobra command. [Run] reads
// stdin via [core/check.Preamble], scans journals
// through [core/ceremony.RecentJournalFiles] and
// [core/ceremony.ScanJournalsForCeremonies], then
// emits the nudge via [write/setup.Nudge] and sends
// a relay notification through [core/nudge].
package check_ceremony
