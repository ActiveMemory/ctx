//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package edit provides scratchpad entry mutation
// logic.
//
// The "ctx pad edit" command modifies an existing
// scratchpad entry in place. This package implements
// four mutation modes, each returning the updated
// entries slice without performing the write.
//
// # Edit Modes
//
// The [Mode] enum and [Opts] struct define the four
// operations:
//
//   - ModeReplace / [Replace]: overwrites the entry
//     at 1-based index n with new text.
//   - ModeAppend / [Append]: appends text after the
//     existing entry. For blob entries, the text is
//     appended to the label portion only.
//   - ModePrepend / [Prepend]: prepends text before
//     the existing entry. For blob entries, the text
//     is prepended to the label portion only.
//   - ModeBlob / [UpdateBlob]: replaces the file
//     content and/or label of a blob entry. Either
//     field can be left empty to keep the existing
//     value. File size is validated against
//     MaxBlobSize.
//
// # Common Pattern
//
// Every function follows the same sequence:
//
//  1. Load current entries via store.ReadEntries.
//  2. Validate the index via validate.Index.
//  3. Apply the mutation.
//  4. Return the updated slice. The caller owns
//     writing the result via store.WriteEntries.
//
// Blob-aware functions use blob.Split and blob.Make
// to preserve the base64 payload while modifying only
// the label.
package edit
