//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package pad provides error constructors for encrypted scratchpad operations.
//
// Error constructors return structured errors with context for
// user-facing messages routed through internal/assets text lookups.
// Exports: [EntryRange], [EditBlobTextConflict], [EditTextConflict], [EditNoMode], [BlobAppendNotAllowed], [BlobPrependNotAllowed].
package pad
