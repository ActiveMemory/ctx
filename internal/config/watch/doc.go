//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package watch defines buffer sizes, parsing
// constants, and default provenance values for the
// ctx watch pipeline.
//
// ctx watch monitors an AI session's live output
// stream, extracting context-update XML tags in
// real time. This package provides the scanner
// buffer configuration, the regex match expectations,
// and the fixed provenance metadata applied to
// machine-generated entries.
//
// # Stream Scanner Buffers
//
//   - [StreamScannerInitCap] (64 KB) — initial
//     buffer capacity for the line scanner.
//   - [StreamScannerMaxSize] (1 MB) — maximum
//     buffer size before the scanner returns an
//     error.
//
// # XML Parsing
//
//   - [ContextUpdateMinGroups] (3) — minimum regex
//     capture groups expected from a context-update
//     match (full match + tag name + content).
//
// # Default Provenance
//
// Watch-originated entries are machine-generated, so
// they receive fixed provenance fields:
//
//   - [ProvenanceSessionID] — default session ID
//     ("watch").
//   - [ProvenanceBranch] — default branch ("watch").
//   - [ProvenanceCommit] — default commit ("watch").
//
// These values let consumers distinguish watch
// entries from human-authored ones.
//
// # Concurrency
//
// All exports are immutable. Safe for any access
// pattern.
package watch
