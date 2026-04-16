//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package merge reads and merges external scratchpad
// files with conflict detection.
//
// The "ctx pad merge" command combines entries from an
// external scratchpad file into the current project
// scratchpad. This package provides the file reading,
// key loading, and conflict detection logic.
//
// # Reading External Files
//
// [ReadFileEntries] reads a scratchpad file and
// attempts decryption first. If a non-nil key is
// provided, it tries crypto.Decrypt on the raw bytes.
// On successful decryption, it parses the plaintext
// into entries. If decryption fails (wrong key or
// unencrypted file), it falls back to parsing the raw
// bytes directly. Empty files return nil.
//
// # Key Loading
//
// [LoadKey] loads an encryption key for merge input
// decryption. When keyFile is non-empty, it loads from
// that path. Otherwise it uses store.KeyPath to find
// the project key. If no key is available, it returns
// nil, which causes ReadFileEntries to skip the
// decryption attempt.
//
// # Blob Conflict Detection
//
// [BuildBlobLabelMap] creates a lookup map from blob
// labels to their full entry strings by scanning
// entries with blob.Split.
//
// [HasBlobConflict] checks whether a new entry has the
// same blob label as an existing entry but different
// content. It updates the label map with the new entry
// regardless. The cmd layer uses this to warn users
// about conflicting blob labels before writing.
//
// # Binary Detection
//
// [HasBinaryEntries] scans entries for non-UTF-8 bytes,
// indicating corrupt or binary data that was not
// properly blob-encoded.
package merge
