//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package export provides scratchpad blob export
// planning.
//
// The "ctx pad export" command extracts binary blob
// entries from the scratchpad and writes them to disk
// as individual files. This package handles the
// planning phase: deciding which blobs to export and
// resolving output paths.
//
// # Export Planning
//
// [Plan] is the primary function. It reads all
// scratchpad entries via store.ReadEntries, filters
// for blob entries using blob.Split, and builds an
// Item for each blob containing the label, decoded
// data, and target output path.
//
// # Collision Avoidance
//
// When force is false and the target file already
// exists, Plan generates an alternative filename by
// prepending a Unix timestamp via tsWithLabel. The
// Item records both the original path (with Exists
// set to true) and the AltName. The cmd layer uses
// these flags to inform the user of renames.
//
// When force is true, the output path is used as-is
// and existing files will be overwritten.
//
// # Item Type
//
// The [Item] struct carries: Label (display name),
// Data (decoded bytes), OutPath (target path), AltName
// (non-empty when collision-renamed), and Exists (true
// when the original path already exists).
//
// # Data Flow
//
// The cmd/pad layer calls Plan, iterates the returned
// items, writes each file to disk, and reports results
// via the write/pad package. Plan never writes files
// itself.
package export
