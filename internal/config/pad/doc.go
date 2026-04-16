//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package pad defines constants for the encrypted
// scratchpad subsystem, which stores short, sensitive
// notes that travel with the project repository.
//
// # File Layout
//
// The scratchpad lives in .context/ as two files:
//
//   - Enc ("scratchpad.enc"): the encrypted file
//     that is committed to the repository.
//   - Md ("scratchpad.md"): the plaintext working
//     copy (gitignored).
//
// During merge conflicts the encrypted file may split
// into EncOurs and EncTheirs variants so the user can
// resolve manually.
//
// # Entry Formatting
//
// Each entry has a stable numeric ID rendered with
// FmtPadEntryID ("[%d] %s"). Merge conflict sides
// are labeled SideOurs ("OURS") and SideTheirs
// ("THEIRS") in the resolution UI.
//
// # Blob Support
//
// Small binary files (up to MaxBlobSize, 64 KB) can
// be attached to scratchpad entries as base64 blobs.
// The blob label and content are separated by BlobSep
// (":::"), and blob entries are tagged with BlobTag
// (" [BLOB]") in list output.
//
// # Tag Filtering
//
// Entries can be tagged with hashtags. TagPrefix ("#")
// and TagPrefixSpace (" #") control how tags are
// rendered. TagNegate ("~") prefixes negated tag
// filters in search queries.
package pad
