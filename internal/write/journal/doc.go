//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package journal provides terminal output for journal
// commands (ctx journal source, site, lock, unlock,
// sync, and related subcommands).
//
// # Site Generation
//
// [InfoSiteGenerated] reports the final build result
// with entry count and output path. [InfoSiteStarting]
// and [InfoSiteBuilding] provide progress feedback.
// [InfoOrphanRemoved] reports cleanup of orphan files.
//
// # Session Import
//
// [ImportSummary] previews what an import will do.
// [ImportedFile] and [SkipFile] report per-file
// results. [ImportFinalSummary] prints aggregate
// counts for new, updated, renamed, and skipped
// files. [ConfirmPrompt] and [Aborted] handle the
// interactive confirmation flow.
//
// # Session Listing
//
// [SessionListHeader] prints the session count.
// [SessionListRow] prints a formatted table row.
// [SessionListFooter] prints the footer with an
// optional --limit hint. [NoSessionsForProject]
// and [NoSessionsWithHint] handle empty results.
// [NoFiltersMatch] handles empty filter results.
// [AmbiguousSessionMatch] and
// [AmbiguousSessionMatchWithHint] resolve ambiguous
// session queries.
//
// # Session Detail
//
// [SessionMetadata] prints the full metadata block
// with identity, timing, and token usage sections.
// [SessionDetail] and [SessionDetailInt] print
// individual metadata lines. [SectionHeader] prints
// a Markdown heading. [ConversationTurn] prints a
// turn header. [TextBlock] and [CodeBlock] render
// content blocks. [ListItem] and [NumberedItem]
// render list entries. [MoreTurns] prints a
// continuation notice. [Hint] prints usage tips.
//
// # Lock / Unlock / Sync
//
// [LockUnlockNone] handles empty entry sets.
// [LockUnlockEntry] confirms per-entry changes.
// [LockUnlockSummary] prints aggregate results.
// [SyncNone], [SyncLocked], [SyncUnlocked], and
// [SyncSummary] handle the sync workflow.
//
// # Types
//
// [SessionInfo] carries pre-formatted session metadata
// for display by [SessionMetadata].
package journal
