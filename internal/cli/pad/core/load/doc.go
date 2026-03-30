//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package load orchestrates scratchpad import from files, stdin,
// or directories.
//
// [Lines] reads line-delimited entries from a file or stdin.
// [Blobs] reads directory contents as binary blob entries.
// Both coordinate parsing (via core/imp), storage (via core/store),
// and user output (via write/pad).
package load
