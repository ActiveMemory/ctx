//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package add provides scratchpad entry creation logic
// with stable ID assignment.
//
// The "ctx pad add" command appends a new text or blob
// entry to the encrypted scratchpad. This package
// handles ID generation and entry construction without
// performing the final write.
//
// # Adding Text Entries
//
// [EntryWithID] loads the current scratchpad entries
// with their IDs via store.ReadEntriesWithIDs, computes
// the next available ID using parse.NextID, appends a
// new parse.Entry with the given text, and returns the
// updated slice together with the assigned ID. The
// caller (cmd layer) is responsible for encrypting and
// writing the result.
//
// # Adding Blob Entries
//
// [BlobWithID] reads a file from disk via SafeReadUserFile,
// validates that its size does not exceed MaxBlobSize,
// encodes it as a blob entry using blob.Make (base64
// with a label prefix), and appends it with a stable
// ID. Like EntryWithID, it returns the updated entries
// and the new ID without writing.
//
// # ID Assignment
//
// Both functions use parse.NextID to find the smallest
// unused integer ID across existing entries. IDs are
// stable: deleting an entry does not reassign IDs to
// remaining entries.
package add
