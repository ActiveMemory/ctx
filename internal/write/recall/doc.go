//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package recall provides terminal output for session history
// commands (ctx recall list, show, import, lock, unlock, sync).
//
// The package is the largest in write/ because recall has many
// output modes:
//
//   - List: [SessionListHeader], [SessionListRow],
//     [SessionListFooter] render a tabular session list.
//   - Show: [SessionMetadata], [SessionDetail], [ConversationTurn],
//     [TextBlock], [CodeBlock], [MoreTurns] render a single session.
//   - Import: [ImportedFile], [SkipFile], [ImportSummary],
//     [ImportFinalSummary] report import progress.
//   - Lock/unlock: [LockUnlockEntry], [LockUnlockSummary],
//     [LockUnlockNone] report lock state changes.
//   - Sync: [JournalSyncLocked], [JournalSyncUnlocked],
//     [JournalSyncSummary] report frontmatter-to-state sync.
//   - Errors: [AmbiguousSessionMatch], [NoSessionsForProject],
//     [NoFiltersMatch] handle lookup failures.
package recall
