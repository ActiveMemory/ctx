//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package pad defines the typed error constructors
// for the encrypted scratchpad (`ctx pad`). These
// errors fire during entry selection, editing mode
// validation, blob operations, merge conflict
// resolution, and scratchpad reads.
//
// # Domain
//
// Errors fall into four categories:
//
//   - **Entry selection**: the requested entry
//     index is out of range, not found by ID, or
//     not a valid number. Constructors:
//     [EntryRange], [EntryNotFound],
//     [InvalidIndex].
//   - **Editing modes**: mutually exclusive edit
//     flags were combined, or no mode was given.
//     Constructors: [EditBlobTextConflict],
//     [EditTextConflict], [EditNoMode].
//   - **Blob operations**: a blob-only flag was
//     used on a text entry, or a file exceeds the
//     size limit. Constructors: [NotBlobEntry],
//     [OutFlagRequiresBlob], [FileTooLarge].
//   - **Conflict resolution**: the scratchpad is
//     not encrypted, or no conflict files exist.
//     Constructors: [ResolveNotEncrypted],
//     [NoConflictFiles], [Read].
//
// # Wrapping Strategy
//
// [Read] wraps its cause with fmt.Errorf %w.
// Validation constructors return plain errors
// because the failures are policy violations, not
// IO errors. All user-facing text is resolved
// through [internal/assets/read/desc].
//
// # Concurrency
//
// Pure constructors. Concurrent callers never race.
package pad
