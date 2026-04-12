//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package pad implements the "ctx pad" command for managing an encrypted
// scratchpad.
//
// The scratchpad stores short, sensitive one-liners that travel with the
// project via git but remain opaque at rest. Entries are encrypted with
// AES-256-GCM using a symmetric key at ~/.ctx/.ctx.key.
//
// File blobs can be stored as entries using the format "label:::base64data".
// The add --file flag ingests a file and shows auto-decoded blob entries.
// Blobs are subject to a 64KB pre-encoding size limit.
//
// A plaintext fallback (.context/scratchpad.md) is available via the
// scratchpad_encrypt config option in .ctxrc.
//
// Subcommands:
//
//   - add:     append a text entry or file blob
//   - show:    display all entries (auto-decodes blobs)
//   - edit:    edit an entry by line number
//   - rm:      remove an entry by line number
//   - mv:      move an entry to a different position
//   - export:  export blob entries as files to a directory
//   - merge:   merge entries from external scratchpad files
//   - resolve: resolve merge conflicts in scratchpad
//   - normalize: reassign entry IDs as 1..N
//   - tag:     list all tags with counts
package pad
