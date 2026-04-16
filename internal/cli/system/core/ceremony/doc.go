//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package ceremony detects missing context ceremonies
// in recent journal entries and generates nudge
// messages to remind agents.
//
// Ceremonies are structured rituals that agents should
// perform: /ctx-remember at session start to load
// context, and /ctx-wrap-up at session end to persist
// learnings. When recent journals lack evidence of
// these ceremonies, the package generates a nudge.
//
// # Journal Scanning
//
// [RecentJournalFiles] reads a journal directory and
// returns the n most recent markdown files, sorted by
// filename descending (newest first).
//
// [ScanJournalsForCeremonies] checks whether the given
// journal files contain references to the remember and
// wrap-up ceremony commands. It scans file contents
// for the ceremony command strings and returns boolean
// flags for each.
//
// # Nudge Emission
//
// [Emit] builds a formatted nudge message box based on
// which ceremonies are missing. It selects a variant
// (both, remember-only, or wrapup-only), loads the
// corresponding message template with an optional
// override from the message system, and wraps it in a
// visual box with a relay prefix.
//
// Returns an empty message when both ceremonies are
// present in recent sessions.
//
// # Data Flow
//
// The hook system calls RecentJournalFiles to get
// recent entries, passes them to
// ScanJournalsForCeremonies, then calls Emit with the
// results. The nudge is relayed to the agent through
// the hook output mechanism.
package ceremony
