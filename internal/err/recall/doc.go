//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package recall provides error constructors for session recall and reindexing.
//
// Error constructors return structured errors with context for
// user-facing messages routed through internal/assets text lookups.
// Exports: [EventLogRead], [StatsGlob], [ReindexFileNotFound], [ReindexFileRead], [ReindexFileWrite], [OpenLogFile].
package recall
